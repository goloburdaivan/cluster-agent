package handlers

import (
	"cluster-agent/internal/api/responses"
	"cluster-agent/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	networkingv1 "k8s.io/api/networking/v1"
)

type IngressHandler struct{ service services.IngressService }

func NewIngressHandler(s services.IngressService) *IngressHandler { return &IngressHandler{service: s} }

func (h *IngressHandler) List(c *gin.Context) {
	namespace := c.Query("namespace")

	data, err := h.service.List(c.Request.Context(), namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, responses.Success(data))
}

func (h *IngressHandler) Get(c *gin.Context) {
	data, err := h.service.Get(c.Request.Context(), c.Param("namespace"), c.Param("name"))
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

func (h *IngressHandler) Create(c *gin.Context) {
	var ingress networkingv1.Ingress

	if err := c.ShouldBindJSON(&ingress); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error()))
		return
	}

	err := h.service.Create(c.Request.Context(), &ingress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, responses.Success(ingress))
}

func (h *IngressHandler) Update(c *gin.Context) {
	var ingress networkingv1.Ingress

	if err := c.ShouldBindJSON(&ingress); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error()))
		return
	}

	err := h.service.Update(c.Request.Context(), &ingress)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, responses.Error(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(ingress))
}

func (h *IngressHandler) Delete(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	err := h.service.Delete(c.Request.Context(), namespace, name)
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
