package provider

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	region = "ap-northeast-1"
)

type EC2API interface {
	AssociateAddress(input *ec2.AssociateAddressInput) (*ec2.AssociateAddressOutput, error)
	DescribeAddresses(input *ec2.DescribeAddressesInput) (*ec2.DescribeAddressesOutput, error)
	DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error)
}

type AWSProvider struct {
	client EC2API
}

func NewAWSProvider() (*AWSProvider, error) {
	s, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return nil, err
	}

	provider := &AWSProvider{
    	client: ec2.New(s),
	}

	return provider, nil
}

func (p *AWSProvider) AssociateAddress(allocationId string, instanceId string) (*ec2.AssociateAddressOutput, error) {
	input := &ec2.AssociateAddressInput{
		AllocationId: aws.String(allocationId),
		InstanceId:   aws.String(instanceId),
	}

	result, err := p.client.AssociateAddress(input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *AWSProvider) DescribeAddresses(input *ec2.DescribeAddressesInput) (*ec2.DescribeAddressesOutput, error) {
	result, err := p.client.DescribeAddresses(input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *AWSProvider) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	result, err := p.client.DescribeInstances(input)
	if err != nil {
		return nil, err
	}
	return result, nil
}
