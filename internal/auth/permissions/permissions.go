package permissions

type Permission string

func (p Permission) String() string {
	return string(p)
}

const (
	// Nodes
	NodesView Permission = "nodes:view"

	// Pods
	PodsView   Permission = "pods:view"
	PodsCreate Permission = "pods:create"
	PodsDelete Permission = "pods:delete"

	// Deployments
	DeploymentsView   Permission = "deployments:view"
	DeploymentsCreate Permission = "deployments:create"
	DeploymentsDelete Permission = "deployments:delete"
	DeploymentsScale  Permission = "deployments:scale"

	// Namespaces
	NamespacesCreate Permission = "namespaces:create"
	NamespacesDelete Permission = "namespaces:delete"

	// Events
	EventsView Permission = "events:view"

	// Topology
	TopologyView Permission = "topology:view"

	// Services
	ServicesView   Permission = "services:view"
	ServicesCreate Permission = "services:create"
	ServicesDelete Permission = "services:delete"

	// Ingresses
	IngressesView   Permission = "ingresses:view"
	IngressesCreate Permission = "ingresses:create"
	IngressesDelete Permission = "ingresses:delete"

	// ConfigMaps
	ConfigMapsView   Permission = "configmaps:view"
	ConfigMapsCreate Permission = "configmaps:create"
	ConfigMapsDelete Permission = "configmaps:delete"

	// Secrets
	SecretsView   Permission = "secrets:view"
	SecretsCreate Permission = "secrets:create"
	SecretsDelete Permission = "secrets:delete"

	// PersistentVolumeClaims
	PVCsView   Permission = "pvcs:view"
	PVCsCreate Permission = "pvcs:create"
	PVCsDelete Permission = "pvcs:delete"

	// Metrics
	NodeMetricsView Permission = "node_metrics:view"
	PodMetricsView  Permission = "pod_metrics:view"
)
