package core

// ContainerPreProcessor 容器前置处理器
type ContainerPreProcessor interface {
	PreProcess(c *Container)
}

// BeanPreProcessor bean的前置处理器
type BeanPreProcessor interface {
	PreProcess(c *Container, bean *Bean)
}

// Initializer bean初始化器
type Initializer interface {
	Init(c *Container)
}

// BeanPostProcessor bean的后置处理器,bean初始化后执行
type BeanPostProcessor interface {
	PostProcess(c *Container, instance interface{})
}

// ContainerPostProcessor 容器后置处理器
type ContainerPostProcessor interface {
	PostProcess(c *Container)
}
