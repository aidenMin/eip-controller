package provider

import (
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Provider interface {
	AssociateAddress(allocationId string, instanceId string) (*ec2.AssociateAddressOutput, error)
	DescribeAddresses(input *ec2.DescribeAddressesInput) (*ec2.DescribeAddressesOutput, error)
	DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error)
}
