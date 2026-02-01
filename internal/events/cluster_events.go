package events

type ClusterChangedEvent struct {
	Namespace string
}

func (e *ClusterChangedEvent) Name() string {
	return "cluster.changed"
}
