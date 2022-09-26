package config

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Orlion/cat-agent/pkg/systemx"
)

type Config struct {
	Domain   string   `yaml:"domain"`
	Hostname string   `yaml:"hostname"`
	Env      string   `yaml:"env"`
	Ip       string   `yaml:"ip"`
	IpHex    string   `yaml:"ip_hex"`
	Servers  []string `yaml:"servers"`
}

type ConfigService struct {
	config      *Config
	routers     []string
	mu          sync.RWMutex
	routersCond *sync.Cond
}

func (c *ConfigService) run() {

}

func (c *ConfigService) shutdown() {

}

func (c *ConfigService) GetDomain() string {
	return c.config.Domain
}

func (c *ConfigService) GetHostname() string {
	return c.config.Hostname
}

func (c *ConfigService) GetEnv() string {
	return c.config.Env
}

func (c *ConfigService) GetIp() string {
	return c.config.Ip
}

func (c *ConfigService) GetIpHex() string {
	return c.config.IpHex
}

func (c *ConfigService) GetServers() []string {
	return c.config.Servers
}

func (c *ConfigService) GetRouters() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.routers
}

func (c *ConfigService) RoutersCondWait() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.routersCond.Wait()
}

var instance *ConfigService

func GetInstance() *ConfigService {
	return instance
}

func Init(conf *Config) error {
	if err := withDefaultConf(conf); err != nil {
		return err
	}

	instance = &ConfigService{}
	instance.mu = sync.RWMutex{}
	instance.routersCond = sync.NewCond(&instance.mu)

	instance.run()
	return nil
}

func Shutdown() {
	instance.shutdown()
}

func withDefaultConf(config *Config) error {
	if config == nil {
		return errors.New("cat config cannot be empty")
	}

	if config.Domain == "" {
		return errors.New("domain cannot be empty")
	}

	var err error
	if config.Hostname, err = systemx.GetHostname(); err != nil {
		config.Hostname = DefaultHostname
	}

	if config.Env == "" {
		config.Env = DefaultEnv
	}

	fmt.Println(config)
	if len(config.Servers) < 1 {
		return errors.New("servers cannot be empty")
	}

	return nil
}
