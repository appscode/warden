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

// Code generated by client-gen. DO NOT EDIT.

package versioned

import (
	"fmt"

	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
	catalogv1alpha1 "kubevault.dev/operator/client/clientset/versioned/typed/catalog/v1alpha1"
	configv1alpha1 "kubevault.dev/operator/client/clientset/versioned/typed/config/v1alpha1"
	enginev1alpha1 "kubevault.dev/operator/client/clientset/versioned/typed/engine/v1alpha1"
	kubevaultv1alpha1 "kubevault.dev/operator/client/clientset/versioned/typed/kubevault/v1alpha1"
	policyv1alpha1 "kubevault.dev/operator/client/clientset/versioned/typed/policy/v1alpha1"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	CatalogV1alpha1() catalogv1alpha1.CatalogV1alpha1Interface
	ConfigV1alpha1() configv1alpha1.ConfigV1alpha1Interface
	EngineV1alpha1() enginev1alpha1.EngineV1alpha1Interface
	KubevaultV1alpha1() kubevaultv1alpha1.KubevaultV1alpha1Interface
	PolicyV1alpha1() policyv1alpha1.PolicyV1alpha1Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	catalogV1alpha1   *catalogv1alpha1.CatalogV1alpha1Client
	configV1alpha1    *configv1alpha1.ConfigV1alpha1Client
	engineV1alpha1    *enginev1alpha1.EngineV1alpha1Client
	kubevaultV1alpha1 *kubevaultv1alpha1.KubevaultV1alpha1Client
	policyV1alpha1    *policyv1alpha1.PolicyV1alpha1Client
}

// CatalogV1alpha1 retrieves the CatalogV1alpha1Client
func (c *Clientset) CatalogV1alpha1() catalogv1alpha1.CatalogV1alpha1Interface {
	return c.catalogV1alpha1
}

// ConfigV1alpha1 retrieves the ConfigV1alpha1Client
func (c *Clientset) ConfigV1alpha1() configv1alpha1.ConfigV1alpha1Interface {
	return c.configV1alpha1
}

// EngineV1alpha1 retrieves the EngineV1alpha1Client
func (c *Clientset) EngineV1alpha1() enginev1alpha1.EngineV1alpha1Interface {
	return c.engineV1alpha1
}

// KubevaultV1alpha1 retrieves the KubevaultV1alpha1Client
func (c *Clientset) KubevaultV1alpha1() kubevaultv1alpha1.KubevaultV1alpha1Interface {
	return c.kubevaultV1alpha1
}

// PolicyV1alpha1 retrieves the PolicyV1alpha1Client
func (c *Clientset) PolicyV1alpha1() policyv1alpha1.PolicyV1alpha1Interface {
	return c.policyV1alpha1
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
// If config's RateLimiter is not set and QPS and Burst are acceptable,
// NewForConfig will generate a rate-limiter in configShallowCopy.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		if configShallowCopy.Burst <= 0 {
			return nil, fmt.Errorf("burst is required to be greater than 0 when RateLimiter is not set and QPS is set to greater than 0")
		}
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.catalogV1alpha1, err = catalogv1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.configV1alpha1, err = configv1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.engineV1alpha1, err = enginev1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.kubevaultV1alpha1, err = kubevaultv1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.policyV1alpha1, err = policyv1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.catalogV1alpha1 = catalogv1alpha1.NewForConfigOrDie(c)
	cs.configV1alpha1 = configv1alpha1.NewForConfigOrDie(c)
	cs.engineV1alpha1 = enginev1alpha1.NewForConfigOrDie(c)
	cs.kubevaultV1alpha1 = kubevaultv1alpha1.NewForConfigOrDie(c)
	cs.policyV1alpha1 = policyv1alpha1.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.catalogV1alpha1 = catalogv1alpha1.New(c)
	cs.configV1alpha1 = configv1alpha1.New(c)
	cs.engineV1alpha1 = enginev1alpha1.New(c)
	cs.kubevaultV1alpha1 = kubevaultv1alpha1.New(c)
	cs.policyV1alpha1 = policyv1alpha1.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
