package jconsul

import (
	"errors"
	"fmt"
	consulApi "github.com/hashicorp/consul/api"
)

const Address = "192.168.33.10"
const Port = 9527

//consul 结构体
type JConsul struct {
	Host string
	Port int64
}

//获取配置
func (c *JConsul) GetConfig() *consulApi.Config {
	address := fmt.Sprintf("http://%s:%d", c.Host, c.Port)
	config := &consulApi.Config{
		Address: address,
	}
	return config
}

//获取客户端
func (c *JConsul) GetClient() (*consulApi.Client, error) {
	config := c.GetConfig()
	return consulApi.NewClient(config)
}

//获取健康检查
func (c *JConsul) GetHealthCheck(checkHost string, checkPort int64, checkPath string) *consulApi.AgentServiceCheck {
	http:=fmt.Sprintf("http://%s:%d%s", checkHost, checkPort, checkPath)
	jlogger.Info(http)
	check := &consulApi.AgentServiceCheck{
		HTTP:                           http,
		Timeout:                        "30s",
		Interval:                       "3s",
		DeregisterCriticalServiceAfter: "30s", //check失败后30秒删除本服务
	}
	return check
}

//获取服务注册信息
func (c *JConsul) GetRegistration(serviceId string, serviceName string, serviceTags []string, serviceAddress string, servicePort int64, check *consulApi.AgentServiceCheck) *consulApi.AgentServiceRegistration {
	//服务注册信息
	registration := new(consulApi.AgentServiceRegistration)
	registration.ID = serviceId
	registration.Name = serviceName
	registration.Tags = serviceTags
	registration.Address = serviceAddress
	registration.Port = int(servicePort)
	//健康检查
	//registration.Check = check
	return registration
}

//服务注册
func (c *JConsul) ServiceRegister(serviceId string, serviceName string, serviceTags []string, serviceAddress string, servicePort int64, checkHost string, checkPort int64, checkPath string) error {
	client, err := c.GetClient()
	if err != nil {
		jlogger.Error("consul new client error:" + err.Error())
		return err
	}
	//健康检查
	check := c.GetHealthCheck(checkHost, checkPort, checkPath)
	//服务注册信息
	registration := c.GetRegistration(serviceId, serviceName, serviceTags, serviceAddress, servicePort, check)
	//注册
	agent := client.Agent()
	if agent == nil {
		jlogger.Error("consul get agent error")
		return errors.New("consul get agent error")
	}
	err = agent.ServiceRegister(registration)
	if err != nil {
		jlogger.Error("consul service register error:" + err.Error())
		return err
	}
	return nil
}

//服务发现
func (c *JConsul) ServiceDiscovery() (map[string]*consulApi.AgentService, error) {
	client, err := c.GetClient()
	if err != nil {
		jlogger.Error("consul new client error:" + err.Error())
		return nil, err
	}
	//发现
	agent := client.Agent()
	if agent == nil {
		jlogger.Error("consul get agent error")
		return nil, errors.New("consul get agent error")
	}
	services, err1 := agent.Services()
	if err1 != nil {
		jlogger.Error("consul get services error:" + err1.Error())
		return nil, err1
	}
	return services, nil
}

//获取实例
func NewJConsul(host string, port int64) *JConsul {
	c := &JConsul{
		Host: host,
		Port: port,
	}
	return c
}
