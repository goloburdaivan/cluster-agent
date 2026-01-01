package rules

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
	"strings"
)

type ResourceNodesRule struct {
}

func (r *ResourceNodesRule) Apply(
	s *models.ClusterSnapshot,
	b *graph.Builder,
) error {

	for _, d := range s.Deployments {
		b.AddNode(node("Deployment", d.Namespace, d.Name))
	}
	for _, s := range s.Services {
		b.AddNode(node("Service", s.Namespace, s.Name))
	}
	for _, ss := range s.StatefulSets {
		b.AddNode(node("StatefulSet", ss.Namespace, ss.Name))
	}
	for _, i := range s.Ingresses {
		b.AddNode(node("Ingress", i.Namespace, i.Name))
	}
	for _, cm := range s.ConfigMaps {
		b.AddNode(node("ConfigMap", cm.Namespace, cm.Name))
	}
	for _, sec := range s.Secrets {
		b.AddNode(node("Secret", sec.Namespace, sec.Name))
	}
	for _, pvc := range s.PVCs {
		b.AddNode(node("PVC", pvc.Namespace, pvc.Name))
	}

	return nil
}

func node(kind, ns, name string) graph.Node {
	return graph.Node{
		ID:   strings.ToLower(kind) + ":" + ns + "/" + name,
		Kind: kind,
		Name: name,
	}
}
