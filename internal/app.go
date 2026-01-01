package internal

import (
	"cluster-agent/internal/consumers"
	"cluster-agent/internal/producers"
	"context"
	"errors"
	"fmt"
	"k8s.io/client-go/informers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cluster-agent/internal/api/handlers"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type App struct {
	Router          *gin.Engine
	Handlers        *handlers.HandlerContainer
	EventCollector  *producers.EventCollector
	EventBatcher    *consumers.EventBatcher
	InformerFactory informers.SharedInformerFactory
}

func NewApp(
	h *handlers.HandlerContainer,
	collector *producers.EventCollector,
	batcher *consumers.EventBatcher,
	factory informers.SharedInformerFactory,
) *App {
	app := &App{
		Router:          gin.Default(),
		Handlers:        h,
		EventCollector:  collector,
		EventBatcher:    batcher,
		InformerFactory: factory,
	}

	app.setRoutes()

	return app
}

func (app *App) setRoutes() {
	v1 := app.Router.Group("/api/v1")
	{
		pods := v1.Group("/pods")
		{
			pods.GET("", app.Handlers.Pod.List)
			pods.GET("/:namespace/:name", app.Handlers.Pod.Get)
			pods.GET("/:namespace/:name/exec", app.Handlers.Terminal.Exec)
		}

		deployments := v1.Group("/deployments")
		{
			deployments.GET("", app.Handlers.Deployment.List)
			deployments.GET("/:namespace/:name", app.Handlers.Deployment.Get)
			deployments.POST("", app.Handlers.Deployment.Create)
			deployments.DELETE("/:namespace/:name", app.Handlers.Deployment.Delete)
			deployments.PATCH("/scale", app.Handlers.Deployment.ScaleDeployment)
		}

		services := v1.Group("/services")
		{
			services.GET("", app.Handlers.Service.List)
		}

		namespace := v1.Group("/namespaces")
		{
			namespace.GET("", app.Handlers.Namespace.List)
		}

		node := v1.Group("/nodes")
		{
			node.GET("", app.Handlers.Node.List)
		}

		topology := v1.Group("/topology")
		{
			topology.GET("", app.Handlers.Topology.Get)
		}
	}
}

func (app *App) Start() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, gCtx := errgroup.WithContext(ctx)

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

	g.Go(func() error {
		log.Println("Starting Event Batcher...")
		app.EventBatcher.Run(gCtx)
		return nil
	})

	g.Go(func() error {
		log.Println("Starting Shared Informer Factory...")
		app.InformerFactory.Start(gCtx.Done())

		log.Println("Waiting for cache sync...")
		results := app.InformerFactory.WaitForCacheSync(gCtx.Done())

		for resType, synced := range results {
			if !synced {
				return fmt.Errorf("failed to sync cache for resource: %v", resType)
			}
		}

		log.Println("All caches synced successfully!")
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Printf("App stopped with error: %v", err)
	} else {
		log.Println("App stopped gracefully")
	}
}
