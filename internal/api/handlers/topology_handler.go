package handlers

import (
	"cluster-agent/internal/api/requests"
	"cluster-agent/internal/services"
	"cluster-agent/internal/services/topology"
	"github.com/gin-gonic/gin"
	"net/http"
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
	var request requests.NamespaceQueryRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	snapshot, err := h.snapshotter.TakeClusterSnapshot(request.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := h.service.BuildFromSnapshot(snapshot)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
