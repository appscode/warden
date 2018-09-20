package framework

import (
	"fmt"
	"time"

	"github.com/appscode/go/crypto/rand"
	. "github.com/onsi/gomega"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	mysqlServiceName    = rand.WithUniqSuffix("test-svc-mysql")
	mysqlDeploymentName = rand.WithUniqSuffix("test-mysql-deploy")
)

func (f *Framework) DeployMySQL() (string, error) {
	label := map[string]string{
		"app": rand.WithUniqSuffix("test-mysql"),
	}

	srv := core.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: f.namespace,
			Name:      mysqlServiceName,
		},
		Spec: core.ServiceSpec{
			Selector: label,
			Ports: []core.ServicePort{
				{
					Name:       "tcp",
					Protocol:   core.ProtocolTCP,
					Port:       3306,
					TargetPort: intstr.FromInt(3306),
				},
			},
		},
	}

	url := fmt.Sprintf("%s.%s.svc:3306", mysqlServiceName, f.namespace)

	mysqlCont := core.Container{
		Name:            "mysql",
		Image:           "mysql:5.6",
		ImagePullPolicy: "IfNotPresent",
		Env: []core.EnvVar{
			{
				Name:  "MYSQL_ROOT_PASSWORD",
				Value: "root",
			},
		},
		Ports: []core.ContainerPort{
			{
				Name:          "mysql",
				Protocol:      core.ProtocolTCP,
				ContainerPort: 3306,
			},
		},
		VolumeMounts: []core.VolumeMount{
			{
				MountPath: "/var/lib/mysql/data/pgdata",
				Name:      "data",
			},
		},
	}

	mysqlDeploy := apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: f.namespace,
			Name:      mysqlDeploymentName,
		},
		Spec: apps.DeploymentSpec{
			Replicas: func(i int32) *int32 { return &i }(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: label,
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: label,
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						mysqlCont,
					},
					Volumes: []core.Volume{
						{
							Name: "data",
							VolumeSource: core.VolumeSource{
								EmptyDir: &core.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}

	err := f.CreateService(srv)
	if err != nil {
		return "", err
	}

	_, err = f.CreateDeployment(mysqlDeploy)
	if err != nil {
		return "", err
	}

	Eventually(func() bool {
		if obj, err := f.KubeClient.AppsV1beta1().Deployments(f.namespace).Get(mysqlDeploy.GetName(), metav1.GetOptions{}); err == nil {
			return *obj.Spec.Replicas == obj.Status.ReadyReplicas
		}
		return false
	}, timeOut, pollingInterval).Should(BeTrue())

	time.Sleep(10 * time.Second)

	return url, nil
}

func (f *Framework) DeleteMySQL() error {
	err := f.DeleteService(mysqlServiceName, f.namespace)
	if err != nil {
		return err
	}

	err = f.DeleteDeployment(mysqlDeploymentName, f.namespace)
	return err
}
