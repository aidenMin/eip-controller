package controller

import (
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/aidenMin/eip-controller/k8s"
	"github.com/aidenMin/eip-controller/resource"
)

func init()  {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

type Controller struct {
	Resource	resource.Resource
	K8s			*k8s.KubeNode
	Interval 	time.Duration
}

func (c *Controller) Run() {
	for {
		c.RunOnce()
		log.Info("finished eip-controller")

		c.Verify()
		log.Info("verified k8s-node-label")
		time.Sleep(c.Interval)
	}
}

func (c *Controller) RunOnce() {
	// 할당 가능한 EIP 가 있는지 확인한다.
	_, err := c.Resource.FindAllocatableEip()
	if err != nil {
		return
	}

	// EIP 에 연동되어 있지 않은 EC2 정보를 가져온다.
	instanceInfo, err := c.Resource.FindAllocatableInstance()
	if err != nil {
		log.Error(err)
		return
	}

	for _, info := range instanceInfo {
		c.AssociateEip(info.InstanceId, info.PrivateDnsName)
	}
}

func (c *Controller) AssociateEip(instanceId, privateDnsName string)  {
	// k8s 노드에 특정 label 이 할당되어 있는 EC2에만 프로세스를 수행한다.
	// label 은 외부에서 지정 가능하며, default 값은 "upbit.com/eip-group"이다.
	_, err := c.K8s.FindByNodeName(privateDnsName)
	if err != nil{
		return
	}

	log.Info("ec2: ", instanceId)

	// 할당 가능한 Eip's allocationId 를 가져온다.
	eipInfo, err := c.Resource.FindAllocatableEip()
	if err != nil {
		return
	}
	log.Infof("allocationId: %s, tagValue: %s", eipInfo.AllocationId, eipInfo.TagValue)

	// Eip 할당!!
	associationId, err := c.Resource.AssociateAddress(eipInfo.AllocationId, instanceId)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("associationId: ", associationId)

	// k8s에 label 설정
	_, err = c.K8s.SetLabel(privateDnsName, "BBBB")
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("associated...")
}

func (c *Controller) Verify() {
	eipInfoList, err := c.Resource.FindAssociatedEip()
	if err != nil {
		log.Error(err)
	}

	for _, eipInfo := range eipInfoList {
		instance, err := c.Resource.FindInstanceById(eipInfo.InstanceId)
		if err != nil {
			log.Error(err)
			continue
		}

		label, err := c.K8s.FindLabelValueByNodeName(instance.PrivateDnsName)
		if err != nil {
			log.Error(err)
			continue
		}
		if label != eipInfo.TagValue {
			c.K8s.SetLabel(instance.PrivateDnsName, eipInfo.TagValue)
			log.Infof("update %s's label: %s -> %s (%s)", instance.InstanceId, label, eipInfo.TagValue, eipInfo.AllocationId)
		}
	}
}
