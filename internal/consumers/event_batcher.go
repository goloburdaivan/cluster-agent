package consumers

import (
	"bytes"
	"cluster-agent/internal/config"
	"context"
	"encoding/json"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"log"
	"math"
	"net"
	"net/http"
	"time"
)

type EventBatcher struct {
	eventsChan chan *corev1.Event
	buffer     []*corev1.Event
	batchSize  int
	interval   time.Duration
	httpClient *http.Client
	cfg        *config.Config
}

func NewEventBatcher(cfg *config.Config) *EventBatcher {
	t := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
	}

	client := &http.Client{
		Transport: t,
		Timeout:   15 * time.Second,
	}

	return &EventBatcher{
		eventsChan: make(chan *corev1.Event, 1000),
		buffer:     make([]*corev1.Event, 0, 100),
		batchSize:  100,
		interval:   5 * time.Second,
		cfg:        cfg,
		httpClient: client,
	}
}

func (b *EventBatcher) Push(event *corev1.Event) {
	// This data is TOO big
	event.ManagedFields = nil

	select {
	case b.eventsChan <- event:
	default:
		log.Println("Event channel full, dropping event")
	}
}

func (b *EventBatcher) Run(ctx context.Context) {
	ticker := time.NewTicker(b.interval)
	defer ticker.Stop()

	for {
		select {
		case event := <-b.eventsChan:
			b.buffer = append(b.buffer, event)
			if len(b.buffer) >= b.batchSize {
				b.flush(ctx)
			}

		case <-ticker.C:
			if len(b.buffer) > 0 {
				b.flush(ctx)
			}

		case <-ctx.Done():
			if len(b.buffer) > 0 {
				b.flush(ctx)
			}
			return
		}
	}
}

func (b *EventBatcher) flush(ctx context.Context) {
	count := len(b.buffer)
	log.Printf("Flushing %d events to Laravel...", count)

	payload, err := json.Marshal(b.buffer)
	if err != nil {
		log.Printf("Failed to marshal events: %v", err)
		b.buffer = b.buffer[:0]
		return
	}

	const maxRetries = 3
	var attempt int

	for attempt = 0; attempt < maxRetries; attempt++ {
		if ctx.Err() != nil {
			log.Println("Flush canceled due to context")
			return
		}

		err = b.sendRequest(ctx, payload)
		if err == nil {
			break
		}

		log.Printf("⚠Attempt %d/%d failed: %v", attempt+1, maxRetries, err)

		backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second

		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			return
		}
	}

	if err != nil {
		log.Printf("Failed to send events after %d attempts. Dropping batch.", maxRetries)
	} else {
		log.Printf("Successfully sent %d events", count)
	}

	b.buffer = b.buffer[:0]
}

func (b *EventBatcher) sendRequest(ctx context.Context, payload []byte) error {
	req, err := http.NewRequestWithContext(ctx, "POST", b.cfg.ApiURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return fmt.Errorf("server error: %d", resp.StatusCode)
	}

	if resp.StatusCode >= 400 {
		log.Printf("Laravel rejected events: %d. Check payload or auth.", resp.StatusCode)
		return nil
	}

	return nil
}
