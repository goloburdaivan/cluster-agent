package internal

import (
	"cluster-agent/internal/api/middleware"
	"cluster-agent/internal/auth/permissions"
	"cluster-agent/internal/consumers"
	"cluster-agent/internal/producers"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"k8s.io/client-go/informers"

	"cluster-agent/internal/api/handlers"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type App struct {
	Router               *gin.Engine
	Handlers             *handlers.HandlerContainer
	EventCollector       *producers.EventCollector
	EventBatcher         *consumers.EventBatcher
	InformerFactory      informers.SharedInformerFactory
	authorizedMiddleware *middleware.AuthorizedMiddleware
}

func NewApp(
	h *handlers.HandlerContainer,
	authorizedMiddleware *middleware.AuthorizedMiddleware,
	collector *producers.EventCollector,
	batcher *consumers.EventBatcher,
	factory informers.SharedInformerFactory,
) *App {
	app := &App{
		Router:               gin.Default(),
		Handlers:             h,
		authorizedMiddleware: authorizedMiddleware,
		EventCollector:       collector,
		EventBatcher:         batcher,
		InformerFactory:      factory,
	}

	app.setRoutes()

	return app
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

			deployments.POST("",
				app.authorizedMiddleware.HasPermission(permissions.DeploymentsCreate),
				app.Handlers.Deployment.Create,
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

		services := v1.Group("/services")
		services.Use(app.authorizedMiddleware.HasPermission(permissions.ServicesView))
		{
			services.GET("", app.Handlers.Service.List)
			services.GET("/:namespace/:name", app.Handlers.Service.Get)
		}

		configmaps := v1.Group("/configmaps")
		configmaps.Use(app.authorizedMiddleware.HasPermission(permissions.ConfigMapsView))
		{
			configmaps.GET("", app.Handlers.ConfigMaps.List)
			configmaps.GET("/:namespace/:name", app.Handlers.ConfigMaps.Get)
		}

		secrets := v1.Group("/secrets")
		secrets.Use(app.authorizedMiddleware.HasPermission(permissions.SecretsView))
		{
			secrets.GET("", app.Handlers.Secrets.List)
			secrets.GET("/:namespace/:name", app.Handlers.Secrets.Get)
		}

		ingresses := v1.Group("/ingresses")
		ingresses.Use(app.authorizedMiddleware.HasPermission(permissions.IngressesView))
		{
			ingresses.GET("", app.Handlers.Ingresses.List)
			ingresses.GET("/:namespace/:name", app.Handlers.Ingresses.Get)
		}

		pvcs := v1.Group("/persistentvolumeclaims")
		pvcs.Use(app.authorizedMiddleware.HasPermission(permissions.PVCsView))
		{
			pvcs.GET("", app.Handlers.Pvcs.List)
			pvcs.GET("/:namespace/:name", app.Handlers.Pvcs.Get)
		}

		namespace := v1.Group("/namespaces")
		{
			namespace.GET("", app.Handlers.Namespace.List)
		}

		node := v1.Group("/nodes")
		node.Use(app.authorizedMiddleware.HasPermission(permissions.NodesView))
		{
			node.GET("", app.Handlers.Node.List)
		}

		topology := v1.Group("/topology")
		topology.Use(app.authorizedMiddleware.HasPermission(permissions.TopologyView))
		{
			topology.GET("", app.Handlers.Topology.Get)
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
