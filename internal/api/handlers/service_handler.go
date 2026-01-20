package handlers

import (
	"cluster-agent/internal/api/responses"
	"cluster-agent/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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
	namespace := c.Query("namespace")

	result, err := handler.service.List(c.Request.Context(), namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(result))
}

func (handler *ServiceHandler) Get(c *gin.Context) {
	data, err := handler.service.Get(c.Request.Context(), c.Param("namespace"), c.Param("name"))
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, responses.Error(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(data))
}
