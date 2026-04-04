package handlers

import (
	"cluster-agent/internal/services"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	writeWait       = 10 * time.Second
	pongWait        = 60 * time.Second
	pingPeriod      = (pongWait * 9) / 10
	metricsInterval = 1 * time.Second
)

type MetricsHandler struct {
	service services.MetricsService
}

func NewMetricsHandler(service services.MetricsService) *MetricsHandler {
	return &MetricsHandler{
		service: service,
	}
}

func (h *MetricsHandler) GetNodeMetrics(c *gin.Context) {
	nodeID := c.Param("id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("Failed to upgrade connection", "error", err)
		return
	}
	defer conn.Close()

	ctx := c.Request.Context()

	done := make(chan struct{})
	go h.readPump(conn, done, nodeID)

	metricsTicker := time.NewTicker(metricsInterval)
	defer metricsTicker.Stop()

	pingTicker := time.NewTicker(pingPeriod)
	defer pingTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Context cancelled, closing WebSocket", "node", nodeID)
			return

		case <-done:
			slog.Info("WebSocket closed by client", "node", nodeID)
			return

		case <-metricsTicker.C:
			metrics, err := h.service.GetNodeMetrics(ctx, nodeID)
			if err != nil {
				slog.Warn("Failed to get node metrics", "node", nodeID, "error", err)
				h.writeError(conn, "Failed to get node metrics: "+err.Error())
				return
			}

			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteJSON(metrics); err != nil {
				slog.Warn("Failed to write metrics", "node", nodeID, "error", err)
				return
			}

		case <-pingTicker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Warn("Failed to write ping", "node", nodeID, "error", err)
				return
			}
		}
	}
}

func (h *MetricsHandler) GetPodMetrics(c *gin.Context) {
	namespace := c.Param("namespace")
	podID := c.Param("id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("Failed to upgrade connection", "error", err)
		return
	}
	defer conn.Close()

	ctx := c.Request.Context()

	done := make(chan struct{})
	go h.readPump(conn, done, namespace+"/"+podID)

	metricsTicker := time.NewTicker(metricsInterval)
	defer metricsTicker.Stop()

	pingTicker := time.NewTicker(pingPeriod)
	defer pingTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Context cancelled, closing WebSocket", "pod", namespace+"/"+podID)
			return

		case <-done:
			slog.Info("WebSocket closed by client", "pod", namespace+"/"+podID)
			return

		case <-metricsTicker.C:
			metrics, err := h.service.GetPodMetrics(ctx, namespace, podID)
			if err != nil {
				slog.Warn("Failed to get pod metrics", "pod", namespace+"/"+podID, "error", err)
				h.writeError(conn, "Failed to get pod metrics: "+err.Error())
				return
			}

			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteJSON(metrics); err != nil {
				slog.Warn("Failed to write metrics", "pod", namespace+"/"+podID, "error", err)
				return
			}

		case <-pingTicker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Warn("Failed to write ping", "pod", namespace+"/"+podID, "error", err)
				return
			}
		}
	}
}

func (h *MetricsHandler) readPump(conn *websocket.Conn, done chan struct{}, resource string) {
	defer close(done)

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Warn("WebSocket closed unexpectedly", "resource", resource, "error", err)
			}
			return
		}
	}
}

func (h *MetricsHandler) writeError(conn *websocket.Conn, message string) {
	conn.SetWriteDeadline(time.Now().Add(writeWait))
	conn.WriteJSON(map[string]string{
		"error": message,
		"type":  "error",
	})
}
