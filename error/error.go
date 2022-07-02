package error

import "errors"

var (
	TypeNotMatchError         = errors.New("Parameter type mismatch")
	NilError                  = errors.New("Parameter nil")
	FactoryMethodReturnsError = errors.New("The number of return values of the factory method is not unique")
	ContainerUpdateError      = errors.New("The container has been initialized and cannot be updated")
	BeanIllegalError          = errors.New("Invalid bean information")
)
