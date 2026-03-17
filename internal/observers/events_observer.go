package observers

import (
	"cluster-agent/internal/consumers"
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

type EventsObserver struct {
	informer cache.SharedIndexInformer
	batcher  *consumers.EventBatcher
}

func NewEventCollector(
	batcher *consumers.EventBatcher,
	informer cache.SharedIndexInformer,
) *EventsObserver {
	collector := &EventsObserver{
		informer: informer,
		batcher:  batcher,
	}

	_, err := informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    collector.handleEvent,
		UpdateFunc: collector.handleEventUpdate,
	})

	if err != nil {
		log.Fatal(err)
	}

	return collector
}

func (e *EventsObserver) handleEvent(obj interface{}) {
	event, ok := obj.(*corev1.Event)
	if !ok {
		return
	}

	log.Printf("New Event: %s/%s - %s", event.InvolvedObject.Kind, event.InvolvedObject.Name, event.Reason)

	e.batcher.Push(event)
}

func (e *EventsObserver) handleEventUpdate(oldObj, newObj interface{}) {
	newEvent, ok := newObj.(*corev1.Event)
	if !ok {
		return
	}
	e.batcher.Push(newEvent)
}
