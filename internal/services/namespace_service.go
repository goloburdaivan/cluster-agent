package services

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NamespaceService interface {
	GetNamespaces(ctx context.Context) ([]string, error)
	GetNamespace(ctx context.Context, name string) (*corev1.Namespace, error)
	CreateNamespace(ctx context.Context, namespace *corev1.Namespace) error
	DeleteNamespace(ctx context.Context, name string) error
}

type namespaceService struct {
	clientset kubernetes.Interface
}

func NewNamespaceService(clientset kubernetes.Interface) NamespaceService {
	return &namespaceService{
		clientset: clientset,
	}
}

func (n *namespaceService) GetNamespaces(ctx context.Context) ([]string, error) {
	list, err := n.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	result := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		result = append(result, item.Name)
	}
	return result, nil
}

func (n *namespaceService) GetNamespace(ctx context.Context, name string) (*corev1.Namespace, error) {
	ns, err := n.clientset.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get namespace: %w", err)
	}
	return ns, nil
}

func (n *namespaceService) CreateNamespace(ctx context.Context, namespace *corev1.Namespace) error {
	_, err := n.clientset.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create namespace: %w", err)
	}
	return nil
}

func (n *namespaceService) DeleteNamespace(ctx context.Context, name string) error {
	err := n.clientset.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete namespace: %w", err)
	}
	return nil
}
