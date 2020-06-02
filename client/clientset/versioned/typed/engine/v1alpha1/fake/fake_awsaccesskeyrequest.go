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

// FakeAWSAccessKeyRequests implements AWSAccessKeyRequestInterface
type FakeAWSAccessKeyRequests struct {
	Fake *FakeEngineV1alpha1
	ns   string
}

var awsaccesskeyrequestsResource = schema.GroupVersionResource{Group: "engine.kubevault.com", Version: "v1alpha1", Resource: "awsaccesskeyrequests"}

var awsaccesskeyrequestsKind = schema.GroupVersionKind{Group: "engine.kubevault.com", Version: "v1alpha1", Kind: "AWSAccessKeyRequest"}

// Get takes name of the aWSAccessKeyRequest, and returns the corresponding aWSAccessKeyRequest object, and an error if there is any.
func (c *FakeAWSAccessKeyRequests) Get(name string, options v1.GetOptions) (result *v1alpha1.AWSAccessKeyRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(awsaccesskeyrequestsResource, c.ns, name), &v1alpha1.AWSAccessKeyRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AWSAccessKeyRequest), err
}

// List takes label and field selectors, and returns the list of AWSAccessKeyRequests that match those selectors.
func (c *FakeAWSAccessKeyRequests) List(opts v1.ListOptions) (result *v1alpha1.AWSAccessKeyRequestList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(awsaccesskeyrequestsResource, awsaccesskeyrequestsKind, c.ns, opts), &v1alpha1.AWSAccessKeyRequestList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.AWSAccessKeyRequestList{ListMeta: obj.(*v1alpha1.AWSAccessKeyRequestList).ListMeta}
	for _, item := range obj.(*v1alpha1.AWSAccessKeyRequestList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested aWSAccessKeyRequests.
func (c *FakeAWSAccessKeyRequests) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(awsaccesskeyrequestsResource, c.ns, opts))

}

// Create takes the representation of a aWSAccessKeyRequest and creates it.  Returns the server's representation of the aWSAccessKeyRequest, and an error, if there is any.
func (c *FakeAWSAccessKeyRequests) Create(aWSAccessKeyRequest *v1alpha1.AWSAccessKeyRequest) (result *v1alpha1.AWSAccessKeyRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(awsaccesskeyrequestsResource, c.ns, aWSAccessKeyRequest), &v1alpha1.AWSAccessKeyRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AWSAccessKeyRequest), err
}

// Update takes the representation of a aWSAccessKeyRequest and updates it. Returns the server's representation of the aWSAccessKeyRequest, and an error, if there is any.
func (c *FakeAWSAccessKeyRequests) Update(aWSAccessKeyRequest *v1alpha1.AWSAccessKeyRequest) (result *v1alpha1.AWSAccessKeyRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(awsaccesskeyrequestsResource, c.ns, aWSAccessKeyRequest), &v1alpha1.AWSAccessKeyRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AWSAccessKeyRequest), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeAWSAccessKeyRequests) UpdateStatus(aWSAccessKeyRequest *v1alpha1.AWSAccessKeyRequest) (*v1alpha1.AWSAccessKeyRequest, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(awsaccesskeyrequestsResource, "status", c.ns, aWSAccessKeyRequest), &v1alpha1.AWSAccessKeyRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AWSAccessKeyRequest), err
}

// Delete takes name of the aWSAccessKeyRequest and deletes it. Returns an error if one occurs.
func (c *FakeAWSAccessKeyRequests) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(awsaccesskeyrequestsResource, c.ns, name), &v1alpha1.AWSAccessKeyRequest{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAWSAccessKeyRequests) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(awsaccesskeyrequestsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.AWSAccessKeyRequestList{})
	return err
}

// Patch applies the patch and returns the patched aWSAccessKeyRequest.
func (c *FakeAWSAccessKeyRequests) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.AWSAccessKeyRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(awsaccesskeyrequestsResource, c.ns, name, pt, data, subresources...), &v1alpha1.AWSAccessKeyRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AWSAccessKeyRequest), err
}
