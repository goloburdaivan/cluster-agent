package permissions

type Permission string

func (p Permission) String() string {
	return string(p)
}

const (
	// Nodes
	NodesView Permission = "nodes:view"

	// Pods
	PodsView Permission = "pods:view"

	// Deployments
	DeploymentsView   Permission = "deployments:view"
	DeploymentsCreate Permission = "deployments:create"
	DeploymentsDelete Permission = "deployments:delete"
	DeploymentsScale  Permission = "deployments:scale"

	// Events
	EventsView Permission = "events:view"

	// Topology
	TopologyView Permission = "topology:view"

	// Services
	ServicesView Permission = "services:view"

	// Ingresses
	IngressesView Permission = "ingresses:view"

	// ConfigMaps
	ConfigMapsView Permission = "configmaps:view"

	// Secrets
	SecretsView Permission = "secrets:view"

	// PersistentVolumeClaims
	PVCsView Permission = "pvcs:view"
)
