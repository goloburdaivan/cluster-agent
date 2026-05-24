package observers

import (
	"cluster-agent/internal/events"
	"context"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

type TopologyObserver struct {
	dispatcher *events.EventDispatcher
	factory    informers.SharedInformerFactory
}

func NewTopologyInvalidator(
	dispatcher *events.EventDispatcher,
	factory informers.SharedInformerFactory,
) *TopologyObserver {
	ti := &TopologyObserver{
		dispatcher: dispatcher,
		factory:    factory,
	}

	ti.registerHandlers()

	return ti
}

func (ti *TopologyObserver) registerHandlers() {
	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    ti.handleObject,
		UpdateFunc: ti.handleUpdate,
		DeleteFunc: ti.handleObject,
	}

	ti.factory.Apps().V1().Deployments().Informer().AddEventHandler(handlers)
	ti.factory.Core().V1().Services().Informer().AddEventHandler(handlers)
	ti.factory.Apps().V1().StatefulSets().Informer().AddEventHandler(handlers)
	ti.factory.Networking().V1().Ingresses().Informer().AddEventHandler(handlers)
	ti.factory.Core().V1().ConfigMaps().Informer().AddEventHandler(handlers)
	ti.factory.Core().V1().Secrets().Informer().AddEventHandler(handlers)
	ti.factory.Core().V1().PersistentVolumeClaims().Informer().AddEventHandler(handlers)
	ti.factory.Core().V1().PersistentVolumes().Informer().AddEventHandler(handlers)
	ti.factory.Networking().V1().NetworkPolicies().Informer().AddEventHandler(handlers)
}

func (ti *TopologyObserver) handleObject(obj interface{}) {
	var namespace string

	switch o := obj.(type) {
	case *appsv1.Deployment:
		namespace = o.Namespace
	case *corev1.Service:
		namespace = o.Namespace
	case *appsv1.StatefulSet:
		namespace = o.Namespace
	case *networkingv1.Ingress:
		namespace = o.Namespace
	case *corev1.ConfigMap:
		namespace = o.Namespace
	case *corev1.Secret:
		namespace = o.Namespace
	case *corev1.PersistentVolumeClaim:
		namespace = o.Namespace
	case *corev1.PersistentVolume:
		namespace = o.Namespace
	case *networkingv1.NetworkPolicy:
		namespace = o.Namespace
	case cache.DeletedFinalStateUnknown:
		ti.handleObject(o.Obj)
		return
	default:
		return
	}

	if namespace != "" {
		err := ti.dispatcher.Dispatch(context.Background(), &events.ClusterChangedEvent{Namespace: namespace})
		if err != nil {
			log.Printf("[TopologyObserver] error during dispatch: %v", err)
		}
	}
}

func (ti *TopologyObserver) handleUpdate(oldObj, newObj interface{}) {
	ti.handleObject(newObj)
}
