package redis

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type pubsub struct {
	ctx    context.Context
	id     string
	pubsub *redis.PubSub
}

func (c *client) Subscribe(ctx context.Context, channels ...string) *pubsub {
	redisPubSub := c.redis.Subscribe(ctx, channels...)
	log.Debug().Msgf("created pubsub for channels: %s", strings.Join(channels, ", "))

	return &pubsub{
		ctx:    ctx,
		pubsub: redisPubSub,
	}
}

func (p *pubsub) Messages() func(func(string) bool) {
	return func(yield func(string) bool) {
		for {
			select {
			case <-p.ctx.Done():
				return
			case msg, ok := <-p.pubsub.Channel():
				if !ok {
					return
				}

				if !yield(msg.Payload) {
					return
				}
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
