package handlers

import (
	"cluster-agent/internal/api/responses"
	"cluster-agent/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	batchv1 "k8s.io/api/batch/v1"
)

type JobHandler struct {
	jobService services.JobService
}

func NewJobHandler(jobService services.JobService) *JobHandler {
	return &JobHandler{
		jobService: jobService,
	}
}

func (handler *JobHandler) List(c *gin.Context) {
	namespace := c.Query("namespace")

	jobs, err := handler.jobService.GetJobs(c.Request.Context(), namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(jobs))
}

func (handler *JobHandler) Get(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	job, err := handler.jobService.GetJob(c.Request.Context(), namespace, name)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, responses.Error(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(job))
}

func (handler *JobHandler) Create(c *gin.Context) {
	var job batchv1.Job

	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error()))
		return
	}

	err := handler.jobService.CreateJob(c.Request.Context(), &job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, responses.Success(job))
}

func (handler *JobHandler) Update(c *gin.Context) {
	var job batchv1.Job

	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error()))
		return
	}

	err := handler.jobService.UpdateJob(c.Request.Context(), &job)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, responses.Error(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.Success(job))
}

func (handler *JobHandler) Delete(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	err := handler.jobService.DeleteJob(c.Request.Context(), namespace, name)
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
