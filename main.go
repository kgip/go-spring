package main

import (
	"fmt"
	"reflect"
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
	//ioc.RegistryBeans()
	//ioc.SetConfigPath("./config.yaml")
	//ioc.SetConfigType("yaml")
	rt := reflect.TypeOf(Mysql{})
	rv := reflect.ValueOf(Mysql{})
	fmt.Println(&rv)
	f := rt.Field(0)
	fmt.Println(f.Name)
	fmt.Println(f.Anonymous)
	fmt.Println(f.Tag.Get("prefix"))
}
