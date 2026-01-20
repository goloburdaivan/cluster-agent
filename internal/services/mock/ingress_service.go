package mock

import (
	"cluster-agent/internal/models"
	"context"

	"github.com/stretchr/testify/mock"
)

type IngressServiceMock struct {
	mock.Mock
}

func (m *IngressServiceMock) List(ctx context.Context, namespace string) ([]models.IngressListInfo, error) {
	args := m.Called(ctx, namespace)
	return args.Get(0).([]models.IngressListInfo), args.Error(1)
}

func (m *IngressServiceMock) Get(ctx context.Context, namespace, name string) (*models.IngressDetails, error) {
	args := m.Called(ctx, namespace, name)
	return args.Get(0).(*models.IngressDetails), args.Error(1)
}
