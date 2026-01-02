package handlers

import (
	"cluster-agent/internal/api/requests"
	"cluster-agent/internal/services"
	"github.com/gin-gonic/gin"
)

type ConfigMapHandler struct{ service services.ConfigMapService }

func NewConfigMapHandler(s services.ConfigMapService) *ConfigMapHandler {
	return &ConfigMapHandler{service: s}
}

func (h *ConfigMapHandler) List(c *gin.Context) {
	var req requests.NamespaceQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	data, err := h.service.List(c.Request.Context(), req.Namespace)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": data})
}

func (h *ConfigMapHandler) Get(c *gin.Context) {
	data, err := h.service.Get(c.Request.Context(), c.Param("namespace"), c.Param("name"))
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, data)
}
