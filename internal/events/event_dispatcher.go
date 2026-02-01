package events

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type (
	Event interface {
		Name() string
	}

	ListenerFunc[T any] func(ctx context.Context, event T) error

	eventHandler func(ctx context.Context, event Event) error
)

type asyncJob struct {
	ctx   context.Context
	event Event
}

type EventDispatcher struct {
	mu sync.RWMutex

	listeners map[string][]eventHandler

	jobQueue chan asyncJob
	wg       sync.WaitGroup
}

func NewEventDispatcher(workerCount int, bufferSize int) *EventDispatcher {
	d := &EventDispatcher{
		listeners: make(map[string][]eventHandler),
		jobQueue:  make(chan asyncJob, bufferSize),
	}

	d.startWorkers(workerCount)

	return d
}

func (d *EventDispatcher) startWorkers(count int) {
	for i := 0; i < count; i++ {
		d.wg.Add(1)
		go func(workerID int) {
			defer d.wg.Done()

			for job := range d.jobQueue {
				if err := d.Dispatch(job.ctx, job.event); err != nil {
					log.Printf("[AsyncWorker-%d] Failed to process %s: %v", workerID, job.event.Name(), err)
				}
			}
		}(i)
	}
}

func (d *EventDispatcher) Shutdown() {
	close(d.jobQueue)
	d.wg.Wait()
}

func RegisterListener[T Event](d *EventDispatcher, listenerFunc ListenerFunc[T]) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var t T
	eventName := t.Name()

	wrapper := func(ctx context.Context, e Event) error {
		typedEvent, ok := e.(T)
		if !ok {
			return fmt.Errorf("failed to cast event to %T", t)
		}
		return listenerFunc(ctx, typedEvent)
	}

	d.listeners[eventName] = append(d.listeners[eventName], wrapper)
}

func (d *EventDispatcher) Dispatch(ctx context.Context, e Event) error {
	d.mu.RLock()
	handlers, ok := d.listeners[e.Name()]
	d.mu.RUnlock()

	if !ok {
		return nil
	}

	for _, handler := range handlers {
		if err := handler(ctx, e); err != nil {
			log.Println("Failed to dispatch event:", err)
		}
	}

	return nil
}

func (d *EventDispatcher) DispatchAsync(ctx context.Context, e Event) {
	select {
	case d.jobQueue <- asyncJob{ctx: ctx, event: e}:
	default:
		log.Printf("EventDispatcher queue is full! Dropping event: %s", e.Name())
	}
}
