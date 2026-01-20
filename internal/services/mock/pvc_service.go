package mock

import (
	"cluster-agent/internal/models"
	"context"

	"github.com/stretchr/testify/mock"
)

type PVCServiceMock struct {
	mock.Mock
}

func (m *PVCServiceMock) List(ctx context.Context, namespace string) ([]models.PVCListInfo, error) {
	args := m.Called(ctx, namespace)
	return args.Get(0).([]models.PVCListInfo), args.Error(1)
}

func (m *PVCServiceMock) Get(ctx context.Context, namespace, name string) (*models.PVCDetails, error) {
	args := m.Called(ctx, namespace, name)
	return args.Get(0).(*models.PVCDetails), args.Error(1)
}
