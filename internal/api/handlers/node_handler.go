package handlers

import (
	"cluster-agent/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NodeHandler struct {
	nodeService services.NodeService
}

func NewNodeHandler(nodeService services.NodeService) *NodeHandler {
	return &NodeHandler{
		nodeService: nodeService,
	}
}

func (h *NodeHandler) List(c *gin.Context) {
	nodes, err := h.nodeService.GetNodes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"nodes": nodes,
	})
}
