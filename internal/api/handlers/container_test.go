package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHandlerContainer(t *testing.T) {
	podHandler := &PodHandler{}
	deploymentHandler := &DeploymentHandler{}
	namespaceHandler := &NamespaceHandler{}
	serviceHandler := &ServiceHandler{}
	nodeHandler := &NodeHandler{}
	terminalHandler := &TerminalHandler{}
	topologyHandler := &TopologyHandler{}
	podLogsHandler := &PodLogsHandler{}
	configMapHandler := &ConfigMapHandler{}
	secretHandler := &SecretHandler{}
	ingressHandler := &IngressHandler{}
	pvcHandler := &PvcHandler{}
	networkInspectorHandler := &NetworkInspectorHandler{}

	container := NewHandlerContainer(
		podHandler,
		deploymentHandler,
		namespaceHandler,
		serviceHandler,
		nodeHandler,
		terminalHandler,
		topologyHandler,
		podLogsHandler,
		configMapHandler,
		secretHandler,
		ingressHandler,
		pvcHandler,
		networkInspectorHandler,
	)

	assert.NotNil(t, container)
	assert.Equal(t, podHandler, container.Pod)
	assert.Equal(t, deploymentHandler, container.Deployment)
	assert.Equal(t, namespaceHandler, container.Namespace)
	assert.Equal(t, serviceHandler, container.Service)
	assert.Equal(t, nodeHandler, container.Node)
	assert.Equal(t, terminalHandler, container.Terminal)
	assert.Equal(t, topologyHandler, container.Topology)
	assert.Equal(t, podLogsHandler, container.PodLogs)
	assert.Equal(t, configMapHandler, container.ConfigMaps)
	assert.Equal(t, secretHandler, container.Secrets)
	assert.Equal(t, ingressHandler, container.Ingresses)
	assert.Equal(t, pvcHandler, container.Pvcs)
	assert.Equal(t, networkInspectorHandler, container.NetworkInspector)
}
