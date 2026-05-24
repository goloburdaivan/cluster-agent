package rules

import (
	"cluster-agent/internal/models"
	"cluster-agent/internal/services/graph"
	"log/slog"
	"strings"
)

type ResourceNodesRule struct {
}

func (r *ResourceNodesRule) Apply(
	s *models.ClusterSnapshot,
	b *graph.Builder,
) error {

	// Workloads
	for _, d := range s.Deployments {
		b.AddNode(node("Deployment", d.Namespace, d.Name))
	}
	for _, ss := range s.StatefulSets {
		b.AddNode(node("StatefulSet", ss.Namespace, ss.Name))
	}
	for _, ds := range s.DaemonSets {
		b.AddNode(node("DaemonSet", ds.Namespace, ds.Name))
	}
	for _, j := range s.Jobs {
		b.AddNode(node("Job", j.Namespace, j.Name))
	}
	for _, cj := range s.CronJobs {
		b.AddNode(node("CronJob", cj.Namespace, cj.Name))
	}

	// Network
	for _, svc := range s.Services {
		b.AddNode(node("Service", svc.Namespace, svc.Name))
	}
	for _, i := range s.Ingresses {
		b.AddNode(node("Ingress", i.Namespace, i.Name))
	}
	slog.Error("debug", "len", len(s.NetworkPolicies))
	for _, np := range s.NetworkPolicies {
		b.AddNode(node("NetworkPolicy", np.Namespace, np.Name))
	}

	// Storage
	for _, pvc := range s.PVCs {
		b.AddNode(node("PVC", pvc.Namespace, pvc.Name))
	}
	for _, pv := range s.PVs {
		b.AddNode(node("PV", "", pv.Name))
	}
	for _, sc := range s.StorageClasses {
		b.AddNode(node("StorageClass", "", sc.Name))
	}

	// Config
	for _, cm := range s.ConfigMaps {
		b.AddNode(node("ConfigMap", cm.Namespace, cm.Name))
	}
	for _, sec := range s.Secrets {
		b.AddNode(node("Secret", sec.Namespace, sec.Name))
	}

	// RBAC
	for _, sa := range s.ServiceAccounts {
		b.AddNode(node("ServiceAccount", sa.Namespace, sa.Name))
	}
	for _, role := range s.Roles {
		b.AddNode(node("Role", role.Namespace, role.Name))
	}
	for _, rb := range s.RoleBindings {
		b.AddNode(node("RoleBinding", rb.Namespace, rb.Name))
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
