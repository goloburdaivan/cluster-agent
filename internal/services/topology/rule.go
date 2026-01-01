package topology

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
)

type Mapper interface {
	Map(snapshot *models.ClusterSnapshot, builder *graph.Builder)
}
