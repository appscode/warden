package controller

import (
	"fmt"

	reg_util "github.com/appscode/kutil/admissionregistration/v1beta1"
	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	"github.com/appscode/kutil/tools/queue"
	"github.com/golang/glog"
	catalogapi "github.com/kubevault/operator/apis/catalog/v1alpha1"
	vaultapi "github.com/kubevault/operator/apis/kubevault/v1alpha1"
	policyapi "github.com/kubevault/operator/apis/policy/v1alpha1"
	cs "github.com/kubevault/operator/client/clientset/versioned"
	vaultinformers "github.com/kubevault/operator/client/informers/externalversions"
	vault_listers "github.com/kubevault/operator/client/listers/kubevault/v1alpha1"
	policy_listers "github.com/kubevault/operator/client/listers/policy/v1alpha1"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	appcat_cs "kmodules.xyz/custom-resources/client/clientset/versioned/typed/appcatalog/v1alpha1"
)

type VaultController struct {
	config
	clientConfig *rest.Config

	// ctxCancels stores vault clusters' contexts that are used to
	// cancel their goroutines when they are deleted
	ctxCancels map[string]CtxWithCancel

	kubeClient       kubernetes.Interface
	extClient        cs.Interface
	appCatalogClient appcat_cs.AppcatalogV1alpha1Interface
	crdClient        crd_cs.ApiextensionsV1beta1Interface
	recorder         record.EventRecorder

	kubeInformerFactory informers.SharedInformerFactory
	extInformerFactory  vaultinformers.SharedInformerFactory

	// for VaultServer
	vsQueue    *queue.Worker
	vsInformer cache.SharedIndexInformer
	vsLister   vault_listers.VaultServerLister

	// for VaultPolicy
	vplcyQueue    *queue.Worker
	vplcyInformer cache.SharedIndexInformer
	vplcyLister   policy_listers.VaultPolicyLister

	// for VaultPolicyBinding
	vplcyBindingQueue    *queue.Worker
	vplcyBindingInformer cache.SharedIndexInformer
	vplcyBindingLister   policy_listers.VaultPolicyBindingLister

	// Contain the currently processing finalizer
	finalizerInfo *mapFinalizer

	// authMethodCtx stores auth method controller contexts that are used to
	// cancel their goroutines when they are not needed
	authMethodCtx map[string]CtxWithCancel
}

func (c *VaultController) ensureCustomResourceDefinitions() error {
	crds := []*crd_api.CustomResourceDefinition{
		vaultapi.VaultServer{}.CustomResourceDefinition(),
		catalogapi.VaultServerVersion{}.CustomResourceDefinition(),
		policyapi.VaultPolicy{}.CustomResourceDefinition(),
		policyapi.VaultPolicyBinding{}.CustomResourceDefinition(),
		appcat.AppBinding{}.CustomResourceDefinition(),
	}
	return crdutils.RegisterCRDs(c.crdClient, crds)
}

func (c *VaultController) Run(stopCh <-chan struct{}) {
	go c.RunInformers(stopCh)

	cancel, _ := reg_util.SyncValidatingWebhookCABundle(c.clientConfig, validatingWebhook)

	<-stopCh

	cancel()
}

func (c *VaultController) RunInformers(stopCh <-chan struct{}) {
	defer runtime.HandleCrash()

	glog.Info("Starting Vault controller")

	c.extInformerFactory.Start(stopCh)

	for _, v := range c.extInformerFactory.WaitForCacheSync(stopCh) {
		if !v {
			runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
			return
		}
	}

	// For VaultServer
	go c.vsQueue.Run(stopCh)

	//For VaultPolicy
	go c.vplcyQueue.Run(stopCh)

	//For VaultPolicyBinding
	go c.vplcyBindingQueue.Run(stopCh)

	<-stopCh
	glog.Info("Stopping Vault operator")
}
