package models

import "time"

type RoleBindingInfo struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Role      string    `json:"role"`
	Age       time.Time `json:"age"`
}
