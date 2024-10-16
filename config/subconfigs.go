package config

type Redis struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
	Channel string `yaml:"channel"`
}

type Consumers struct {
	Count int `yaml:"count"`
}
