package core

type PriorityProvider interface {
	GetPriority() int
}

type BeanNameProvider interface {
	GetBeanName() string
}
