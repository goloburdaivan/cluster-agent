package handlers

import (
	"cluster-agent/internal/api/requests"
	"cluster-agent/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
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
	var request requests.NamespaceQueryRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	pods, err := handler.podService.GetPods(c.Request.Context(), request.Namespace)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": pods,
	})
}

func (handler *PodHandler) Get(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	pod, err := handler.podService.GetPod(c.Request.Context(), namespace, name)
	if err != nil {
		c.JSON(404, gin.H{
			"error": "Pod not found: " + err.Error(),
		})
		return
	}

	c.JSON(200, pod)
}
