package rules

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
)

type DeploymentConfigRule struct{}

func (r *DeploymentConfigRule) Apply(
	s *models.ClusterSnapshot,
	b *graph.Builder,
) error {

	for _, d := range s.Deployments {
		dNodeID := "deployment:" + d.Namespace + "/" + d.Name

		for _, cm := range s.ConfigMaps {
			if cm.Namespace != d.Namespace {
				continue
			}

			if usedInDeployment(
				d.Spec.Template.Spec,
				cm.Name,
				"ConfigMap",
			) {
				b.AddEdge(graph.Edge{
					Source: dNodeID,
					Target: "configmap:" + cm.Namespace + "/" + cm.Name,
				})
			}
		}
	}
	return nil
}
