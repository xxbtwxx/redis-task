package consumer

import (
	"fmt"
)

type messageProvider interface {
	Messages() func(func(string) bool)
}

type consumer struct {
	id              string
	messageProvider messageProvider
}

func newConsumer(id string, msgProvider messageProvider) *consumer {
	return &consumer{
		id:              id,
		messageProvider: msgProvider,
	}
}

func (c *consumer) Consume(doneCallback func(id string) error) {
	go func() {
		for msg := range c.messageProvider.Messages() {
			fmt.Printf("%s consumed %s\n", c.id, msg)
		}

		doneCallback(c.id)
	}()
}
