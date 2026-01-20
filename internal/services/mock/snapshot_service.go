package mock

import (
	"cluster-agent/internal/models"

	"github.com/stretchr/testify/mock"
)

type SnapshotServiceMock struct {
	mock.Mock
}

func (m *SnapshotServiceMock) TakeClusterSnapshot(namespace string) (*models.ClusterSnapshot, error) {
	args := m.Called(namespace)
	return args.Get(0).(*models.ClusterSnapshot), args.Error(1)
}
