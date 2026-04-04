package services

import (
	"cluster-agent/internal/models"
	"context"
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type RoleBindingService interface {
	List(ctx context.Context, namespace string) ([]models.RoleBindingInfo, error)
	Get(ctx context.Context, namespace, name string) (*rbacv1.RoleBinding, error)
	Create(ctx context.Context, rb *rbacv1.RoleBinding) error
	Update(ctx context.Context, rb *rbacv1.RoleBinding) error
	Delete(ctx context.Context, namespace, name string) error
}

type roleBindingService struct {
	clientset kubernetes.Interface
}

func NewRoleBindingService(c kubernetes.Interface) RoleBindingService {
	return &roleBindingService{
		clientset: c,
	}
}

func (s *roleBindingService) List(ctx context.Context, namespace string) ([]models.RoleBindingInfo, error) {
	list, err := s.clientset.RbacV1().RoleBindings(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list rolebindings: %w", err)
	}

	result := make([]models.RoleBindingInfo, 0, len(list.Items))
	for _, item := range list.Items {
		result = append(result, models.RoleBindingInfo{
			Name:      item.Name,
			Namespace: item.Namespace,
			Role:      item.RoleRef.Name,
			Age:       item.CreationTimestamp.Time,
		})
	}
	return result, nil
}

func (s *roleBindingService) Get(ctx context.Context, namespace, name string) (*rbacv1.RoleBinding, error) {
	rb, err := s.clientset.RbacV1().RoleBindings(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get rolebinding: %w", err)
	}
	return rb, nil
}

func (s *roleBindingService) Create(ctx context.Context, rb *rbacv1.RoleBinding) error {
	_, err := s.clientset.RbacV1().RoleBindings(rb.Namespace).Create(ctx, rb, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create rolebinding: %w", err)
	}
	return nil
}

func (s *roleBindingService) Update(ctx context.Context, rb *rbacv1.RoleBinding) error {
	_, err := s.clientset.RbacV1().RoleBindings(rb.Namespace).Update(ctx, rb, metav1.UpdateOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update rolebinding: %w", err)
	}
	return nil
}

func (s *roleBindingService) Delete(ctx context.Context, namespace, name string) error {
	err := s.clientset.RbacV1().RoleBindings(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete rolebinding: %w", err)
	}
	return nil
}
