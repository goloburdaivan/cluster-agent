package mock

import (
	"cluster-agent/internal/models"
	"context"

	"github.com/stretchr/testify/mock"
)

type ConfigMapServiceMock struct {
	mock.Mock
}

func (m *ConfigMapServiceMock) List(ctx context.Context, namespace string) ([]models.ConfigMapListInfo, error) {
	args := m.Called(ctx, namespace)
	return args.Get(0).([]models.ConfigMapListInfo), args.Error(1)
}

func (m *ConfigMapServiceMock) Get(ctx context.Context, namespace, name string) (*models.ConfigMapDetails, error) {
	args := m.Called(ctx, namespace, name)
	return args.Get(0).(*models.ConfigMapDetails), args.Error(1)
}
