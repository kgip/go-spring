package ioc

import (
	"github.com/kgip/go-spring/configuration"
	"github.com/kgip/go-spring/core"
	errors "github.com/kgip/go-spring/error"
)

type ModuleRegister interface {
	Register()
}

func RegisterModules(registers ...ModuleRegister) {
	for _, register := range registers {
		if register != nil {
			register.Register()
		}
	}
}

func verifyBean(bean *core.Bean) bool {
	//bean名称不能为空
	if bean.GetName() == "" {
		return false
	}
	//model和factoryMethod不能同时为空
	if bean.GetModel() == nil && bean.GetFactoryMethod() == nil {
		return false
	}
	return true
}

func RegisterBeans(beans ...*core.Bean) {
	for _, bean := range beans {
		if !verifyBean(bean) {
			panic(errors.BeanIllegalError)
		}
	}
	for _, bean := range beans {
		container.AddBean(bean)
	}
}

func registerBeans(action func(o interface{}) *core.Bean, o ...interface{}) {
	if len(o) > 0 {
		beans := make([]*core.Bean, len(o))
		for i := 0; i < len(o); i++ {
			beans[i] = action(o)
		}
		RegisterBeans(beans...)
	}
}

func RegisterSimpleBean(model ...interface{}) {
	registerBeans(func(o interface{}) *core.Bean {
		return core.NewBean(model)
	}, model...)
}

func RegisterSimpleFactoryBean(factoryMethod ...interface{}) {
	registerBeans(func(o interface{}) *core.Bean {
		return core.NewFactoryBean(factoryMethod)
	}, factoryMethod...)
}

func RegisterBeanPreProcessors(processors ...core.BeanPreProcessor) {
	for _, processor := range processors {
		if processor == nil {
			panic(errors.NilError)
		}
		container.AddBeanPreProcessor(processor)
	}
}

func RegisterBeanPostProcessors(processors ...core.BeanPostProcessor) {
	for _, processor := range processors {
		if processor == nil {
			panic(errors.NilError)
		}
		container.AddBeanPostProcessor(processor)
	}
}

func RegisterPreProcessors(processors ...core.ContainerPreProcessor) {
	for _, processor := range processors {
		if processor == nil {
			panic(errors.NilError)
		}
		container.AddContainerPreProcessor(processor)
	}
}

func RegisterPostProcessors(processors ...core.ContainerPostProcessor) {
	for _, processor := range processors {
		if processor == nil {
			panic(errors.NilError)
		}
		container.AddContainerPostProcessor(processor)
	}
}

func SetConfiguration(provider configuration.Provider) {
	if provider == nil {
		panic(errors.NilError)
	}
	container.SetConfiguration(provider)
}

func setConfigInfo(action func(config *configuration.Configuration)) bool {
	config := container.GetConfiguration()
	if config == nil {
		panic(errors.NilError)
	}
	if defaultConfig, ok := config.(*configuration.Configuration); ok {
		action(defaultConfig)
		return true
	}
	return false
}

func SetConfigPath(path string) bool {
	return setConfigInfo(func(config *configuration.Configuration) {
		config.SetPath(path)
	})
}

func SetConfigType(configType string) bool {
	return setConfigInfo(func(config *configuration.Configuration) {
		config.SetConfigType(configType)
	})
}

func SetConfigRefresh(refresh bool) bool {
	return setConfigInfo(func(config *configuration.Configuration) {
		config.SetRefresh(refresh)
	})
}

func Start() {
	container.Init()
}
