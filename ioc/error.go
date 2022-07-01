package ioc

import "errors"

var (
	TypeNotMatchError         = errors.New("parameter type mismatch")
	NilError                  = errors.New("parameter nil")
	FactoryMethodReturnsError = errors.New("The number of return values of the factory method is not unique")
	ContainerUpdateError      = errors.New("The container has been initialized and cannot be updated")
)
