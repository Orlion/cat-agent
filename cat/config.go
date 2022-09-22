package cat

import (
	"errors"
	"fmt"

	"github.com/Orlion/cat-agent/pkg/systemx"
)

type Config struct {
	Domain        string   `yaml:"domain"`
	Hostname      string   `yaml:"hostname"`
	Env           string   `yaml:"env"`
	Ip            string   `yaml:"ip"`
	IpHex         string   `yaml:"ip_hex"`
	RouterServers []string `yaml:"router-servers"`
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
		config.Hostname = defaultHostname
	}

	if config.Env == "" {
		config.Env = defaultEnv
	}

	fmt.Println(config)
	if len(config.RouterServers) < 1 {
		return errors.New("router servers cannot be empty")
	}

	return nil
}
