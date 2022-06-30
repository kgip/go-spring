package ioc

// ContainerPreProcessor 容器前置处理器
type ContainerPreProcessor func(c *Container)

// BeanPreProcessor bean的前置处理器
type BeanPreProcessor func(c *Container, value interface{})

type BeanPreProcessorInterface interface {
	PreProcess(c *Container)
}

// Initializer bean初始化器
type Initializer interface {
	Init(c *Container)
}

// BeanPostProcessor bean的后置处理器,bean初始化后执行
type BeanPostProcessor func(c *Container, value interface{})

type BeanPostProcessorInterface interface {
	PostProcess(c *Container)
}

// ContainerPostProcessor 容器后置处理器
type ContainerPostProcessor func(c *Container)
