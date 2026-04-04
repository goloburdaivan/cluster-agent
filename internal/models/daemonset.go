package models

import "time"

type DaemonSetInfo struct {
	Name             string    `json:"name"`
	Namespace        string    `json:"namespace"`
	DesiredScheduled int32     `json:"desiredScheduled"`
	CurrentScheduled int32     `json:"currentScheduled"`
	NumberReady      int32     `json:"numberReady"`
	Age              time.Time `json:"age"`
}
