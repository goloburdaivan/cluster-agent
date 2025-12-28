//go:build wireinject
// +build wireinject

package main

import (
	"cluster-agent/internal"
	"cluster-agent/internal/api/handlers"
	"cluster-agent/internal/k8s"
	"cluster-agent/internal/services"
	"github.com/google/wire"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func ProvideK8sInterface(client *k8s.Client) kubernetes.Interface {
	return client.GetClientset()
}

func ProvideRestConfig(client *k8s.Client) *rest.Config {
	return client.GetConfig()
}

func InitializeApp() (*internal.App, error) {
	wire.Build(
		k8s.NewClient,

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

		internal.NewApp,
	)
	return &internal.App{}, nil
}
