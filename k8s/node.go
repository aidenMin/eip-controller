package k8s

import (
	"errors"
	"fmt"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)


type KubeNode struct {
	client 		kubernetes.Clientset
	nodeLabel	string
}

func NewKubeNode(client kubernetes.Clientset, nodeLabel string) (*KubeNode, error) {
	return &KubeNode{
		client: 	client,
		nodeLabel:	nodeLabel,
	}, nil
}

func (kube *KubeNode) SetLabel(nodeName string, labelValue string) (map[string]string, error) {
	node, err := kube.FindByNodeName(nodeName)
	if err != nil {
		return nil, err
	}

	node.Labels[kube.nodeLabel] = labelValue
	node.SetLabels(node.Labels)

	result, err := kube.client.CoreV1().Nodes().Update(node)
	if err != nil {
		return nil, err
	}
	return result.Labels, nil
}

func (kube *KubeNode) FindByNodeName(nodeName string) (*v1.Node, error) {
	input := metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", nodeName),
		LabelSelector: kube.nodeLabel,
	}

	nodes, err := kube.client.CoreV1().Nodes().List(input)
	if err != nil {
		return nil, err
	}

	if len(nodes.Items) == 0 {
		return nil, errors.New("not found node")
	}

	return &nodes.Items[0], nil
}


func (kube *KubeNode) FindLabelValueByNodeName(nodeName string) (string, error) {
	node, err := kube.FindByNodeName(nodeName)
	if err != nil {
		return "", err
	}

	return node.Labels[kube.nodeLabel], nil
}