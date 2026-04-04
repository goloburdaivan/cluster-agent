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

type RoleService interface {
	List(ctx context.Context, namespace string) ([]models.RoleInfo, error)
	Get(ctx context.Context, namespace, name string) (*rbacv1.Role, error)
	Create(ctx context.Context, role *rbacv1.Role) error
	Update(ctx context.Context, role *rbacv1.Role) error
	Delete(ctx context.Context, namespace, name string) error
}

type roleService struct {
	clientset kubernetes.Interface
}

func NewRoleService(c kubernetes.Interface) RoleService {
	return &roleService{
		clientset: c,
	}
}

func (s *roleService) List(ctx context.Context, namespace string) ([]models.RoleInfo, error) {
	list, err := s.clientset.RbacV1().Roles(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	result := make([]models.RoleInfo, 0, len(list.Items))
	for _, item := range list.Items {
		result = append(result, models.RoleInfo{
			Name:      item.Name,
			Namespace: item.Namespace,
			Age:       item.CreationTimestamp.Time,
		})
	}
	return result, nil
}

func (s *roleService) Get(ctx context.Context, namespace, name string) (*rbacv1.Role, error) {
	role, err := s.clientset.RbacV1().Roles(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	return role, nil
}

func (s *roleService) Create(ctx context.Context, role *rbacv1.Role) error {
	_, err := s.clientset.RbacV1().Roles(role.Namespace).Create(ctx, role, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}
	return nil
}

func (s *roleService) Update(ctx context.Context, role *rbacv1.Role) error {
	_, err := s.clientset.RbacV1().Roles(role.Namespace).Update(ctx, role, metav1.UpdateOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update role: %w", err)
	}
	return nil
}

func (s *roleService) Delete(ctx context.Context, namespace, name string) error {
	err := s.clientset.RbacV1().Roles(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete role: %w", err)
	}
	return nil
}
