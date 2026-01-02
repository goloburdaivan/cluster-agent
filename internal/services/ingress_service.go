package services

import (
	"cluster-agent/internal/models"
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type IngressService interface {
	List(ctx context.Context, namespace string) ([]models.IngressListInfo, error)
	Get(ctx context.Context, namespace, name string) (*models.IngressDetails, error)
}

type ingressService struct{ clientset kubernetes.Interface }

func NewIngressService(c kubernetes.Interface) IngressService { return &ingressService{c} }

func (s *ingressService) List(ctx context.Context, namespace string) ([]models.IngressListInfo, error) {
	list, err := s.clientset.NetworkingV1().Ingresses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	result := make([]models.IngressListInfo, 0, len(list.Items))
	for _, item := range list.Items {
		var rules []string
		for _, r := range item.Spec.Rules {
			rules = append(rules, r.Host)
		}

		lb := ""
		if len(item.Status.LoadBalancer.Ingress) > 0 {
			lb = item.Status.LoadBalancer.Ingress[0].IP
			if lb == "" {
				lb = item.Status.LoadBalancer.Ingress[0].Hostname
			}
		}

		result = append(result, models.IngressListInfo{
			Name: item.Name, Namespace: item.Namespace, LoadBalancer: lb, Rules: rules, Age: item.CreationTimestamp.Time,
		})
	}
	return result, nil
}

func (s *ingressService) Get(ctx context.Context, namespace, name string) (*models.IngressDetails, error) {
	item, err := s.clientset.NetworkingV1().Ingresses(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var rules []string
	for _, r := range item.Spec.Rules {
		rules = append(rules, r.Host)
	}
	lb := ""
	if len(item.Status.LoadBalancer.Ingress) > 0 {
		lb = item.Status.LoadBalancer.Ingress[0].IP
		if lb == "" {
			lb = item.Status.LoadBalancer.Ingress[0].Hostname
		}
	}

	return &models.IngressDetails{
		IngressListInfo: models.IngressListInfo{Name: item.Name, Namespace: item.Namespace, LoadBalancer: lb, Rules: rules, Age: item.CreationTimestamp.Time},
		Spec:            item.Spec, UID: string(item.UID),
	}, nil
}
