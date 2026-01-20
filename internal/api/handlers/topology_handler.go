package handlers

import (
	"cluster-agent/internal/api/responses"
	"cluster-agent/internal/services"
	"cluster-agent/internal/services/topology"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TopologyHandler struct {
	service     topology.Service
	snapshotter services.SnapshotService
}

func NewTopologyHandler(
	service topology.Service,
	snapshotter services.SnapshotService,
) *TopologyHandler {
	return &TopologyHandler{
		service:     service,
		snapshotter: snapshotter,
	}
}

func (h *TopologyHandler) Get(c *gin.Context) {
	namespace := c.Query("namespace")

	snapshot, err := h.snapshotter.TakeClusterSnapshot(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	result, err := h.service.BuildFromSnapshot(c.Request.Context(), snapshot)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(result))
}
