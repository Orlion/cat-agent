package config

import (
	"errors"
	"io/ioutil"

	catconfig "github.com/Orlion/cat-agent/cat/config"
	"github.com/Orlion/cat-agent/log"
	"github.com/Orlion/cat-agent/server"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Cat    *catconfig.Config `yaml:"cat"`
	Server *server.Config    `yaml:"server"`
	Log    *log.Config       `yaml:"log"`
}

func ParseConfig(filename string) (config *Config, err error) {
	if filename == "" {
		err = errors.New("please enter a configuration file name")
		return
	}

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
