package resource

import (
	"errors"
	"github.com/aidenMin/eip-controller/provider"
	//"fmt"

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

func (awsec2 *AWSEC2) FindAllNotAssociatedEC2InstanceToEip() (map[string]string, error) {
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

func (awsec2 *AWSEC2) FindNotAllocatedEipAllocationId() (string, error) {
  result, err := awsec2.provider.DescribeAddresses(&ec2.DescribeAddressesInput{})
  if err != nil {
    return "", err
  }

  addr := FindAddress(result.Addresses)
  if addr == nil {
    return "", errors.New("NotFoundAvailableEipException")
  }
  return *addr.AllocationId, nil
}

func (awsec2 *AWSEC2) FindEipGroupNameByAllocationId(allocationId string) (string, error) {
  input := &ec2.DescribeAddressesInput{
      Filters: []*ec2.Filter{
          {
              Name: aws.String("allocation-id"),
              Values: []*string{
                  aws.String(allocationId),
              },
          },
      },
  }
  result, err := awsec2.provider.DescribeAddresses(input)
  if err != nil {
    return "", err
  }

	tag := FindTag(result.Addresses)
	if tag == nil {
		return "", errors.New("NotFoundEipGroupException")
	}
	return *tag.Value, nil
}

func (awsec2 *AWSEC2) AssociateEipToEC2(allocationId string, instanceId string) (string, error) {
  result, err := awsec2.provider.AssociateAddress(allocationId, instanceId)
  if err != nil {
    return "", err
  }
  return *result.AssociationId, nil
}

func FindAddress(addresses []*ec2.Address) *ec2.Address {
	for _, addr := range addresses {

		// Eip에 연결된 EC2인스턴스가 있다면 패스
		if addr.AssociationId != nil {
			continue
		}

		for _, tag := range addr.Tags {
			if *tag.Key == string(tagName) {
				return addr
			}
		}
	}
	return nil
}

func FindTag(addresses []*ec2.Address) *ec2.Tag {
	for _, addr := range addresses {
		for _, tag := range addr.Tags {
			if *tag.Key == tagName {
				return tag
			}
		}
	}
	return nil
}

func ExtractData(reservations []*ec2.Reservation) map[string]string {
	var data = make(map[string]string)
	for _, reservation := range reservations {
		for _, instance := range reservation.Instances {
			//fmt.Println(instance)
			data[*instance.InstanceId] = *instance.PrivateDnsName
		}
	}

	return data
}
