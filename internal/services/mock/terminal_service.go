package mock

import (
	"github.com/stretchr/testify/mock"
	"k8s.io/client-go/tools/remotecommand"
)

type TerminalServiceMock struct {
	mock.Mock
}

func (m *TerminalServiceMock) GetAuthExecutor(namespace, podName, container string) (remotecommand.Executor, error) {
	args := m.Called(namespace, podName, container)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(remotecommand.Executor), args.Error(1)
}
