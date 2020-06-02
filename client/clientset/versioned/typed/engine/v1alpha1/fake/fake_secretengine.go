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
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1alpha1 "kubevault.dev/operator/apis/engine/v1alpha1"
)

// FakeSecretEngines implements SecretEngineInterface
type FakeSecretEngines struct {
	Fake *FakeEngineV1alpha1
	ns   string
}

var secretenginesResource = schema.GroupVersionResource{Group: "engine.kubevault.com", Version: "v1alpha1", Resource: "secretengines"}

var secretenginesKind = schema.GroupVersionKind{Group: "engine.kubevault.com", Version: "v1alpha1", Kind: "SecretEngine"}

// Get takes name of the secretEngine, and returns the corresponding secretEngine object, and an error if there is any.
func (c *FakeSecretEngines) Get(name string, options v1.GetOptions) (result *v1alpha1.SecretEngine, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(secretenginesResource, c.ns, name), &v1alpha1.SecretEngine{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SecretEngine), err
}

// List takes label and field selectors, and returns the list of SecretEngines that match those selectors.
func (c *FakeSecretEngines) List(opts v1.ListOptions) (result *v1alpha1.SecretEngineList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(secretenginesResource, secretenginesKind, c.ns, opts), &v1alpha1.SecretEngineList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.SecretEngineList{ListMeta: obj.(*v1alpha1.SecretEngineList).ListMeta}
	for _, item := range obj.(*v1alpha1.SecretEngineList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested secretEngines.
func (c *FakeSecretEngines) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(secretenginesResource, c.ns, opts))

}

// Create takes the representation of a secretEngine and creates it.  Returns the server's representation of the secretEngine, and an error, if there is any.
func (c *FakeSecretEngines) Create(secretEngine *v1alpha1.SecretEngine) (result *v1alpha1.SecretEngine, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(secretenginesResource, c.ns, secretEngine), &v1alpha1.SecretEngine{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SecretEngine), err
}

// Update takes the representation of a secretEngine and updates it. Returns the server's representation of the secretEngine, and an error, if there is any.
func (c *FakeSecretEngines) Update(secretEngine *v1alpha1.SecretEngine) (result *v1alpha1.SecretEngine, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(secretenginesResource, c.ns, secretEngine), &v1alpha1.SecretEngine{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SecretEngine), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeSecretEngines) UpdateStatus(secretEngine *v1alpha1.SecretEngine) (*v1alpha1.SecretEngine, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(secretenginesResource, "status", c.ns, secretEngine), &v1alpha1.SecretEngine{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SecretEngine), err
}

// Delete takes name of the secretEngine and deletes it. Returns an error if one occurs.
func (c *FakeSecretEngines) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(secretenginesResource, c.ns, name), &v1alpha1.SecretEngine{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSecretEngines) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(secretenginesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.SecretEngineList{})
	return err
}

// Patch applies the patch and returns the patched secretEngine.
func (c *FakeSecretEngines) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.SecretEngine, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(secretenginesResource, c.ns, name, pt, data, subresources...), &v1alpha1.SecretEngine{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SecretEngine), err
}
