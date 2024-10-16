//go:generate mockgen -source=$GOFILE -destination=consumer_mock.go -package=consumer -exclude_interfaces=messageProvider

package consumer

import (
	"context"

	"github.com/rs/zerolog/log"
)

type (
	messageProvider interface {
		Messages() func(func(string) bool)
	}

	processor interface {
		Process(context.Context, string, string) error
	}
)

type consumer struct {
	id              string
	messageProvider messageProvider
	processor       processor
}

func newConsumer(
	id string,
	msgProvider messageProvider,
	processor processor,
) *consumer {
	return &consumer{
		id:              id,
		messageProvider: msgProvider,
		processor:       processor,
	}
}

func (c *consumer) Consume(doneCallback func()) {
	go func() {
		log.Debug().Msgf("started consumer: %s", c.id)
		for msg := range c.messageProvider.Messages() {
			err := c.processor.Process(context.Background(), c.id, msg)
			if err != nil {
				log.Error().Err(err).Str("consumer id", c.id).Msg("failed to process message")
			}
		}

		doneCallback()
	}()
}
