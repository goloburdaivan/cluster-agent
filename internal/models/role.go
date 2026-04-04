package models

import "time"

type RoleInfo struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Age       time.Time `json:"age"`
}
