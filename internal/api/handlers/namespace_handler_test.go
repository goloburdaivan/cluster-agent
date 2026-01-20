package handlers

import (
	"cluster-agent/internal/services/mock"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestNamespaceHandler_List(t *testing.T) {
	type testCase struct {
		name          string
		mockBehavior  func(m *mock.NamespaceServiceMock)
		expectedCode  int
		expectedError string
	}

	tests := []testCase{
		{
			name: "Success",
			mockBehavior: func(m *mock.NamespaceServiceMock) {
				m.On("GetNamespaces", testifyMock.Anything).
					Return([]string{"default", "kube-system"}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Error",
			mockBehavior: func(m *mock.NamespaceServiceMock) {
				m.On("GetNamespaces", testifyMock.Anything).
					Return([]string(nil), assert.AnError)
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "assert.AnError general error for testing",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mock.NamespaceServiceMock)
			tc.mockBehavior(svc)

			handler := NewNamespaceHandler(svc)
			r := setupRouter()
			r.GET("/namespaces", handler.List)

			w := performRequest(r, "GET", "/namespaces", nil)

			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedError != "" {
				assert.Contains(t, w.Body.String(), tc.expectedError)
			}

			svc.AssertExpectations(t)
		})
	}
}
