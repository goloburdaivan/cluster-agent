package handlers

import (
	"cluster-agent/internal/api/responses"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func parseResponse[T any](t *testing.T, w *httptest.ResponseRecorder) responses.Response[T] {
	var response responses.Response[T]
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "failed to unmarshal response")
	return response
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}
