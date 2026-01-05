package services

import (
	"cluster-agent/internal/models"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/listers/core/v1"
)

type PVCService interface {
	List(ctx context.Context, namespace string) ([]models.PVCListInfo, error)
	Get(ctx context.Context, namespace, name string) (*models.PVCDetails, error)
}

type pvcService struct {
	clientset kubernetes.Interface
	podLister v1.PodLister
}

func NewPVCService(
	c kubernetes.Interface,
	podLister v1.PodLister,
) PVCService {
	return &pvcService{
		clientset: c,
		podLister: podLister,
	}
}

func (s *pvcService) List(ctx context.Context, namespace string) ([]models.PVCListInfo, error) {
	list, err := s.clientset.CoreV1().PersistentVolumeClaims(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed list pvcs: %w", err)
	}

	result := make([]models.PVCListInfo, 0, len(list.Items))
	for _, item := range list.Items {
		result = append(result, s.mapToListInfo(&item))
	}
	return result, nil
}

func (s *pvcService) Get(ctx context.Context, namespace, name string) (*models.PVCDetails, error) {
	item, err := s.clientset.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed get pvc: %w", err)
	}

	mountedPods, err := s.getMountedPods(namespace, name)

	if err != nil {
		return nil, fmt.Errorf("failed get mounted pods: %w", err)
	}

	return &models.PVCDetails{
		PVCListInfo: s.mapToListInfo(item),

		AccessModes:  item.Spec.AccessModes,
		StorageClass: item.Spec.StorageClassName,
		VolumeName:   item.Spec.VolumeName,
		UID:          string(item.UID),
		MountedBy:    mountedPods,
	}, nil
}

func (s *pvcService) mapToListInfo(item *corev1.PersistentVolumeClaim) models.PVCListInfo {
	capacity := ""
	if cap, ok := item.Status.Capacity[corev1.ResourceStorage]; ok {
		capacity = cap.String()
	}

	return models.PVCListInfo{
		Name:      item.Name,
		Namespace: item.Namespace,
		Status:    string(item.Status.Phase),
		Capacity:  capacity,
		Age:       item.CreationTimestamp.Time,
	}
}

func (s *pvcService) getMountedPods(namespace string, claimName string) ([]string, error) {
	result := make([]string, 0)
	pods, err := s.podLister.Pods(namespace).List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("failed list pods: %w", err)
	}

	for _, pod := range pods {
		for _, vol := range pod.Spec.Volumes {
			if vol.PersistentVolumeClaim != nil && vol.PersistentVolumeClaim.ClaimName == claimName {
				result = append(result, pod.Name)
			}
		}
	}

	return result, nil
}
