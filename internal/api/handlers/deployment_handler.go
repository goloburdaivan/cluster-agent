package handlers

import (
	"cluster-agent/internal/api/requests"
	"cluster-agent/internal/models"
	"cluster-agent/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
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
