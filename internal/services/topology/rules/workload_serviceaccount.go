package rules

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
)

type WorkloadServiceAccountRule struct{}

func (r *WorkloadServiceAccountRule) Apply(
	s *models.ClusterSnapshot,
	b *graph.Builder,
) error {

	for _, d := range s.Deployments {
		if d.Spec.Template.Spec.ServiceAccountName != "" {
			b.AddEdge(edge(
				id("Deployment", d.Namespace, d.Name),
				id("ServiceAccount", d.Namespace, d.Spec.Template.Spec.ServiceAccountName),
			))
		}
	}

	for _, ss := range s.StatefulSets {
		if ss.Spec.Template.Spec.ServiceAccountName != "" {
			b.AddEdge(edge(
				id("StatefulSet", ss.Namespace, ss.Name),
				id("ServiceAccount", ss.Namespace, ss.Spec.Template.Spec.ServiceAccountName),
			))
		}
	}

	for _, ds := range s.DaemonSets {
		if ds.Spec.Template.Spec.ServiceAccountName != "" {
			b.AddEdge(edge(
				id("DaemonSet", ds.Namespace, ds.Name),
				id("ServiceAccount", ds.Namespace, ds.Spec.Template.Spec.ServiceAccountName),
			))
		}
	}

	for _, j := range s.Jobs {
		if j.Spec.Template.Spec.ServiceAccountName != "" {
			b.AddEdge(edge(
				id("Job", j.Namespace, j.Name),
				id("ServiceAccount", j.Namespace, j.Spec.Template.Spec.ServiceAccountName),
			))
		}
	}

	return nil
}
