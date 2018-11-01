package main

import (
	"fmt"
	"os"

	"github.com/aidenMin/eip-controller/client"
	"github.com/aidenMin/eip-controller/controller"
	"github.com/aidenMin/eip-controller/config"
	"github.com/aidenMin/eip-controller/provider"
	"github.com/aidenMin/eip-controller/resource"
	"github.com/aidenMin/eip-controller/source"
)

func main() {
	fmt.Println("================================")

	cfg := config.NewConfig()
	if err := cfg.ParseFlags(os.Args[1:]); err != nil {
		fmt.Printf("flag parsing error: %v\n", err)
	}

	p, err := provider.NewAWSProvider()
	if err != nil {
		fmt.Println("Error creating session", err)
		return
	}

	r, err := resource.NewAWSEC2(p)
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	k, err := client.NewKubeClient("")
	if err != nil {
		fmt.Println("Error creating session", err)
		return
	}

	s, err := source.NewKubeNode(k)
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	ctrl := controller.Controller{
		Resource: r,
		Source:   s,
		Interval: cfg.Interval,
	}
	ctrl.Run()
}


