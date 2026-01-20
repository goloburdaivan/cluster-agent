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
	NewPodLogsHandler,
	NewConfigMapHandler,
	NewSecretHandler,
	NewIngressHandler,
	NewPvcHandler,
	NewNetworkInspectorHandler,
)

type HandlerContainer struct {
	Pod              *PodHandler
	Deployment       *DeploymentHandler
	Namespace        *NamespaceHandler
	Service          *ServiceHandler
	Node             *NodeHandler
	Terminal         *TerminalHandler
	Topology         *TopologyHandler
	PodLogs          *PodLogsHandler
	ConfigMaps       *ConfigMapHandler
	Secrets          *SecretHandler
	Ingresses        *IngressHandler
	Pvcs             *PvcHandler
	NetworkInspector *NetworkInspectorHandler
}

func NewHandlerContainer(
	pod *PodHandler,
	deployment *DeploymentHandler,
	namespace *NamespaceHandler,
	service *ServiceHandler,
	node *NodeHandler,
	terminal *TerminalHandler,
	topology *TopologyHandler,
	podLogs *PodLogsHandler,
	configmaps *ConfigMapHandler,
	secrets *SecretHandler,
	ingresses *IngressHandler,
	pvcs *PvcHandler,
	networkInspector *NetworkInspectorHandler,
) *HandlerContainer {
	return &HandlerContainer{
		Pod:              pod,
		Deployment:       deployment,
		Namespace:        namespace,
		Service:          service,
		Node:             node,
		Terminal:         terminal,
		Topology:         topology,
		PodLogs:          podLogs,
		ConfigMaps:       configmaps,
		Secrets:          secrets,
		Ingresses:        ingresses,
		Pvcs:             pvcs,
		NetworkInspector: networkInspector,
	}
}
