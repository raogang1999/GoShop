package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

type Registry struct {
	Host string
	Port int
}

type RegistryClient interface {
	Register(addr string, port int, name string, tags []string, id string) (err error)
	DeRegister(serviceID string) error
}

func NewRegistry(host string, port int) RegistryClient {
	/*
		host: 注册中心的ip
		port: 注册中心的端口
	*/
	return &Registry{
		Host: host,
		Port: port,
	}
}

func (self *Registry) Register(addr string, port int, name string, tags []string, id string) (err error) {
	/*
		addr: 指本服务运行在哪个host上
		port: 本服务运行在哪个端口上
		name: 服务名称
		tags: 服务标签
		id: 服务id
	*/
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("http://%s:%d", self.Host, self.Port) //consul的ip和port

	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}

	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", addr, port),
		Timeout:                        "3s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	registration := new(api.AgentServiceRegistration)
	registration.Name = name
	registration.ID = id
	registration.Port = port
	registration.Tags = tags
	registration.Address = addr
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	return nil
}

func (self *Registry) DeRegister(serviceID string) error {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("http://%s:%d", self.Host, self.Port) //consul的ip和port

	client, err := api.NewClient(config)
	if err != nil {
		return err
	}
	err = client.Agent().ServiceDeregister(serviceID)
	return err
}
