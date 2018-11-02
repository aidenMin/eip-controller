package source

import (
	"fmt"
	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

)


type KubeNode struct {
	client kubernetes.Clientset
}

func NewKubeNode(client kubernetes.Clientset) (*KubeNode, error) {
	return &KubeNode{
		client: client,
	}, nil
}

func (kube *KubeNode) SetLabel(nodeName string, labelKey string, labelValue string) (map[string]string, error) {
	node := kube.FindNode(nodeName)

	glog.Infoln("nodeName: %s\n, labelKey: %s\n, labelValue: %s\n", nodeName, labelKey, labelValue)

	node.Labels[labelKey] = labelValue
	node.SetLabels(node.Labels)

	result, err := kube.client.CoreV1().Nodes().Update(node)
	if err != nil {
		return nil, err
	}
	return result.Labels, nil
}

func (kube *KubeNode) FindNode(nodeName string) *v1.Node {
	input := metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", nodeName),
		LabelSelector: "EipGroup",
	}

	nodes, err := kube.client.CoreV1().Nodes().List(input)
	if err != nil {
		return nil
	}

	if len(nodes.Items) != 1 {
		return nil
	}

	return &nodes.Items[0]
}