package ioc

import (
	"github.com/kgip/go-spring/configuration"
	"github.com/kgip/go-spring/core"
	"log"
)

const (
	defaultConfigPath = "./config.yaml"
	defaultConfigType = "yaml"
	refreshConfig     = false
)

var (
	container *core.Container
)

func init() {
	logger := log.Default()
	configurationProvider := configuration.NewConfiguration(defaultConfigPath, defaultConfigType, refreshConfig, logger)
	container = core.NewContainer(configurationProvider, logger)
	RegisterBeanPreProcessors()
	RegisterBeanPostProcessors(&core.AssignBeanPostProcessor{})
	RegisterPreProcessors()
	RegisterPostProcessors()
}
