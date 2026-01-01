package rules

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
	"strings"
)

type ServiceDiscoveryRule struct{}

func (r *ServiceDiscoveryRule) Apply(
	s *models.ClusterSnapshot,
	b *graph.Builder,
) error {
	servicesMap := make(map[string]string)
	for _, svc := range s.Services {
		servicesMap[svc.Name] = id("Service", svc.Namespace, svc.Name)
	}

	for _, d := range s.Deployments {
		deployNodeID := id("Deployment", d.Namespace, d.Name)

		for _, container := range d.Spec.Template.Spec.Containers {
			for _, env := range container.Env {
				envValue := strings.TrimSpace(env.Value)

				if len(envValue) < 3 {
					continue
				}

				for svcName, svcID := range servicesMap {
					if strings.Contains(envValue, svcName) {
						b.AddEdge(graph.Edge{
							Source: deployNodeID,
							Target: svcID,
						})
					}
				}
			}
		}
	}

	return nil
}
