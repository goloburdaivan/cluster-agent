package services

import (
	"cluster-agent/internal/models"
	"context"
	"fmt"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DeploymentService interface {
	GetDeployments(ctx context.Context, namespace string) ([]models.DeploymentInfo, error)
	CreateDeployment(ctx context.Context, deployment *v1.Deployment) error
	ScaleDeployment(ctx context.Context, params models.ScaleDeploymentParams) error
}

type deploymentService struct {
	clientset kubernetes.Interface
}

func NewDeploymentService(clientset kubernetes.Interface) DeploymentService {
	return &deploymentService{
		clientset: clientset,
	}
}

func (d *deploymentService) GetDeployments(ctx context.Context, namespace string) ([]models.DeploymentInfo, error) {
	list, err := d.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments in namespace %s: %w", namespace, err)
	}

	result := make([]models.DeploymentInfo, 0, len(list.Items))
	for _, item := range list.Items {
		deployInfo := models.DeploymentInfo{
			Name:            item.Name,
			Namespace:       item.Namespace,
			Replicas:        *item.Spec.Replicas,
			ReadyReplicas:   item.Status.ReadyReplicas,
			UpdatedReplicas: item.Status.UpdatedReplicas,
			Status:          calculateDeployStatus(item.Status),
		}
		result = append(result, deployInfo)
	}

	return result, nil
}

func (d *deploymentService) CreateDeployment(ctx context.Context, deployment *v1.Deployment) error {
	_, err := d.clientset.AppsV1().Deployments(deployment.Namespace).Create(ctx, deployment, metav1.CreateOptions{})

	if err != nil {
		return fmt.Errorf("failed to create deployment: %w", err)
	}

	return nil
}

func (d *deploymentService) ScaleDeployment(ctx context.Context, params models.ScaleDeploymentParams) error {
	return executeWithRetry("scale deployment", func() error {
		scale, err := d.clientset.AppsV1().Deployments(params.Namespace).GetScale(ctx, params.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		scale.Spec.Replicas = params.Replicas
		_, err = d.clientset.AppsV1().Deployments(params.Namespace).UpdateScale(ctx, params.Name, scale, metav1.UpdateOptions{})
		return err
	})
}
