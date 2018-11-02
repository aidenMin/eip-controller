package controller

import (
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/aidenMin/eip-controller/resource"
	"github.com/aidenMin/eip-controller/source"
)

func init()  {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

type Controller struct {
	Resource	resource.Resource
	Source		*source.KubeNode
	Interval 	time.Duration
}

func (c *Controller) Run() {
	for {
		c.RunOnce()
		time.Sleep(c.Interval)
	}
}

func (c *Controller) RunOnce() {
	instanceMap, err := c.Resource.FindAllNotAssociatedEC2InstanceToEip()
	if err != nil {
		log.Fatal(err)
	}

	for instanceId, privateDnsName := range instanceMap {
		c.AssociateEip(instanceId, privateDnsName)
		log.Info("-----")
	}
}

func (c *Controller) AssociateEip(instanceId, privateDnsName string)  {
	log.Info("EC2 InstancesId:", instanceId)
	allocationId, err := c.Resource.FindNotAllocatedEipAllocationId()
	if err != nil {
		log.Panic(err)
	}
	log.Info("EIP AllocationId:", allocationId)

	associationId, err := c.Resource.AssociateEipToEC2(allocationId, instanceId)
	if err != nil {
		log.Panic(err)
	}
	log.Info("AssociationId:", associationId)

	eipGroupName, err := c.Resource.FindEipGroupNameByAllocationId(allocationId)
	if err != nil {
		log.Panic(err)
	}
	log.Info("EipGroupName:", eipGroupName)

	result, err := c.Source.SetLabel(privateDnsName, "EipGroup", eipGroupName)
	if err != nil {
		log.Panic(err)
	}
	log.Info("Complete:", result)
}
