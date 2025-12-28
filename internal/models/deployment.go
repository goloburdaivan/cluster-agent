package models

type DeploymentStatus string

const (
	DeploymentStatusHealthy     DeploymentStatus = "Healthy"
	DeploymentStatusProgressing DeploymentStatus = "Progressing"
	DeploymentStatusDegraded    DeploymentStatus = "Degraded"
)

type DeploymentInfo struct {
	Name            string           `json:"name"`
	Namespace       string           `json:"namespace"`
	Replicas        int32            `json:"replicas"`
	ReadyReplicas   int32            `json:"ready_replicas"`
	UpdatedReplicas int32            `json:"updated_replicas"`
	Status          DeploymentStatus `json:"status"`
}

type ScaleDeploymentParams struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Replicas  int32  `json:"replicas" binding:"required,min=0"`
}
