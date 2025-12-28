package services

import (
	"cluster-agent/internal/models"
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NodeService interface {
	GetNodes(ctx context.Context) ([]models.Node, error)
}

type nodeService struct {
	clientset kubernetes.Interface
}

func NewNodeService(clientset kubernetes.Interface) NodeService {
	return &nodeService{
		clientset: clientset,
	}
}

func (n *nodeService) GetNodes(ctx context.Context) ([]models.Node, error) {
	nodes, err := n.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	result := make([]models.Node, 0, len(nodes.Items))

	for _, node := range nodes.Items {
		result = append(result, models.Node{
			Name:    node.Name,
			Status:  getNodeStatus(&node),
			Role:    getNodeRole(&node),
			Version: node.Status.NodeInfo.KubeletVersion,
		})
	}

	return result, nil
}
