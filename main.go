package main

import (
	"context"
	"os"
	"os/signal"
	"redis-task/config"
	"redis-task/consumer"
	"redis-task/redis"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
	shutdownSig := make(chan os.Signal, 1)
	signal.Notify(shutdownSig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP)

	teardowns := []func(){}
	defer func() {
		for _, teardown := range teardowns {
			if teardown != nil {
				teardown()
			}
		}
	}()

	cfg, err := config.Load("./config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	redisClient, err := redis.New(cfg.Redis)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to redis")
	}
	teardowns = append(teardowns, redisClient.Teardown)

	pubsub := redisClient.Subscribe(context.Background(), cfg.Redis.Channel)

	consumerManager := consumer.NewManager(cfg.Consumers, pubsub)
	consumerManager.Start()
	teardowns = append(teardowns, consumerManager.Teardown)

	<-shutdownSig
}
