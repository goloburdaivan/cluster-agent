package services

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type TerminalService interface {
	GetAuthExecutor(namespace, podName, container string) (remotecommand.Executor, error)
}

type terminalService struct {
	clientset kubernetes.Interface
	config    *rest.Config
}

func NewTerminalService(clientset kubernetes.Interface, config *rest.Config) TerminalService {
	return &terminalService{
		clientset: clientset,
		config:    config,
	}
}

func (t *terminalService) GetAuthExecutor(namespace, podName, container string) (remotecommand.Executor, error) {
	req := t.clientset.CoreV1().RESTClient().Get().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		Param("container", container).
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "true").
		Param("command", "/bin/sh")

	return remotecommand.NewWebSocketExecutor(t.config, "GET", req.URL().String())
}
