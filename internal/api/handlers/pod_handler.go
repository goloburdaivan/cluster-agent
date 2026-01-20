package handlers

import (
	"cluster-agent/internal/api/responses"
	"cluster-agent/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PodHandler struct {
	podService services.PodService
}

func NewPodHandler(podService services.PodService) *PodHandler {
	return &PodHandler{
		podService: podService,
	}
}

func (handler *PodHandler) List(c *gin.Context) {
	namespace := c.Query("namespace")

	pods, err := handler.podService.GetPods(c.Request.Context(), namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(pods))
}

func (handler *PodHandler) Get(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	pod, err := handler.podService.GetPod(c.Request.Context(), namespace, name)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, responses.Error(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(pod))
}
