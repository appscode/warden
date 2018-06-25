package storage

import (
	api "github.com/kubevault/operator/apis/core/v1alpha1"
	"github.com/kubevault/operator/pkg/vault/storage/etcd"
	"github.com/kubevault/operator/pkg/vault/storage/gcs"
	"github.com/kubevault/operator/pkg/vault/storage/inmem"
	"github.com/kubevault/operator/pkg/vault/storage/s3"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

type Storage interface {
	Apply(pt *corev1.PodTemplateSpec) error
	GetSecrets(namespace string) ([]corev1.Secret, error)
	GetStorageConfig() (string, error)
}

func NewStorage(s *api.BackendStorageSpec) (Storage, error) {
	if s.Inmem {
		return inmem.NewOptions()
	} else if s.Etcd != nil {
		return etcd.NewOptions(*s.Etcd)
	} else if s.Gcs != nil {
		return gcs.NewOptions(*s.Gcs)
	} else if s.S3 != nil {
		return s3.NewOptions(*s.S3)
	} else {
		return nil, errors.New("invalid storage backend")
	}
}
