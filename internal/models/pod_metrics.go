package models

import "time"

type ContainerMetrics struct {
	Name        string    `json:"name"`
	Timestamp   time.Time `json:"timestamp"`
	CPUUsage    float64   `json:"cpu_usage"`
	MemoryUsage float64   `json:"memory_usage"`
}

type PodMetrics struct {
	Name             string             `json:"name"`
	Namespace        string             `json:"namespace"`
	Timestamp        time.Time          `json:"timestamp"`
	ContainerMetrics []ContainerMetrics `json:"container_metrics"`
}
