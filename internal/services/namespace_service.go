package services

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NamespaceService interface {
	GetNamespaces(ctx context.Context) ([]string, error)
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
