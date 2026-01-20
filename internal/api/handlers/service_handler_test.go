package handlers

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services"
	"cluster-agent/internal/services/mock"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestServiceHandler_List(t *testing.T) {
	type testCase struct {
		name          string
		queryString   string
		mockBehavior  func(m *mock.KubernetesServiceServiceMock)
		expectedCode  int
		expectedError string
	}

	tests := []testCase{
		{
			name:        "Success",
			queryString: "?namespace=default",
			mockBehavior: func(m *mock.KubernetesServiceServiceMock) {
				m.On("List", testifyMock.Anything, "default").
					Return([]models.ServiceInfo{{Name: "svc1"}}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "Success without namespace",
			queryString: "",
			mockBehavior: func(m *mock.KubernetesServiceServiceMock) {
				m.On("List", testifyMock.Anything, "").
					Return([]models.ServiceInfo{{Name: "svc1"}}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "Internal error",
			queryString: "?namespace=default",
			mockBehavior: func(m *mock.KubernetesServiceServiceMock) {
				m.On("List", testifyMock.Anything, "default").
					Return([]models.ServiceInfo(nil), assert.AnError)
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "assert.AnError general error for testing",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mock.KubernetesServiceServiceMock)
			tc.mockBehavior(svc)

			handler := NewServiceHandler(svc)
			r := setupRouter()
			r.GET("/services", handler.List)

			w := performRequest(r, "GET", "/services"+tc.queryString, nil)

			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedError != "" {
				assert.Contains(t, w.Body.String(), tc.expectedError)
			}

			svc.AssertExpectations(t)
		})
	}
}

func TestServiceHandler_Get(t *testing.T) {
	type testCase struct {
		name         string
		namespace    string
		serviceName  string
		mockBehavior func(m *mock.KubernetesServiceServiceMock)
		expectedCode int
		expectError  bool
	}

	tests := []testCase{
		{
			name:        "Success",
			namespace:   "default",
			serviceName: "my-service",
			mockBehavior: func(m *mock.KubernetesServiceServiceMock) {
				m.On("Get", testifyMock.Anything, "default", "my-service").
					Return(&models.ServiceDetails{ServiceInfo: models.ServiceInfo{Name: "my-service"}}, nil)
			},
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:        "Not Found",
			namespace:   "default",
			serviceName: "missing",
			mockBehavior: func(m *mock.KubernetesServiceServiceMock) {
				m.On("Get", testifyMock.Anything, "default", "missing").
					Return((*models.ServiceDetails)(nil), services.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectError:  true,
		},
		{
			name:        "Internal error",
			namespace:   "default",
			serviceName: "error-service",
			mockBehavior: func(m *mock.KubernetesServiceServiceMock) {
				m.On("Get", testifyMock.Anything, "default", "error-service").
					Return((*models.ServiceDetails)(nil), assert.AnError)
			},
			expectedCode: http.StatusInternalServerError,
			expectError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mock.KubernetesServiceServiceMock)
			tc.mockBehavior(svc)

			handler := NewServiceHandler(svc)
			r := setupRouter()
			r.GET("/:namespace/:name", handler.Get)

			url := fmt.Sprintf("/%s/%s", tc.namespace, tc.serviceName)
			w := performRequest(r, "GET", url, nil)

			assert.Equal(t, tc.expectedCode, w.Code)
			svc.AssertExpectations(t)
		})
	}
}
