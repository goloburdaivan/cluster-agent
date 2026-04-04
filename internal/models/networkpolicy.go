package models

import "time"

type NetworkPolicyInfo struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	PodSelector map[string]string `json:"podSelector"`
	Age         time.Time         `json:"age"`
}
