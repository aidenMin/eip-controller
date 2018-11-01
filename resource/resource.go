package resource

type Resource interface {
	FindAllNotAssociatedEC2InstanceToEip() (map[string]string, error)
	FindNotAllocatedEipAllocationId() (string, error)
	FindEipGroupNameByAllocationId(allocationId string) (string, error)
	AssociateEipToEC2(allocationId string, instanceId string) (string, error)
}
