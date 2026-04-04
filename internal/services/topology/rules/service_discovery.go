package rules

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
	"strings"
)

type ServiceDiscoveryRule struct{}

func (r *ServiceDiscoveryRule) Apply(
	s *models.ClusterSnapshot,
	b *graph.Builder,
) error {
	servicesMap := make(map[string]string)
	for _, svc := range s.Services {
		servicesMap[svc.Name] = id("Service", svc.Namespace, svc.Name)
	}

	for _, d := range s.Deployments {
		deployNodeID := id("Deployment", d.Namespace, d.Name)

		for _, container := range d.Spec.Template.Spec.Containers {
			for _, env := range container.Env {
				envValue := strings.TrimSpace(env.Value)

				if len(envValue) < 3 {
					continue
				}

				for svcName, svcID := range servicesMap {
					if strings.Contains(envValue, svcName) {
						b.AddEdge(edge(deployNodeID, svcID))
					}
				}
			}
		}
	}

	for _, ss := range s.StatefulSets {
		ssNodeID := id("StatefulSet", ss.Namespace, ss.Name)

		for _, container := range ss.Spec.Template.Spec.Containers {
			for _, env := range container.Env {
				envValue := strings.TrimSpace(env.Value)

				if len(envValue) < 3 {
					continue
				}

				for svcName, svcID := range servicesMap {
					if strings.Contains(envValue, svcName) {
						b.AddEdge(edge(ssNodeID, svcID))
					}
				}
			}
		}
	}

	for _, ds := range s.DaemonSets {
		dsNodeID := id("DaemonSet", ds.Namespace, ds.Name)

		for _, container := range ds.Spec.Template.Spec.Containers {
			for _, env := range container.Env {
				envValue := strings.TrimSpace(env.Value)

				if len(envValue) < 3 {
					continue
				}

				for svcName, svcID := range servicesMap {
					if strings.Contains(envValue, svcName) {
						b.AddEdge(edge(dsNodeID, svcID))
					}
				}
			}
		}
	}

	for _, j := range s.Jobs {
		jNodeID := id("Job", j.Namespace, j.Name)

		for _, container := range j.Spec.Template.Spec.Containers {
			for _, env := range container.Env {
				envValue := strings.TrimSpace(env.Value)

				if len(envValue) < 3 {
					continue
				}

				for svcName, svcID := range servicesMap {
					if strings.Contains(envValue, svcName) {
						b.AddEdge(edge(jNodeID, svcID))
					}
				}
			}
		}
	}

	return nil
}
