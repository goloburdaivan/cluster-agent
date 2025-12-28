package services

import (
	"cluster-agent/internal/models"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodService interface {
	GetPods(ctx context.Context, namespace string) ([]models.PodListInfo, error)
	GetPod(ctx context.Context, namespace, name string) (*models.PodDetails, error)
}

type podService struct {
	clientset kubernetes.Interface
}

func NewPodService(clientset kubernetes.Interface) PodService {
	return &podService{
		clientset: clientset,
	}
}

func (p *podService) GetPods(ctx context.Context, namespace string) ([]models.PodListInfo, error) {
	list, err := p.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, fmt.Errorf("failed to list pods in namespace %s: %w", namespace, err)
	}

	result := make([]models.PodListInfo, 0, len(list.Items))

	for _, item := range list.Items {
		ownerName := ""
		if len(item.OwnerReferences) > 0 {
			ownerName = item.OwnerReferences[0].Name
		}

		podInfo := models.PodListInfo{
			Name:      item.Name,
			Namespace: item.Namespace,
			IP:        item.Status.PodIP,
			OwnerName: ownerName,
		}

		if item.DeletionTimestamp != nil {
			podInfo.Status = models.PodStatusTerminating
		} else {
			podInfo.Status = models.PodStatus(item.Status.Phase)
		}

		result = append(result, podInfo)
	}

	return result, nil
}

func (p *podService) GetPod(ctx context.Context, namespace, name string) (*models.PodDetails, error) {
	rawPod, err := p.clientset.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return p.mapToDetails(rawPod), nil
}

func (p *podService) mapToDetails(pod *corev1.Pod) *models.PodDetails {
	details := &models.PodDetails{
		PodListInfo: models.PodListInfo{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Status:    models.PodStatus(pod.Status.Phase),
			IP:        pod.Status.PodIP,
		},
		Created: pod.CreationTimestamp.Time,
		Node:    pod.Spec.NodeName,
		Labels:  pod.Labels,
		UID:     string(pod.UID),

		PodIPs:            pod.Status.PodIPs,
		Containers:        pod.Spec.Containers,
		ContainerStatuses: pod.Status.ContainerStatuses,
		Volumes:           pod.Spec.Volumes,
		Conditions:        pod.Status.Conditions,
	}

	if len(pod.OwnerReferences) > 0 {
		details.OwnerName = pod.OwnerReferences[0].Name
	}

	var totalRestarts int32
	for _, status := range pod.Status.ContainerStatuses {
		totalRestarts += status.RestartCount
	}
	details.Restarts = totalRestarts

	return details
}
