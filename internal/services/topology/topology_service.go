package topology

import (
	"cluster-agent/internal/cache"
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
	"cluster-agent/internal/services/topology/rules"
	"context"
	"errors"
	"log"
)

type (
	Service interface {
		BuildFromSnapshot(ctx context.Context, snapshot *models.ClusterSnapshot) (*graph.Graph, error)
	}

	TopologyCacheStorage interface {
		Get(ctx context.Context, namespace string) (*graph.Graph, error)
		Set(ctx context.Context, namespace string, g *graph.Graph) error
	}
)

type topologyService struct {
	rules []Rule
	cache TopologyCacheStorage
}

func NewTopologyService(
	topologyCache TopologyCacheStorage,
) Service {
	return &topologyService{
		cache: topologyCache,
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

func (s *topologyService) BuildFromSnapshot(ctx context.Context, snapshot *models.ClusterSnapshot) (*graph.Graph, error) {
	cachedTopology, err := s.cache.Get(ctx, snapshot.Namespace)

	if err == nil {
		return cachedTopology, nil
	}

	if !errors.Is(err, cache.ErrNotFound) {
		log.Printf("Warning: failed to read topology from cache: %v", err)
	}

	builder := graph.NewGraphBuilder()

	for _, rule := range s.rules {
		if err := rule.Apply(snapshot, builder); err != nil {
			return nil, err
		}
	}

	topology := builder.Build()

	if err := s.cache.Set(ctx, snapshot.Namespace, topology); err != nil {
		log.Printf("Warning: failed to save topology to cache: %v", err)
	}

	return topology, nil
}
