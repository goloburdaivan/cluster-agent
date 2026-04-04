package services

import (
	"cluster-agent/internal/models"
	"context"
	"fmt"

	networkingv1 "k8s.io/api/networking/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NetworkPolicyService interface {
	List(ctx context.Context, namespace string) ([]models.NetworkPolicyInfo, error)
	Get(ctx context.Context, namespace, name string) (*networkingv1.NetworkPolicy, error)
	Create(ctx context.Context, np *networkingv1.NetworkPolicy) error
	Update(ctx context.Context, np *networkingv1.NetworkPolicy) error
	Delete(ctx context.Context, namespace, name string) error
}

type networkPolicyService struct {
	clientset kubernetes.Interface
}

func NewNetworkPolicyService(c kubernetes.Interface) NetworkPolicyService {
	return &networkPolicyService{
		clientset: c,
	}
}

func (s *networkPolicyService) List(ctx context.Context, namespace string) ([]models.NetworkPolicyInfo, error) {
	list, err := s.clientset.NetworkingV1().NetworkPolicies(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list networkpolicies: %w", err)
	}

	result := make([]models.NetworkPolicyInfo, 0, len(list.Items))
	for _, item := range list.Items {
		result = append(result, models.NetworkPolicyInfo{
			Name:      item.Name,
			Namespace: item.Namespace,
			PodSelector: item.Spec.PodSelector.MatchLabels,
			Age:       item.CreationTimestamp.Time,
		})
	}
	return result, nil
}

func (s *networkPolicyService) Get(ctx context.Context, namespace, name string) (*networkingv1.NetworkPolicy, error) {
	np, err := s.clientset.NetworkingV1().NetworkPolicies(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get networkpolicy: %w", err)
	}
	return np, nil
}

func (s *networkPolicyService) Create(ctx context.Context, np *networkingv1.NetworkPolicy) error {
	_, err := s.clientset.NetworkingV1().NetworkPolicies(np.Namespace).Create(ctx, np, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create networkpolicy: %w", err)
	}
	return nil
}

func (s *networkPolicyService) Update(ctx context.Context, np *networkingv1.NetworkPolicy) error {
	_, err := s.clientset.NetworkingV1().NetworkPolicies(np.Namespace).Update(ctx, np, metav1.UpdateOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update networkpolicy: %w", err)
	}
	return nil
}

func (s *networkPolicyService) Delete(ctx context.Context, namespace, name string) error {
	err := s.clientset.NetworkingV1().NetworkPolicies(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete networkpolicy: %w", err)
	}
	return nil
}
