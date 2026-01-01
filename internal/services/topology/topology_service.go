package topology

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
	"cluster-agent/internal/services/topology/rules"
)

type Service interface {
	BuildFromSnapshot(snapshot *models.ClusterSnapshot) (*graph.Graph, error)
}

type topologyService struct {
	rules []Rule
}

func NewTopologyService() Service {
	return &topologyService{
		rules: []Rule{
			&rules.ResourceNodesRule{},

			// Network layer
			&rules.WorkloadServiceRule{},
			&rules.IngressServiceRule{},
			&rules.ServiceDiscoveryRule{},

			// Storage / config
			&rules.WorkloadPVCRule{},
			&rules.DeploymentConfigRule{},
			&rules.DeploymentSecretRule{},
		},
	}
}

func (s *topologyService) BuildFromSnapshot(snapshot *models.ClusterSnapshot) (*graph.Graph, error) {
	builder := graph.NewGraphBuilder()

	for _, rule := range s.rules {
		if err := rule.Apply(snapshot, builder); err != nil {
			return nil, err
		}
	}

	return builder.Build(), nil
}
