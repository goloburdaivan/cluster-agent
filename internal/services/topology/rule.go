package topology

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
)

type Rule interface {
	Apply(
		snapshot *models.ClusterSnapshot,
		builder *graph.Builder,
	) error
}
