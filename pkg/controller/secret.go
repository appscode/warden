package controller

import (
	"fmt"
	"path"
	"time"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

func (c *VaultController) runSecretWatcher() {
	for c.processNextSecret() {
	}
}

func (c *VaultController) processNextSecret() bool {
	// Wait until there is a new item in the working queue
	key, quit := c.sQueue.Get()
	if quit {
		return false
	}
	// Tell the queue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two secrets with the same key are never processed in
	// parallel.
	defer c.sQueue.Done(key)

	// Invoke the method containing the business logic
	err := c.syncSecretToVault(key.(string))
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		c.sQueue.Forget(key)
		return true
	}

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if c.sQueue.NumRequeues(key) < c.options.MaxNumRequeues {
		glog.Infof("Error syncing secret %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		c.sQueue.AddRateLimited(key)
		return true
	}

	c.sQueue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
	glog.Infof("Dropping secret %q out of the queue: %v", key, err)
	return true
}

// syncToStdout is the business logic of the controller. In this controller it simply prints
// information about the secret to stdout. In case an error happened, it has to simply return the error.
// The retry logic should not be part of the business logic.
func (c *VaultController) syncSecretToVault(key string) error {
	obj, exists, err := c.sIndexer.GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}
	if !exists {
		// Below we will warm up our cache with a Secret, so that we will see a delete for one secret
		fmt.Printf("Secret %s does not exist anymore\n", key)

		namespace, name, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			return err
		}
		c.vaultClient.Logical().Delete(path.Join(c.options.SecretBackend(), namespace, name))
	} else {
		secret := obj.(*v1.Secret)
		fmt.Printf("Sync/Add/Update for Secret %s\n", secret.GetName())

		data := map[string]interface{}{}
		for k, v := range secret.Data {
			data[k] = string(v)
		}
		data["m:uid"] = secret.ObjectMeta.UID
		data["m:resourceVersion"] = secret.ObjectMeta.ResourceVersion
		data["m:generation"] = secret.ObjectMeta.Generation
		data["m:creationTimestamp"] = secret.ObjectMeta.CreationTimestamp.UTC().Format(time.RFC3339)
		if secret.ObjectMeta.DeletionTimestamp != nil {
			data["m:deletionTimestamp"] = secret.ObjectMeta.DeletionTimestamp.UTC().Format(time.RFC3339)
		}
		for k, v := range secret.ObjectMeta.Labels {
			data["l:"+k] = v
		}
		for k, v := range secret.ObjectMeta.Annotations {
			data["a:"+k] = v
		}
		_, err = c.vaultClient.Logical().Write(path.Join(c.options.SecretBackend(), secret.Namespace, secret.Name), data)
		return err
	}
	return nil
}
