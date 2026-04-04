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

type ConfigMapService interface {
	List(ctx context.Context, namespace string) ([]models.ConfigMapListInfo, error)
	Get(ctx context.Context, namespace, name string) (*models.ConfigMapDetails, error)
	Create(ctx context.Context, configMap *corev1.ConfigMap) error
	Update(ctx context.Context, configMap *corev1.ConfigMap) error
	Delete(ctx context.Context, namespace, name string) error
}

type configMapService struct {
	clientset kubernetes.Interface
}

func NewConfigMapService(c kubernetes.Interface) ConfigMapService {
	return &configMapService{
		clientset: c,
	}
}

func (s *configMapService) List(ctx context.Context, namespace string) ([]models.ConfigMapListInfo, error) {
	list, err := s.clientset.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed list cm: %w", err)
	}

	result := make([]models.ConfigMapListInfo, 0, len(list.Items))
	for _, item := range list.Items {
		result = append(result, s.mapToListInfo(&item))
	}
	return result, nil
}

func (s *configMapService) Get(ctx context.Context, namespace, name string) (*models.ConfigMapDetails, error) {
	item, err := s.clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get configmap: %w", err)
	}

	return &models.ConfigMapDetails{
		ConfigMapListInfo: s.mapToListInfo(item),
		Data:              item.Data,
		Labels:            item.Labels,
		Annotations:       item.Annotations,
		Immutable:         item.Immutable,
		UID:               string(item.UID),
	}, nil
}

func (s *configMapService) Create(ctx context.Context, configMap *corev1.ConfigMap) error {
	_, err := s.clientset.CoreV1().ConfigMaps(configMap.Namespace).Create(ctx, configMap, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create configmap: %w", err)
	}
	return nil
}

func (s *configMapService) Update(ctx context.Context, configMap *corev1.ConfigMap) error {
	_, err := s.clientset.CoreV1().ConfigMaps(configMap.Namespace).Update(ctx, configMap, metav1.UpdateOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update configmap: %w", err)
	}
	return nil
}

func (s *configMapService) Delete(ctx context.Context, namespace, name string) error {
	err := s.clientset.CoreV1().ConfigMaps(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete configmap: %w", err)
	}
	return nil
}

func (s *configMapService) mapToListInfo(item *corev1.ConfigMap) models.ConfigMapListInfo {
	keys := make([]string, 0, len(item.Data)+len(item.BinaryData))

	for k := range item.Data {
		keys = append(keys, k)
	}
	for k := range item.BinaryData {
		keys = append(keys, k)
	}

	return models.ConfigMapListInfo{
		Name:      item.Name,
		Namespace: item.Namespace,
		Keys:      keys,
		Age:       item.CreationTimestamp.Time,
	}
}
