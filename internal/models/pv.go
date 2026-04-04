package models

import (
	"time"

	corev1 "k8s.io/api/core/v1"
)

type PVInfo struct {
	Name         string                  `json:"name"`
	Capacity     string                  `json:"capacity"`
	Status       string                  `json:"status"`
	Claim        *corev1.ObjectReference `json:"claim,omitempty"`
	StorageClass string                  `json:"storageClass"`
	Age          time.Time               `json:"age"`
}
