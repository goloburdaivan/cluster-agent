package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type NamespaceServiceMock struct {
	mock.Mock
}

func (m *NamespaceServiceMock) GetNamespaces(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}
