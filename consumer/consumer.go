package consumer

import (
	"context"

	"github.com/rs/zerolog/log"
)

type messageProvider interface {
	Messages() func(func(string) bool)
}

type processor interface {
	Process(context.Context, string, string) error
}

type consumer struct {
	id              string
	messageProvider messageProvider
	processor       processor
}

func newConsumer(id string, msgProvider messageProvider, processor processor) *consumer {
	return &consumer{
		id:              id,
		messageProvider: msgProvider,
		processor:       processor,
	}
}

func (c *consumer) Consume(doneCallback func(id string) error) {
	go func() {
		for msg := range c.messageProvider.Messages() {
			err := c.processor.Process(context.Background(), c.id, msg)
			if err != nil {
				log.Error().Err(err).Msgf("failed to process message")
			}
		}

		err := doneCallback(c.id)
		if err != nil {
			log.Error().Err(err).Msgf("failed to execute consumer %s done callback", c.id)
		}
	}()
}
