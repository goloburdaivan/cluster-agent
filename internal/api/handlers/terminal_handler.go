package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"cluster-agent/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"k8s.io/client-go/tools/remotecommand"
)

type TerminalHandler struct {
	service services.TerminalService
}

func NewTerminalHandler(service services.TerminalService) *TerminalHandler {
	return &TerminalHandler{
		service: service,
	}
}

type TerminalMessage struct {
	Op   string `json:"op"`
	Data string `json:"data,omitempty"`
	Rows uint16 `json:"rows,omitempty"`
	Cols uint16 `json:"cols,omitempty"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *TerminalHandler) Exec(c *gin.Context) {
	namespace := c.Param("namespace")
	podName := c.Param("name")
	container := c.Query("container")

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	defer ws.Close()

	exec, err := h.service.GetAuthExecutor(namespace, podName, container)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error init executor: "+err.Error()))
		return
	}

	handler := &WsPtyHandler{
		conn:        ws,
		resizeEvent: make(chan remotecommand.TerminalSize),
	}

	err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:             handler,
		Stdout:            handler,
		Stderr:            handler,
		Tty:               true,
		TerminalSizeQueue: handler,
	})

	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("\r\nSession ended with error: "+err.Error()))
	}
}

type WsPtyHandler struct {
	conn        *websocket.Conn
	resizeEvent chan remotecommand.TerminalSize
	writeMutex  sync.Mutex
	readBuf     bytes.Buffer
}

func (t *WsPtyHandler) Write(p []byte) (int, error) {
	t.writeMutex.Lock()
	defer t.writeMutex.Unlock()

	err := t.conn.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (t *WsPtyHandler) Read(p []byte) (int, error) {
	if t.readBuf.Len() > 0 {
		return t.readBuf.Read(p)
	}

	for {
		_, message, err := t.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				return 0, io.EOF
			}

			return 0, err
		}

		var msg TerminalMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		switch msg.Op {
		case "resize":
			select {
			case t.resizeEvent <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}:
			default:
			}
			continue

		case "stdin":
			t.readBuf.WriteString(msg.Data)
			return t.readBuf.Read(p)

		default:
			continue
		}
	}
}

func (t *WsPtyHandler) Next() *remotecommand.TerminalSize {
	size := <-t.resizeEvent
	return &size
}
