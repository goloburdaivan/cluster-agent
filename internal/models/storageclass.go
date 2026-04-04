package models

import "time"

type StorageClassInfo struct {
	Name        string    `json:"name"`
	Provisioner string    `json:"provisioner"`
	IsDefault   bool      `json:"isDefault"`
	Age         time.Time `json:"age"`
}
