package rules

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
)

type WorkloadPVCRule struct {
}

func (r *WorkloadPVCRule) Apply(
	s *models.ClusterSnapshot,
	b *graph.Builder,
) error {

	for _, ss := range s.StatefulSets {
		ssID := id("StatefulSet", ss.Namespace, ss.Name)

		for _, v := range ss.Spec.Template.Spec.Volumes {
			if v.PersistentVolumeClaim != nil {
				b.AddEdge(edge(
					ssID,
					id("PVC", ss.Namespace, v.PersistentVolumeClaim.ClaimName),
				))
			}
		}
	}

	return nil
}
