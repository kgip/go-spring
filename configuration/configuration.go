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
	viper      *viper.Viper
	logger     *log.Logger
}

func NewConfiguration(path string, configType string, refresh bool, logger *log.Logger) *Configuration {
	return &Configuration{path: path, refresh: refresh, configType: configType, logger: logger, viper: viper.New()}
}

func (c *Configuration) SetPath(path string) {
	c.path = path
}

func (c *Configuration) SetConfigType(configType string) {
	c.configType = configType
}

func (c *Configuration) SetRefresh(refresh bool) {
	c.refresh = refresh
}

func (c *Configuration) Load() {
	c.logger.Printf("load configuration from path %s", c.path)
	c.viper.SetConfigFile(c.path)
	c.viper.SetConfigType(c.configType)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if c.refresh {
		c.viper.OnConfigChange(func(e fsnotify.Event) {
			c.logger.Println("config file changed")
			c.configs = viper.AllSettings()
		})
	}
	c.configs = viper.AllSettings()
	c.logger.Println("initialize config complete")
}

func (c *Configuration) GetConfig(configKey string) interface{} {
	splitedKeys := strings.Split(configKey, ".")
	configMap := c.configs
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
