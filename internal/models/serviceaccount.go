package models

import "time"

type ServiceAccountInfo struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Secrets   int       `json:"secrets"`
	Age       time.Time `json:"age"`
}
