package handlers

import (
	"cluster-agent/internal/api/responses"
	"cluster-agent/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	rbacv1 "k8s.io/api/rbac/v1"
)

type RoleBindingHandler struct{ service services.RoleBindingService }

func NewRoleBindingHandler(s services.RoleBindingService) *RoleBindingHandler {
	return &RoleBindingHandler{service: s}
}

func (h *RoleBindingHandler) List(c *gin.Context) {
	namespace := c.Query("namespace")

	data, err := h.service.List(c.Request.Context(), namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, responses.Success(data))
}

func (h *RoleBindingHandler) Get(c *gin.Context) {
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

func (h *RoleBindingHandler) Create(c *gin.Context) {
	var rb rbacv1.RoleBinding

	if err := c.ShouldBindJSON(&rb); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error()))
		return
	}

	err := h.service.Create(c.Request.Context(), &rb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, responses.Success(rb))
}

func (h *RoleBindingHandler) Update(c *gin.Context) {
	var rb rbacv1.RoleBinding

	if err := c.ShouldBindJSON(&rb); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error()))
		return
	}

	err := h.service.Update(c.Request.Context(), &rb)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, responses.Error(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(rb))
}

func (h *RoleBindingHandler) Delete(c *gin.Context) {
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
