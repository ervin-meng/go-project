package global

import "github.com/ervin-meng/go-stitch-monster/infrastructure/config"

type ServiceConfig struct {
	Name   string              `mapstructure:"name" json:"name"`
	IP     string              `mapstructure:"ip" json:"ip"`
	Port   int                 `mapstructure:"port" json:"port"`
	Mysql  config.MysqlConfig  `mapstruct:"mysql" json:"mysql"`
	Consul config.ConsulConfig `mapstruct:"consul" json:"consul"`
}

var Config *ServiceConfig
