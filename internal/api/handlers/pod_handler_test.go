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

func TestPodHandler_List(t *testing.T) {
	type testCase struct {
		name          string
		queryString   string
		mockBehavior  func(m *mock.PodServiceMock)
		expectedCode  int
		expectedError string
	}

	tests := []testCase{
		{
			name:        "Success",
			queryString: "?namespace=default",
			mockBehavior: func(m *mock.PodServiceMock) {
				m.On("GetPods", testifyMock.Anything, "default").
					Return([]models.PodListInfo{{Name: "pod1"}}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "Success without namespace",
			queryString: "",
			mockBehavior: func(m *mock.PodServiceMock) {
				m.On("GetPods", testifyMock.Anything, "").
					Return([]models.PodListInfo{{Name: "pod1"}}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "Internal error",
			queryString: "?namespace=default",
			mockBehavior: func(m *mock.PodServiceMock) {
				m.On("GetPods", testifyMock.Anything, "default").
					Return([]models.PodListInfo(nil), assert.AnError)
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "assert.AnError general error for testing",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mock.PodServiceMock)
			tc.mockBehavior(svc)

			handler := NewPodHandler(svc)
			r := setupRouter()
			r.GET("/pods", handler.List)

			w := performRequest(r, "GET", "/pods"+tc.queryString, nil)

			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedError != "" {
				assert.Contains(t, w.Body.String(), tc.expectedError)
			}

			svc.AssertExpectations(t)
		})
	}
}

func TestPodHandler_Get(t *testing.T) {
	type testCase struct {
		name         string
		namespace    string
		podName      string
		mockBehavior func(m *mock.PodServiceMock)
		expectedCode int
		expectError  bool
	}

	tests := []testCase{
		{
			name:      "Success",
			namespace: "default",
			podName:   "my-pod",
			mockBehavior: func(m *mock.PodServiceMock) {
				m.On("GetPod", testifyMock.Anything, "default", "my-pod").
					Return(&models.PodDetails{PodListInfo: models.PodListInfo{Name: "my-pod"}}, nil)
			},
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:      "Not Found",
			namespace: "default",
			podName:   "missing",
			mockBehavior: func(m *mock.PodServiceMock) {
				m.On("GetPod", testifyMock.Anything, "default", "missing").
					Return((*models.PodDetails)(nil), services.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectError:  true,
		},
		{
			name:      "Internal error",
			namespace: "default",
			podName:   "error-pod",
			mockBehavior: func(m *mock.PodServiceMock) {
				m.On("GetPod", testifyMock.Anything, "default", "error-pod").
					Return((*models.PodDetails)(nil), assert.AnError)
			},
			expectedCode: http.StatusInternalServerError,
			expectError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mock.PodServiceMock)
			tc.mockBehavior(svc)

			handler := NewPodHandler(svc)
			r := setupRouter()
			r.GET("/:namespace/:name", handler.Get)

			url := fmt.Sprintf("/%s/%s", tc.namespace, tc.podName)
			w := performRequest(r, "GET", url, nil)

			assert.Equal(t, tc.expectedCode, w.Code)
			svc.AssertExpectations(t)
		})
	}
}
