package client

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

func NewKubeClient(kubeConfig string) kubernetes.Clientset {
	if kubeConfig == "" {
		if _, err := os.Stat(clientcmd.RecommendedHomeFile); err == nil {
			kubeConfig = clientcmd.RecommendedHomeFile
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		log.Panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err)
	}

	return *clientset
}
