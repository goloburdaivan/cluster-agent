package models

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CronJobInfo struct {
	Name         string       `json:"name"`
	Namespace    string       `json:"namespace"`
	Schedule     string       `json:"schedule"`
	Suspended    bool         `json:"suspended"`
	Active       int          `json:"active"`
	LastSchedule *metav1.Time `json:"lastSchedule,omitempty"`
	Age          time.Time    `json:"age"`
}
