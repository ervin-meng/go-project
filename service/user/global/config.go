package global

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Db       string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type ConsulConfig struct {
	IP   string `mapstructure:"ip" json:"ip"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServiceConfig struct {
	Name   string       `mapstructure:"name" json:"name"`
	IP     string       `mapstructure:"ip" json:"ip"`
	Port   int          `mapstructure:"port" json:"port"`
	Mysql  MysqlConfig  `mapstruct:"mysql" json:"mysql"`
	Consul ConsulConfig `mapstruct:"consul" json:"consul"`
}

type NacosConfig struct {
	IP          string `mapstructure:"ip" json:"ip"`
	Port        int    `mapstructure:"port" json:"port"`
	NamespaceId string `mapstructure:"namespaceId" json:"namespaceId"`
	DataId      string `mapstructure:"dataId" json:"dataId"`
	Group       string `mapstructure:"group" json:"group"`
}

var Config *ServiceConfig
