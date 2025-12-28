package models

type ServiceType string

const (
	ServiceTypeClusterIP    ServiceType = "ClusterIP"
	ServiceTypeNodePort     ServiceType = "NodePort"
	ServiceTypeLoadBalancer ServiceType = "LoadBalancer"
)

type ServiceInfo struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Type      ServiceType       `json:"type"`
	ClusterIP string            `json:"cluster_ip"`
	Selector  map[string]string `json:"selector"`
	Ports     []int32           `json:"ports"`
}
