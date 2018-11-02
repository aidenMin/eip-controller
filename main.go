package main

import (
	"github.com/aidenMin/eip-controller/client"
	"github.com/aidenMin/eip-controller/config"
	"github.com/aidenMin/eip-controller/controller"
	"github.com/aidenMin/eip-controller/provider"
	"github.com/aidenMin/eip-controller/resource"
	"github.com/aidenMin/eip-controller/source"
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

	r, err := resource.NewAWSEC2(p)
	if err != nil {
		log.Panic(err)
	}


	k := client.NewKubeClient(cfg.KubeConfig)
	s, err := source.NewKubeNode(k, cfg.NodeLabel)
	if err != nil {
		log.Panic(err)
	}

	ctrl := controller.Controller{
		Resource: r,
		Source:   s,
		Interval: cfg.Interval,
	}

	ctrl.Run()
}


