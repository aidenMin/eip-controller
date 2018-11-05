package resource

import (
	"errors"
	"github.com/aidenMin/eip-controller/provider"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type AWSEC2 struct {
	provider 	provider.Provider
	eipTag		string
}

func NewAWSEC2(provider provider.Provider, eipTag string) (*AWSEC2, error) {
	return &AWSEC2{
		provider: 	provider,
		eipTag:		eipTag,
	}, nil
}

func (awsec2 *AWSEC2) FindAllocatableInstance() ([]InstanceInfo, error) {
	// EIP가 할당되지 않은 인스턴스를 가져온다.
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("network-interface.addresses.association.ip-owner-id"),
				Values: []*string{
					aws.String("amazon"),
				},
			},
		},
  	}

	result, err := awsec2.provider.DescribeInstances(input)
	if err != nil {
		return nil, err
	}
	return NormalizeData(result.Reservations), nil
}

func (awsec2 *AWSEC2) FindInstanceById(instanceId string) (*InstanceInfo, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-id"),
				Values: []*string{
					aws.String(instanceId),
				},
			},
		},
	}

	result, err := awsec2.provider.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	data := NormalizeData(result.Reservations)
	if len(data) == 0 {
		return nil, errors.New("not found instance")
	}

	return &data[0], nil
}

func (awsec2 *AWSEC2) FindAssociatedEip() ([]EipInfo, error) {
	var data []EipInfo
	result, err := awsec2.provider.DescribeAddresses(&ec2.DescribeAddressesInput{})
	if err != nil {
		return nil, err
	}

	for _, addr := range result.Addresses {
		if addr.AssociationId != nil {
			tag := FilterTagByTagName(addr, awsec2.eipTag)

			var tagValue = ""
			if tag != nil {
				tagValue = *tag.Value
			}

			data = append(data, EipInfo{
				AllocationId: *addr.AllocationId,
				InstanceId: *addr.InstanceId,
				TagValue: tagValue,
			})
		}
	}
	return data, nil
}

func (awsec2 *AWSEC2) FindAllocatableEip() (*EipInfo, error) {
	input := &ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag-key"),
				Values: []*string{
					aws.String(awsec2.eipTag),
				},
			},
		},
	}
	result, err := awsec2.provider.DescribeAddresses(input)
	if err != nil {
		return nil, err
	}

	addr := FilterNotAssociatedAddress(result.Addresses)
	if addr == nil {
		return nil, errors.New("not found address")
	}

	tag := FilterTagByTagName(addr, awsec2.eipTag)
	if tag == nil {
		return nil, errors.New("not found tag in address")
	}

	return &EipInfo{
		AllocationId: *addr.AllocationId,
		TagValue: *tag.Value,
	}, nil
}

func (awsec2 *AWSEC2) AssociateAddress(allocationId string, instanceId string) (string, error) {
  result, err := awsec2.provider.AssociateAddress(allocationId, instanceId)
  if err != nil {
    return "", err
  }
  return *result.AssociationId, nil
}

func FilterNotAssociatedAddress(addresses []*ec2.Address) *ec2.Address {
	for _, addr := range addresses {
		if addr.AssociationId == nil {
			return addr
		}
	}
	return nil
}

func FilterTagByTagName(address *ec2.Address, tagName string) *ec2.Tag {
	for _, tag := range address.Tags {
		if *tag.Key == tagName {
			return tag
		}
	}
	return nil
}

func NormalizeData(reservations []*ec2.Reservation) []InstanceInfo {
	var data []InstanceInfo
	for _, reservation := range reservations {
		for _, instance := range reservation.Instances {
			data = append(data, InstanceInfo{
				InstanceId: *instance.InstanceId,
				PrivateDnsName: *instance.PrivateDnsName,
			})
		}
	}

	return data
}
