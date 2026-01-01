package rules

import (
	"cluster-agent/internal/services/graph"
	corev1 "k8s.io/api/core/v1"
	"strings"
)

func labelsMatch(
	selector map[string]string,
	labels map[string]string,
) bool {

	if len(selector) == 0 {
		return false
	}

	for k, v := range selector {
		if labels[k] != v {
			return false
		}
	}
	return true
}

func id(kind, namespace, name string) string {
	return strings.ToLower(kind) + ":" + namespace + "/" + name
}

func edge(source, target string) graph.Edge {
	return graph.Edge{
		Source: source,
		Target: target,
	}
}

func usedInDeployment(
	podSpec corev1.PodSpec,
	resourceName string,
	kind string,
) bool {

	for _, v := range podSpec.Volumes {

		if kind == "ConfigMap" &&
			v.ConfigMap != nil &&
			v.ConfigMap.Name == resourceName {
			return true
		}

		if kind == "Secret" &&
			v.Secret != nil &&
			v.Secret.SecretName == resourceName {
			return true
		}
	}

	containers := append(
		podSpec.InitContainers,
		podSpec.Containers...,
	)

	for _, c := range containers {

		for _, ef := range c.EnvFrom {

			if kind == "ConfigMap" &&
				ef.ConfigMapRef != nil &&
				ef.ConfigMapRef.Name == resourceName {
				return true
			}

			if kind == "Secret" &&
				ef.SecretRef != nil &&
				ef.SecretRef.Name == resourceName {
				return true
			}
		}

		for _, e := range c.Env {

			if e.ValueFrom == nil {
				continue
			}

			if kind == "ConfigMap" &&
				e.ValueFrom.ConfigMapKeyRef != nil &&
				e.ValueFrom.ConfigMapKeyRef.Name == resourceName {
				return true
			}

			if kind == "Secret" &&
				e.ValueFrom.SecretKeyRef != nil &&
				e.ValueFrom.SecretKeyRef.Name == resourceName {
				return true
			}
		}
	}

	return false
}
