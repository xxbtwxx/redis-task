package redis

import (
	"context"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type listManager interface {
	LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	LRem(ctx context.Context, key string, count int64, value interface{}) *redis.IntCmd
	// TODO: define the other list command when/if needed
}

// TODO: rework it to store generic types
type list struct {
	id      string
	key     string
	manager listManager
}

func (c *client) List(key string) *list {
	listUUID := uuid.NewString()

	list := &list{
		id:      listUUID,
		key:     key,
		manager: c.redis,
	}

	return list
}

func (l *list) Add(ctx context.Context, value string) error {
	return l.manager.LPush(ctx, l.key, value).Err()
}

func (l *list) Remove(ctx context.Context, value string) error {
	return l.manager.LRem(ctx, l.key, 0, value).Err()
}
