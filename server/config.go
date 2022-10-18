package server

type Config struct {
	Addr               string `yaml:"addr"`
	ReadTimeoutMillis  int    `yaml:"read_timeout_millis"`
	WriteTimeoutMillis int    `yaml:"write_timeout_millis"`
}

func withDefaultConf(config *Config) {
	if config.Addr == "" {
		config.Addr = "127.0.0.1:2280"
	}

	if config.ReadTimeoutMillis < 1 {
		config.ReadTimeoutMillis = 100
	}

	if config.WriteTimeoutMillis < 1 {
		config.WriteTimeoutMillis = 100
	}
}
