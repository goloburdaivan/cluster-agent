package mock

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
	"context"

	"github.com/stretchr/testify/mock"
)

type TopologyServiceMock struct {
	mock.Mock
}

func (m *TopologyServiceMock) BuildFromSnapshot(ctx context.Context, snapshot *models.ClusterSnapshot) (*graph.Graph, error) {
	args := m.Called(ctx, snapshot)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*graph.Graph), args.Error(1)
}
