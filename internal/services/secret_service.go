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

type SecretService interface {
	List(ctx context.Context, namespace string) ([]models.SecretListInfo, error)
	Get(ctx context.Context, namespace, name string) (*models.SecretDetails, error)
}

type secretService struct {
	clientset kubernetes.Interface
}

func NewSecretService(c kubernetes.Interface) SecretService {
	return &secretService{
		clientset: c,
	}
}

func (s *secretService) List(ctx context.Context, namespace string) ([]models.SecretListInfo, error) {
	list, err := s.clientset.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed list secrets: %w", err)
	}

	result := make([]models.SecretListInfo, 0, len(list.Items))
	for _, item := range list.Items {
		result = append(result, s.mapToListInfo(&item))
	}
	return result, nil
}

func (s *secretService) Get(ctx context.Context, namespace, name string) (*models.SecretDetails, error) {
	item, err := s.clientset.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}

	return &models.SecretDetails{
		SecretListInfo: s.mapToListInfo(item),
		Data:           item.Data,
		UID:            string(item.UID),
		Labels:         item.Labels,
		Annotations:    item.Annotations,
		Immutable:      item.Immutable,
	}, nil
}

func (s *secretService) mapToListInfo(item *corev1.Secret) models.SecretListInfo {
	keys := make([]string, 0, len(item.Data))
	for k := range item.Data {
		keys = append(keys, k)
	}

	return models.SecretListInfo{
		Name:      item.Name,
		Namespace: item.Namespace,
		Type:      string(item.Type),
		Keys:      keys,
		Age:       item.CreationTimestamp.Time,
	}
}
