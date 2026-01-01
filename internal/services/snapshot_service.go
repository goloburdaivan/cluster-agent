package services

import (
	"cluster-agent/internal/models"
	"fmt"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	appsv1 "k8s.io/client-go/listers/apps/v1"
	corev1 "k8s.io/client-go/listers/core/v1"
	v1 "k8s.io/client-go/listers/networking/v1"
)

type SnapshotService interface {
	TakeClusterSnapshot(namespace string) (*models.ClusterSnapshot, error)
}

type snapshotService struct {
	deploymentLister  appsv1.DeploymentLister
	serviceLister     corev1.ServiceLister
	statefulSetLister appsv1.StatefulSetLister
	ingressLister     v1.IngressLister
	configMapLister   corev1.ConfigMapLister
	secretLister      corev1.SecretLister
	pvcLister         corev1.PersistentVolumeClaimLister
}

func NewSnapshotService(
	factory informers.SharedInformerFactory,
) SnapshotService {
	return &snapshotService{
		deploymentLister:  factory.Apps().V1().Deployments().Lister(),
		serviceLister:     factory.Core().V1().Services().Lister(),
		statefulSetLister: factory.Apps().V1().StatefulSets().Lister(),
		ingressLister:     factory.Networking().V1().Ingresses().Lister(),
		configMapLister:   factory.Core().V1().ConfigMaps().Lister(),
		secretLister:      factory.Core().V1().Secrets().Lister(),
		pvcLister:         factory.Core().V1().PersistentVolumeClaims().Lister(),
	}
}

func (s snapshotService) TakeClusterSnapshot(namespace string) (*models.ClusterSnapshot, error) {
	var snapshot models.ClusterSnapshot

	deployments, err := s.deploymentLister.Deployments(namespace).List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	services, err := s.serviceLister.Services(namespace).List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	statefulSets, err := s.statefulSetLister.StatefulSets(namespace).List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("failed to list statefulSets: %w", err)
	}

	ingresses, err := s.ingressLister.Ingresses(namespace).List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("failed to list ingresses: %w", err)
	}

	configMaps, err := s.configMapLister.ConfigMaps(namespace).List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("failed to list configMaps: %w", err)
	}

	secrets, err := s.secretLister.Secrets(namespace).List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}

	pvcs, err := s.pvcLister.PersistentVolumeClaims(namespace).List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("failed to list pvcs: %w", err)
	}

	snapshot.Deployments = deployments
	snapshot.Services = services
	snapshot.StatefulSets = statefulSets
	snapshot.Ingresses = ingresses
	snapshot.ConfigMaps = configMaps
	snapshot.Secrets = secrets
	snapshot.PVCs = pvcs

	return &snapshot, nil
}
