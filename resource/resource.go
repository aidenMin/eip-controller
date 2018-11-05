package resource

type InstanceInfo struct {
	InstanceId 		string
	PrivateDnsName	string
}

type EipInfo struct {
	AllocationId	string
	InstanceId      string
	TagValue		string
}

type Resource interface {
	FindAllocatableInstance() ([]InstanceInfo, error)
	FindAllocatableEip() (*EipInfo, error)
	FindAssociatedEip() ([]EipInfo, error)
	FindInstanceById(instanceId string) (*InstanceInfo, error)
	AssociateAddress(allocationId string, instanceId string) (string, error)
}
