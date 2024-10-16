package config

type Redis struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
	Channel string `yaml:"channel"`
}

type Consumers struct {
	Count    int    `yaml:"count"`
	ListName string `yaml:"list_name"`
}

type Processor struct {
	ProcessedEventsStream string `yaml:"processed_events_stream"`
}
