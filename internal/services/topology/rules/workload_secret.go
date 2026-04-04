package rules

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
)

type WorkloadSecretRule struct{}

func (r *WorkloadSecretRule) Apply(
	s *models.ClusterSnapshot,
	b *graph.Builder,
) error {

	for _, d := range s.Deployments {
		dNodeID := id("Deployment", d.Namespace, d.Name)

		for _, sec := range s.Secrets {
			if sec.Namespace != d.Namespace {
				continue
			}

			if usedInWorkload(
				d.Spec.Template.Spec,
				sec.Name,
				"Secret",
			) {
				b.AddEdge(edge(
					dNodeID,
					id("Secret", sec.Namespace, sec.Name),
				))
			}
		}
	}

	for _, ss := range s.StatefulSets {
		ssNodeID := id("StatefulSet", ss.Namespace, ss.Name)

		for _, sec := range s.Secrets {
			if sec.Namespace != ss.Namespace {
				continue
			}

			if usedInWorkload(
				ss.Spec.Template.Spec,
				sec.Name,
				"Secret",
			) {
				b.AddEdge(edge(
					ssNodeID,
					id("Secret", sec.Namespace, sec.Name),
				))
			}
		}
	}

	for _, ds := range s.DaemonSets {
		dsNodeID := id("DaemonSet", ds.Namespace, ds.Name)

		for _, sec := range s.Secrets {
			if sec.Namespace != ds.Namespace {
				continue
			}

			if usedInWorkload(
				ds.Spec.Template.Spec,
				sec.Name,
				"Secret",
			) {
				b.AddEdge(edge(
					dsNodeID,
					id("Secret", sec.Namespace, sec.Name),
				))
			}
		}
	}

	for _, j := range s.Jobs {
		jNodeID := id("Job", j.Namespace, j.Name)

		for _, sec := range s.Secrets {
			if sec.Namespace != j.Namespace {
				continue
			}

			if usedInWorkload(
				j.Spec.Template.Spec,
				sec.Name,
				"Secret",
			) {
				b.AddEdge(edge(
					jNodeID,
					id("Secret", sec.Namespace, sec.Name),
				))
			}
		}
	}

	return nil
}
