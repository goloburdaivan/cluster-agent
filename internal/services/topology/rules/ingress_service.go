package rules

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
)

type IngressServiceRule struct {
}

func (r *IngressServiceRule) Apply(
	s *models.ClusterSnapshot,
	b *graph.Builder,
) error {

	for _, ing := range s.Ingresses {
		ingID := id("Ingress", ing.Namespace, ing.Name)

		for _, rule := range ing.Spec.Rules {
			if rule.HTTP == nil {
				continue
			}

			for _, path := range rule.HTTP.Paths {
				if path.Backend.Service == nil {
					continue
				}

				svcName := path.Backend.Service.Name
				b.AddEdge(edge(
					ingID,
					id("Service", ing.Namespace, svcName),
				))
			}
		}
	}

	return nil
}
