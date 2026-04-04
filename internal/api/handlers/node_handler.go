package handlers

import (
	"cluster-agent/internal/api/responses"
	"cluster-agent/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(nodes))
}

func (h *NodeHandler) Get(c *gin.Context) {
	name := c.Param("name")

	node, err := h.nodeService.GetNode(c.Request.Context(), name)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, responses.Error(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(node))
}
