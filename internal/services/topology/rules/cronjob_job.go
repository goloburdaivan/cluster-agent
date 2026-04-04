package rules

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
)

type CronJobJobRule struct{}

func (r *CronJobJobRule) Apply(
	s *models.ClusterSnapshot,
	b *graph.Builder,
) error {

	// CronJob -> Job (based on owner references)
	for _, job := range s.Jobs {
		for _, owner := range job.OwnerReferences {
			if owner.Kind == "CronJob" {
				b.AddEdge(edge(
					id("CronJob", job.Namespace, owner.Name),
					id("Job", job.Namespace, job.Name),
				))
			}
		}
	}

	return nil
}
