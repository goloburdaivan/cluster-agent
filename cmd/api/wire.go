//go:build wireinject
// +build wireinject

package main

import (
	"cluster-agent/internal"
	"cluster-agent/internal/api/handlers"
	"cluster-agent/internal/api/listeners"
	"cluster-agent/internal/api/middleware"
	cache2 "cluster-agent/internal/cache"
	"cluster-agent/internal/config"
	"cluster-agent/internal/consumers"
	"cluster-agent/internal/events"
	"cluster-agent/internal/k8s"
	"cluster-agent/internal/observers"
	"cluster-agent/internal/services"
	"cluster-agent/internal/services/topology"
	"time"

	"github.com/google/wire"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
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

func ProvidePodLister(factory informers.SharedInformerFactory) corelisters.PodLister {
	return factory.Core().V1().Pods().Lister()
}

func ProvideEventDispatcher() *events.EventDispatcher {
	return events.NewEventDispatcher(5, 100)
}

func ProvideDynamicClient(cfg *rest.Config) (dynamic.Interface, error) {
	return dynamic.NewForConfig(cfg)
}

func ProvideMetricsClient(cfg *rest.Config) (metricsclientset.Interface, error) {
	return metricsclientset.NewForConfig(cfg)
}

func InitializeApp() (*internal.App, func(), error) {
	wire.Build(
		config.NewConfig,
		k8s.NewClient,

		ProvideInformerFactory,
		ProvideEventInformer,
		ProvidePodLister,
		ProvideK8sInterface,
		ProvideRestConfig,
		ProvideMetricsClient,
		ProvideDynamicClient,

		cache2.NewRedisClient,
		cache2.NewTopologyCache,
		wire.Bind(new(topology.TopologyCacheStorage), new(*cache2.TopologyCache)),
		wire.Bind(new(listeners.TopologyCacheInvalidator), new(*cache2.TopologyCache)),

		listeners.NewTopologyCacheListener,

		handlers.HandlerSet,
		middleware.NewAuthorizedMiddleware,

		// Services
		services.NewDeploymentService,
		services.NewNamespaceService,
		services.NewNodeService,
		services.NewPodService,
		services.NewServiceService,
		services.NewTerminalService,
		services.NewSnapshotService,
		services.NewPodLogsService,
		services.NewIngressService,
		services.NewPVCService,
		services.NewConfigMapService,
		services.NewSecretService,
		services.NewNetworkInspectorService,
		services.NewDaemonSetService,
		services.NewJobService,
		services.NewCronJobService,
		services.NewPVService,
		services.NewStorageClassService,
		services.NewNetworkPolicyService,
		services.NewServiceAccountService,
		services.NewRoleService,
		services.NewRoleBindingService,
		services.NewMetricsService,
		topology.NewTopologyService,

		consumers.NewEventBatcher,
		observers.NewEventCollector,
		observers.NewTopologyInvalidator,
		ProvideEventDispatcher,

		internal.NewApp,
	)

	return &internal.App{}, nil, nil
}
