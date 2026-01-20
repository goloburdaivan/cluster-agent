package handlers

import (
	"cluster-agent/internal/api/responses"
	"cluster-agent/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(result))
}
