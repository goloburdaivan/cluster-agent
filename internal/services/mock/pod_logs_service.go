package mock

import (
	"context"
	"io"

	"github.com/stretchr/testify/mock"
)

type PodLogsServiceMock struct {
	mock.Mock
}

func (m *PodLogsServiceMock) StreamLogs(ctx context.Context, namespace, podName, containerName string) (io.ReadCloser, error) {
	args := m.Called(ctx, namespace, podName, containerName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}
