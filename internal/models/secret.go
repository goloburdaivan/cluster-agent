package models

import "time"

type SecretListInfo struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Type      string    `json:"type"`
	Keys      []string  `json:"keys"`
	Age       time.Time `json:"age"`
}

type SecretDetails struct {
	SecretListInfo
	Data   map[string][]byte `json:"data"`
	UID    string            `json:"uid"`
	Labels map[string]string `json:"labels"`

	Annotations map[string]string `json:"annotations"`
	Immutable   *bool             `json:"immutable,omitempty"`
}
