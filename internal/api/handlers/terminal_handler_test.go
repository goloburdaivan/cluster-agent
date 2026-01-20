package handlers

import (
	"bytes"
	"cluster-agent/internal/services/mock"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/remotecommand"
)

func TestNewTerminalHandler(t *testing.T) {
	svc := new(mock.TerminalServiceMock)
	handler := NewTerminalHandler(svc)
	assert.NotNil(t, handler)
	assert.Equal(t, svc, handler.service)
}

func TestTerminalHandler_Exec_UpgradeFailure(t *testing.T) {
	svc := new(mock.TerminalServiceMock)
	handler := NewTerminalHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/:namespace/:name/exec", handler.Exec)

	req, _ := http.NewRequest("GET", "/default/my-pod/exec", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	svc.AssertExpectations(t)
}

func TestTerminalHandler_Exec_ServiceError(t *testing.T) {
	svc := new(mock.TerminalServiceMock)
	svc.On("GetAuthExecutor", "default", "my-pod", "main").
		Return(nil, assert.AnError)

	handler := NewTerminalHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/:namespace/:name/exec", handler.Exec)

	server := httptest.NewServer(r)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/default/my-pod/exec?container=main"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open websocket connection: %v", err)
	}
	defer ws.Close()

	_, message, err := ws.ReadMessage()
	assert.NoError(t, err)
	assert.Contains(t, string(message), "Error init executor")

	svc.AssertExpectations(t)
}

type MockExecutor struct {
	streamErr error
	output    string
}

func (m *MockExecutor) Stream(options remotecommand.StreamOptions) error {
	return m.streamErr
}

func (m *MockExecutor) StreamWithContext(ctx context.Context, options remotecommand.StreamOptions) error {
	if m.output != "" {
		options.Stdout.Write([]byte(m.output))
	}
	return m.streamErr
}

func TestTerminalHandler_Exec_StreamError(t *testing.T) {
	mockExec := &MockExecutor{streamErr: assert.AnError}

	svc := new(mock.TerminalServiceMock)
	svc.On("GetAuthExecutor", "default", "my-pod", "main").
		Return(mockExec, nil)

	handler := NewTerminalHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/:namespace/:name/exec", handler.Exec)

	server := httptest.NewServer(r)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/default/my-pod/exec?container=main"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open websocket connection: %v", err)
	}
	defer ws.Close()

	ws.SetReadDeadline(time.Now().Add(2 * time.Second))
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		if strings.Contains(string(message), "Session ended with error") {
			break
		}
	}

	svc.AssertExpectations(t)
}

func TestTerminalHandler_Exec_Success(t *testing.T) {
	mockExec := &MockExecutor{output: "hello terminal"}

	svc := new(mock.TerminalServiceMock)
	svc.On("GetAuthExecutor", "default", "my-pod", "main").
		Return(mockExec, nil)

	handler := NewTerminalHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/:namespace/:name/exec", handler.Exec)

	server := httptest.NewServer(r)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/default/my-pod/exec?container=main"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open websocket connection: %v", err)
	}
	defer ws.Close()

	ws.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, message, err := ws.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, "hello terminal", string(message))

	svc.AssertExpectations(t)
}

func TestWsPtyHandler_Write(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	var capturedHandler *WsPtyHandler

	r.GET("/test", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		capturedHandler = &WsPtyHandler{
			conn:        ws,
			resizeEvent: make(chan remotecommand.TerminalSize),
		}

		n, err := capturedHandler.Write([]byte("test output"))
		assert.NoError(t, err)
		assert.Equal(t, 11, n)

		ws.Close()
	})

	server := httptest.NewServer(r)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/test"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open websocket connection: %v", err)
	}
	defer ws.Close()

	ws.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, message, err := ws.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, "test output", string(message))
}

func TestWsPtyHandler_Write_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/test", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		handler := &WsPtyHandler{
			conn:        ws,
			resizeEvent: make(chan remotecommand.TerminalSize),
		}

		ws.Close()

		_, err = handler.Write([]byte("test"))
		assert.Error(t, err)
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

func TestWsPtyHandler_Read_FromBuffer(t *testing.T) {
	handler := &WsPtyHandler{
		conn:        nil,
		resizeEvent: make(chan remotecommand.TerminalSize),
		readBuf:     *bytes.NewBufferString("buffered data"),
	}

	buf := make([]byte, 20)
	n, err := handler.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, 13, n)
	assert.Equal(t, "buffered data", string(buf[:n]))
}

func TestWsPtyHandler_Read_StdinMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	done := make(chan struct{})

	r.GET("/test", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		handler := &WsPtyHandler{
			conn:        ws,
			resizeEvent: make(chan remotecommand.TerminalSize, 1),
		}

		buf := make([]byte, 100)
		n, err := handler.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, "user input", string(buf[:n]))

		close(done)
	})

	server := httptest.NewServer(r)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/test"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open websocket connection: %v", err)
	}
	defer ws.Close()

	msg := TerminalMessage{Op: "stdin", Data: "user input"}
	data, _ := json.Marshal(msg)
	ws.WriteMessage(websocket.TextMessage, data)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for read")
	}
}

func TestWsPtyHandler_Read_ResizeMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	done := make(chan struct{})

	r.GET("/test", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		handler := &WsPtyHandler{
			conn:        ws,
			resizeEvent: make(chan remotecommand.TerminalSize, 1),
		}

		go func() {
			buf := make([]byte, 100)
			n, err := handler.Read(buf)
			assert.NoError(t, err)
			assert.Equal(t, "data", string(buf[:n]))
			close(done)
		}()

		select {
		case size := <-handler.resizeEvent:
			assert.Equal(t, uint16(80), size.Width)
			assert.Equal(t, uint16(24), size.Height)
		case <-time.After(2 * time.Second):
			t.Fatal("timeout waiting for resize event")
		}
	})

	server := httptest.NewServer(r)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/test"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open websocket connection: %v", err)
	}
	defer ws.Close()

	resizeMsg := TerminalMessage{Op: "resize", Cols: 80, Rows: 24}
	data, _ := json.Marshal(resizeMsg)
	ws.WriteMessage(websocket.TextMessage, data)

	time.Sleep(100 * time.Millisecond)
	stdinMsg := TerminalMessage{Op: "stdin", Data: "data"}
	data, _ = json.Marshal(stdinMsg)
	ws.WriteMessage(websocket.TextMessage, data)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for read")
	}
}

func TestWsPtyHandler_Read_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	done := make(chan struct{})

	r.GET("/test", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		handler := &WsPtyHandler{
			conn:        ws,
			resizeEvent: make(chan remotecommand.TerminalSize, 1),
		}

		buf := make([]byte, 100)
		n, err := handler.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, "valid", string(buf[:n]))

		close(done)
	})

	server := httptest.NewServer(r)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/test"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open websocket connection: %v", err)
	}
	defer ws.Close()

	ws.WriteMessage(websocket.TextMessage, []byte("not valid json"))

	msg := TerminalMessage{Op: "stdin", Data: "valid"}
	data, _ := json.Marshal(msg)
	ws.WriteMessage(websocket.TextMessage, data)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for read")
	}
}

func TestWsPtyHandler_Read_UnknownOp(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	done := make(chan struct{})

	r.GET("/test", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		handler := &WsPtyHandler{
			conn:        ws,
			resizeEvent: make(chan remotecommand.TerminalSize, 1),
		}

		buf := make([]byte, 100)
		n, err := handler.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, "valid", string(buf[:n]))

		close(done)
	})

	server := httptest.NewServer(r)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/test"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open websocket connection: %v", err)
	}
	defer ws.Close()

	unknownMsg := TerminalMessage{Op: "unknown", Data: "test"}
	data, _ := json.Marshal(unknownMsg)
	ws.WriteMessage(websocket.TextMessage, data)

	msg := TerminalMessage{Op: "stdin", Data: "valid"}
	data, _ = json.Marshal(msg)
	ws.WriteMessage(websocket.TextMessage, data)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for read")
	}
}

func TestWsPtyHandler_Read_ConnectionClosed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	done := make(chan struct{})

	r.GET("/test", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		handler := &WsPtyHandler{
			conn:        ws,
			resizeEvent: make(chan remotecommand.TerminalSize, 1),
		}

		buf := make([]byte, 100)
		_, err = handler.Read(buf)
		assert.Equal(t, io.EOF, err)

		close(done)
	})

	server := httptest.NewServer(r)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/test"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open websocket connection: %v", err)
	}

	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	ws.Close()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for read")
	}
}

func TestWsPtyHandler_Next(t *testing.T) {
	handler := &WsPtyHandler{
		resizeEvent: make(chan remotecommand.TerminalSize, 1),
	}

	expectedSize := remotecommand.TerminalSize{Width: 120, Height: 40}
	handler.resizeEvent <- expectedSize

	size := handler.Next()
	assert.NotNil(t, size)
	assert.Equal(t, expectedSize.Width, size.Width)
	assert.Equal(t, expectedSize.Height, size.Height)
}

func TestWsPtyHandler_Read_NonCloseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	done := make(chan struct{})

	r.GET("/test", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		handler := &WsPtyHandler{
			conn:        ws,
			resizeEvent: make(chan remotecommand.TerminalSize, 1),
		}

		buf := make([]byte, 100)
		_, err = handler.Read(buf)
		assert.Error(t, err)
		assert.NotEqual(t, io.EOF, err)

		close(done)
	})

	server := httptest.NewServer(r)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/test"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open websocket connection: %v", err)
	}

	ws.UnderlyingConn().Close()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for read")
	}
}

func TestWsPtyHandler_Read_ResizeChannelFull(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	done := make(chan struct{})

	r.GET("/test", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		handler := &WsPtyHandler{
			conn:        ws,
			resizeEvent: make(chan remotecommand.TerminalSize),
		}

		go func() {
			buf := make([]byte, 100)
			n, err := handler.Read(buf)
			assert.NoError(t, err)
			assert.Equal(t, "data", string(buf[:n]))
			close(done)
		}()
	})

	server := httptest.NewServer(r)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/test"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open websocket connection: %v", err)
	}
	defer ws.Close()

	resizeMsg := TerminalMessage{Op: "resize", Cols: 80, Rows: 24}
	data, _ := json.Marshal(resizeMsg)
	ws.WriteMessage(websocket.TextMessage, data)

	time.Sleep(50 * time.Millisecond)
	stdinMsg := TerminalMessage{Op: "stdin", Data: "data"}
	data, _ = json.Marshal(stdinMsg)
	ws.WriteMessage(websocket.TextMessage, data)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for read")
	}
}
