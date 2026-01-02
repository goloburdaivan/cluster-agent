package services

import (
	"cluster-agent/internal/models"
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PVCService interface {
	List(ctx context.Context, namespace string) ([]models.PVCListInfo, error)
}

type pvcService struct{ clientset kubernetes.Interface }

func NewPVCService(c kubernetes.Interface) PVCService { return &pvcService{c} }

func (s *pvcService) List(ctx context.Context, namespace string) ([]models.PVCListInfo, error) {
	list, err := s.clientset.CoreV1().PersistentVolumeClaims(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	result := make([]models.PVCListInfo, 0, len(list.Items))
	for _, item := range list.Items {
		capacity := item.Status.Capacity.Storage().String()
		result = append(result, models.PVCListInfo{
			Name: item.Name, Namespace: item.Namespace, Status: string(item.Status.Phase), Capacity: capacity, Age: item.CreationTimestamp.Time,
		})
	}
	return result, nil
}
