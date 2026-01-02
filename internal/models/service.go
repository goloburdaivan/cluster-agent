package models

import (
	corev1 "k8s.io/api/core/v1"
	"time"
)

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

type ServiceDetails struct {
	ServiceInfo
	ExternalIPs []string             `json:"external_ips"`
	FullPorts   []corev1.ServicePort `json:"full_ports"`
	UID         string               `json:"uid"`
	Age         time.Time            `json:"age"`
}
