package handlers

import (
	"cluster-agent/internal/api/responses"
	"cluster-agent/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
)

type DaemonSetHandler struct {
	daemonSetService services.DaemonSetService
}

func NewDaemonSetHandler(daemonSetService services.DaemonSetService) *DaemonSetHandler {
	return &DaemonSetHandler{
		daemonSetService: daemonSetService,
	}
}

func (handler *DaemonSetHandler) List(c *gin.Context) {
	namespace := c.Query("namespace")

	daemonsets, err := handler.daemonSetService.GetDaemonSets(c.Request.Context(), namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(daemonsets))
}

func (handler *DaemonSetHandler) Get(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	daemonset, err := handler.daemonSetService.GetDaemonSet(c.Request.Context(), namespace, name)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, responses.Error(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(daemonset))
}

func (handler *DaemonSetHandler) Create(c *gin.Context) {
	var daemonset appsv1.DaemonSet

	if err := c.ShouldBindJSON(&daemonset); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error()))
		return
	}

	err := handler.daemonSetService.CreateDaemonSet(c.Request.Context(), &daemonset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, responses.Success(daemonset))
}

func (handler *DaemonSetHandler) Update(c *gin.Context) {
	var daemonset appsv1.DaemonSet

	if err := c.ShouldBindJSON(&daemonset); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error()))
		return
	}

	err := handler.daemonSetService.UpdateDaemonSet(c.Request.Context(), &daemonset)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, responses.Error(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(daemonset))
}

func (handler *DaemonSetHandler) Delete(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	err := handler.daemonSetService.DeleteDaemonSet(c.Request.Context(), namespace, name)
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
