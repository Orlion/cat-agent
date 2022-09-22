package log

type Config struct {
	Level    int8   `yaml:"level"`
	Filename string `yaml:"filename"`
}
