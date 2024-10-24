package main

import (
	"context"
	"os"
	"os/signal"
	"redis-task/config"
	"redis-task/consumer"
	"redis-task/metrics"
	"redis-task/processor"
	"redis-task/redis"
	"syscall"

	_ "net/http/pprof"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	shutdownSig := make(chan os.Signal, 1)
	signal.Notify(shutdownSig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	wait := make(chan struct{}, 1)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-shutdownSig
		cancel()
		wait <- struct{}{}
	}()

	metricsSrv := metrics.New()
	metricsSrv.Start()

	teardowns := []func(){
		metricsSrv.Stop,
	}
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

	logLevel, err := zerolog.ParseLevel(cfg.Log.Level)
	if err != nil {
		log.Error().Err(err).Msgf("failed to parse log level: %s, defaulting to `debug`", cfg.Log.Level)
		logLevel = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	redisClient, err := redis.New(cfg.Redis)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to redis")
		shutdownSig <- syscall.SIGABRT
	}

	pubsub := redisClient.Subscribe(ctx, cfg.Redis.Channel)
	teardowns = append(teardowns, pubsub.Teardown)

	processedMsgStream := redisClient.Stream(cfg.Processor.ProcessedEventsStream)
	processor := processor.New(processedMsgStream)

	consumerListManager := redisClient.List(cfg.Consumers.ListName)
	consumerManager, err := consumer.NewManager(cfg.Consumers, pubsub, consumerListManager, processor)
	if err != nil {
		log.Error().Err(err).Msg("consumer manager initialization failed")
		shutdownSig <- syscall.SIGABRT
	}

	consumerManager.Start()
	teardowns = append(teardowns, consumerManager.Teardown, redisClient.Teardown)

	<-wait
}
