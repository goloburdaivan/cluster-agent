package rules

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
)

type DeploymentServiceRule struct {
}

func (r *DeploymentServiceRule) Apply(
	s *models.ClusterSnapshot,
	b *graph.Builder,
) error {

	for _, d := range s.Deployments {
		dNode := graph.Node{
			ID:   "deployment:" + d.Namespace + "/" + d.Name,
			Kind: "Deployment",
			Name: d.Name,
		}
		b.AddNode(dNode)

		for _, svc := range s.Services {
			if svc.Namespace != d.Namespace {
				continue
			}

			if labelsMatch(
				svc.Spec.Selector,
				d.Spec.Template.Labels,
			) {
				svcNode := graph.Node{
					ID:   "service:" + svc.Namespace + "/" + svc.Name,
					Kind: "Service",
					Name: svc.Name,
				}
				b.AddNode(svcNode)

				b.AddEdge(graph.Edge{
					Source: dNode.ID,
					Target: svcNode.ID,
				})
			}
		}
	}
	return nil
}
