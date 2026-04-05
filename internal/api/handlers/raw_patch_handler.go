package handlers

import (
	"cluster-agent/internal/api/responses"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

type RawPatchHandler struct {
	client dynamic.Interface
}

func NewRawPatchHandler(client dynamic.Interface) *RawPatchHandler {
	return &RawPatchHandler{client: client}
}

func (h *RawPatchHandler) PatchNamespaced(gvr schema.GroupVersionResource) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.Error("failed to read request body"))
			return
		}

		result, err := h.client.Resource(gvr).Namespace(namespace).Patch(
			c.Request.Context(), name, types.MergePatchType, body, metav1.PatchOptions{},
		)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, responses.Error(err.Error()))
			return
		}

		c.JSON(http.StatusOK, responses.Success(result.Object))
	}
}

func (h *RawPatchHandler) PatchClusterScoped(gvr schema.GroupVersionResource) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.Error("failed to read request body"))
			return
		}

		result, err := h.client.Resource(gvr).Patch(
			c.Request.Context(), name, types.MergePatchType, body, metav1.PatchOptions{},
		)

		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, responses.Error(err.Error()))
			return
		}

		if metadata, ok := result.Object["metadata"].(map[string]interface{}); ok {
			delete(metadata, "managedFields")
		}

		c.JSON(http.StatusOK, responses.Success(result.Object))
	}
}

func (h *RawPatchHandler) GetNamespaced(gvr schema.GroupVersionResource) gin.HandlerFunc {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		name := c.Param("name")

		result, err := h.client.Resource(gvr).Namespace(namespace).Get(
			c.Request.Context(), name, metav1.GetOptions{},
		)

		if err != nil {
			c.JSON(http.StatusNotFound, responses.Error(err.Error()))
			return
		}

		if metadata, ok := result.Object["metadata"].(map[string]interface{}); ok {
			delete(metadata, "managedFields")
		}

		c.JSON(http.StatusOK, responses.Success(result.Object))
	}
}

func (h *RawPatchHandler) GetClusterScoped(gvr schema.GroupVersionResource) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")

		result, err := h.client.Resource(gvr).Get(
			c.Request.Context(), name, metav1.GetOptions{},
		)
		if err != nil {
			c.JSON(http.StatusNotFound, responses.Error(err.Error()))
			return
		}

		c.JSON(http.StatusOK, responses.Success(result.Object))
	}
}
