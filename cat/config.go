package cat

import (
	"errors"

	"github.com/Orlion/cat-agent/pkg/systemx"
)

type Config struct {
	CatServerVersion string
	Domain           string
	Hostname         string
	Env              string
	Ip               string
	IpHex            string
	routerServers    []string
}

func withDefaultConf(config *Config) error {
	if len(config.CatServerVersion) < 1 {
		config.CatServerVersion = defaultCatServerVersion
	}

	if config.Domain == "" {
		return errors.New("domain cannot be empty.")
	}

	var err error
	if config.Hostname, err = systemx.GetHostname(); err != nil {
		config.Hostname = defaultHostname
	}

	if config.Env == "" {
		config.Env = defaultEnv
	}

	if len(config.routerServers) < 1 {
		return errors.New("router servers cannot be empty.")
	}

	return nil
}
