package controller

import (
	"fmt"
	"testing"

	api "github.com/kubevault/operator/apis/kubevault/v1alpha1"
	cfake "github.com/kubevault/operator/client/clientset/versioned/fake"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kfake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/record"
	appcatfake "kmodules.xyz/custom-resources/client/clientset/versioned/fake"
)

type vaultFake struct {
	sr                *core.Secret
	cm                *core.ConfigMap
	dp                *appsv1.Deployment
	sa                *core.ServiceAccount
	svc               *core.Service
	roles             []rbac.Role
	pt                *core.PodTemplateSpec
	cnt               core.Container
	ErrInGetServerTLS bool
	ErrInGetConfig    bool
	ErrInApply        bool
}

func (v *vaultFake) GetServerTLS() (*core.Secret, error) {
	if v.ErrInGetServerTLS {
		return nil, fmt.Errorf("error")
	}
	return v.sr, nil
}
func (v *vaultFake) GetConfig() (*core.ConfigMap, error) {
	if v.ErrInGetConfig {
		return nil, fmt.Errorf("error")
	}
	return v.cm, nil
}
func (v *vaultFake) Apply(pt *core.PodTemplateSpec) error {
	if v.ErrInApply {
		return fmt.Errorf("error")
	}
	return nil
}
func (v *vaultFake) GetService() *core.Service {
	return v.svc
}
func (v *vaultFake) GetDeployment(pt *core.PodTemplateSpec) *appsv1.Deployment {
	return v.dp
}
func (v *vaultFake) GetServiceAccounts() []core.ServiceAccount {
	return []core.ServiceAccount{*v.sa}
}
func (v *vaultFake) GetRBACRolesAndRoleBindings() ([]rbac.Role, []rbac.RoleBinding) {
	return v.roles, nil
}
func (v *vaultFake) GetRBACClusterRoleBinding() rbac.ClusterRoleBinding {
	return rbac.ClusterRoleBinding{}
}
func (v *vaultFake) GetPodTemplate(c core.Container, saName string) *core.PodTemplateSpec {
	return v.pt
}
func (v *vaultFake) GetContainer() core.Container {
	return v.cnt
}

func TestReconcileVault(t *testing.T) {
	vfk := vaultFake{
		sr: &core.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sr-test",
				Namespace: "test",
			},
			Data: map[string][]byte{
				"ca.crt":     []byte("ca"),
				"server.crt": []byte("srv"),
				"server.key": []byte("srv"),
			},
		},
		cm: &core.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sa-test",
				Namespace: "test",
			},
		},
		sa: &core.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sa-test",
				Namespace: "test",
			},
		},
		cnt: core.Container{},
		pt:  &core.PodTemplateSpec{},
		dp: &appsv1.Deployment{
			ObjectMeta: getVaultObjectMeta(1),
		},
		svc: &core.Service{
			ObjectMeta: getVaultObjectMeta(1),
		},
		roles: []rbac.Role{},
	}

	testData := []struct {
		name      string
		vs        *api.VaultServer
		vfake     *vaultFake
		expectErr bool
	}{
		{
			name: "no error",
			vs: &api.VaultServer{
				ObjectMeta: getVaultObjectMeta(1),
			},

			vfake:     &vfk,
			expectErr: false,
		},
		{
			name: "failed to create vault tls",
			vs: &api.VaultServer{
				ObjectMeta: getVaultObjectMeta(1),
			},

			vfake:     func(v vaultFake) *vaultFake { v.ErrInGetServerTLS = true; return &v }(vfk),
			expectErr: true,
		},
		{
			name: "failed to create vault config",
			vs: &api.VaultServer{
				ObjectMeta: getVaultObjectMeta(1),
			},

			vfake:     func(v vaultFake) *vaultFake { v.ErrInGetConfig = true; return &v }(vfk),
			expectErr: true,
		},
		{
			name: "failed to deploy vault",
			vs: &api.VaultServer{
				ObjectMeta: getVaultObjectMeta(1),
			},

			vfake:     func(v vaultFake) *vaultFake { v.ErrInApply = true; return &v }(vfk),
			expectErr: true,
		},
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			vaultCtrl := VaultController{
				kubeClient:       kfake.NewSimpleClientset(),
				recorder:         record.NewFakeRecorder(0),
				ctxCancels:       map[string]CtxWithCancel{},
				extClient:        cfake.NewSimpleClientset(),
				appCatalogClient: appcatfake.NewSimpleClientset().AppcatalogV1alpha1(),
			}

			// to ignore monitorAndUpdateStatus
			vaultCtrl.ctxCancels[test.vs.Name] = CtxWithCancel{}

			err := vaultCtrl.reconcileVault(test.vs, test.vfake)
			if test.expectErr {
				assert.NotNil(t, err, "error must be non-empty")
			} else {
				assert.Nil(t, err, "error must be nil")

				_, err := vaultCtrl.kubeClient.AppsV1().Deployments(test.vs.Namespace).Get(test.vs.Name, metav1.GetOptions{})
				assert.Nil(t, err, "deployment for vaultserver should exist")
			}
		})
	}
}

func TestDeployVault(t *testing.T) {
	vfk := vaultFake{
		sa: &core.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sa-test",
				Namespace: "test",
			},
		},
		cnt: core.Container{},
		pt:  &core.PodTemplateSpec{},
		dp: &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "dp-test",
				Namespace: "test",
			},
		},
		svc: &core.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "svc-test",
				Namespace: "test",
			},
		},
		roles:      []rbac.Role{},
		ErrInApply: false,
	}
	testData := []struct {
		name      string
		vfake     *vaultFake
		vs        *api.VaultServer
		expectErr bool
	}{
		{
			name:  "no error",
			vfake: &vfk,
			vs: &api.VaultServer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "svc-test",
					Namespace: "test",
				},
			},
			expectErr: false,
		},
		{
			name:  "expected error",
			vfake: func(v vaultFake) *vaultFake { v.ErrInApply = true; return &v }(vfk),
			vs: &api.VaultServer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "svc-test",
					Namespace: "test",
				},
			},
			expectErr: true,
		},
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			vaultCtrl := VaultController{
				kubeClient: kfake.NewSimpleClientset(),
			}

			err := vaultCtrl.DeployVault(test.vs, test.vfake)
			if test.expectErr {
				assert.NotNil(t, err, "error must be non-empty")
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestCreateRoleAndRoleBinding(t *testing.T) {
	vaultCtrl := VaultController{
		kubeClient: kfake.NewSimpleClientset(),
	}

	demoRole := rbac.Role{
		Rules: []rbac.PolicyRule{
			{
				APIGroups: []string{core.GroupName},
				Resources: []string{"secret"},
				Verbs:     []string{"*"},
			},
		},
	}

	vs := &api.VaultServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "try",
			Namespace: "default",
		},
	}

	testData := []struct {
		testName           string
		preCreatedRole     []rbac.Role
		roles              []rbac.Role
		expectErr          bool
		expectRoles        []string
		expectRoleBindings []string
	}{
		{
			"create 2 rbac role and rolebinding",
			nil,
			[]rbac.Role{
				func(r *rbac.Role) rbac.Role { r.SetName("test1"); r.SetNamespace(vs.Namespace); return *r }(&demoRole),
				func(r *rbac.Role) rbac.Role { r.SetName("test2"); r.SetNamespace(vs.Namespace); return *r }(&demoRole),
			},
			false,
			[]string{"test1", "test2"},
			[]string{"test1", "test2"},
		},
		{
			"create 1 rbac role and rolebinding, but role already exists",
			[]rbac.Role{
				func(r *rbac.Role) rbac.Role { r.SetName("test3"); r.SetNamespace(vs.Namespace); return *r }(&demoRole),
			},
			[]rbac.Role{
				func(r *rbac.Role) rbac.Role { r.SetName("test3"); r.SetNamespace(vs.Namespace); return *r }(&demoRole),
			},
			false,
			[]string{"test3"},
			[]string{"test3"},
		},
	}

	for _, test := range testData {
		t.Run(test.testName, func(t *testing.T) {
			for _, r := range test.preCreatedRole {
				_, err := vaultCtrl.kubeClient.RbacV1().Roles(vs.Namespace).Create(&r)
				assert.Nil(t, err)
			}

			err := ensureRoleAndRoleBinding(vaultCtrl.kubeClient, vs, test.roles, nil)
			if test.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			for _, r := range test.expectRoles {
				_, err := vaultCtrl.kubeClient.RbacV1().Roles(vs.Namespace).Get(r, metav1.GetOptions{})
				assert.Nil(t, err, fmt.Sprintf("role(%s) should exists", r))
			}

			for _, rb := range test.expectRoleBindings {
				_, err := vaultCtrl.kubeClient.RbacV1().RoleBindings(vs.Namespace).Get(rb, metav1.GetOptions{})
				assert.Nil(t, err, fmt.Sprintf("rolebinding (%s) should exists", rb))
			}
		})
	}
}
