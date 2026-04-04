package handlers

import (
	"cluster-agent/internal/api/responses"
	"cluster-agent/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
)

type PvcHandler struct{ service services.PVCService }

func NewPvcHandler(s services.PVCService) *PvcHandler { return &PvcHandler{service: s} }

func (h *PvcHandler) List(c *gin.Context) {
	namespace := c.Query("namespace")

	data, err := h.service.List(c.Request.Context(), namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, responses.Success(data))
}

func (h *PvcHandler) Get(c *gin.Context) {
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

func (h *PvcHandler) Create(c *gin.Context) {
	var pvc corev1.PersistentVolumeClaim

	if err := c.ShouldBindJSON(&pvc); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error()))
		return
	}

	err := h.service.Create(c.Request.Context(), &pvc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, responses.Success(pvc))
}

func (h *PvcHandler) Update(c *gin.Context) {
	var pvc corev1.PersistentVolumeClaim

	if err := c.ShouldBindJSON(&pvc); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error()))
		return
	}

	err := h.service.Update(c.Request.Context(), &pvc)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, responses.Error(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(pvc))
}

func (h *PvcHandler) Delete(c *gin.Context) {
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
