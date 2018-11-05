package main

import (
	"github.com/aidenMin/eip-controller/config"
	"github.com/aidenMin/eip-controller/controller"
	"github.com/aidenMin/eip-controller/k8s"
	"github.com/aidenMin/eip-controller/provider"
	"github.com/aidenMin/eip-controller/resource"
	log "github.com/sirupsen/logrus"
	"os"
)

func init()  {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {
	cfg := config.NewConfig()
	if err := cfg.ParseFlags(os.Args[1:]); err != nil {
		log.Panic(err)
	}

	logLevel, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Panic(err)
	}
	log.SetLevel(logLevel)

	doProcess(cfg)
}

func doProcess(cfg *config.Config) {
	p, err := provider.NewAWSProvider()
	if err != nil {
		log.Panic(err)
	}

	r, err := resource.NewAWSEC2(p, cfg.EipTag)
	if err != nil {
		log.Panic(err)
	}

	c, err := k8s.NewClient(cfg.KubeConfig)
	if err != nil {
		log.Panic(err)
	}

	n, err := k8s.NewKubeNode(*c, cfg.NodeLabel)
	if err != nil {
		log.Panic(err)
	}

	ctrl := controller.Controller{
		Resource: 	r,
		K8s:  	 	n,
		Interval: 	cfg.Interval,
	}

	ctrl.Run()
}


