package handlers

import (
	"cluster-agent/internal/services"
	"cluster-agent/internal/services/mock"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestNetworkInspectorHandler_GetConnections(t *testing.T) {
	type testCase struct {
		name          string
		namespace     string
		podName       string
		container     string
		mockBehavior  func(m *mock.NetworkInspectorServiceMock)
		expectedCode  int
		expectedError string
	}

	tests := []testCase{
		{
			name:      "Success",
			namespace: "default",
			podName:   "my-pod",
			container: "main",
			mockBehavior: func(m *mock.NetworkInspectorServiceMock) {
				m.On("GetPodNetworkConnections", testifyMock.Anything, "default", "my-pod", "main").
					Return([]services.TCPSocketEntry{{LocalAddress: "127.0.0.1", LocalPort: 8080}}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:      "Success without container",
			namespace: "default",
			podName:   "my-pod",
			container: "",
			mockBehavior: func(m *mock.NetworkInspectorServiceMock) {
				m.On("GetPodNetworkConnections", testifyMock.Anything, "default", "my-pod", "").
					Return([]services.TCPSocketEntry{}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:      "Internal error",
			namespace: "default",
			podName:   "my-pod",
			container: "main",
			mockBehavior: func(m *mock.NetworkInspectorServiceMock) {
				m.On("GetPodNetworkConnections", testifyMock.Anything, "default", "my-pod", "main").
					Return([]services.TCPSocketEntry(nil), assert.AnError)
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "failed to get network connections",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mock.NetworkInspectorServiceMock)
			tc.mockBehavior(svc)

			handler := NewNetworkInspectorHandler(svc)
			r := setupRouter()
			r.GET("/:namespace/:name/connections", handler.GetConnections)

			url := fmt.Sprintf("/%s/%s/connections", tc.namespace, tc.podName)
			if tc.container != "" {
				url += "?container=" + tc.container
			}
			w := performRequest(r, "GET", url, nil)

			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedError != "" {
				assert.Contains(t, w.Body.String(), tc.expectedError)
			}

			svc.AssertExpectations(t)
		})
	}
}
