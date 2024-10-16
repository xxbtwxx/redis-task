package redis

import (
	"context"
	"fmt"
	"redis-task/config"

	"github.com/redis/go-redis/v9"
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

func (c *client) Teardown() error {
	return c.redis.Close()
}
