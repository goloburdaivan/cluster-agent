package handlers

import (
	"cluster-agent/internal/api/responses"
	"cluster-agent/internal/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NetworkInspectorHandler struct {
	service services.NetworkInspectorService
}

func NewNetworkInspectorHandler(service services.NetworkInspectorService) *NetworkInspectorHandler {
	return &NetworkInspectorHandler{
		service: service,
	}
}

func (h *NetworkInspectorHandler) GetConnections(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	container := c.Query("container")

	connections, err := h.service.GetPodNetworkConnections(c.Request.Context(), namespace, name, container)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(fmt.Sprintf("failed to get network connections: %v", err)))
		return
	}

	c.JSON(http.StatusOK, responses.Success(connections))
}
