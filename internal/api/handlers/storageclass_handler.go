package handlers

import (
	"cluster-agent/internal/api/responses"
	"cluster-agent/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	storagev1 "k8s.io/api/storage/v1"
)

type StorageClassHandler struct{ service services.StorageClassService }

func NewStorageClassHandler(s services.StorageClassService) *StorageClassHandler {
	return &StorageClassHandler{service: s}
}

func (h *StorageClassHandler) List(c *gin.Context) {
	data, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, responses.Success(data))
}

func (h *StorageClassHandler) Get(c *gin.Context) {
	name := c.Param("name")

	data, err := h.service.Get(c.Request.Context(), name)
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

func (h *StorageClassHandler) Create(c *gin.Context) {
	var sc storagev1.StorageClass

	if err := c.ShouldBindJSON(&sc); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error()))
		return
	}

	err := h.service.Create(c.Request.Context(), &sc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, responses.Success(sc))
}

func (h *StorageClassHandler) Update(c *gin.Context) {
	var sc storagev1.StorageClass

	if err := c.ShouldBindJSON(&sc); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error()))
		return
	}

	err := h.service.Update(c.Request.Context(), &sc)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, responses.Error(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(sc))
}

func (h *StorageClassHandler) Delete(c *gin.Context) {
	name := c.Param("name")

	err := h.service.Delete(c.Request.Context(), name)
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
