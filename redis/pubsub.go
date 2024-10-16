package redis

import (
	"context"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type pubsub struct {
	id     string
	pubsub *redis.PubSub
}

func (c *client) Subscribe(ctx context.Context, channels ...string) *pubsub {
	pubsubUUID := uuid.NewString()
	redisPubSub := c.redis.Subscribe(ctx, channels...)

	log.Info().Msgf("created pubsub with id: %s", pubsubUUID)

	return &pubsub{
		id:     pubsubUUID,
		pubsub: redisPubSub,
	}
}

func (p *pubsub) Messages() func(func(string) bool) {
	return func(yield func(string) bool) {
		for msg := range p.pubsub.Channel() {
			if !yield(msg.Payload) {
				return
			}
		}
	}
}

func (p *pubsub) Teardown() {
	err := p.pubsub.Close()
	if err != nil {
		log.Error().Err(err).Msgf("failed to close pubsub: %s", p.id)
	}
}
