package global

import "github.com/ervin-meng/go-stitch-monster/infrastructure/config"

type UserServiceConfig struct {
	IP   string `mapstructure:"ip" json:"ip"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type ServiceConfig struct {
	User UserServiceConfig `mapstructure:"user" json:"user"`
}

type ApiConfig struct {
	Name    string              `mapstructure:"name" json:"name"`
	IP      string              `mapstructure:"ip" json:"ip"`
	Port    int                 `mapstructure:"port" json:"port"`
	Service ServiceConfig       `mapstructure:"service" json:"service"`
	Consul  config.ConsulConfig `mapstructure:"consul" json:"consul"`
}

var Config *ApiConfig
