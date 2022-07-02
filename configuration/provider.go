package configuration

type Provider interface {
	Load()
	GetConfig(configKey string) interface{}
}
