package config

import (
	"errors"
	"fmt"

	"github.com/Orlion/cat-agent/pkg/systemx"
)

var config *Config

type Config struct {
	Domain        string   `yaml:"domain"`
	Hostname      string   `yaml:"hostname"`
	Env           string   `yaml:"env"`
	Ip            string   `yaml:"ip"`
	IpHex         string   `yaml:"ip_hex"`
	RouterServers []string `yaml:"router-servers"`
}

func GetDomain() string {
	return config.Domain
}

func Init(conf *Config) error {
	if err := withDefaultConf(conf); err != nil {
		return err
	}

	config = conf

	return nil
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
	if len(config.RouterServers) < 1 {
		return errors.New("router servers cannot be empty")
	}

	return nil
}
