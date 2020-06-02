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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1alpha1 "kubevault.dev/operator/apis/engine/v1alpha1"
)

// AWSAccessKeyRequestLister helps list AWSAccessKeyRequests.
type AWSAccessKeyRequestLister interface {
	// List lists all AWSAccessKeyRequests in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.AWSAccessKeyRequest, err error)
	// AWSAccessKeyRequests returns an object that can list and get AWSAccessKeyRequests.
	AWSAccessKeyRequests(namespace string) AWSAccessKeyRequestNamespaceLister
	AWSAccessKeyRequestListerExpansion
}

// aWSAccessKeyRequestLister implements the AWSAccessKeyRequestLister interface.
type aWSAccessKeyRequestLister struct {
	indexer cache.Indexer
}

// NewAWSAccessKeyRequestLister returns a new AWSAccessKeyRequestLister.
func NewAWSAccessKeyRequestLister(indexer cache.Indexer) AWSAccessKeyRequestLister {
	return &aWSAccessKeyRequestLister{indexer: indexer}
}

// List lists all AWSAccessKeyRequests in the indexer.
func (s *aWSAccessKeyRequestLister) List(selector labels.Selector) (ret []*v1alpha1.AWSAccessKeyRequest, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.AWSAccessKeyRequest))
	})
	return ret, err
}

// AWSAccessKeyRequests returns an object that can list and get AWSAccessKeyRequests.
func (s *aWSAccessKeyRequestLister) AWSAccessKeyRequests(namespace string) AWSAccessKeyRequestNamespaceLister {
	return aWSAccessKeyRequestNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// AWSAccessKeyRequestNamespaceLister helps list and get AWSAccessKeyRequests.
type AWSAccessKeyRequestNamespaceLister interface {
	// List lists all AWSAccessKeyRequests in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.AWSAccessKeyRequest, err error)
	// Get retrieves the AWSAccessKeyRequest from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.AWSAccessKeyRequest, error)
	AWSAccessKeyRequestNamespaceListerExpansion
}

// aWSAccessKeyRequestNamespaceLister implements the AWSAccessKeyRequestNamespaceLister
// interface.
type aWSAccessKeyRequestNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all AWSAccessKeyRequests in the indexer for a given namespace.
func (s aWSAccessKeyRequestNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.AWSAccessKeyRequest, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.AWSAccessKeyRequest))
	})
	return ret, err
}

// Get retrieves the AWSAccessKeyRequest from the indexer for a given namespace and name.
func (s aWSAccessKeyRequestNamespaceLister) Get(name string) (*v1alpha1.AWSAccessKeyRequest, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("awsaccesskeyrequest"), name)
	}
	return obj.(*v1alpha1.AWSAccessKeyRequest), nil
}
