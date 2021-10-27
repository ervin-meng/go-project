package register

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

const (
	HTTPService = iota
	RPCService
)

var Client *api.Client

func Init(ip string, port int) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", ip, port)
	Client, _ = api.NewClient(cfg)
}

func ServiceRegister(serviceType int, id string, name string, ip string, port int) {

	var check *api.AgentServiceCheck

	switch serviceType {
	case RPCService:
		check = &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", ip, port),
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "10s",
		}
	case HTTPService:
		fallthrough
	default:
		check = &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/health", ip, port),
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "10s",
		}
	}

	reg := new(api.AgentServiceRegistration)
	reg.ID = id
	reg.Name = name
	reg.Address = ip
	reg.Port = port
	reg.Tags = []string{name}

	reg.Check = check

	err := Client.Agent().ServiceRegister(reg)

	if err != nil {
		panic(err)
	}
}

func ServiceDeregister(id string) error {
	return Client.Agent().ServiceDeregister(id)
}
