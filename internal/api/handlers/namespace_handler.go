package handlers

import (
	"cluster-agent/internal/api/responses"
	"cluster-agent/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
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

func (handler *NamespaceHandler) Get(c *gin.Context) {
	name := c.Param("name")

	result, err := handler.namespaceService.GetNamespace(c.Request.Context(), name)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, responses.Error(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(result))
}

func (handler *NamespaceHandler) Create(c *gin.Context) {
	var namespace corev1.Namespace

	if err := c.ShouldBindJSON(&namespace); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error()))
		return
	}

	err := handler.namespaceService.CreateNamespace(c.Request.Context(), &namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, responses.Success(namespace))
}

func (handler *NamespaceHandler) Delete(c *gin.Context) {
	name := c.Param("name")

	err := handler.namespaceService.DeleteNamespace(c.Request.Context(), name)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, responses.Error(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success("OK"))
}
