package models

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

type ClusterSnapshot struct {
	Deployments  []*appsv1.Deployment
	Services     []*corev1.Service
	StatefulSets []*appsv1.StatefulSet
	Ingresses    []*networkingv1.Ingress
	ConfigMaps   []*corev1.ConfigMap
	Secrets      []*corev1.Secret
	PVCs         []*corev1.PersistentVolumeClaim
	Namespace    string
}
