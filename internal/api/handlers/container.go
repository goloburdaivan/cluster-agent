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
	NewTopologyHandler,
)

type HandlerContainer struct {
	Pod        *PodHandler
	Deployment *DeploymentHandler
	Namespace  *NamespaceHandler
	Service    *ServiceHandler
	Node       *NodeHandler
	Terminal   *TerminalHandler
	Topology   *TopologyHandler
}

func NewHandlerContainer(
	pod *PodHandler,
	deployment *DeploymentHandler,
	namespace *NamespaceHandler,
	service *ServiceHandler,
	node *NodeHandler,
	terminal *TerminalHandler,
	topology *TopologyHandler,
) *HandlerContainer {
	return &HandlerContainer{
		Pod:        pod,
		Deployment: deployment,
		Namespace:  namespace,
		Service:    service,
		Node:       node,
		Terminal:   terminal,
		Topology:   topology,
	}
}
