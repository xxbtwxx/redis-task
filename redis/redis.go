package redis

import (
	"context"
	"fmt"
	"redis-task/config"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type client struct {
	redis *redis.Client
}

func New(cfg *config.Redis) (*client, error) {
	redisClient := redis.NewClient(
		&redis.Options{
			Addr: fmt.Sprintf("%s:%s", cfg.Address, cfg.Port),
		},
	)

	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return &client{
		redis: redisClient,
	}, nil
}

func (c *client) Teardown() {
	err := c.redis.Close()
	if err != nil {
		log.Error().Err(err).Msg("failed to close redis connection")
	}
}
