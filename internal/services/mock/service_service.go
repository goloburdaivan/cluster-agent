package mock

import (
	"cluster-agent/internal/models"
	"context"

	"github.com/stretchr/testify/mock"
)

type KubernetesServiceServiceMock struct {
	mock.Mock
}

func (m *KubernetesServiceServiceMock) List(ctx context.Context, namespace string) ([]models.ServiceInfo, error) {
	args := m.Called(ctx, namespace)
	return args.Get(0).([]models.ServiceInfo), args.Error(1)
}

func (m *KubernetesServiceServiceMock) Get(ctx context.Context, namespace, name string) (*models.ServiceDetails, error) {
	args := m.Called(ctx, namespace, name)
	return args.Get(0).(*models.ServiceDetails), args.Error(1)
}
