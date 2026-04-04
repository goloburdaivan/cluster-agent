package rules

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
)

type StorageHierarchyRule struct{}

func (r *StorageHierarchyRule) Apply(
	s *models.ClusterSnapshot,
	b *graph.Builder,
) error {

	// PVC -> PV
	for _, pvc := range s.PVCs {
		if pvc.Spec.VolumeName != "" {
			b.AddEdge(edge(
				id("PVC", pvc.Namespace, pvc.Name),
				id("PV", "", pvc.Spec.VolumeName),
			))
		}
	}

	// PV -> StorageClass
	for _, pv := range s.PVs {
		if pv.Spec.StorageClassName != "" {
			b.AddEdge(edge(
				id("PV", "", pv.Name),
				id("StorageClass", "", pv.Spec.StorageClassName),
			))
		}
	}

	// PVC -> StorageClass (direct reference)
	for _, pvc := range s.PVCs {
		if pvc.Spec.StorageClassName != nil && *pvc.Spec.StorageClassName != "" {
			b.AddEdge(edge(
				id("PVC", pvc.Namespace, pvc.Name),
				id("StorageClass", "", *pvc.Spec.StorageClassName),
			))
		}
	}

	return nil
}
