package global

type UserServiceConfig struct {
	IP   string `mapstructure:"ip" json:"ip"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type ServiceConfig struct {
	User UserServiceConfig `mapstructure:"user" json:"user"`
}

type ConsulConfig struct {
	IP   string `mapstructure:"ip" json:"ip"`
	Port int    `mapstructure:"port" json:"port"`
}

type ApiConfig struct {
	Name    string        `mapstructure:"name" json:"name"`
	Port    int           `mapstructure:"port" json:"port"`
	Service ServiceConfig `mapstructure:"service" json:"service"`
	Consul  ConsulConfig  `mapstructure:"consul" json:"consul"`
}

type NacosConfig struct {
	IP          string `mapstructure:"ip" json:"ip"`
	Port        int    `mapstructure:"port" json:"port"`
	NamespaceId string `mapstructure:"namespaceId" json:"namespaceId"`
	DataId      string `mapstructure:"dataId" json:"dataId"`
	Group       string `mapstructure:"group" json:"group"`
}

var Config *ApiConfig
