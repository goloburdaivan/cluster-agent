package services

import (
	"cluster-agent/internal/models"
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubernetesServiceService interface {
	GetServices(ctx context.Context, namespace string) ([]models.ServiceInfo, error)
}

type service struct {
	clientset kubernetes.Interface
}

func NewServiceService(clientset kubernetes.Interface) KubernetesServiceService {
	return &service{
		clientset: clientset,
	}
}

func (s *service) GetServices(ctx context.Context, namespace string) ([]models.ServiceInfo, error) {
	list, err := s.clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	result := make([]models.ServiceInfo, 0, len(list.Items))
	for _, item := range list.Items {
		ports := make([]int32, 0)
		for _, p := range item.Spec.Ports {
			ports = append(ports, p.Port)
		}

		result = append(result, models.ServiceInfo{
			Name:      item.Name,
			Namespace: item.Namespace,
			Type:      models.ServiceType(item.Spec.Type),
			ClusterIP: item.Spec.ClusterIP,
			Selector:  item.Spec.Selector,
			Ports:     ports,
		})
	}
	return result, nil
}
