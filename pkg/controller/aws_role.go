package controller

import (
	"fmt"
	"time"

	"github.com/appscode/go/encoding/json/types"
	core_util "github.com/appscode/kutil/core/v1"
	meta_util "github.com/appscode/kutil/meta"
	"github.com/appscode/kutil/tools/queue"
	"github.com/golang/glog"
	"github.com/kubevault/operator/apis"
	vsapis "github.com/kubevault/operator/apis"
	api "github.com/kubevault/operator/apis/secretengine/v1alpha1"
	patchutil "github.com/kubevault/operator/client/clientset/versioned/typed/secretengine/v1alpha1/util"
	"github.com/kubevault/operator/pkg/vault/role/aws"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	AWSRolePhaseSuccess    api.AWSRolePhase = "Success"
	AWSRoleConditionFailed                  = "Failed"
	AWSRoleFinalizer                        = "awsrole.secretengine.kubevault.com"
)

func (c *VaultController) initAWSRoleWatcher() {
	c.awsRoleInformer = c.extInformerFactory.Secretengine().V1alpha1().AWSRoles().Informer()
	c.awsRoleQueue = queue.New(api.ResourceKindAWSRole, c.MaxNumRequeues, c.NumThreads, c.runAWSRoleInjector)
	c.awsRoleInformer.AddEventHandler(queue.NewObservableHandler(c.awsRoleQueue.GetQueue(), apis.EnableStatusSubresource))
	c.awsRoleLister = c.extInformerFactory.Secretengine().V1alpha1().AWSRoles().Lister()
}

func (c *VaultController) runAWSRoleInjector(key string) error {
	obj, exist, err := c.awsRoleInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exist {
		glog.Warningf("AWSRole %s does not exist anymore", key)

	} else {
		awsRole := obj.(*api.AWSRole).DeepCopy()

		glog.Infof("Sync/Add/Update for AWSRole %s/%s", awsRole.Namespace, awsRole.Name)

		if awsRole.DeletionTimestamp != nil {
			if core_util.HasFinalizer(awsRole.ObjectMeta, AWSRoleFinalizer) {
				go c.runAWSRoleFinalizer(awsRole, finalizerTimeout, finalizerInterval)
			}
		} else {
			if !core_util.HasFinalizer(awsRole.ObjectMeta, AWSRoleFinalizer) {
				// Add finalizer
				_, _, err := patchutil.PatchAWSRole(c.extClient.SecretengineV1alpha1(), awsRole, func(role *api.AWSRole) *api.AWSRole {
					role.ObjectMeta = core_util.AddFinalizer(role.ObjectMeta, AWSRoleFinalizer)
					return role
				})
				if err != nil {
					return errors.Wrapf(err, "failed to set AWSRole finalizer for %s/%s", awsRole.Namespace, awsRole.Name)
				}
			}

			awsRClient, err := aws.NewAWSRole(c.kubeClient, c.appCatalogClient, awsRole)
			if err != nil {
				return err
			}

			err = c.reconcileAWSRole(awsRClient, awsRole)
			if err != nil {
				return errors.Wrapf(err, "for AWSRole %s/%s:", awsRole.Namespace, awsRole.Name)
			}
		}
	}
	return nil
}

// Will do:
//	For vault:
//	  - enable the aws secrets engine if it is not already enabled
//	  - configure Vault AWS secret engine
// 	  - configure a AWS role
//    - sync role
func (c *VaultController) reconcileAWSRole(awsRClient aws.AWSRoleInterface, awsRole *api.AWSRole) error {
	status := awsRole.Status
	// enable the aws secrets engine if it is not already enabled
	err := awsRClient.EnableAWS()
	if err != nil {
		status.Conditions = []api.AWSRoleCondition{
			{
				Type:    AWSRoleConditionFailed,
				Status:  corev1.ConditionTrue,
				Reason:  "FailedToEnableAWS",
				Message: err.Error(),
			},
		}

		err2 := c.updatedAWSRoleStatus(&status, awsRole)
		if err2 != nil {
			return errors.Wrap(err2, "failed to update status")
		}
		return errors.Wrap(err, "failed to enable database secret engine")
	}

	// create aws config
	err = awsRClient.CreateConfig()
	if err != nil {
		status.Conditions = []api.AWSRoleCondition{
			{
				Type:    AWSRoleConditionFailed,
				Status:  corev1.ConditionTrue,
				Reason:  "FailedToCreateAWSConfig",
				Message: err.Error(),
			},
		}

		err2 := c.updatedAWSRoleStatus(&status, awsRole)
		if err2 != nil {
			return errors.Wrap(err2, "failed to update status")
		}
		return errors.Wrap(err, "failed to create database connection config")
	}

	// create role
	err = awsRClient.CreateRole()
	if err != nil {
		status.Conditions = []api.AWSRoleCondition{
			{
				Type:    AWSRoleConditionFailed,
				Status:  corev1.ConditionTrue,
				Reason:  "FailedToCreateRole",
				Message: err.Error(),
			},
		}

		err2 := c.updatedAWSRoleStatus(&status, awsRole)
		if err2 != nil {
			return errors.Wrap(err2, "failed to update status")
		}
		return errors.Wrap(err, "failed to create role")
	}

	status.Conditions = []api.AWSRoleCondition{}
	status.Phase = AWSRolePhaseSuccess
	status.ObservedGeneration = types.NewIntHash(awsRole.Generation, meta_util.GenerationHash(awsRole))

	err = c.updatedAWSRoleStatus(&status, awsRole)
	if err != nil {
		return errors.Wrapf(err, "failed to update AWSRole status")
	}
	return nil
}

func (c *VaultController) updatedAWSRoleStatus(status *api.AWSRoleStatus, awsRole *api.AWSRole) error {
	_, err := patchutil.UpdateAWSRoleStatus(c.extClient.SecretengineV1alpha1(), awsRole, func(s *api.AWSRoleStatus) *api.AWSRoleStatus {
		s = status
		return s
	}, vsapis.EnableStatusSubresource)
	return err
}

func (c *VaultController) runAWSRoleFinalizer(awsRole *api.AWSRole, timeout time.Duration, interval time.Duration) {
	if awsRole == nil {
		glog.Infoln("AWSRole is nil")
		return
	}

	id := getAWSRoleId(awsRole)
	if c.finalizerInfo.IsAlreadyProcessing(id) {
		// already processing
		return
	}

	glog.Infof("Processing finalizer for AWSRole %s/%s", awsRole.Namespace, awsRole.Name)
	// Add key to finalizerInfo, it will prevent other go routine to processing for this AWSRole
	c.finalizerInfo.Add(id)

	stopCh := time.After(timeout)
	finalizationDone := false
	timeOutOccured := false
	attempt := 0

	for {
		glog.Infof("AWSRole %s/%s finalizer: attempt %d\n", awsRole.Namespace, awsRole.Name, attempt)

		select {
		case <-stopCh:
			timeOutOccured = true
		default:
		}

		if timeOutOccured {
			break
		}

		if !finalizationDone {
			d, err := aws.NewAWSRole(c.kubeClient, c.appCatalogClient, awsRole)
			if err != nil {
				glog.Errorf("AWSRole %s/%s finalizer: %v", awsRole.Namespace, awsRole.Name, err)
			} else {
				err = c.finalizeAWSRole(d, awsRole)
				if err != nil {
					glog.Errorf("AWSRole %s/%s finalizer: %v", awsRole.Namespace, awsRole.Name, err)
				} else {
					finalizationDone = true
				}
			}
		}

		if finalizationDone {
			err := c.removeAWSRoleFinalizer(awsRole)
			if err != nil {
				glog.Errorf("AWSRole %s/%s finalizer: removing finalizer %v", awsRole.Namespace, awsRole.Name, err)
			} else {
				break
			}
		}

		select {
		case <-stopCh:
			timeOutOccured = true
		case <-time.After(interval):
		}
		attempt++
	}

	err := c.removeAWSRoleFinalizer(awsRole)
	if err != nil {
		glog.Errorf("AWSRole %s/%s finalizer: removing finalizer %v", awsRole.Namespace, awsRole.Name, err)
	} else {
		glog.Infof("Removed finalizer for AWSRole %s/%s", awsRole.Namespace, awsRole.Name)
	}

	// Delete key from finalizer info as processing is done
	c.finalizerInfo.Delete(id)
}

// Do:
//	- delete role in vault
func (c *VaultController) finalizeAWSRole(awsRClient aws.AWSRoleInterface, awsRole *api.AWSRole) error {
	err := awsRClient.DeleteRole(awsRole.RoleName())
	if err != nil {
		return errors.Wrap(err, "failed to database role")
	}
	return nil
}

func (c *VaultController) removeAWSRoleFinalizer(awsRole *api.AWSRole) error {
	m, err := c.extClient.SecretengineV1alpha1().AWSRoles(awsRole.Namespace).Get(awsRole.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	// remove finalizer
	_, _, err = patchutil.PatchAWSRole(c.extClient.SecretengineV1alpha1(), m, func(role *api.AWSRole) *api.AWSRole {
		role.ObjectMeta = core_util.RemoveFinalizer(role.ObjectMeta, AWSRoleFinalizer)
		return role
	})
	return err
}

func getAWSRoleId(awsRole *api.AWSRole) string {
	return fmt.Sprintf("%s/%s/%s", api.ResourceAWSRole, awsRole.Namespace, awsRole.Name)
}
