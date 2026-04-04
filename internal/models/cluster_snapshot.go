package models

import (
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
)

type ClusterSnapshot struct {
	// Workloads
	Deployments  []*appsv1.Deployment
	StatefulSets []*appsv1.StatefulSet
	DaemonSets   []*appsv1.DaemonSet
	Jobs         []*batchv1.Job
	CronJobs     []*batchv1.CronJob

	// Network
	Services        []*corev1.Service
	Ingresses       []*networkingv1.Ingress
	NetworkPolicies []*networkingv1.NetworkPolicy

	// Storage
	PVCs            []*corev1.PersistentVolumeClaim
	PVs             []*corev1.PersistentVolume
	StorageClasses  []*storagev1.StorageClass

	// Config
	ConfigMaps []*corev1.ConfigMap
	Secrets    []*corev1.Secret

	// Security/RBAC
	ServiceAccounts []*corev1.ServiceAccount
	Roles           []*rbacv1.Role
	RoleBindings    []*rbacv1.RoleBinding

	Namespace string
}
