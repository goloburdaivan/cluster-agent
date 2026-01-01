package rules

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
)

type DeploymentSecretRule struct{}

func (r *DeploymentSecretRule) Apply(
	s *models.ClusterSnapshot,
	b *graph.Builder,
) error {

	for _, d := range s.Deployments {
		dNodeID := "deployment:" + d.Namespace + "/" + d.Name

		for _, sec := range s.Secrets {
			if sec.Namespace != d.Namespace {
				continue
			}

			if usedInDeployment(
				d.Spec.Template.Spec,
				sec.Name,
				"Secret",
			) {
				b.AddEdge(graph.Edge{
					Source: dNodeID,
					Target: "secret:" + sec.Namespace + "/" + sec.Name,
				})
			}
		}
	}
	return nil
}
