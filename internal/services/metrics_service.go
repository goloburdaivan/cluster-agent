package services

import (
	"cluster-agent/internal/models"
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

type MetricsService interface {
	GetNodeMetrics(ctx context.Context, id string) (models.NodeMetrics, error)
	GetPodMetrics(ctx context.Context, namespace, id string) (models.PodMetrics, error)
}

type metricsService struct {
	client     metricsclientset.Interface
	kubeClient kubernetes.Interface
}

func NewMetricsService(
	client metricsclientset.Interface,
	kubeClient kubernetes.Interface,
) MetricsService {
	return &metricsService{
		client:     client,
		kubeClient: kubeClient,
	}
}

func (m *metricsService) GetNodeMetrics(ctx context.Context, id string) (models.NodeMetrics, error) {
	metrics, err := m.client.MetricsV1beta1().NodeMetricses().Get(ctx, id, metav1.GetOptions{})
	if err != nil {
		return models.NodeMetrics{}, err
	}

	node, err := m.kubeClient.CoreV1().Nodes().Get(ctx, id, metav1.GetOptions{})
	if err != nil {
		return models.NodeMetrics{}, err
	}

	cpuCapacity := node.Status.Allocatable.Cpu().AsApproximateFloat64()
	cpuUsage := metrics.Usage.Cpu().AsApproximateFloat64()
	cpuPercent := (cpuUsage / cpuCapacity) * 100

	memCapacity := float64(node.Status.Allocatable.Memory().Value())
	memUsage := float64(metrics.Usage.Memory().Value())
	memPercent := (memUsage / memCapacity) * 100

	return models.NodeMetrics{
		Name:        metrics.Name,
		Timestamp:   metrics.Timestamp.Time,
		CPUUsage:    cpuPercent,
		MemoryUsage: memPercent,
	}, nil
}

func (m *metricsService) GetPodMetrics(ctx context.Context, namespace, id string) (models.PodMetrics, error) {
	metrics, err := m.client.MetricsV1beta1().PodMetricses(namespace).Get(ctx, id, metav1.GetOptions{})
	if err != nil {
		return models.PodMetrics{}, err
	}

	timestamp := metrics.Timestamp.Time
	containerMetrics := make([]models.ContainerMetrics, 0, len(metrics.Containers))

	for _, container := range metrics.Containers {
		containerMetrics = append(containerMetrics, models.ContainerMetrics{
			Name:        container.Name,
			Timestamp:   timestamp,
			CPUUsage:    container.Usage.Cpu().AsApproximateFloat64(),
			MemoryUsage: float64(container.Usage.Memory().Value()) / (1024.0 * 1024.0),
		})
	}

	return models.PodMetrics{
		Name:             metrics.Name,
		Namespace:        metrics.Namespace,
		Timestamp:        timestamp,
		ContainerMetrics: containerMetrics,
	}, nil
}
