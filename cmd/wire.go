//go:build wireinject
// +build wireinject

package main

import (
	"cluster-agent/internal"
	"cluster-agent/internal/api/handlers"
	"cluster-agent/internal/config"
	"cluster-agent/internal/consumers"
	"cluster-agent/internal/k8s"
	"cluster-agent/internal/producers"
	"cluster-agent/internal/services"
	"cluster-agent/internal/services/topology"
	"github.com/google/wire"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"time"
)

func ProvideK8sInterface(client *k8s.Client) kubernetes.Interface {
	return client.GetClientset()
}

func ProvideRestConfig(client *k8s.Client) *rest.Config {
	return client.GetConfig()
}

func ProvideInformerFactory(clientset kubernetes.Interface) informers.SharedInformerFactory {
	return informers.NewSharedInformerFactory(clientset, 12*time.Hour)
}

func ProvideEventInformer(factory informers.SharedInformerFactory) cache.SharedIndexInformer {
	return factory.Core().V1().Events().Informer()
}

func InitializeApp() (*internal.App, error) {
	wire.Build(
		k8s.NewClient,

		ProvideInformerFactory,
		ProvideEventInformer,
		ProvideK8sInterface,
		ProvideRestConfig,

		handlers.HandlerSet,

		// Services
		services.NewDeploymentService,
		services.NewNamespaceService,
		services.NewNodeService,
		services.NewPodService,
		services.NewServiceService,
		services.NewTerminalService,
		services.NewSnapshotService,
		topology.NewTopologyService,

		config.NewConfig,
		consumers.NewEventBatcher,
		producers.NewEventCollector,

		internal.NewApp,
	)
	return &internal.App{}, nil
}
