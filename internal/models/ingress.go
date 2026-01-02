package models

import (
	v1 "k8s.io/api/networking/v1"
	"time"
)

type IngressListInfo struct {
	Name         string    `json:"name"`
	Namespace    string    `json:"namespace"`
	LoadBalancer string    `json:"load_balancer"`
	Rules        []string  `json:"rules"`
	Age          time.Time `json:"age"`
}

type IngressDetails struct {
	IngressListInfo
	Spec v1.IngressSpec `json:"spec"`
	UID  string         `json:"uid"`
}
