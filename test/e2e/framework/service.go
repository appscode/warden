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

package framework

import (
	"context"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	meta_util "kmodules.xyz/client-go/meta"
)

func (f *Framework) CreateService(obj core.Service) error {
	_, err := f.KubeClient.CoreV1().Services(obj.Namespace).Create(context.TODO(), &obj, metav1.CreateOptions{})
	return err
}

func (f *Framework) DeleteService(name, namespace string) error {
	return f.KubeClient.CoreV1().Services(namespace).Delete(context.TODO(), name, meta_util.DeleteInForeground())
}
