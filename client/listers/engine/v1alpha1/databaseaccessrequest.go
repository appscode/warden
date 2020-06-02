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

// DatabaseAccessRequestLister helps list DatabaseAccessRequests.
type DatabaseAccessRequestLister interface {
	// List lists all DatabaseAccessRequests in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.DatabaseAccessRequest, err error)
	// DatabaseAccessRequests returns an object that can list and get DatabaseAccessRequests.
	DatabaseAccessRequests(namespace string) DatabaseAccessRequestNamespaceLister
	DatabaseAccessRequestListerExpansion
}

// databaseAccessRequestLister implements the DatabaseAccessRequestLister interface.
type databaseAccessRequestLister struct {
	indexer cache.Indexer
}

// NewDatabaseAccessRequestLister returns a new DatabaseAccessRequestLister.
func NewDatabaseAccessRequestLister(indexer cache.Indexer) DatabaseAccessRequestLister {
	return &databaseAccessRequestLister{indexer: indexer}
}

// List lists all DatabaseAccessRequests in the indexer.
func (s *databaseAccessRequestLister) List(selector labels.Selector) (ret []*v1alpha1.DatabaseAccessRequest, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.DatabaseAccessRequest))
	})
	return ret, err
}

// DatabaseAccessRequests returns an object that can list and get DatabaseAccessRequests.
func (s *databaseAccessRequestLister) DatabaseAccessRequests(namespace string) DatabaseAccessRequestNamespaceLister {
	return databaseAccessRequestNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// DatabaseAccessRequestNamespaceLister helps list and get DatabaseAccessRequests.
type DatabaseAccessRequestNamespaceLister interface {
	// List lists all DatabaseAccessRequests in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.DatabaseAccessRequest, err error)
	// Get retrieves the DatabaseAccessRequest from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.DatabaseAccessRequest, error)
	DatabaseAccessRequestNamespaceListerExpansion
}

// databaseAccessRequestNamespaceLister implements the DatabaseAccessRequestNamespaceLister
// interface.
type databaseAccessRequestNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all DatabaseAccessRequests in the indexer for a given namespace.
func (s databaseAccessRequestNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.DatabaseAccessRequest, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.DatabaseAccessRequest))
	})
	return ret, err
}

// Get retrieves the DatabaseAccessRequest from the indexer for a given namespace and name.
func (s databaseAccessRequestNamespaceLister) Get(name string) (*v1alpha1.DatabaseAccessRequest, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("databaseaccessrequest"), name)
	}
	return obj.(*v1alpha1.DatabaseAccessRequest), nil
}
