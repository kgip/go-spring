package ioc

import "log"

// Container bean容器
type Container struct {
	beans         map[string]*Bean
	Configuration interface{}
	Logger        *log.Logger
}

func (c *Container) GetBean(name string) interface{} {
	if bean := c.beans[name]; bean != nil {

	}
	return nil
}

func (c *Container) AddBean(bean *Bean) {
	c.beans[bean.name] = bean
}

// Bean 表示一个对象
type Bean struct {
	name               string
	value              interface{}
	isSingleton        bool //是否单例
	beanPreProcessors  []BeanPreProcessor
	beanPostProcessors []BeanPostProcessor
}
