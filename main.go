package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go-spring/ioc"
	"log"
)

type Mysql struct {
	MysqlConfig `prefix:"mysql"`
}

type MysqlConfig struct {
	Path      string            `configKey:"path"`      // 服务器地址:端口
	Dbname    string            `configKey:"dbname"`    // 数据库名
	Username  string            `configKey:"username."` // 数据库用户名
	Password  string            `configKey:"password"`  // 数据库密码
	SubConfig map[string]string `prefix:"sub-config"`
}

type MysqlAllConfig struct {
	Path      string            `configKey:"mysql.path"`     // 服务器地址:端口
	Dbname    string            `configKey:"mysql.dbname"`   // 数据库名
	Username  string            `configKey:"mysql.username"` // 数据库用户名
	Password  string            `configKey:"mysql.password"` // 数据库密码
	SubConfig map[string]string `prefix:"mysql.sub-config"`
}

func main() {
	viper := viper.New()
	viper.SetConfigFile("./config.yaml")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("config file changed")
		viper.AllSettings()
	})
	fmt.Println(viper.AllSettings())
	log.Println("finished initializing config")
	ioc.RegistryBeans(nil, nil, nil)
}
