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

func TestPvcHandler_List(t *testing.T) {
	type testCase struct {
		name          string
		queryString   string
		mockBehavior  func(m *mock.PVCServiceMock)
		expectedCode  int
		expectedError string
	}

	tests := []testCase{
		{
			name:        "Success",
			queryString: "?namespace=default",
			mockBehavior: func(m *mock.PVCServiceMock) {
				m.On("List", testifyMock.Anything, "default").
					Return([]models.PVCListInfo{{Name: "pvc1"}}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "Success without namespace",
			queryString: "",
			mockBehavior: func(m *mock.PVCServiceMock) {
				m.On("List", testifyMock.Anything, "").
					Return([]models.PVCListInfo{{Name: "pvc1"}}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "Internal error",
			queryString: "?namespace=default",
			mockBehavior: func(m *mock.PVCServiceMock) {
				m.On("List", testifyMock.Anything, "default").
					Return([]models.PVCListInfo(nil), assert.AnError)
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "assert.AnError general error for testing",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mock.PVCServiceMock)
			tc.mockBehavior(svc)

			handler := NewPvcHandler(svc)
			r := setupRouter()
			r.GET("/pvcs", handler.List)

			w := performRequest(r, "GET", "/pvcs"+tc.queryString, nil)

			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedError != "" {
				assert.Contains(t, w.Body.String(), tc.expectedError)
			}

			svc.AssertExpectations(t)
		})
	}
}

func TestPvcHandler_Get(t *testing.T) {
	type testCase struct {
		name         string
		namespace    string
		pvcName      string
		mockBehavior func(m *mock.PVCServiceMock)
		expectedCode int
		expectError  bool
	}

	tests := []testCase{
		{
			name:      "Success",
			namespace: "default",
			pvcName:   "my-pvc",
			mockBehavior: func(m *mock.PVCServiceMock) {
				m.On("Get", testifyMock.Anything, "default", "my-pvc").
					Return(&models.PVCDetails{PVCListInfo: models.PVCListInfo{Name: "my-pvc"}}, nil)
			},
			expectedCode: http.StatusOK,
			expectError:  false,
		},
		{
			name:      "Not Found",
			namespace: "default",
			pvcName:   "missing",
			mockBehavior: func(m *mock.PVCServiceMock) {
				m.On("Get", testifyMock.Anything, "default", "missing").
					Return((*models.PVCDetails)(nil), services.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
			expectError:  true,
		},
		{
			name:      "Internal error",
			namespace: "default",
			pvcName:   "error-pvc",
			mockBehavior: func(m *mock.PVCServiceMock) {
				m.On("Get", testifyMock.Anything, "default", "error-pvc").
					Return((*models.PVCDetails)(nil), assert.AnError)
			},
			expectedCode: http.StatusInternalServerError,
			expectError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mock.PVCServiceMock)
			tc.mockBehavior(svc)

			handler := NewPvcHandler(svc)
			r := setupRouter()
			r.GET("/:namespace/:name", handler.Get)

			url := fmt.Sprintf("/%s/%s", tc.namespace, tc.pvcName)
			w := performRequest(r, "GET", url, nil)

			assert.Equal(t, tc.expectedCode, w.Code)
			svc.AssertExpectations(t)
		})
	}
}
