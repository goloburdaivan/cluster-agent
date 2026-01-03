package models

import "time"

type ConfigMapListInfo struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Keys      []string  `json:"keys"`
	Age       time.Time `json:"age"`
}

type ConfigMapDetails struct {
	ConfigMapListInfo
	Data        map[string]string `json:"data"`
	UID         string            `json:"uid"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	Immutable   *bool             `json:"immutable,omitempty"`
}
