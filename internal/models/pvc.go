package models

import (
	corev1 "k8s.io/api/core/v1"
	"time"
)

type PVCListInfo struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Status    string    `json:"status"`
	Capacity  string    `json:"capacity"`
	Age       time.Time `json:"age"`
}

type PVCDetails struct {
	PVCListInfo
	AccessModes  []corev1.PersistentVolumeAccessMode `json:"access_modes"`
	StorageClass *string                             `json:"storage_class"`
	VolumeName   string                              `json:"volume_name"`
	UID          string                              `json:"uid"`
	MountedBy    []string                            `json:"mounted_by"`
}
