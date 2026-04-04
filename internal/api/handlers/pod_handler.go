package handlers

import (
	"cluster-agent/internal/api/responses"
	"cluster-agent/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
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

func (handler *PodHandler) Create(c *gin.Context) {
	var pod corev1.Pod

	if err := c.ShouldBindJSON(&pod); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error()))
		return
	}

	err := handler.podService.CreatePod(c.Request.Context(), &pod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, responses.Success(pod))
}

func (handler *PodHandler) Update(c *gin.Context) {
	var pod corev1.Pod

	if err := c.ShouldBindJSON(&pod); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error()))
		return
	}

	err := handler.podService.UpdatePod(c.Request.Context(), &pod)
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

func (handler *PodHandler) Delete(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	err := handler.podService.DeletePod(c.Request.Context(), namespace, name)
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
