package internal

import (
	"cluster-agent/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

type App struct {
	Router   *gin.Engine
	Handlers *handlers.HandlerContainer
}

func NewApp(h *handlers.HandlerContainer) *App {
	app := &App{
		Router:   gin.Default(),
		Handlers: h,
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
	}
}

func (app *App) Start() {
	app.Router.Run(":8080")
}
