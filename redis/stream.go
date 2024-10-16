package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type streamManager interface {
	XAdd(ctx context.Context, a *redis.XAddArgs) *redis.StringCmd
}

type stream struct {
	name          string
	streamManager streamManager
}

func (c *client) Stream(name string) *stream {
	return &stream{
		name:          name,
		streamManager: c.redis,
	}
}

func (s *stream) Add(ctx context.Context, value interface{}) error {
	return s.streamManager.XAdd(ctx, &redis.XAddArgs{
		Stream: s.name,
		Values: map[string]interface{}{
			"data": value,
		},
	}).Err()
}
