package error

import (
	"fmt"
)

type IocError struct {
	message string
	detail  string
}

func (e *IocError) Error() string {
	return fmt.Sprintf(`{"message": "%s", "detail": "%s"}`, e.message, e.detail)
}

func (e IocError) Detail(detail string) IocError {
	e.detail = detail
	return e
}

var (
	TypeNotMatchError           = &IocError{message: "Parameter type mismatch"}
	NilError                    = &IocError{message: "Parameter nil"}
	FactoryMethodReturnsError   = &IocError{message: "The number of return values of the factory method is not unique"}
	ContainerUpdateError        = &IocError{message: "The container has been initialized and cannot be updated"}
	BeanIllegalError            = &IocError{message: "Invalid bean information"}
	NameEmptyError              = &IocError{message: "Bean name can't be empty"}
	CircularReferenceError      = &IocError{message: "Cannot depend on the factory bean being created"}
	UnknownBeanNameError        = &IocError{message: "Unknown bean name"}
	ConfigKeyError              = &IocError{message: "config key error"}
	UnknownConfigKeySubKeyError = &IocError{message: "unknown config key sub key"}
	ConfigKeySubKeyResolveError = &IocError{message: "config key sub key resolve failed"}
)
