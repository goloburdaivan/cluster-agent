package services

import (
	"cluster-agent/internal/models"
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PVService interface {
	List(ctx context.Context) ([]models.PVInfo, error)
	Get(ctx context.Context, name string) (*corev1.PersistentVolume, error)
	Create(ctx context.Context, pv *corev1.PersistentVolume) error
	Update(ctx context.Context, pv *corev1.PersistentVolume) error
	Delete(ctx context.Context, name string) error
}

type pvService struct {
	clientset kubernetes.Interface
}

func NewPVService(c kubernetes.Interface) PVService {
	return &pvService{
		clientset: c,
	}
}

func (s *pvService) List(ctx context.Context) ([]models.PVInfo, error) {
	list, err := s.clientset.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pvs: %w", err)
	}

	result := make([]models.PVInfo, 0, len(list.Items))
	for _, item := range list.Items {
		capacity := ""
		if cap, ok := item.Spec.Capacity[corev1.ResourceStorage]; ok {
			capacity = cap.String()
		}

		result = append(result, models.PVInfo{
			Name:         item.Name,
			Capacity:     capacity,
			Status:       string(item.Status.Phase),
			Claim:        item.Spec.ClaimRef,
			StorageClass: item.Spec.StorageClassName,
			Age:          item.CreationTimestamp.Time,
		})
	}
	return result, nil
}

func (s *pvService) Get(ctx context.Context, name string) (*corev1.PersistentVolume, error) {
	pv, err := s.clientset.CoreV1().PersistentVolumes().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get pv: %w", err)
	}
	return pv, nil
}

func (s *pvService) Create(ctx context.Context, pv *corev1.PersistentVolume) error {
	_, err := s.clientset.CoreV1().PersistentVolumes().Create(ctx, pv, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create pv: %w", err)
	}
	return nil
}

func (s *pvService) Update(ctx context.Context, pv *corev1.PersistentVolume) error {
	_, err := s.clientset.CoreV1().PersistentVolumes().Update(ctx, pv, metav1.UpdateOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update pv: %w", err)
	}
	return nil
}

func (s *pvService) Delete(ctx context.Context, name string) error {
	err := s.clientset.CoreV1().PersistentVolumes().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete pv: %w", err)
	}
	return nil
}
