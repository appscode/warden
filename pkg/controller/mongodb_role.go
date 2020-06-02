/*
Copyright The KubeVault Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"time"

	"kubevault.dev/operator/apis"
	api "kubevault.dev/operator/apis/engine/v1alpha1"
	patchutil "kubevault.dev/operator/client/clientset/versioned/typed/engine/v1alpha1/util"
	"kubevault.dev/operator/pkg/vault/role/database"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	kmapi "kmodules.xyz/client-go/api/v1"
	core_util "kmodules.xyz/client-go/core/v1"
	"kmodules.xyz/client-go/tools/queue"
)

const (
	MongoDBRolePhaseSuccess api.MongoDBRolePhase = "Success"
	finalizerInterval                            = 5 * time.Second
	finalizerTimeout                             = 30 * time.Second
)

func (c *VaultController) initMongoDBRoleWatcher() {
	c.mgRoleInformer = c.extInformerFactory.Engine().V1alpha1().MongoDBRoles().Informer()
	c.mgRoleQueue = queue.New(api.ResourceKindMongoDBRole, c.MaxNumRequeues, c.NumThreads, c.runMongoDBRoleInjector)
	c.mgRoleInformer.AddEventHandler(queue.NewReconcilableHandler(c.mgRoleQueue.GetQueue()))
	c.mgRoleLister = c.extInformerFactory.Engine().V1alpha1().MongoDBRoles().Lister()
}

func (c *VaultController) runMongoDBRoleInjector(key string) error {
	obj, exist, err := c.mgRoleInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exist {
		glog.Warningf("MongoDBRole %s does not exist anymore", key)

	} else {
		role := obj.(*api.MongoDBRole).DeepCopy()

		glog.Infof("Sync/Add/Update for MongoDBRole %s/%s", role.Namespace, role.Name)

		if role.DeletionTimestamp != nil {
			if core_util.HasFinalizer(role.ObjectMeta, apis.Finalizer) {
				go c.runMongoDBRoleFinalizer(role, finalizerTimeout, finalizerInterval)
			}
		} else {
			if !core_util.HasFinalizer(role.ObjectMeta, apis.Finalizer) {
				// Add finalizer
				_, _, err := patchutil.PatchMongoDBRole(context.TODO(), c.extClient.EngineV1alpha1(), role, func(role *api.MongoDBRole) *api.MongoDBRole {
					role.ObjectMeta = core_util.AddFinalizer(role.ObjectMeta, apis.Finalizer)
					return role
				}, metav1.PatchOptions{})
				if err != nil {
					return errors.Wrapf(err, "failed to set MongoDBRole finalizer for %s/%s", role.Namespace, role.Name)
				}
			}

			dbRClient, err := database.NewDatabaseRoleForMongodb(c.kubeClient, c.appCatalogClient, role)
			if err != nil {
				return err
			}

			err = c.reconcileMongoDBRole(dbRClient, role)
			if err != nil {
				return errors.Wrapf(err, "for MongoDBRole %s/%s:", role.Namespace, role.Name)
			}
		}
	}
	return nil
}

// Will do:
//	For vault:
// 	  - configure a role that maps a name in Vault to an SQL statement to execute to create the database credential.
//    - sync role
//	  - revoke previous lease of all the respective mongodbRoleBinding and reissue a new lease
func (c *VaultController) reconcileMongoDBRole(dbRClient database.DatabaseRoleInterface, role *api.MongoDBRole) error {
	// create role
	err := dbRClient.CreateRole()
	if err != nil {
		_, err2 := patchutil.UpdateMongoDBRoleStatus(
			context.TODO(),
			c.extClient.EngineV1alpha1(),
			role.ObjectMeta,
			func(status *api.MongoDBRoleStatus) *api.MongoDBRoleStatus {
				status.Conditions = kmapi.SetCondition(status.Conditions, kmapi.Condition{
					Type:    kmapi.ConditionFailure,
					Status:  kmapi.ConditionTrue,
					Reason:  "FailedToCreateRole",
					Message: err.Error(),
				})
				return status
			},
			metav1.UpdateOptions{},
		)
		return utilerrors.NewAggregate([]error{err2, errors.Wrap(err, "failed to create role")})
	}

	_, err = patchutil.UpdateMongoDBRoleStatus(
		context.TODO(),
		c.extClient.EngineV1alpha1(),
		role.ObjectMeta,
		func(status *api.MongoDBRoleStatus) *api.MongoDBRoleStatus {
			status.Phase = MongoDBRolePhaseSuccess
			status.ObservedGeneration = role.Generation
			status.Conditions = kmapi.SetCondition(status.Conditions, kmapi.Condition{
				Type:    kmapi.ConditionAvailable,
				Status:  kmapi.ConditionTrue,
				Reason:  "Provisioned",
				Message: "role is ready to use",
			})
			return status
		},
		metav1.UpdateOptions{},
	)
	return err
}

func (c *VaultController) runMongoDBRoleFinalizer(role *api.MongoDBRole, timeout time.Duration, interval time.Duration) {
	if role == nil {
		glog.Infoln("MongoDBRole is nil")
		return
	}

	id := getMongoDBRoleId(role)
	if c.finalizerInfo.IsAlreadyProcessing(id) {
		// already processing
		return
	}

	glog.Infof("Processing finalizer for MongoDBRole %s/%s", role.Namespace, role.Name)
	// Add key to finalizerInfo, it will prevent other go routine to processing for this MongoDBRole
	c.finalizerInfo.Add(id)

	stopCh := time.After(timeout)
	finalizationDone := false
	timeOutOccured := false
	attempt := 0

	for {
		glog.Infof("MongoDBRole %s/%s finalizer: attempt %d\n", role.Namespace, role.Name, attempt)

		select {
		case <-stopCh:
			timeOutOccured = true
		default:
		}

		if timeOutOccured {
			break
		}

		if !finalizationDone {
			d, err := database.NewDatabaseRoleForMongodb(c.kubeClient, c.appCatalogClient, role)
			if err != nil {
				glog.Errorf("MongoDBRole %s/%s finalizer: %v", role.Namespace, role.Name, err)
			} else {
				err = c.finalizeMongoDBRole(d, role)
				if err != nil {
					glog.Errorf("MongoDBRole %s/%s finalizer: %v", role.Namespace, role.Name, err)
				} else {
					finalizationDone = true
				}
			}
		}

		if finalizationDone {
			err := c.removeMongoDBRoleFinalizer(role)
			if err != nil {
				glog.Errorf("MongoDBRole %s/%s finalizer: removing finalizer %v", role.Namespace, role.Name, err)
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

	err := c.removeMongoDBRoleFinalizer(role)
	if err != nil {
		glog.Errorf("MongoDBRole %s/%s finalizer: removing finalizer %v", role.Namespace, role.Name, err)
	} else {
		glog.Infof("Removed finalizer for MongoDBRole %s/%s", role.Namespace, role.Name)
	}

	// Delete key from finalizer info as processing is done
	c.finalizerInfo.Delete(id)
}

// Do:
//	- delete role in vault
//	- revoke lease of all the corresponding mongodbRoleBinding
func (c *VaultController) finalizeMongoDBRole(dbRClient database.DatabaseRoleInterface, role *api.MongoDBRole) error {
	err := dbRClient.DeleteRole(role.RoleName())
	if err != nil {
		return errors.Wrap(err, "failed to database role")
	}
	return nil
}

func (c *VaultController) removeMongoDBRoleFinalizer(mRole *api.MongoDBRole) error {
	m, err := c.extClient.EngineV1alpha1().MongoDBRoles(mRole.Namespace).Get(context.TODO(), mRole.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	// remove finalizer
	_, _, err = patchutil.PatchMongoDBRole(context.TODO(), c.extClient.EngineV1alpha1(), m, func(role *api.MongoDBRole) *api.MongoDBRole {
		role.ObjectMeta = core_util.RemoveFinalizer(role.ObjectMeta, apis.Finalizer)
		return role
	}, metav1.PatchOptions{})
	return err
}

func getMongoDBRoleId(mRole *api.MongoDBRole) string {
	return fmt.Sprintf("%s/%s/%s", api.ResourceMongoDBRole, mRole.Namespace, mRole.Name)
}
