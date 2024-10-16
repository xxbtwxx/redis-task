package main

import (
	"redis-task/config"

	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.Load("./config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	_ = cfg
}
