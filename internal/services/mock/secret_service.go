package mock

import (
	"cluster-agent/internal/models"
	"context"

	"github.com/stretchr/testify/mock"
)

type SecretServiceMock struct {
	mock.Mock
}

func (m *SecretServiceMock) List(ctx context.Context, namespace string) ([]models.SecretListInfo, error) {
	args := m.Called(ctx, namespace)
	return args.Get(0).([]models.SecretListInfo), args.Error(1)
}

func (m *SecretServiceMock) Get(ctx context.Context, namespace, name string) (*models.SecretDetails, error) {
	args := m.Called(ctx, namespace, name)
	return args.Get(0).(*models.SecretDetails), args.Error(1)
}
