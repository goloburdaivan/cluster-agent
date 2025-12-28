package handlers

import "github.com/google/wire"

var HandlerSet = wire.NewSet(
	NewPodHandler,
	NewDeploymentHandler,
	NewNamespaceHandler,
	NewServiceHandler,
	NewNodeHandler,
	NewTerminalHandler,
	NewHandlerContainer,
)

type HandlerContainer struct {
	Pod        *PodHandler
	Deployment *DeploymentHandler
	Namespace  *NamespaceHandler
	Service    *ServiceHandler
	Node       *NodeHandler
	Terminal   *TerminalHandler
}

func NewHandlerContainer(
	pod *PodHandler,
	deployment *DeploymentHandler,
	namespace *NamespaceHandler,
	service *ServiceHandler,
	node *NodeHandler,
	terminal *TerminalHandler,
) *HandlerContainer {
	return &HandlerContainer{
		Pod:        pod,
		Deployment: deployment,
		Namespace:  namespace,
		Service:    service,
		Node:       node,
		Terminal:   terminal,
	}
}
