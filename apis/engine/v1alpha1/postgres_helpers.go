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

package v1alpha1

import (
	"fmt"

	"kubevault.dev/operator/api/crds"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/yaml"
)

func (_ PostgresRole) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	data := crds.MustAsset("engine.kubevault.com_postgresroles.yaml")
	var out apiextensions.CustomResourceDefinition
	utilruntime.Must(yaml.Unmarshal(data, &out))
	return &out
}

const DefaultPostgresDatabasePlugin = "postgresql-database-plugin"

func (r PostgresRole) RoleName() string {
	cluster := "-"
	if r.ClusterName != "" {
		cluster = r.ClusterName
	}
	return fmt.Sprintf("k8s.%s.%s.%s", cluster, r.Namespace, r.Name)
}

func (r PostgresRole) IsValid() error {
	return nil
}

func (p *PostgresConfiguration) SetDefaults() {
	if p == nil {
		return
	}

	// If user doesn't specify the list of AllowedRoles
	// It is set to "*" (allow all)
	if p.AllowedRoles == nil || len(p.AllowedRoles) == 0 {
		p.AllowedRoles = []string{"*"}
	}

	if p.PluginName == "" {
		p.PluginName = DefaultPostgresDatabasePlugin
	}
}
