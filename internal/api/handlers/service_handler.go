package handlers

import (
	"cluster-agent/internal/api/requests"
	"cluster-agent/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ServiceHandler struct {
	service services.KubernetesServiceService
}

func NewServiceHandler(service services.KubernetesServiceService) *ServiceHandler {
	return &ServiceHandler{
		service: service,
	}
}

func (handler *ServiceHandler) List(c *gin.Context) {
	var request requests.NamespaceQueryRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := handler.service.List(c.Request.Context(), request.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

func (handler *ServiceHandler) Get(c *gin.Context) {
	data, err := handler.service.Get(c.Request.Context(), c.Param("namespace"), c.Param("name"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}
