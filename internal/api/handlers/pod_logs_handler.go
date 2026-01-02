package handlers

import (
	"bufio"
	"cluster-agent/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"time"
)

type PodLogsHandler struct {
	service services.PodLogsService
}

func NewPodLogsHandler(service services.PodLogsService) *PodLogsHandler {
	return &PodLogsHandler{
		service: service,
	}
}

func (p *PodLogsHandler) StreamLogs(c *gin.Context) {
	namespace := c.Param("namespace")
	podName := c.Param("name")
	container := c.Query("container")

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	defer ws.Close()

	logsStream, err := p.service.StreamLogs(c.Request.Context(), namespace, podName, container)

	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error init logs stream: "+err.Error()))
		return
	}

	defer logsStream.Close()

	processLogsStream(ws, logsStream)
}

func processLogsStream(conn *websocket.Conn, stream io.ReadCloser) {
	buffer := bufio.NewReader(stream)

	for {
		str, err := buffer.ReadString('\n')

		if err != nil {
			break
		}

		conn.SetWriteDeadline(time.Now().Add(time.Second * 5))
		if err := conn.WriteMessage(websocket.TextMessage, []byte(str)); err != nil {
			break
		}
	}
}
