package internal

import (
	"cluster-agent/internal/api/middleware"
	"cluster-agent/internal/auth/permissions"
	"cluster-agent/internal/consumers"
	"cluster-agent/internal/events"
	"cluster-agent/internal/observers"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"

	"cluster-agent/internal/api/handlers"
	"cluster-agent/internal/api/listeners"
	"cluster-agent/internal/services/topology"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type App struct {
	Router                *gin.Engine
	Handlers              *handlers.HandlerContainer
	EventCollector        *observers.EventsObserver
	TopologyInvalidator   *observers.TopologyObserver
	TopologyCacheListener *listeners.TopologyCacheListener
	EventBatcher          *consumers.EventBatcher
	EventDispatcher       *events.EventDispatcher
	TopologyService       topology.Service
	InformerFactory       informers.SharedInformerFactory
	authorizedMiddleware  *middleware.AuthorizedMiddleware
}

func NewApp(
	h *handlers.HandlerContainer,
	authorizedMiddleware *middleware.AuthorizedMiddleware,
	collector *observers.EventsObserver,
	invalidator *observers.TopologyObserver,
	topologyCacheListener *listeners.TopologyCacheListener,
	batcher *consumers.EventBatcher,
	dispatcher *events.EventDispatcher,
	topologyService topology.Service,
	factory informers.SharedInformerFactory,
) *App {
	app := &App{
		Router:                gin.Default(),
		Handlers:              h,
		authorizedMiddleware:  authorizedMiddleware,
		EventCollector:        collector,
		TopologyInvalidator:   invalidator,
		TopologyCacheListener: topologyCacheListener,
		EventBatcher:          batcher,
		EventDispatcher:       dispatcher,
		TopologyService:       topologyService,
		InformerFactory:       factory,
	}

	app.setRoutes()
	app.setupEventListeners()

	return app
}

func (app *App) setupEventListeners() {
	events.RegisterListener(app.EventDispatcher, app.TopologyCacheListener.Handle)
}

func (app *App) setRoutes() {
	v1 := app.Router.Group("/api/v1")
	v1.Use(app.authorizedMiddleware.Handle())
	{
		pods := v1.Group("/pods")
		{
			pods.GET("",
				app.authorizedMiddleware.HasPermission(permissions.PodsView),
				app.Handlers.Pod.List,
			)
			pods.GET("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.PodsView),
				app.Handlers.Pod.Get,
			)
			pods.GET("/:namespace/:name/raw",
				app.authorizedMiddleware.HasPermission(permissions.PodsView),
				app.Handlers.RawPatch.GetNamespaced(schema.GroupVersionResource{Version: "v1", Resource: "pods"}),
			)
			pods.POST("",
				app.authorizedMiddleware.HasPermission(permissions.PodsCreate),
				app.Handlers.Pod.Create,
			)
			pods.PUT("",
				app.authorizedMiddleware.HasPermission(permissions.PodsCreate),
				app.Handlers.Pod.Update,
			)
			pods.PATCH("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.PodsCreate),
				app.Handlers.RawPatch.PatchNamespaced(schema.GroupVersionResource{Version: "v1", Resource: "pods"}),
			)
			pods.DELETE("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.PodsDelete),
				app.Handlers.Pod.Delete,
			)

			pods.GET("/:namespace/:name/logs",
				app.authorizedMiddleware.HasPermission(permissions.PodsView),
				app.Handlers.PodLogs.StreamLogs,
			)
			pods.GET("/:namespace/:name/exec",
				app.authorizedMiddleware.HasPermission(permissions.PodsView),
				app.Handlers.Terminal.Exec,
			)
			pods.GET("/:namespace/:name/network",
				app.authorizedMiddleware.HasPermission(permissions.PodsView),
				app.Handlers.NetworkInspector.GetConnections,
			)
		}

		deployments := v1.Group("/deployments")
		{
			deployments.GET("",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsView),
				app.Handlers.Deployment.List,
			)

			deployments.GET("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsView),
				app.Handlers.Deployment.Get,
			)
			deployments.GET("/:namespace/:name/raw",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsView),
				app.Handlers.RawPatch.GetNamespaced(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}),
			)

			deployments.POST("",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsCreate),
				app.Handlers.Deployment.Create,
			)

			deployments.PUT("",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsCreate),
				app.Handlers.Deployment.Update,
			)
			deployments.PATCH("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsCreate),
				app.Handlers.RawPatch.PatchNamespaced(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}),
			)

			deployments.DELETE("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsDelete),
				app.Handlers.Deployment.Delete,
			)

			deployments.PATCH("/scale",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsScale),
				app.Handlers.Deployment.ScaleDeployment,
			)
		}

		daemonsets := v1.Group("/daemonsets")
		{
			daemonsets.GET("",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsView),
				app.Handlers.DaemonSet.List,
			)

			daemonsets.GET("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsView),
				app.Handlers.DaemonSet.Get,
			)
			daemonsets.GET("/:namespace/:name/raw",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsView),
				app.Handlers.RawPatch.GetNamespaced(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "daemonsets"}),
			)

			daemonsets.POST("",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsCreate),
				app.Handlers.DaemonSet.Create,
			)

			daemonsets.PUT("",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsCreate),
				app.Handlers.DaemonSet.Update,
			)
			daemonsets.PATCH("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsCreate),
				app.Handlers.RawPatch.PatchNamespaced(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "daemonsets"}),
			)

			daemonsets.DELETE("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsDelete),
				app.Handlers.DaemonSet.Delete,
			)
		}

		jobs := v1.Group("/jobs")
		{
			jobs.GET("",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsView),
				app.Handlers.Job.List,
			)

			jobs.GET("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsView),
				app.Handlers.Job.Get,
			)
			jobs.GET("/:namespace/:name/raw",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsView),
				app.Handlers.RawPatch.GetNamespaced(schema.GroupVersionResource{Group: "batch", Version: "v1", Resource: "jobs"}),
			)

			jobs.POST("",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsCreate),
				app.Handlers.Job.Create,
			)

			jobs.PUT("",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsCreate),
				app.Handlers.Job.Update,
			)
			jobs.PATCH("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsCreate),
				app.Handlers.RawPatch.PatchNamespaced(schema.GroupVersionResource{Group: "batch", Version: "v1", Resource: "jobs"}),
			)

			jobs.DELETE("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsDelete),
				app.Handlers.Job.Delete,
			)
		}

		cronjobs := v1.Group("/cronjobs")
		{
			cronjobs.GET("",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsView),
				app.Handlers.CronJob.List,
			)

			cronjobs.GET("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsView),
				app.Handlers.CronJob.Get,
			)
			cronjobs.GET("/:namespace/:name/raw",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsView),
				app.Handlers.RawPatch.GetNamespaced(schema.GroupVersionResource{Group: "batch", Version: "v1", Resource: "cronjobs"}),
			)

			cronjobs.POST("",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsCreate),
				app.Handlers.CronJob.Create,
			)

			cronjobs.PUT("",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsCreate),
				app.Handlers.CronJob.Update,
			)
			cronjobs.PATCH("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsCreate),
				app.Handlers.RawPatch.PatchNamespaced(schema.GroupVersionResource{Group: "batch", Version: "v1", Resource: "cronjobs"}),
			)

			cronjobs.DELETE("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsDelete),
				app.Handlers.CronJob.Delete,
			)
		}

		services := v1.Group("/services")
		{
			services.GET("",
				app.authorizedMiddleware.HasPermission(permissions.ServicesView),
				app.Handlers.Service.List,
			)
			services.GET("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.ServicesView),
				app.Handlers.Service.Get,
			)
			services.GET("/:namespace/:name/raw",
				app.authorizedMiddleware.HasPermission(permissions.ServicesView),
				app.Handlers.RawPatch.GetNamespaced(schema.GroupVersionResource{Version: "v1", Resource: "services"}),
			)
			services.POST("",
				app.authorizedMiddleware.HasPermission(permissions.ServicesCreate),
				app.Handlers.Service.Create,
			)
			services.PUT("",
				app.authorizedMiddleware.HasPermission(permissions.ServicesCreate),
				app.Handlers.Service.Update,
			)
			services.PATCH("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.ServicesCreate),
				app.Handlers.RawPatch.PatchNamespaced(schema.GroupVersionResource{Version: "v1", Resource: "services"}),
			)
			services.DELETE("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.ServicesDelete),
				app.Handlers.Service.Delete,
			)
		}

		configmaps := v1.Group("/configmaps")
		{
			configmaps.GET("",
				app.authorizedMiddleware.HasPermission(permissions.ConfigMapsView),
				app.Handlers.ConfigMaps.List,
			)
			configmaps.GET("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.ConfigMapsView),
				app.Handlers.ConfigMaps.Get,
			)
			configmaps.GET("/:namespace/:name/raw",
				app.authorizedMiddleware.HasPermission(permissions.ConfigMapsView),
				app.Handlers.RawPatch.GetNamespaced(schema.GroupVersionResource{Version: "v1", Resource: "configmaps"}),
			)
			configmaps.POST("",
				app.authorizedMiddleware.HasPermission(permissions.ConfigMapsCreate),
				app.Handlers.ConfigMaps.Create,
			)
			configmaps.PUT("",
				app.authorizedMiddleware.HasPermission(permissions.ConfigMapsCreate),
				app.Handlers.ConfigMaps.Update,
			)
			configmaps.PATCH("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.ConfigMapsCreate),
				app.Handlers.RawPatch.PatchNamespaced(schema.GroupVersionResource{Version: "v1", Resource: "configmaps"}),
			)
			configmaps.DELETE("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.ConfigMapsDelete),
				app.Handlers.ConfigMaps.Delete,
			)
		}

		secrets := v1.Group("/secrets")
		{
			secrets.GET("",
				app.authorizedMiddleware.HasPermission(permissions.SecretsView),
				app.Handlers.Secrets.List,
			)
			secrets.GET("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.SecretsView),
				app.Handlers.Secrets.Get,
			)
			secrets.GET("/:namespace/:name/raw",
				app.authorizedMiddleware.HasPermission(permissions.SecretsView),
				app.Handlers.RawPatch.GetNamespaced(schema.GroupVersionResource{Version: "v1", Resource: "secrets"}),
			)
			secrets.POST("",
				app.authorizedMiddleware.HasPermission(permissions.SecretsCreate),
				app.Handlers.Secrets.Create,
			)
			secrets.PUT("",
				app.authorizedMiddleware.HasPermission(permissions.SecretsCreate),
				app.Handlers.Secrets.Update,
			)
			secrets.PATCH("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.SecretsCreate),
				app.Handlers.RawPatch.PatchNamespaced(schema.GroupVersionResource{Version: "v1", Resource: "secrets"}),
			)
			secrets.DELETE("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.SecretsDelete),
				app.Handlers.Secrets.Delete,
			)
		}

		ingresses := v1.Group("/ingresses")
		{
			ingresses.GET("",
				app.authorizedMiddleware.HasPermission(permissions.IngressesView),
				app.Handlers.Ingresses.List,
			)
			ingresses.GET("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.IngressesView),
				app.Handlers.Ingresses.Get,
			)
			ingresses.GET("/:namespace/:name/raw",
				app.authorizedMiddleware.HasPermission(permissions.IngressesView),
				app.Handlers.RawPatch.GetNamespaced(schema.GroupVersionResource{Group: "networking.k8s.io", Version: "v1", Resource: "ingresses"}),
			)
			ingresses.POST("",
				app.authorizedMiddleware.HasPermission(permissions.IngressesCreate),
				app.Handlers.Ingresses.Create,
			)
			ingresses.PUT("",
				app.authorizedMiddleware.HasPermission(permissions.IngressesCreate),
				app.Handlers.Ingresses.Update,
			)
			ingresses.PATCH("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.IngressesCreate),
				app.Handlers.RawPatch.PatchNamespaced(schema.GroupVersionResource{Group: "networking.k8s.io", Version: "v1", Resource: "ingresses"}),
			)
			ingresses.DELETE("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.IngressesDelete),
				app.Handlers.Ingresses.Delete,
			)
		}

		pvcs := v1.Group("/persistentvolumeclaims")
		{
			pvcs.GET("",
				app.authorizedMiddleware.HasPermission(permissions.PVCsView),
				app.Handlers.Pvcs.List,
			)
			pvcs.GET("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.PVCsView),
				app.Handlers.Pvcs.Get,
			)
			pvcs.GET("/:namespace/:name/raw",
				app.authorizedMiddleware.HasPermission(permissions.PVCsView),
				app.Handlers.RawPatch.GetNamespaced(schema.GroupVersionResource{Version: "v1", Resource: "persistentvolumeclaims"}),
			)
			pvcs.POST("",
				app.authorizedMiddleware.HasPermission(permissions.PVCsCreate),
				app.Handlers.Pvcs.Create,
			)
			pvcs.PUT("",
				app.authorizedMiddleware.HasPermission(permissions.PVCsCreate),
				app.Handlers.Pvcs.Update,
			)
			pvcs.PATCH("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.PVCsCreate),
				app.Handlers.RawPatch.PatchNamespaced(schema.GroupVersionResource{Version: "v1", Resource: "persistentvolumeclaims"}),
			)
			pvcs.DELETE("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.PVCsDelete),
				app.Handlers.Pvcs.Delete,
			)
		}

		namespace := v1.Group("/namespaces")
		{
			namespace.GET("", app.Handlers.Namespace.List)
			namespace.GET("/:name", app.Handlers.Namespace.Get)
			namespace.POST("",
				app.authorizedMiddleware.HasPermission(permissions.NamespacesCreate),
				app.Handlers.Namespace.Create,
			)
			namespace.DELETE("/:name",
				app.authorizedMiddleware.HasPermission(permissions.NamespacesDelete),
				app.Handlers.Namespace.Delete,
			)
		}

		node := v1.Group("/nodes")
		{
			node.GET("",
				app.authorizedMiddleware.HasPermission(permissions.NodesView),
				app.Handlers.Node.List,
			)
			node.GET("/:name",
				app.authorizedMiddleware.HasPermission(permissions.NodesView),
				app.Handlers.Node.Get,
			)
		}

		pvs := v1.Group("/persistentvolumes")
		{
			pvs.GET("",
				app.authorizedMiddleware.HasPermission(permissions.PVCsView),
				app.Handlers.PV.List,
			)
			pvs.GET("/:name",
				app.authorizedMiddleware.HasPermission(permissions.PVCsView),
				app.Handlers.PV.Get,
			)
			pvs.GET("/:name/raw",
				app.authorizedMiddleware.HasPermission(permissions.PVCsView),
				app.Handlers.RawPatch.GetClusterScoped(schema.GroupVersionResource{Version: "v1", Resource: "persistentvolumes"}),
			)
			pvs.POST("",
				app.authorizedMiddleware.HasPermission(permissions.PVCsCreate),
				app.Handlers.PV.Create,
			)
			pvs.PUT("",
				app.authorizedMiddleware.HasPermission(permissions.PVCsCreate),
				app.Handlers.PV.Update,
			)
			pvs.PATCH("/:name",
				app.authorizedMiddleware.HasPermission(permissions.PVCsCreate),
				app.Handlers.RawPatch.PatchClusterScoped(schema.GroupVersionResource{Version: "v1", Resource: "persistentvolumes"}),
			)
			pvs.DELETE("/:name",
				app.authorizedMiddleware.HasPermission(permissions.PVCsDelete),
				app.Handlers.PV.Delete,
			)
		}

		storageclasses := v1.Group("/storageclasses")
		{
			storageclasses.GET("",
				app.authorizedMiddleware.HasPermission(permissions.PVCsView),
				app.Handlers.StorageClass.List,
			)
			storageclasses.GET("/:name",
				app.authorizedMiddleware.HasPermission(permissions.PVCsView),
				app.Handlers.StorageClass.Get,
			)
			storageclasses.GET("/:name/raw",
				app.authorizedMiddleware.HasPermission(permissions.PVCsView),
				app.Handlers.RawPatch.GetClusterScoped(schema.GroupVersionResource{Group: "storage.k8s.io", Version: "v1", Resource: "storageclasses"}),
			)
			storageclasses.POST("",
				app.authorizedMiddleware.HasPermission(permissions.PVCsCreate),
				app.Handlers.StorageClass.Create,
			)
			storageclasses.PUT("",
				app.authorizedMiddleware.HasPermission(permissions.PVCsCreate),
				app.Handlers.StorageClass.Update,
			)
			storageclasses.PATCH("/:name",
				app.authorizedMiddleware.HasPermission(permissions.PVCsCreate),
				app.Handlers.RawPatch.PatchClusterScoped(schema.GroupVersionResource{Group: "storage.k8s.io", Version: "v1", Resource: "storageclasses"}),
			)
			storageclasses.DELETE("/:name",
				app.authorizedMiddleware.HasPermission(permissions.PVCsDelete),
				app.Handlers.StorageClass.Delete,
			)
		}

		networkpolicies := v1.Group("/networkpolicies")
		{
			networkpolicies.GET("",
				app.authorizedMiddleware.HasPermission(permissions.ServicesView),
				app.Handlers.NetworkPolicy.List,
			)
			networkpolicies.GET("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.ServicesView),
				app.Handlers.NetworkPolicy.Get,
			)
			networkpolicies.GET("/:namespace/:name/raw",
				app.authorizedMiddleware.HasPermission(permissions.ServicesView),
				app.Handlers.RawPatch.GetNamespaced(schema.GroupVersionResource{Group: "networking.k8s.io", Version: "v1", Resource: "networkpolicies"}),
			)
			networkpolicies.POST("",
				app.authorizedMiddleware.HasPermission(permissions.ServicesCreate),
				app.Handlers.NetworkPolicy.Create,
			)
			networkpolicies.PUT("",
				app.authorizedMiddleware.HasPermission(permissions.ServicesCreate),
				app.Handlers.NetworkPolicy.Update,
			)
			networkpolicies.PATCH("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.ServicesCreate),
				app.Handlers.RawPatch.PatchNamespaced(schema.GroupVersionResource{Group: "networking.k8s.io", Version: "v1", Resource: "networkpolicies"}),
			)
			networkpolicies.DELETE("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.ServicesDelete),
				app.Handlers.NetworkPolicy.Delete,
			)
		}

		serviceaccounts := v1.Group("/serviceaccounts")
		{
			serviceaccounts.GET("",
				app.authorizedMiddleware.HasPermission(permissions.SecretsView),
				app.Handlers.ServiceAccount.List,
			)
			serviceaccounts.GET("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.SecretsView),
				app.Handlers.ServiceAccount.Get,
			)
			serviceaccounts.GET("/:namespace/:name/raw",
				app.authorizedMiddleware.HasPermission(permissions.SecretsView),
				app.Handlers.RawPatch.GetNamespaced(schema.GroupVersionResource{Version: "v1", Resource: "serviceaccounts"}),
			)
			serviceaccounts.POST("",
				app.authorizedMiddleware.HasPermission(permissions.SecretsCreate),
				app.Handlers.ServiceAccount.Create,
			)
			serviceaccounts.PUT("",
				app.authorizedMiddleware.HasPermission(permissions.SecretsCreate),
				app.Handlers.ServiceAccount.Update,
			)
			serviceaccounts.PATCH("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.SecretsCreate),
				app.Handlers.RawPatch.PatchNamespaced(schema.GroupVersionResource{Version: "v1", Resource: "serviceaccounts"}),
			)
			serviceaccounts.DELETE("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.SecretsDelete),
				app.Handlers.ServiceAccount.Delete,
			)
		}

		roles := v1.Group("/roles")
		{
			roles.GET("",
				app.authorizedMiddleware.HasPermission(permissions.SecretsView),
				app.Handlers.Role.List,
			)
			roles.GET("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.SecretsView),
				app.Handlers.Role.Get,
			)
			roles.GET("/:namespace/:name/raw",
				app.authorizedMiddleware.HasPermission(permissions.SecretsView),
				app.Handlers.RawPatch.GetNamespaced(schema.GroupVersionResource{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "roles"}),
			)
			roles.POST("",
				app.authorizedMiddleware.HasPermission(permissions.SecretsCreate),
				app.Handlers.Role.Create,
			)
			roles.PUT("",
				app.authorizedMiddleware.HasPermission(permissions.SecretsCreate),
				app.Handlers.Role.Update,
			)
			roles.PATCH("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.SecretsCreate),
				app.Handlers.RawPatch.PatchNamespaced(schema.GroupVersionResource{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "roles"}),
			)
			roles.DELETE("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.SecretsDelete),
				app.Handlers.Role.Delete,
			)
		}

		rolebindings := v1.Group("/rolebindings")
		{
			rolebindings.GET("",
				app.authorizedMiddleware.HasPermission(permissions.SecretsView),
				app.Handlers.RoleBinding.List,
			)
			rolebindings.GET("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.SecretsView),
				app.Handlers.RoleBinding.Get,
			)
			rolebindings.GET("/:namespace/:name/raw",
				app.authorizedMiddleware.HasPermission(permissions.SecretsView),
				app.Handlers.RawPatch.GetNamespaced(schema.GroupVersionResource{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "rolebindings"}),
			)
			rolebindings.POST("",
				app.authorizedMiddleware.HasPermission(permissions.SecretsCreate),
				app.Handlers.RoleBinding.Create,
			)
			rolebindings.PUT("",
				app.authorizedMiddleware.HasPermission(permissions.SecretsCreate),
				app.Handlers.RoleBinding.Update,
			)
			rolebindings.PATCH("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.SecretsCreate),
				app.Handlers.RawPatch.PatchNamespaced(schema.GroupVersionResource{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "rolebindings"}),
			)
			rolebindings.DELETE("/:namespace/:name",
				app.authorizedMiddleware.HasPermission(permissions.SecretsDelete),
				app.Handlers.RoleBinding.Delete,
			)
		}

		topology := v1.Group("/topology")
		topology.Use(app.authorizedMiddleware.HasPermission(permissions.TopologyView))
		{
			topology.GET("", app.Handlers.Topology.Get)
		}

		metrics := v1.Group("/metrics/stream")
		{
			metrics.GET("/nodes/:id", app.Handlers.Metrics.GetNodeMetrics)
			metrics.GET("/pods/:namespace/:id", app.Handlers.Metrics.GetPodMetrics)
		}
	}
}

func (app *App) Start() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		log.Println("Starting Event Batcher...")
		app.EventBatcher.Run(gCtx)
		return nil
	})

	log.Println("Starting Shared Informer Factory...")
	app.InformerFactory.Start(ctx.Done())

	log.Println("Waiting for cache sync...")
	results := app.InformerFactory.WaitForCacheSync(ctx.Done())
	for resType, synced := range results {
		if !synced {
			log.Fatalf("failed to sync cache for resource: %v", resType)
		}
	}
	log.Println("All caches synced successfully!")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: app.Router,
	}

	g.Go(func() error {
		log.Println("Starting HTTP Server on :8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server error: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		<-gCtx.Done()
		log.Println("Shutting down HTTP Server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return srv.Shutdown(shutdownCtx)
	})

	if err := g.Wait(); err != nil {
		log.Printf("App stopped with error: %v", err)
	} else {
		log.Println("App stopped gracefully")
	}
}
