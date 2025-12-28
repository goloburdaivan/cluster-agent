package services

import (
	"cluster-agent/internal/models"
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/util/retry"
)

func executeWithRetry(operationName string, f func() error) error {
	err := retry.RetryOnConflict(retry.DefaultRetry, f)
	if err != nil {
		return fmt.Errorf("operation '%s' failed after retries: %w", operationName, err)
	}
	return nil
}

func calculateDeployStatus(status appsv1.DeploymentStatus) models.DeploymentStatus {
	if status.ReadyReplicas == status.Replicas {
		return models.DeploymentStatusHealthy
	}

	if status.ReadyReplicas < status.Replicas {
		return models.DeploymentStatusProgressing
	}

	return models.DeploymentStatusDegraded
}

func getNodeRole(node *corev1.Node) string {
	var roles []string

	for label := range node.Labels {
		if strings.HasPrefix(label, "node-role.kubernetes.io/") {
			role := strings.TrimPrefix(label, "node-role.kubernetes.io/")
			if role == "" {
				continue
			}
			roles = append(roles, role)
		}

		if label == "kubernetes.io/role" {
			roles = append(roles, node.Labels[label])
		}
	}

	if len(roles) == 0 {
		return "worker"
	}

	return strings.Join(roles, ", ")
}

func getNodeStatus(node *corev1.Node) models.NodeStatus {
	status := models.NodeStatusUnknown

	for _, cond := range node.Status.Conditions {
		if cond.Type == corev1.NodeReady {
			if cond.Status == corev1.ConditionTrue {
				status = models.NodeStatusReady
			} else {
				status = models.NodeStatusNotReady
			}
			break
		}
	}

	return status
}
