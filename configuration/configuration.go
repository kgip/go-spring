package configuration

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

type Configuration struct {
	configs    map[string]interface{}
	path       string
	refresh    bool //是否刷新配置
	configType string
	logger     *log.Logger
}

func NewConfiguration(path string, refresh bool, configType string, logger *log.Logger) *Configuration {
	config := &Configuration{path: path, refresh: refresh, configType: configType, logger: logger}
	//加载配置
	config.Load()
	return config
}

func (configuration *Configuration) AppendConfig() {

}

func (configuration *Configuration) Load() {
	viper := viper.New()
	viper.SetConfigFile(configuration.path)
	viper.SetConfigType(configuration.configType)
	configuration.logger.Println("start load config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if configuration.refresh {
		viper.OnConfigChange(func(e fsnotify.Event) {
			configuration.logger.Println("config file changed")
			configuration.configs = viper.AllSettings()
		})
	}
	configuration.configs = viper.AllSettings()
	configuration.logger.Println("finished initializing config")
}

func (configuration *Configuration) GetConfig(configKey string) interface{} {
	return configuration.configs[configKey]
}
