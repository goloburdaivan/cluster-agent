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

type ServiceAccountService interface {
	List(ctx context.Context, namespace string) ([]models.ServiceAccountInfo, error)
	Get(ctx context.Context, namespace, name string) (*corev1.ServiceAccount, error)
	Create(ctx context.Context, sa *corev1.ServiceAccount) error
	Update(ctx context.Context, sa *corev1.ServiceAccount) error
	Delete(ctx context.Context, namespace, name string) error
}

type serviceAccountService struct {
	clientset kubernetes.Interface
}

func NewServiceAccountService(c kubernetes.Interface) ServiceAccountService {
	return &serviceAccountService{
		clientset: c,
	}
}

func (s *serviceAccountService) List(ctx context.Context, namespace string) ([]models.ServiceAccountInfo, error) {
	list, err := s.clientset.CoreV1().ServiceAccounts(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list serviceaccounts: %w", err)
	}

	result := make([]models.ServiceAccountInfo, 0, len(list.Items))
	for _, item := range list.Items {
		result = append(result, models.ServiceAccountInfo{
			Name:      item.Name,
			Namespace: item.Namespace,
			Secrets:   len(item.Secrets),
			Age:       item.CreationTimestamp.Time,
		})
	}
	return result, nil
}

func (s *serviceAccountService) Get(ctx context.Context, namespace, name string) (*corev1.ServiceAccount, error) {
	sa, err := s.clientset.CoreV1().ServiceAccounts(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get serviceaccount: %w", err)
	}
	return sa, nil
}

func (s *serviceAccountService) Create(ctx context.Context, sa *corev1.ServiceAccount) error {
	_, err := s.clientset.CoreV1().ServiceAccounts(sa.Namespace).Create(ctx, sa, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create serviceaccount: %w", err)
	}
	return nil
}

func (s *serviceAccountService) Update(ctx context.Context, sa *corev1.ServiceAccount) error {
	_, err := s.clientset.CoreV1().ServiceAccounts(sa.Namespace).Update(ctx, sa, metav1.UpdateOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update serviceaccount: %w", err)
	}
	return nil
}

func (s *serviceAccountService) Delete(ctx context.Context, namespace, name string) error {
	err := s.clientset.CoreV1().ServiceAccounts(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete serviceaccount: %w", err)
	}
	return nil
}
