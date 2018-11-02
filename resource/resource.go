package resource

type Resource interface {
	FindAllocatableInstance() (map[string]string, error)
	FindAllocatableEip() (string, string, error)
	AssociateAddress(allocationId string, instanceId string) (string, error)
}
