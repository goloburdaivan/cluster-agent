package handlers

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/mock"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestNodeHandler_List(t *testing.T) {
	type testCase struct {
		name          string
		mockBehavior  func(m *mock.NodeServiceMock)
		expectedCode  int
		expectedError string
	}

	tests := []testCase{
		{
			name: "Success",
			mockBehavior: func(m *mock.NodeServiceMock) {
				m.On("GetNodes", testifyMock.Anything).
					Return([]models.Node{{Name: "node1", Status: models.NodeStatusReady}}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Error",
			mockBehavior: func(m *mock.NodeServiceMock) {
				m.On("GetNodes", testifyMock.Anything).
					Return([]models.Node(nil), assert.AnError)
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "assert.AnError general error for testing",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mock.NodeServiceMock)
			tc.mockBehavior(svc)

			handler := NewNodeHandler(svc)
			r := setupRouter()
			r.GET("/nodes", handler.List)

			w := performRequest(r, "GET", "/nodes", nil)

			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedError != "" {
				assert.Contains(t, w.Body.String(), tc.expectedError)
			}

			svc.AssertExpectations(t)
		})
	}
}
