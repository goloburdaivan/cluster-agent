package rules

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
)

type WorkloadConfigRule struct{}

func (r *WorkloadConfigRule) Apply(
	s *models.ClusterSnapshot,
	b *graph.Builder,
) error {

	for _, d := range s.Deployments {
		dNodeID := id("Deployment", d.Namespace, d.Name)

		for _, cm := range s.ConfigMaps {
			if cm.Namespace != d.Namespace {
				continue
			}

			if usedInWorkload(
				d.Spec.Template.Spec,
				cm.Name,
				"ConfigMap",
			) {
				b.AddEdge(edge(
					dNodeID,
					id("ConfigMap", cm.Namespace, cm.Name),
				))
			}
		}
	}

	for _, ss := range s.StatefulSets {
		ssNodeID := id("StatefulSet", ss.Namespace, ss.Name)

		for _, cm := range s.ConfigMaps {
			if cm.Namespace != ss.Namespace {
				continue
			}

			if usedInWorkload(
				ss.Spec.Template.Spec,
				cm.Name,
				"ConfigMap",
			) {
				b.AddEdge(edge(
					ssNodeID,
					id("ConfigMap", cm.Namespace, cm.Name),
				))
			}
		}
	}

	for _, ds := range s.DaemonSets {
		dsNodeID := id("DaemonSet", ds.Namespace, ds.Name)

		for _, cm := range s.ConfigMaps {
			if cm.Namespace != ds.Namespace {
				continue
			}

			if usedInWorkload(
				ds.Spec.Template.Spec,
				cm.Name,
				"ConfigMap",
			) {
				b.AddEdge(edge(
					dsNodeID,
					id("ConfigMap", cm.Namespace, cm.Name),
				))
			}
		}
	}

	for _, j := range s.Jobs {
		jNodeID := id("Job", j.Namespace, j.Name)

		for _, cm := range s.ConfigMaps {
			if cm.Namespace != j.Namespace {
				continue
			}

			if usedInWorkload(
				j.Spec.Template.Spec,
				cm.Name,
				"ConfigMap",
			) {
				b.AddEdge(edge(
					jNodeID,
					id("ConfigMap", cm.Namespace, cm.Name),
				))
			}
		}
	}

	return nil
}
