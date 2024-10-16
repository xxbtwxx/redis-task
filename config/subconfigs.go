package config

type Redis struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
}

type Consumers struct {
	Count int `yaml:"count"`
}
