package configuration

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"strings"
)

type Configuration struct {
	Configs    map[string]interface{}
	Path       string
	Refresh    bool //是否刷新配置
	ConfigType string
	Logger     *log.Logger
}

func (configuration *Configuration) Load() {
	viper := viper.New()
	viper.SetConfigFile(configuration.Path)
	viper.SetConfigType(configuration.ConfigType)
	configuration.Logger.Println("start load config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if configuration.Refresh {
		viper.OnConfigChange(func(e fsnotify.Event) {
			configuration.Logger.Println("config file changed")
			configuration.Configs = viper.AllSettings()
		})
	}
	configuration.Configs = viper.AllSettings()
	configuration.Logger.Println("finished initializing config")
}

func (configuration *Configuration) GetConfig(configKey string) interface{} {
	splitedKeys := strings.Split(configKey, ".")
	configMap := configuration.Configs
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
