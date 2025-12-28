package handlers

import (
	"cluster-agent/internal/api/requests"
	"cluster-agent/internal/models"
	"cluster-agent/internal/services"
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
	var query requests.NamespaceQueryRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	deployments, err := handler.deploymentService.GetDeployments(c.Request.Context(), query.Namespace)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": deployments,
	})
}

func (handler *DeploymentHandler) Create(c *gin.Context) {
	var deployment v1.Deployment

	if err := c.ShouldBindJSON(&deployment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	err := handler.deploymentService.CreateDeployment(c.Request.Context(), &deployment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create deployment",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Deployment created",
	})
}

func (handler *DeploymentHandler) ScaleDeployment(c *gin.Context) {
	var request models.ScaleDeploymentParams
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	err := handler.deploymentService.ScaleDeployment(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Deployment scaled",
	})
}
