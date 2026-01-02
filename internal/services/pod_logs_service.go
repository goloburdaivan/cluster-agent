package services

import (
	"context"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type PodLogsService interface {
	StreamLogs(ctx context.Context, namespace string, podName string, containerName string) (io.ReadCloser, error)
}

type podLogsService struct {
	client kubernetes.Interface
}

func NewPodLogsService(client kubernetes.Interface) PodLogsService {
	return &podLogsService{
		client: client,
	}
}

func (p podLogsService) StreamLogs(ctx context.Context, namespace string, podName string, containerName string) (io.ReadCloser, error) {
	logOpts := &corev1.PodLogOptions{
		Container:  containerName,
		Follow:     true,
		TailLines:  new(int64),
		Timestamps: true,
	}

	*logOpts.TailLines = 100

	req := p.client.CoreV1().Pods(namespace).GetLogs(podName, logOpts)

	stream, err := req.Stream(ctx)

	if err != nil {
		return nil, err
	}

	return stream, nil
}
