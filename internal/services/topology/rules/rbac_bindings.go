package rules

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
)

type RBACBindingsRule struct{}

func (r *RBACBindingsRule) Apply(
	s *models.ClusterSnapshot,
	b *graph.Builder,
) error {

	// RoleBinding -> Role
	for _, rb := range s.RoleBindings {
		if rb.RoleRef.Kind == "Role" {
			b.AddEdge(edge(
				id("RoleBinding", rb.Namespace, rb.Name),
				id("Role", rb.Namespace, rb.RoleRef.Name),
			))
		}
	}

	// RoleBinding -> ServiceAccount
	for _, rb := range s.RoleBindings {
		for _, subject := range rb.Subjects {
			if subject.Kind == "ServiceAccount" {
				b.AddEdge(edge(
					id("RoleBinding", rb.Namespace, rb.Name),
					id("ServiceAccount", subject.Namespace, subject.Name),
				))
			}
		}
	}

	return nil
}
