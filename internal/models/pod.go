package models

import (
	corev1 "k8s.io/api/core/v1"
	"time"
)

type PodStatus string

const (
	PodStatusTerminating PodStatus = "Terminating"
)

type PodListInfo struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Status    PodStatus `json:"status"`
	IP        string    `json:"ip"`
	OwnerName string    `json:"owner_name"`
	Restarts  int32     `json:"restarts"`
}

type PodDetails struct {
	PodListInfo
	Created           time.Time                `json:"created"`
	Node              string                   `json:"node"`
	Labels            map[string]string        `json:"labels"`
	UID               string                   `json:"uid"`
	PodIPs            []corev1.PodIP           `json:"pod_ips"`
	Containers        []corev1.Container       `json:"containers"`
	ContainerStatuses []corev1.ContainerStatus `json:"container_statuses"`
	Volumes           []corev1.Volume          `json:"volumes"`
	Conditions        []corev1.PodCondition    `json:"conditions"`
}
