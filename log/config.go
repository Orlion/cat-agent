package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level    zapcore.Level `yaml:"level"`
	Filename string        `yaml:"filename"`
}

func withDefaultConf(config *Config) *Config {
	if config == nil {
		config = &Config{
			Level: zap.ErrorLevel,
		}

		return config
	}

	if config.Level < zap.DebugLevel {
		config.Level = zap.ErrorLevel
	}

	return config
}
