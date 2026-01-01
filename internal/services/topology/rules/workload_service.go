package rules

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
)

type WorkloadServiceRule struct {
}

func (r *WorkloadServiceRule) Apply(
	s *models.ClusterSnapshot,
	b *graph.Builder,
) error {

	for _, d := range s.Deployments {
		workloadID := id("Deployment", d.Namespace, d.Name)

		for _, svc := range s.Services {
			if svc.Namespace != d.Namespace {
				continue
			}

			if labelsMatch(svc.Spec.Selector, d.Spec.Template.Labels) {
				b.AddEdge(edge(workloadID, id("Service", svc.Namespace, svc.Name)))
			}
		}
	}

	for _, ss := range s.StatefulSets {
		workloadID := id("StatefulSet", ss.Namespace, ss.Name)

		for _, svc := range s.Services {
			if svc.Namespace != ss.Namespace {
				continue
			}

			if labelsMatch(svc.Spec.Selector, ss.Spec.Template.Labels) {
				b.AddEdge(edge(workloadID, id("Service", svc.Namespace, svc.Name)))
			}
		}
	}

	return nil
}
