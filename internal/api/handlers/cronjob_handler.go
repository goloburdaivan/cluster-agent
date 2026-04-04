package handlers

import (
	"cluster-agent/internal/api/responses"
	"cluster-agent/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	batchv1 "k8s.io/api/batch/v1"
)

type CronJobHandler struct {
	cronJobService services.CronJobService
}

func NewCronJobHandler(cronJobService services.CronJobService) *CronJobHandler {
	return &CronJobHandler{
		cronJobService: cronJobService,
	}
}

func (handler *CronJobHandler) List(c *gin.Context) {
	namespace := c.Query("namespace")

	cronjobs, err := handler.cronJobService.GetCronJobs(c.Request.Context(), namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(cronjobs))
}

func (handler *CronJobHandler) Get(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	cronjob, err := handler.cronJobService.GetCronJob(c.Request.Context(), namespace, name)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, responses.Error(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(cronjob))
}

func (handler *CronJobHandler) Create(c *gin.Context) {
	var cronjob batchv1.CronJob

	if err := c.ShouldBindJSON(&cronjob); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error()))
		return
	}

	err := handler.cronJobService.CreateCronJob(c.Request.Context(), &cronjob)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, responses.Success(cronjob))
}

func (handler *CronJobHandler) Update(c *gin.Context) {
	var cronjob batchv1.CronJob

	if err := c.ShouldBindJSON(&cronjob); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error()))
		return
	}

	err := handler.cronJobService.UpdateCronJob(c.Request.Context(), &cronjob)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, responses.Error(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(cronjob))
}

func (handler *CronJobHandler) Delete(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	err := handler.cronJobService.DeleteCronJob(c.Request.Context(), namespace, name)
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
