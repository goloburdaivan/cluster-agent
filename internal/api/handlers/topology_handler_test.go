package handlers

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
	"cluster-agent/internal/services/mock"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestTopologyHandler_Get(t *testing.T) {
	type testCase struct {
		name                 string
		queryString          string
		mockSnapshotBehavior func(m *mock.SnapshotServiceMock)
		mockTopologyBehavior func(m *mock.TopologyServiceMock)
		expectedCode         int
		expectedError        string
	}

	tests := []testCase{
		{
			name:        "Success",
			queryString: "?namespace=default",
			mockSnapshotBehavior: func(m *mock.SnapshotServiceMock) {
				m.On("TakeClusterSnapshot", "default").
					Return(&models.ClusterSnapshot{}, nil)
			},
			mockTopologyBehavior: func(m *mock.TopologyServiceMock) {
				m.On("BuildFromSnapshot", testifyMock.Anything, testifyMock.Anything).
					Return(&graph.Graph{Nodes: []graph.Node{}, Edges: []graph.Edge{}}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "Success without namespace",
			queryString: "",
			mockSnapshotBehavior: func(m *mock.SnapshotServiceMock) {
				m.On("TakeClusterSnapshot", "").
					Return(&models.ClusterSnapshot{}, nil)
			},
			mockTopologyBehavior: func(m *mock.TopologyServiceMock) {
				m.On("BuildFromSnapshot", testifyMock.Anything, testifyMock.Anything).
					Return(&graph.Graph{Nodes: []graph.Node{}, Edges: []graph.Edge{}}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:        "Snapshot error",
			queryString: "?namespace=default",
			mockSnapshotBehavior: func(m *mock.SnapshotServiceMock) {
				m.On("TakeClusterSnapshot", "default").
					Return((*models.ClusterSnapshot)(nil), assert.AnError)
			},
			mockTopologyBehavior: func(m *mock.TopologyServiceMock) {
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "assert.AnError general error for testing",
		},
		{
			name:        "Build topology error",
			queryString: "?namespace=default",
			mockSnapshotBehavior: func(m *mock.SnapshotServiceMock) {
				m.On("TakeClusterSnapshot", "default").
					Return(&models.ClusterSnapshot{}, nil)
			},
			mockTopologyBehavior: func(m *mock.TopologyServiceMock) {
				m.On("BuildFromSnapshot", testifyMock.Anything, testifyMock.Anything).
					Return((*graph.Graph)(nil), assert.AnError)
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "assert.AnError general error for testing",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			snapshotSvc := new(mock.SnapshotServiceMock)
			topologySvc := new(mock.TopologyServiceMock)

			tc.mockSnapshotBehavior(snapshotSvc)
			tc.mockTopologyBehavior(topologySvc)

			handler := NewTopologyHandler(topologySvc, snapshotSvc)
			r := setupRouter()
			r.GET("/topology", handler.Get)

			w := performRequest(r, "GET", "/topology"+tc.queryString, nil)

			assert.Equal(t, tc.expectedCode, w.Code)

			if tc.expectedError != "" {
				assert.Contains(t, w.Body.String(), tc.expectedError)
			}

			snapshotSvc.AssertExpectations(t)
			topologySvc.AssertExpectations(t)
		})
	}
}
