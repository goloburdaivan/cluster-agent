package producers

import (
	"cluster-agent/internal/consumers"
	"context"
	"log"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type EventCollector struct {
	clientset kubernetes.Interface
	factory   informers.SharedInformerFactory
	informer  cache.SharedIndexInformer
	batcher   *consumers.EventBatcher
}

func NewEventCollector(
	clientset kubernetes.Interface,
	batcher *consumers.EventBatcher,
) *EventCollector {
	factory := informers.NewSharedInformerFactoryWithOptions(
		clientset,
		time.Hour*12,
		informers.WithNamespace(metav1.NamespaceAll),
	)

	eventInformer := factory.Core().V1().Events().Informer()

	collector := &EventCollector{
		clientset: clientset,
		factory:   factory,
		informer:  eventInformer,
		batcher:   batcher,
	}

	_, err := eventInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    collector.handleEvent,
		UpdateFunc: collector.handleEventUpdate,
	})

	if err != nil {
		log.Fatal(err)
	}

	return collector
}

func (e *EventCollector) Start(ctx context.Context) {
	log.Println("Starting Event Collector...")

	e.factory.Start(ctx.Done())

	log.Println("Waiting for cache sync...")
	if !cache.WaitForCacheSync(ctx.Done(), e.informer.HasSynced) {
		log.Println("Timed out waiting for caches to sync")
		return
	}

	log.Println("Event Collector synced! Listening for events...")

	<-ctx.Done()
	log.Println("Stopping Event Collector")
}

func (e *EventCollector) handleEvent(obj interface{}) {
	event, ok := obj.(*corev1.Event)
	if !ok {
		return
	}

	log.Printf("New Event: %s/%s - %s", event.InvolvedObject.Kind, event.InvolvedObject.Name, event.Reason)

	e.batcher.Push(event)
}

func (e *EventCollector) handleEventUpdate(oldObj, newObj interface{}) {
	newEvent, ok := newObj.(*corev1.Event)
	if !ok {
		return
	}
	e.batcher.Push(newEvent)
}
