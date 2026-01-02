package services

import (
	"cluster-agent/internal/models"
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubernetesServiceService interface {
	List(ctx context.Context, namespace string) ([]models.ServiceInfo, error)
	Get(ctx context.Context, namespace, name string) (*models.ServiceDetails, error)
}

type service struct {
	clientset kubernetes.Interface
}

func NewServiceService(clientset kubernetes.Interface) KubernetesServiceService {
	return &service{
		clientset: clientset,
	}
}

func (s *service) List(ctx context.Context, namespace string) ([]models.ServiceInfo, error) {
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

func (s *service) Get(ctx context.Context, namespace, name string) (*models.ServiceDetails, error) {
	svc, err := s.clientset.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	ports := make([]int32, 0, len(svc.Spec.Ports))
	for _, p := range svc.Spec.Ports {
		ports = append(ports, p.Port)
	}

	var extIPs []string
	extIPs = append(extIPs, svc.Spec.ExternalIPs...)
	for _, ing := range svc.Status.LoadBalancer.Ingress {
		if ing.IP != "" {
			extIPs = append(extIPs, ing.IP)
		} else if ing.Hostname != "" {
			extIPs = append(extIPs, ing.Hostname)
		}
	}

	return &models.ServiceDetails{
		ServiceInfo: models.ServiceInfo{
			Name:      svc.Name,
			Namespace: svc.Namespace,
			Type:      models.ServiceType(svc.Spec.Type),
			ClusterIP: svc.Spec.ClusterIP,
			Selector:  svc.Spec.Selector,
			Ports:     ports,
		},
		FullPorts:   svc.Spec.Ports,
		ExternalIPs: extIPs,
		UID:         string(svc.UID),
		Age:         svc.CreationTimestamp.Time,
	}, nil
}
