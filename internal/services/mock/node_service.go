package mock

import (
	"cluster-agent/internal/models"
	"context"

	"github.com/stretchr/testify/mock"
)

type NodeServiceMock struct {
	mock.Mock
}

func (m *NodeServiceMock) GetNodes(ctx context.Context) ([]models.Node, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Node), args.Error(1)
}
