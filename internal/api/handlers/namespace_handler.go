package handlers

import (
	"cluster-agent/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NamespaceHandler struct {
	namespaceService services.NamespaceService
}

func NewNamespaceHandler(namespaceService services.NamespaceService) *NamespaceHandler {
	return &NamespaceHandler{
		namespaceService: namespaceService,
	}
}

func (handler *NamespaceHandler) List(c *gin.Context) {
	result, err := handler.namespaceService.GetNamespaces(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, result)
}
