package controller

import (
	"errors"
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

type EipInfo struct {
	allocationId	string
	tagValue		string
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
	// 할당 가능한 EIP 가 있는지 확인한다.
	_, err := c.FindAllocatableEip()
	if err != nil {
		log.Warning(err)
		return
	}

	// EIP 에 연동되어 있지 않은 EC2 정보를 가져온다.
	// 반환값은 다읨의 형식을 따른다.
	// { instanceId string: privateDnsName string }
	instanceMap, err := c.Resource.FindAllocatableInstance()
	if err != nil {
		log.Error(err)
		return
	}

	for instanceId, privateDnsName := range instanceMap {
		c.AssociateEip(instanceId, privateDnsName)
	}
}

func (c *Controller) AssociateEip(instanceId, privateDnsName string)  {
	// k8s 노드에 특정 label 이 할당되어 있는 EC2에만 프로세스를 수행한다.
	// label 은 외부에서 지정 가능하며, default 값은 "upbit.com/eip-group"이다.
	_, err := c.Source.FindNodeByLabelName(privateDnsName)
	if err != nil{
		return
	}

	log.Info("ec2: ", instanceId)

	// 할당 가능한 Eip's allocationId 를 가져온다.
	eipInfo, err := c.FindAllocatableEip()
	if err != nil {
		return
	}
	log.Infof("allocationId: %s, tagValue: %s", eipInfo.allocationId, eipInfo.tagValue)

	// Eip 할당!!
	associationId, err := c.Resource.AssociateAddress(eipInfo.allocationId, instanceId)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("associationId: ", associationId)

	// k8s에 label 설정
	_, err = c.Source.SetLabel(privateDnsName, eipInfo.tagValue)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("associated...")
}

func (c *Controller) FindAllocatableEip() (*EipInfo, error) {
	allocationId, tagValue, err := c.Resource.FindAllocatableEip()
	if err != nil {
		return nil, err
	}

	if allocationId == "" || tagValue == "" {
		return nil, errors.New("not found allocatable eip")
	}

	return &EipInfo{
		allocationId: allocationId,
		tagValue: tagValue}, nil
}
