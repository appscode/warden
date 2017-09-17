package controller

import (
	"fmt"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/pkg/api/v1"
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
	err := c.syncSToStdout(key.(string))
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
func (c *VaultController) syncSToStdout(key string) error {
	obj, exists, err := c.sIndexer.GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		// Below we will warm up our cache with a Secret, so that we will see a delete for one secret
		fmt.Printf("Secret %s does not exist anymore\n", key)
	} else {
		// Note that you also have to check the uid if you have a local controlled resource, which
		// is dependent on the actual instance, to detect that a Secret was recreated with the same name
		fmt.Printf("Sync/Add/Update for Secret %s\n", obj.(*v1.Secret).GetName())
	}
	return nil
}
