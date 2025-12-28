package models

type NodeStatus string

const (
	NodeStatusReady    NodeStatus = "Ready"
	NodeStatusNotReady NodeStatus = "NotReady"
	NodeStatusUnknown  NodeStatus = "Unknown"
)

type Node struct {
	Name     string     `json:"name"`
	Status   NodeStatus `json:"status"`
	Role     string     `json:"role"`
	Version  string     `json:"version"`
	CpuUsage float64    `json:"cpu_usage"`
	MemUsage float64    `json:"mem_usage"`
}
