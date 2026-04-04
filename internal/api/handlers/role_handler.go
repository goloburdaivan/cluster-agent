package handlers

import (
	"cluster-agent/internal/api/responses"
	"cluster-agent/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	rbacv1 "k8s.io/api/rbac/v1"
)

type RoleHandler struct{ service services.RoleService }

func NewRoleHandler(s services.RoleService) *RoleHandler {
	return &RoleHandler{service: s}
}

func (h *RoleHandler) List(c *gin.Context) {
	namespace := c.Query("namespace")

	data, err := h.service.List(c.Request.Context(), namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}
	c.JSON(http.StatusOK, responses.Success(data))
}

func (h *RoleHandler) Get(c *gin.Context) {
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

func (h *RoleHandler) Create(c *gin.Context) {
	var role rbacv1.Role

	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error()))
		return
	}

	err := h.service.Create(c.Request.Context(), &role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, responses.Success(role))
}

func (h *RoleHandler) Update(c *gin.Context) {
	var role rbacv1.Role

	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error()))
		return
	}

	err := h.service.Update(c.Request.Context(), &role)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, responses.Error(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(role))
}

func (h *RoleHandler) Delete(c *gin.Context) {
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
