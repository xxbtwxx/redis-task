package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Redis     *Redis     `yaml:"redis"`
	Consumers *Consumers `yaml:"consumers"`
	Processor *Processor `yaml:"processor"`
}

func Load(path string) (*Config, error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(fileData, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
