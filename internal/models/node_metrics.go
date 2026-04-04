package models

import "time"

type NodeMetrics struct {
	Name        string    `json:"name"`
	Timestamp   time.Time `json:"timestamp"`
	CPUUsage    float64   `json:"cpu_usage"`
	MemoryUsage float64   `json:"memory_usage"`
}
