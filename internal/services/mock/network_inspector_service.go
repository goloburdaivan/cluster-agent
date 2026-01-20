package mock

import (
	"cluster-agent/internal/services"
	"context"

	"github.com/stretchr/testify/mock"
)

type NetworkInspectorServiceMock struct {
	mock.Mock
}

func (m *NetworkInspectorServiceMock) GetPodNetworkConnections(ctx context.Context, namespace, podName, container string) ([]services.TCPSocketEntry, error) {
	args := m.Called(ctx, namespace, podName, container)
	return args.Get(0).([]services.TCPSocketEntry), args.Error(1)
}
