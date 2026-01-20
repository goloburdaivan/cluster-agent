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

func TestConfigMapHandler_List(t *testing.T) {
	type testCase struct {
		name          string
		queryString   string
		mockBehavior  func(m *mock.ConfigMapServiceMock)
		expectedCode  int
		expectedError string
	}

	tests := []testCase{
		{
			name:        "Success",
			queryString: "?namespace=default",
			mockBehavior: func(m *mock.ConfigMapServiceMock) {
				m.On("List", testifyMock.Anything, "default").
					Return([]models.ConfigMapListInfo{{Name: "config1"}}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "Success without namespace",
			queryString: "",
			mockBehavior: func(m *mock.ConfigMapServiceMock) {
				m.On("List", testifyMock.Anything, "").
					Return([]models.ConfigMapListInfo{{Name: "config1"}}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "Internal error",
			queryString: "?namespace=default",
			mockBehavior: func(m *mock.ConfigMapServiceMock) {
				m.On("List", testifyMock.Anything, "default").
					Return([]models.ConfigMapListInfo(nil), assert.AnError)
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "assert.AnError general error for testing",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mock.ConfigMapServiceMock)
			tc.mockBehavior(svc)

			handler := NewConfigMapHandler(svc)
			r := setupRouter()
			r.GET("/configmaps", handler.List)

			w := performRequest(r, "GET", "/configmaps"+tc.queryString, nil)

			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedError != "" {
				assert.Contains(t, w.Body.String(), tc.expectedError)
			}

			svc.AssertExpectations(t)
		})
	}
}

func TestConfigMapHandler_Get(t *testing.T) {
	type testCase struct {
		name         string
		namespace    string
		configName   string
		mockBehavior func(m *mock.ConfigMapServiceMock)
		expectedCode int
		expectError  bool
	}

	tests := []testCase{
		{
			name:       "Success",
			namespace:  "default",
			configName: "my-config",
			mockBehavior: func(m *mock.ConfigMapServiceMock) {
				m.On("Get", testifyMock.Anything, "default", "my-config").
					Return(&models.ConfigMapDetails{ConfigMapListInfo: models.ConfigMapListInfo{Name: "my-config"}}, nil)
			},
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:       "Not Found",
			namespace:  "default",
			configName: "missing",
			mockBehavior: func(m *mock.ConfigMapServiceMock) {
				m.On("Get", testifyMock.Anything, "default", "missing").
					Return((*models.ConfigMapDetails)(nil), services.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectError:  true,
		},
		{
			name:       "Internal error",
			namespace:  "default",
			configName: "error-config",
			mockBehavior: func(m *mock.ConfigMapServiceMock) {
				m.On("Get", testifyMock.Anything, "default", "error-config").
					Return((*models.ConfigMapDetails)(nil), assert.AnError)
			},
			expectedCode: http.StatusInternalServerError,
			expectError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mock.ConfigMapServiceMock)
			tc.mockBehavior(svc)

			handler := NewConfigMapHandler(svc)
			r := setupRouter()
			r.GET("/:namespace/:name", handler.Get)

			url := fmt.Sprintf("/%s/%s", tc.namespace, tc.configName)
			w := performRequest(r, "GET", url, nil)

			assert.Equal(t, tc.expectedCode, w.Code)
			svc.AssertExpectations(t)
		})
	}
}
