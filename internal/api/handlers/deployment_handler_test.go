package handlers

import (
	"bytes"
	"cluster-agent/internal/models"
	"cluster-agent/internal/services"
	"cluster-agent/internal/services/mock"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"testing"
)

func TestListDeployments(t *testing.T) {
	type testCase struct {
		name          string
		queryString   string
		mockBehavior  func(m *mock.DeploymentServiceMock)
		expectedCode  int
		expectedError string
		expectedData  []models.DeploymentInfo
	}

	tests := []testCase{
		{
			name:        "Success",
			queryString: "?namespace=default",
			mockBehavior: func(m *mock.DeploymentServiceMock) {
				m.On("GetDeployments", testifyMock.Anything, "default").
					Return([]models.DeploymentInfo{{Name: "Test"}}, nil)
			},
			expectedCode: http.StatusOK,
			expectedData: []models.DeploymentInfo{{Name: "Test"}},
		},
		{
			name:        "Kubernetes internal error",
			queryString: "?namespace=default",
			mockBehavior: func(m *mock.DeploymentServiceMock) {
				m.On("GetDeployments", testifyMock.Anything, "default").
					Return([]models.DeploymentInfo(nil), assert.AnError)
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "assert.AnError general error for testing",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mock.DeploymentServiceMock)
			tc.mockBehavior(svc)

			handler := NewDeploymentHandler(svc)
			r := setupRouter()
			r.GET("/deployments", handler.List)

			w := performRequest(r, "GET", "/deployments"+tc.queryString, nil)

			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedError != "" {
				assert.Contains(t, w.Body.String(), tc.expectedError)
			} else {
				resp := parseResponse[[]models.DeploymentInfo](t, w)
				assert.Equal(t, tc.expectedData, resp.Data)
			}

			svc.AssertExpectations(t)
		})
	}
}

func TestGetDeployment(t *testing.T) {
	type testCase struct {
		name           string
		namespace      string
		deploymentName string
		mockBehavior   func(m *mock.DeploymentServiceMock)
		expectedCode   int
		expectError    bool
	}

	tests := []testCase{
		{
			name:           "Success",
			namespace:      "default",
			deploymentName: "nginx",
			mockBehavior: func(m *mock.DeploymentServiceMock) {
				m.On("GetDeployment", testifyMock.Anything, "default", "nginx").
					Return(&v1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: "nginx"}}, nil)
			},
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:           "Not Found",
			namespace:      "default",
			deploymentName: "missing",
			mockBehavior: func(m *mock.DeploymentServiceMock) {
				m.On("GetDeployment", testifyMock.Anything, "default", "missing").
					Return((*v1.Deployment)(nil), services.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectError:  true,
		},
		{
			name:           "Kubernetes internal error",
			namespace:      "default",
			deploymentName: "nginx",
			mockBehavior: func(m *mock.DeploymentServiceMock) {
				m.On("GetDeployment", testifyMock.Anything, "default", "nginx").
					Return((*v1.Deployment)(nil), assert.AnError)
			},
			expectedCode: http.StatusInternalServerError,
			expectError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mock.DeploymentServiceMock)
			tc.mockBehavior(svc)

			r := setupRouter()
			r.GET("/:namespace/:name", NewDeploymentHandler(svc).Get)

			url := fmt.Sprintf("/%s/%s", tc.namespace, tc.deploymentName)
			w := performRequest(r, "GET", url, nil)

			assert.Equal(t, tc.expectedCode, w.Code)

			if !tc.expectError {
				response := parseResponse[*v1.Deployment](t, w)
				assert.Equal(t, tc.deploymentName, response.Data.Name)
			}

			svc.AssertExpectations(t)
		})
	}
}

func TestCreateDeployment(t *testing.T) {
	type testCase struct {
		name         string
		inputBody    string
		mockBehavior func(m *mock.DeploymentServiceMock)
		expectedCode int
	}

	tests := []testCase{
		{
			name:      "Success",
			inputBody: `{"metadata": {"name": "test-app"}}`,
			mockBehavior: func(m *mock.DeploymentServiceMock) {
				m.On("CreateDeployment", testifyMock.Anything, testifyMock.MatchedBy(func(d *v1.Deployment) bool {
					return d.Name == "test-app"
				})).Return(nil)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:      "Bad Request (Invalid JSON)",
			inputBody: `{invalid-json}`,
			mockBehavior: func(m *mock.DeploymentServiceMock) {
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:      "Service Error",
			inputBody: `{"metadata": {"name": "test-app"}}`,
			mockBehavior: func(m *mock.DeploymentServiceMock) {
				m.On("CreateDeployment", testifyMock.Anything, testifyMock.Anything).
					Return(errors.New("creation failed"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mock.DeploymentServiceMock)
			tc.mockBehavior(svc)

			r := setupRouter()
			r.POST("/deployments", NewDeploymentHandler(svc).Create)

			w := performRequest(r, "POST", "/deployments", bytes.NewBufferString(tc.inputBody))

			assert.Equal(t, tc.expectedCode, w.Code)
			svc.AssertExpectations(t)
		})
	}
}

func TestDeleteDeployment(t *testing.T) {
	type testCase struct {
		name           string
		namespace      string
		deploymentName string
		mockBehavior   func(m *mock.DeploymentServiceMock)
		expectedCode   int
		expectError    bool
	}

	tests := []testCase{
		{
			name:           "Success",
			namespace:      "default",
			deploymentName: "old-app",
			mockBehavior: func(m *mock.DeploymentServiceMock) {
				m.On("DeleteDeployment", testifyMock.Anything, "default", "old-app").
					Return(nil)
			},
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:           "Not Found",
			namespace:      "default",
			deploymentName: "missing-app",
			mockBehavior: func(m *mock.DeploymentServiceMock) {
				m.On("DeleteDeployment", testifyMock.Anything, "default", "missing-app").
					Return(services.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectError:  true,
		},
		{
			name:           "Kubernetes Internal Error",
			namespace:      "default",
			deploymentName: "broken-app",
			mockBehavior: func(m *mock.DeploymentServiceMock) {
				m.On("DeleteDeployment", testifyMock.Anything, "default", "broken-app").
					Return(assert.AnError)
			},
			expectedCode: http.StatusInternalServerError,
			expectError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mock.DeploymentServiceMock)
			tc.mockBehavior(svc)

			r := setupRouter()
			r.DELETE("/:namespace/:name", NewDeploymentHandler(svc).Delete)

			url := fmt.Sprintf("/%s/%s", tc.namespace, tc.deploymentName)
			w := performRequest(r, "DELETE", url, nil)

			assert.Equal(t, tc.expectedCode, w.Code)

			if !tc.expectError {
				response := parseResponse[string](t, w)
				assert.Equal(t, "OK", response.Data)
			}

			svc.AssertExpectations(t)
		})
	}
}

func TestScaleDeployment(t *testing.T) {
	type testCase struct {
		name         string
		inputBody    string
		mockBehavior func(m *mock.DeploymentServiceMock)
		expectedCode int
		expectError  bool
	}

	tests := []testCase{
		{
			name:      "Success",
			inputBody: `{"namespace": "default", "name": "app", "replicas": 3}`,
			mockBehavior: func(m *mock.DeploymentServiceMock) {
				expectedParams := models.ScaleDeploymentParams{
					Namespace: "default",
					Name:      "app",
					Replicas:  3,
				}
				m.On("ScaleDeployment", testifyMock.Anything, expectedParams).Return(nil)
			},
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:      "Bad Request (Invalid JSON)",
			inputBody: `Wait, this is not JSON`,
			mockBehavior: func(m *mock.DeploymentServiceMock) {
			},
			expectedCode: http.StatusBadRequest,
			expectError:  true,
		},
		{
			name:      "Kubernetes Internal Error",
			inputBody: `{"namespace": "default", "name": "app", "replicas": 3}`,
			mockBehavior: func(m *mock.DeploymentServiceMock) {
				m.On("ScaleDeployment", testifyMock.Anything, testifyMock.Anything).
					Return(assert.AnError)
			},
			expectedCode: http.StatusInternalServerError,
			expectError:  true,
		},
		{
			name:      "Deployment Not Found",
			inputBody: `{"namespace": "missing", "name": "app", "replicas": 3}`,
			mockBehavior: func(m *mock.DeploymentServiceMock) {
				expectedParams := models.ScaleDeploymentParams{
					Namespace: "missing",
					Name:      "app",
					Replicas:  3,
				}
				m.On("ScaleDeployment", testifyMock.Anything, expectedParams).
					Return(services.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mock.DeploymentServiceMock)
			tc.mockBehavior(svc)

			r := setupRouter()
			r.POST("/scale", NewDeploymentHandler(svc).ScaleDeployment)

			w := performRequest(r, "POST", "/scale", bytes.NewBufferString(tc.inputBody))

			assert.Equal(t, tc.expectedCode, w.Code)

			if !tc.expectError {
				response := parseResponse[string](t, w)
				assert.Equal(t, "OK", response.Data)
			}

			svc.AssertExpectations(t)
		})
	}
}
