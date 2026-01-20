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

func TestIngressHandler_List(t *testing.T) {
	type testCase struct {
		name          string
		queryString   string
		mockBehavior  func(m *mock.IngressServiceMock)
		expectedCode  int
		expectedError string
	}

	tests := []testCase{
		{
			name:        "Success",
			queryString: "?namespace=default",
			mockBehavior: func(m *mock.IngressServiceMock) {
				m.On("List", testifyMock.Anything, "default").
					Return([]models.IngressListInfo{{Name: "ingress1"}}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "Success without namespace",
			queryString: "",
			mockBehavior: func(m *mock.IngressServiceMock) {
				m.On("List", testifyMock.Anything, "").
					Return([]models.IngressListInfo{{Name: "ingress1"}}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "Internal error",
			queryString: "?namespace=default",
			mockBehavior: func(m *mock.IngressServiceMock) {
				m.On("List", testifyMock.Anything, "default").
					Return([]models.IngressListInfo(nil), assert.AnError)
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "assert.AnError general error for testing",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mock.IngressServiceMock)
			tc.mockBehavior(svc)

			handler := NewIngressHandler(svc)
			r := setupRouter()
			r.GET("/ingresses", handler.List)

			w := performRequest(r, "GET", "/ingresses"+tc.queryString, nil)

			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedError != "" {
				assert.Contains(t, w.Body.String(), tc.expectedError)
			}

			svc.AssertExpectations(t)
		})
	}
}

func TestIngressHandler_Get(t *testing.T) {
	type testCase struct {
		name         string
		namespace    string
		ingressName  string
		mockBehavior func(m *mock.IngressServiceMock)
		expectedCode int
		expectError  bool
	}

	tests := []testCase{
		{
			name:        "Success",
			namespace:   "default",
			ingressName: "my-ingress",
			mockBehavior: func(m *mock.IngressServiceMock) {
				m.On("Get", testifyMock.Anything, "default", "my-ingress").
					Return(&models.IngressDetails{IngressListInfo: models.IngressListInfo{Name: "my-ingress"}}, nil)
			},
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:        "Not Found",
			namespace:   "default",
			ingressName: "missing",
			mockBehavior: func(m *mock.IngressServiceMock) {
				m.On("Get", testifyMock.Anything, "default", "missing").
					Return((*models.IngressDetails)(nil), services.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectError:  true,
		},
		{
			name:        "Internal error",
			namespace:   "default",
			ingressName: "error-ingress",
			mockBehavior: func(m *mock.IngressServiceMock) {
				m.On("Get", testifyMock.Anything, "default", "error-ingress").
					Return((*models.IngressDetails)(nil), assert.AnError)
			},
			expectedCode: http.StatusInternalServerError,
			expectError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mock.IngressServiceMock)
			tc.mockBehavior(svc)

			handler := NewIngressHandler(svc)
			r := setupRouter()
			r.GET("/:namespace/:name", handler.Get)

			url := fmt.Sprintf("/%s/%s", tc.namespace, tc.ingressName)
			w := performRequest(r, "GET", url, nil)

			assert.Equal(t, tc.expectedCode, w.Code)
			svc.AssertExpectations(t)
		})
	}
}
