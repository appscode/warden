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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/kubevault/operator/client/clientset/versioned/typed/engine/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeEngineV1alpha1 struct {
	*testing.Fake
}

func (c *FakeEngineV1alpha1) AWSAccessKeyRequests(namespace string) v1alpha1.AWSAccessKeyRequestInterface {
	return &FakeAWSAccessKeyRequests{c, namespace}
}

func (c *FakeEngineV1alpha1) AWSRoles(namespace string) v1alpha1.AWSRoleInterface {
	return &FakeAWSRoles{c, namespace}
}

func (c *FakeEngineV1alpha1) GCPRoles(namespace string) v1alpha1.GCPRoleInterface {
	return &FakeGCPRoles{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeEngineV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
