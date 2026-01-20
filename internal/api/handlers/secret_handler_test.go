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

func TestSecretHandler_List(t *testing.T) {
	type testCase struct {
		name          string
		queryString   string
		mockBehavior  func(m *mock.SecretServiceMock)
		expectedCode  int
		expectedError string
	}

	tests := []testCase{
		{
			name:        "Success",
			queryString: "?namespace=default",
			mockBehavior: func(m *mock.SecretServiceMock) {
				m.On("List", testifyMock.Anything, "default").
					Return([]models.SecretListInfo{{Name: "secret1"}}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "Success without namespace",
			queryString: "",
			mockBehavior: func(m *mock.SecretServiceMock) {
				m.On("List", testifyMock.Anything, "").
					Return([]models.SecretListInfo{{Name: "secret1"}}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "Internal error",
			queryString: "?namespace=default",
			mockBehavior: func(m *mock.SecretServiceMock) {
				m.On("List", testifyMock.Anything, "default").
					Return([]models.SecretListInfo(nil), assert.AnError)
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "assert.AnError general error for testing",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mock.SecretServiceMock)
			tc.mockBehavior(svc)

			handler := NewSecretHandler(svc)
			r := setupRouter()
			r.GET("/secrets", handler.List)

			w := performRequest(r, "GET", "/secrets"+tc.queryString, nil)

			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedError != "" {
				assert.Contains(t, w.Body.String(), tc.expectedError)
			}

			svc.AssertExpectations(t)
		})
	}
}

func TestSecretHandler_Get(t *testing.T) {
	type testCase struct {
		name         string
		namespace    string
		secretName   string
		mockBehavior func(m *mock.SecretServiceMock)
		expectedCode int
		expectError  bool
	}

	tests := []testCase{
		{
			name:       "Success",
			namespace:  "default",
			secretName: "my-secret",
			mockBehavior: func(m *mock.SecretServiceMock) {
				m.On("Get", testifyMock.Anything, "default", "my-secret").
					Return(&models.SecretDetails{SecretListInfo: models.SecretListInfo{Name: "my-secret"}}, nil)
			},
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:       "Not Found",
			namespace:  "default",
			secretName: "missing",
			mockBehavior: func(m *mock.SecretServiceMock) {
				m.On("Get", testifyMock.Anything, "default", "missing").
					Return((*models.SecretDetails)(nil), services.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectError:  true,
		},
		{
			name:       "Internal error",
			namespace:  "default",
			secretName: "error-secret",
			mockBehavior: func(m *mock.SecretServiceMock) {
				m.On("Get", testifyMock.Anything, "default", "error-secret").
					Return((*models.SecretDetails)(nil), assert.AnError)
			},
			expectedCode: http.StatusInternalServerError,
			expectError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mock.SecretServiceMock)
			tc.mockBehavior(svc)

			handler := NewSecretHandler(svc)
			r := setupRouter()
			r.GET("/:namespace/:name", handler.Get)

			url := fmt.Sprintf("/%s/%s", tc.namespace, tc.secretName)
			w := performRequest(r, "GET", url, nil)

			assert.Equal(t, tc.expectedCode, w.Code)
			svc.AssertExpectations(t)
		})
	}
}
