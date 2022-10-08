package log

type Config struct {
	StdoutLevel string `yaml:"stdout_level"`
	Level       string `yaml:"level"`
	Filename    string `yaml:"filename"`
}

func withDefaultConf(config *Config) *Config {
	if config == nil {
		config = &Config{
			StdoutLevel: "info",
			Level:       "error",
		}

		return config
	}

	if config.StdoutLevel == "" {
		config.StdoutLevel = "info"
	}

	if config.Level == "" {
		config.Level = "error"
	}

	return config
}
