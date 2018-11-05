package k8s

import (
log "github.com/sirupsen/logrus"
"k8s.io/client-go/kubernetes"
"k8s.io/client-go/tools/clientcmd"
"os"
)

func init()  {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

func NewClient(kubeConfig string) (*kubernetes.Clientset, error) {
	if kubeConfig == "" {
		if _, err := os.Stat(clientcmd.RecommendedHomeFile); err == nil {
			kubeConfig = clientcmd.RecommendedHomeFile
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
