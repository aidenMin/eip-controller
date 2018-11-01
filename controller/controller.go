package controller

import (
	"fmt"
	"time"

	"github.com/aidenMin/eip-controller/resource"
	"github.com/aidenMin/eip-controller/source"
)

type Controller struct {
	Resource	resource.Resource
	Source		*source.KubeNode
	Interval 	time.Duration
}

func (c *Controller) Run() {
	instanceMap, err := c.Resource.FindAllNotAssociatedEC2InstanceToEip()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for instanceId, privateDnsName := range instanceMap {
		c.AssociateEip(instanceId, privateDnsName)
	}
}

func (c *Controller) AssociateEip(instanceId, privateDnsName string)  {
	fmt.Printf("[EC2 InstancesId: %v]\n", instanceId)
	allocationId, err := c.Resource.FindNotAllocatedEipAllocationId()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("EIP AllocationId:", allocationId)

	associationId, err := c.Resource.AssociateEipToEC2(allocationId, instanceId)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("AssociationId:", associationId)

	eipGroupName, err := c.Resource.FindEipGroupNameByAllocationId(allocationId)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("EipGroupName:", eipGroupName)

	result, err := c.Source.SetLabel(privateDnsName, "EipGroup", eipGroupName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Complete:", result)
}
