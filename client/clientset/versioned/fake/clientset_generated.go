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

package fake

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/testing"
	clientset "kubevault.dev/operator/client/clientset/versioned"
	approlev1alpha1 "kubevault.dev/operator/client/clientset/versioned/typed/approle/v1alpha1"
	fakeapprolev1alpha1 "kubevault.dev/operator/client/clientset/versioned/typed/approle/v1alpha1/fake"
	catalogv1alpha1 "kubevault.dev/operator/client/clientset/versioned/typed/catalog/v1alpha1"
	fakecatalogv1alpha1 "kubevault.dev/operator/client/clientset/versioned/typed/catalog/v1alpha1/fake"
	configv1alpha1 "kubevault.dev/operator/client/clientset/versioned/typed/config/v1alpha1"
	fakeconfigv1alpha1 "kubevault.dev/operator/client/clientset/versioned/typed/config/v1alpha1/fake"
	enginev1alpha1 "kubevault.dev/operator/client/clientset/versioned/typed/engine/v1alpha1"
	fakeenginev1alpha1 "kubevault.dev/operator/client/clientset/versioned/typed/engine/v1alpha1/fake"
	kubevaultv1alpha1 "kubevault.dev/operator/client/clientset/versioned/typed/kubevault/v1alpha1"
	fakekubevaultv1alpha1 "kubevault.dev/operator/client/clientset/versioned/typed/kubevault/v1alpha1/fake"
	policyv1alpha1 "kubevault.dev/operator/client/clientset/versioned/typed/policy/v1alpha1"
	fakepolicyv1alpha1 "kubevault.dev/operator/client/clientset/versioned/typed/policy/v1alpha1/fake"
)

// NewSimpleClientset returns a clientset that will respond with the provided objects.
// It's backed by a very simple object tracker that processes creates, updates and deletions as-is,
// without applying any validations and/or defaults. It shouldn't be considered a replacement
// for a real clientset and is mostly useful in simple unit tests.
func NewSimpleClientset(objects ...runtime.Object) *Clientset {
	o := testing.NewObjectTracker(scheme, codecs.UniversalDecoder())
	for _, obj := range objects {
		if err := o.Add(obj); err != nil {
			panic(err)
		}
	}

	cs := &Clientset{tracker: o}
	cs.discovery = &fakediscovery.FakeDiscovery{Fake: &cs.Fake}
	cs.AddReactor("*", "*", testing.ObjectReaction(o))
	cs.AddWatchReactor("*", func(action testing.Action) (handled bool, ret watch.Interface, err error) {
		gvr := action.GetResource()
		ns := action.GetNamespace()
		watch, err := o.Watch(gvr, ns)
		if err != nil {
			return false, nil, err
		}
		return true, watch, nil
	})

	return cs
}

// Clientset implements clientset.Interface. Meant to be embedded into a
// struct to get a default implementation. This makes faking out just the method
// you want to test easier.
type Clientset struct {
	testing.Fake
	discovery *fakediscovery.FakeDiscovery
	tracker   testing.ObjectTracker
}

func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	return c.discovery
}

func (c *Clientset) Tracker() testing.ObjectTracker {
	return c.tracker
}

var _ clientset.Interface = &Clientset{}

// ApproleV1alpha1 retrieves the ApproleV1alpha1Client
func (c *Clientset) ApproleV1alpha1() approlev1alpha1.ApproleV1alpha1Interface {
	return &fakeapprolev1alpha1.FakeApproleV1alpha1{Fake: &c.Fake}
}

// CatalogV1alpha1 retrieves the CatalogV1alpha1Client
func (c *Clientset) CatalogV1alpha1() catalogv1alpha1.CatalogV1alpha1Interface {
	return &fakecatalogv1alpha1.FakeCatalogV1alpha1{Fake: &c.Fake}
}

// ConfigV1alpha1 retrieves the ConfigV1alpha1Client
func (c *Clientset) ConfigV1alpha1() configv1alpha1.ConfigV1alpha1Interface {
	return &fakeconfigv1alpha1.FakeConfigV1alpha1{Fake: &c.Fake}
}

// EngineV1alpha1 retrieves the EngineV1alpha1Client
func (c *Clientset) EngineV1alpha1() enginev1alpha1.EngineV1alpha1Interface {
	return &fakeenginev1alpha1.FakeEngineV1alpha1{Fake: &c.Fake}
}

// KubevaultV1alpha1 retrieves the KubevaultV1alpha1Client
func (c *Clientset) KubevaultV1alpha1() kubevaultv1alpha1.KubevaultV1alpha1Interface {
	return &fakekubevaultv1alpha1.FakeKubevaultV1alpha1{Fake: &c.Fake}
}

// PolicyV1alpha1 retrieves the PolicyV1alpha1Client
func (c *Clientset) PolicyV1alpha1() policyv1alpha1.PolicyV1alpha1Interface {
	return &fakepolicyv1alpha1.FakePolicyV1alpha1{Fake: &c.Fake}
}
