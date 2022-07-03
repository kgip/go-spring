package configuration

type Provider interface {
	Load()
	GetConfig(configKey string) interface{}
}

// Storage 实现该接口的类被视为配置类
type Storage interface {
	ConfigurationPrefix() string
}
