package main

import (
	"context"
	"os"
	"os/signal"
	"redis-task/config"
	"redis-task/redis"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
	shutdownSig := make(chan os.Signal, 1)
	signal.Notify(shutdownSig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP)

	cfg, err := config.Load("./config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	redisClient, err := redis.New(cfg.Redis)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to redis")
	}
	defer redisClient.Teardown()

	pubsub := redisClient.Subscribe(context.Background(), cfg.Redis.Channel)

	_ = pubsub

	<-shutdownSig
}
