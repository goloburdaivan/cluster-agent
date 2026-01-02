package services

import (
	"cluster-agent/internal/models"
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ConfigMapService interface {
	List(ctx context.Context, namespace string) ([]models.ConfigMapListInfo, error)
	Get(ctx context.Context, namespace, name string) (*models.ConfigMapDetails, error)
}

type configMapService struct {
	clientset kubernetes.Interface
}

func NewConfigMapService(c kubernetes.Interface) ConfigMapService {
	return &configMapService{
		c,
	}
}

func (s *configMapService) List(ctx context.Context, namespace string) ([]models.ConfigMapListInfo, error) {
	list, err := s.clientset.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed list cm: %w", err)
	}
	result := make([]models.ConfigMapListInfo, 0, len(list.Items))
	for _, item := range list.Items {
		keys := make([]string, 0, len(item.Data))
		for k := range item.Data {
			keys = append(keys, k)
		}

		result = append(result, models.ConfigMapListInfo{
			Name: item.Name, Namespace: item.Namespace, Keys: keys, Age: item.CreationTimestamp.Time,
		})
	}
	return result, nil
}

func (s *configMapService) Get(ctx context.Context, namespace, name string) (*models.ConfigMapDetails, error) {
	item, err := s.clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(item.Data))
	for k := range item.Data {
		keys = append(keys, k)
	}

	return &models.ConfigMapDetails{
		ConfigMapListInfo: models.ConfigMapListInfo{Name: item.Name, Namespace: item.Namespace, Keys: keys, Age: item.CreationTimestamp.Time},
		Data:              item.Data, UID: string(item.UID),
	}, nil
}
