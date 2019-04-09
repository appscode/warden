/*
Copyright 2019 The Kube Vault Authors.

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
	v1alpha1 "github.com/kubevault/operator/apis/engine/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// GCPAccessKeyRequestLister helps list GCPAccessKeyRequests.
type GCPAccessKeyRequestLister interface {
	// List lists all GCPAccessKeyRequests in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.GCPAccessKeyRequest, err error)
	// GCPAccessKeyRequests returns an object that can list and get GCPAccessKeyRequests.
	GCPAccessKeyRequests(namespace string) GCPAccessKeyRequestNamespaceLister
	GCPAccessKeyRequestListerExpansion
}

// gCPAccessKeyRequestLister implements the GCPAccessKeyRequestLister interface.
type gCPAccessKeyRequestLister struct {
	indexer cache.Indexer
}

// NewGCPAccessKeyRequestLister returns a new GCPAccessKeyRequestLister.
func NewGCPAccessKeyRequestLister(indexer cache.Indexer) GCPAccessKeyRequestLister {
	return &gCPAccessKeyRequestLister{indexer: indexer}
}

// List lists all GCPAccessKeyRequests in the indexer.
func (s *gCPAccessKeyRequestLister) List(selector labels.Selector) (ret []*v1alpha1.GCPAccessKeyRequest, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.GCPAccessKeyRequest))
	})
	return ret, err
}

// GCPAccessKeyRequests returns an object that can list and get GCPAccessKeyRequests.
func (s *gCPAccessKeyRequestLister) GCPAccessKeyRequests(namespace string) GCPAccessKeyRequestNamespaceLister {
	return gCPAccessKeyRequestNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// GCPAccessKeyRequestNamespaceLister helps list and get GCPAccessKeyRequests.
type GCPAccessKeyRequestNamespaceLister interface {
	// List lists all GCPAccessKeyRequests in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.GCPAccessKeyRequest, err error)
	// Get retrieves the GCPAccessKeyRequest from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.GCPAccessKeyRequest, error)
	GCPAccessKeyRequestNamespaceListerExpansion
}

// gCPAccessKeyRequestNamespaceLister implements the GCPAccessKeyRequestNamespaceLister
// interface.
type gCPAccessKeyRequestNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all GCPAccessKeyRequests in the indexer for a given namespace.
func (s gCPAccessKeyRequestNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.GCPAccessKeyRequest, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.GCPAccessKeyRequest))
	})
	return ret, err
}

// Get retrieves the GCPAccessKeyRequest from the indexer for a given namespace and name.
func (s gCPAccessKeyRequestNamespaceLister) Get(name string) (*v1alpha1.GCPAccessKeyRequest, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("gcpaccesskeyrequest"), name)
	}
	return obj.(*v1alpha1.GCPAccessKeyRequest), nil
}
