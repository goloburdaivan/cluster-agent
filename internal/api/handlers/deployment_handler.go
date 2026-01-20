package handlers

import (
	"cluster-agent/internal/api/responses"
	"cluster-agent/internal/models"
	"cluster-agent/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/apps/v1"
)

type DeploymentHandler struct {
	deploymentService services.DeploymentService
}

func NewDeploymentHandler(deploymentService services.DeploymentService) *DeploymentHandler {
	return &DeploymentHandler{
		deploymentService: deploymentService,
	}
}

func (handler *DeploymentHandler) List(c *gin.Context) {
	namespace := c.Query("namespace")

	deployments, err := handler.deploymentService.GetDeployments(c.Request.Context(), namespace)

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(deployments))
}

func (handler *DeploymentHandler) Get(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	deployment, err := handler.deploymentService.GetDeployment(c.Request.Context(), namespace, name)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, responses.Error(err.Error()))
			return
		}

		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(deployment))
}

func (handler *DeploymentHandler) Create(c *gin.Context) {
	var deployment v1.Deployment

	if err := c.ShouldBindJSON(&deployment); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error()))
		return
	}

	err := handler.deploymentService.CreateDeployment(c.Request.Context(), &deployment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, responses.Success(deployment))
}

func (handler *DeploymentHandler) Delete(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := handler.deploymentService.DeleteDeployment(c.Request.Context(), namespace, name)
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

func (handler *DeploymentHandler) ScaleDeployment(c *gin.Context) {
	var request models.ScaleDeploymentParams
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error()))
		return
	}

	err := handler.deploymentService.ScaleDeployment(c.Request.Context(), request)
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
