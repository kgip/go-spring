package configuration

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"strings"
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
	splitedKeys := strings.Split(configKey, ".")
	configMap := configuration.configs
	for i := 0; i < len(splitedKeys); i++ {
		isMatched := false
		for mapKey, value := range configMap {
			if mapKey == splitedKeys[i] {
				if i >= len(splitedKeys) {
					return value
				} else {
					isMatched = true
					configMap = value.(map[string]interface{})
					break
				}
			}
		}
		if !isMatched {
			break
		}
	}
	return nil
}
