package mock

import (
	"cluster-agent/internal/models"
	"context"

	"github.com/stretchr/testify/mock"
)

type PodServiceMock struct {
	mock.Mock
}

func (m *PodServiceMock) GetPods(ctx context.Context, namespace string) ([]models.PodListInfo, error) {
	args := m.Called(ctx, namespace)
	return args.Get(0).([]models.PodListInfo), args.Error(1)
}

func (m *PodServiceMock) GetPod(ctx context.Context, namespace, name string) (*models.PodDetails, error) {
	args := m.Called(ctx, namespace, name)
	return args.Get(0).(*models.PodDetails), args.Error(1)
}
