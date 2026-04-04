package services

import (
	"cluster-agent/internal/models"
	"context"
	"fmt"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/client-go/kubernetes"
)

type StorageClassService interface {
	List(ctx context.Context) ([]models.StorageClassInfo, error)
	Get(ctx context.Context, name string) (*storagev1.StorageClass, error)
	Create(ctx context.Context, sc *storagev1.StorageClass) error
	Update(ctx context.Context, sc *storagev1.StorageClass) error
	Delete(ctx context.Context, name string) error
}

type storageClassService struct {
	clientset kubernetes.Interface
}

func NewStorageClassService(c kubernetes.Interface) StorageClassService {
	return &storageClassService{
		clientset: c,
	}
}

func (s *storageClassService) List(ctx context.Context) ([]models.StorageClassInfo, error) {
	list, err := s.clientset.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list storageclasses: %w", err)
	}

	result := make([]models.StorageClassInfo, 0, len(list.Items))
	for _, item := range list.Items {
		isDefault := false
		if item.Annotations["storageclass.kubernetes.io/is-default-class"] == "true" {
			isDefault = true
		}

		result = append(result, models.StorageClassInfo{
			Name:        item.Name,
			Provisioner: item.Provisioner,
			IsDefault:   isDefault,
			Age:         item.CreationTimestamp.Time,
		})
	}
	return result, nil
}

func (s *storageClassService) Get(ctx context.Context, name string) (*storagev1.StorageClass, error) {
	sc, err := s.clientset.StorageV1().StorageClasses().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get storageclass: %w", err)
	}
	return sc, nil
}

func (s *storageClassService) Create(ctx context.Context, sc *storagev1.StorageClass) error {
	_, err := s.clientset.StorageV1().StorageClasses().Create(ctx, sc, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create storageclass: %w", err)
	}
	return nil
}

func (s *storageClassService) Update(ctx context.Context, sc *storagev1.StorageClass) error {
	_, err := s.clientset.StorageV1().StorageClasses().Update(ctx, sc, metav1.UpdateOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update storageclass: %w", err)
	}
	return nil
}

func (s *storageClassService) Delete(ctx context.Context, name string) error {
	err := s.clientset.StorageV1().StorageClasses().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete storageclass: %w", err)
	}
	return nil
}
