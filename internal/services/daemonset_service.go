package services

import (
	"cluster-agent/internal/models"
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DaemonSetService interface {
	GetDaemonSets(ctx context.Context, namespace string) ([]models.DaemonSetInfo, error)
	GetDaemonSet(ctx context.Context, namespace, name string) (*appsv1.DaemonSet, error)
	CreateDaemonSet(ctx context.Context, daemonset *appsv1.DaemonSet) error
	UpdateDaemonSet(ctx context.Context, daemonset *appsv1.DaemonSet) error
	DeleteDaemonSet(ctx context.Context, namespace, name string) error
}

type daemonSetService struct {
	clientset kubernetes.Interface
}

func NewDaemonSetService(clientset kubernetes.Interface) DaemonSetService {
	return &daemonSetService{
		clientset: clientset,
	}
}

func (d *daemonSetService) GetDaemonSets(ctx context.Context, namespace string) ([]models.DaemonSetInfo, error) {
	list, err := d.clientset.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list daemonsets in namespace %s: %w", namespace, err)
	}

	result := make([]models.DaemonSetInfo, 0, len(list.Items))
	for _, item := range list.Items {
		result = append(result, models.DaemonSetInfo{
			Name:             item.Name,
			Namespace:        item.Namespace,
			DesiredScheduled: item.Status.DesiredNumberScheduled,
			CurrentScheduled: item.Status.CurrentNumberScheduled,
			NumberReady:      item.Status.NumberReady,
			Age:              item.CreationTimestamp.Time,
		})
	}

	return result, nil
}

func (d *daemonSetService) GetDaemonSet(ctx context.Context, namespace, name string) (*appsv1.DaemonSet, error) {
	daemonset, err := d.clientset.AppsV1().DaemonSets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get daemonset: %w", err)
	}

	daemonset.Kind = "DaemonSet"
	daemonset.APIVersion = "apps/v1"
	daemonset.ManagedFields = nil

	return daemonset, nil
}

func (d *daemonSetService) CreateDaemonSet(ctx context.Context, daemonset *appsv1.DaemonSet) error {
	_, err := d.clientset.AppsV1().DaemonSets(daemonset.Namespace).Create(ctx, daemonset, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create daemonset: %w", err)
	}
	return nil
}

func (d *daemonSetService) UpdateDaemonSet(ctx context.Context, daemonset *appsv1.DaemonSet) error {
	_, err := d.clientset.AppsV1().DaemonSets(daemonset.Namespace).Update(ctx, daemonset, metav1.UpdateOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update daemonset: %w", err)
	}
	return nil
}

func (d *daemonSetService) DeleteDaemonSet(ctx context.Context, namespace, name string) error {
	err := d.clientset.AppsV1().DaemonSets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete daemonset: %w", err)
	}
	return nil
}
