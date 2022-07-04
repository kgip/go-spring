package main

import "github.com/kgip/go-spring/ioc"

type Mysql struct {
	MysqlConfig `prefix:"mysql"`
}

type MysqlConfig struct {
	Path      string            `configKey:"value:path default:10.4.48.44:3306"` // 服务器地址:端口
	Dbname    string            `configKey:"dbname"`                             // 数据库名
	Username  string            `configKey:"username"`                           // 数据库用户名
	Password  string            `configKey:"password"`                           // 数据库密码
	SubConfig map[string]string `prefix:"sub-config"`
}

type MysqlAllConfig struct {
	Path      string            `configKey:"mysql.path"`     // 服务器地址:端口
	Dbname    string            `configKey:"mysql.dbname"`   // 数据库名
	Username  string            `configKey:"mysql.username"` // 数据库用户名
	Password  string            `configKey:"mysql.password"` // 数据库密码
	SubConfig map[string]string `prefix:"mysql.sub-config"`
}

func (MysqlAllConfig) ConfigurationPrefix() string {
	return "mysql"
}

func main() {
	ioc.RegisterModules()
	ioc.RegisterSimpleBean()
	ioc.RegisterSimpleFactoryBean()
	ioc.Start()
}
