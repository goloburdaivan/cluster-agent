package rules

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
)

type NetworkPolicyRule struct {
}

func (r *NetworkPolicyRule) Apply(
	s *models.ClusterSnapshot,
	b *graph.Builder,
) error {

	for _, np := range s.NetworkPolicies {
		npID := id("NetworkPolicy", np.Namespace, np.Name)
		selector := np.Spec.PodSelector.MatchLabels

		for _, d := range s.Deployments {
			if d.Namespace != np.Namespace {
				continue
			}
			if labelsMatch(selector, d.Spec.Template.Labels) {
				b.AddEdge(edge(npID, id("Deployment", d.Namespace, d.Name)))
			}
		}

		for _, ss := range s.StatefulSets {
			if ss.Namespace != np.Namespace {
				continue
			}
			if labelsMatch(selector, ss.Spec.Template.Labels) {
				b.AddEdge(edge(npID, id("StatefulSet", ss.Namespace, ss.Name)))
			}
		}

		for _, ds := range s.DaemonSets {
			if ds.Namespace != np.Namespace {
				continue
			}
			if labelsMatch(selector, ds.Spec.Template.Labels) {
				b.AddEdge(edge(npID, id("DaemonSet", ds.Namespace, ds.Name)))
			}
		}

		for _, j := range s.Jobs {
			if j.Namespace != np.Namespace {
				continue
			}
			if labelsMatch(selector, j.Spec.Template.Labels) {
				b.AddEdge(edge(npID, id("Job", j.Namespace, j.Name)))
			}
		}
	}

	return nil
}
