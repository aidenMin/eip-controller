package resource

import (
	"github.com/aidenMin/eip-controller/provider"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	tagName = "EipGroup"
)

type AWSEC2 struct {
	provider provider.Provider
}

func NewAWSEC2(provider provider.Provider) (*AWSEC2, error) {
	return &AWSEC2{
		provider: provider,
	}, nil
}

func (awsec2 *AWSEC2) FindAllocatableInstance() (map[string]string, error) {
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

  return ExtractData(result.Reservations), nil
}

func (awsec2 *AWSEC2) FindAllocatableEip() (string, string, error) {
	input := &ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag-key"),
				Values: []*string{
					aws.String(tagName),
				},
			},
		},
	}
	result, err := awsec2.provider.DescribeAddresses(input)
	if err != nil {
		return "", "", err
	}

	addr := FilterNotAssociatedAddress(result.Addresses)
	if addr == nil {
		return "", "", nil
	}

	tag := FilterTagByTagName(addr, tagName)
	if tag == nil {
		return *addr.AllocationId, "", nil
	}

	return *addr.AllocationId, *tag.Value, nil
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

func ExtractData(reservations []*ec2.Reservation) map[string]string {
	var data = make(map[string]string)
	for _, reservation := range reservations {
		for _, instance := range reservation.Instances {
			data[*instance.InstanceId] = *instance.PrivateDnsName
		}
	}

	return data
}
