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
	log.Info("================================")

	cfg := config.NewConfig()
	if err := cfg.ParseFlags(os.Args[1:]); err != nil {
		log.Panic("flag parsing error:", err)
	}

	p, err := provider.NewAWSProvider()
	if err != nil {
		log.Panic("Error creating session", err)
	}

	r, err := resource.NewAWSEC2(p)
	if err != nil {
		log.Panic("Error", err)
	}

	k, err := client.NewKubeClient("")
	if err != nil {
		log.Panic("Error creating session", err)
	}

	s, err := source.NewKubeNode(k)
	if err != nil {
		log.Panic("Error", err)
	}

	ctrl := controller.Controller{
		Resource: r,
		Source:   s,
		Interval: cfg.Interval,
	}
	ctrl.Run()
}


