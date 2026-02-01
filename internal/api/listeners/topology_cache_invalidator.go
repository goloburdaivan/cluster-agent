package listeners

import (
	"cluster-agent/internal/events"
	"context"
)

type TopologyCacheInvalidator interface {
	Invalidate(ctx context.Context, namespace string) error
}

type TopologyCacheListener struct {
	invalidator TopologyCacheInvalidator
}

func NewTopologyCacheListener(invalidator TopologyCacheInvalidator) *TopologyCacheListener {
	return &TopologyCacheListener{
		invalidator: invalidator,
	}
}

func (l *TopologyCacheListener) Handle(ctx context.Context, e *events.ClusterChangedEvent) error {
	return l.invalidator.Invalidate(ctx, e.Namespace)
}
