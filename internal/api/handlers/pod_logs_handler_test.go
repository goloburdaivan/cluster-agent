package handlers

import (
	"bytes"
	"cluster-agent/internal/services/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
)

func TestNewPodLogsHandler(t *testing.T) {
	svc := new(mock.PodLogsServiceMock)
	handler := NewPodLogsHandler(svc)
	assert.NotNil(t, handler)
	assert.Equal(t, svc, handler.service)
}

func TestPodLogsHandler_StreamLogs_UpgradeFailure(t *testing.T) {
	svc := new(mock.PodLogsServiceMock)
	handler := NewPodLogsHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/:namespace/:name/logs", handler.StreamLogs)

	req, _ := http.NewRequest("GET", "/default/my-pod/logs", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	svc.AssertExpectations(t)
}

func TestPodLogsHandler_StreamLogs_ServiceError(t *testing.T) {
	svc := new(mock.PodLogsServiceMock)
	svc.On("StreamLogs", testifyMock.Anything, "default", "my-pod", "main").
		Return(nil, assert.AnError)

	handler := NewPodLogsHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/:namespace/:name/logs", handler.StreamLogs)

	server := httptest.NewServer(r)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/default/my-pod/logs?container=main"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open websocket connection: %v", err)
	}
	defer ws.Close()

	_, message, err := ws.ReadMessage()
	assert.NoError(t, err)
	assert.Contains(t, string(message), "Error init logs stream")

	svc.AssertExpectations(t)
}

func TestPodLogsHandler_StreamLogs_Success(t *testing.T) {
	logContent := "log line 1\nlog line 2\n"
	mockReader := io.NopCloser(bytes.NewBufferString(logContent))

	svc := new(mock.PodLogsServiceMock)
	svc.On("StreamLogs", testifyMock.Anything, "default", "my-pod", "main").
		Return(mockReader, nil)

	handler := NewPodLogsHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/:namespace/:name/logs", handler.StreamLogs)

	server := httptest.NewServer(r)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/default/my-pod/logs?container=main"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open websocket connection: %v", err)
	}
	defer ws.Close()

	ws.SetReadDeadline(time.Now().Add(2 * time.Second))

	_, message1, err := ws.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, "log line 1\n", string(message1))

	_, message2, err := ws.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, "log line 2\n", string(message2))

	svc.AssertExpectations(t)
}

func TestProcessLogsStream_WriteError(t *testing.T) {
	logContent := "log line 1\nlog line 2\n"
	mockReader := io.NopCloser(bytes.NewBufferString(logContent))

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		ws.Close()
		processLogsStream(ws, mockReader)
	})

	server := httptest.NewServer(r)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/test"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open websocket connection: %v", err)
	}
	ws.Close()
}
