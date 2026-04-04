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

	for _, d := range s.Deployments {
		dID := id("Deployment", d.Namespace, d.Name)

		for _, v := range d.Spec.Template.Spec.Volumes {
			if v.PersistentVolumeClaim != nil {
				b.AddEdge(edge(
					dID,
					id("PVC", d.Namespace, v.PersistentVolumeClaim.ClaimName),
				))
			}
		}
	}

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

	for _, ds := range s.DaemonSets {
		dsID := id("DaemonSet", ds.Namespace, ds.Name)

		for _, v := range ds.Spec.Template.Spec.Volumes {
			if v.PersistentVolumeClaim != nil {
				b.AddEdge(edge(
					dsID,
					id("PVC", ds.Namespace, v.PersistentVolumeClaim.ClaimName),
				))
			}
		}
	}

	for _, j := range s.Jobs {
		jID := id("Job", j.Namespace, j.Name)

		for _, v := range j.Spec.Template.Spec.Volumes {
			if v.PersistentVolumeClaim != nil {
				b.AddEdge(edge(
					jID,
					id("PVC", j.Namespace, v.PersistentVolumeClaim.ClaimName),
				))
			}
		}
	}

	return nil
}
