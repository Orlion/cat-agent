package config

import (
	"io/ioutil"

	"github.com/Orlion/cat-agent/cat"
	"github.com/Orlion/cat-agent/log"
	"github.com/Orlion/cat-agent/server"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Cat    *cat.Config
	Server *server.Config
	Log    *log.Config
}

func ParseConfig(filename string) (config *Config, err error) {
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	config = new(Config)
	err = yaml.Unmarshal(fileData, config)
	if err != nil {
		return
	}

	return
}
