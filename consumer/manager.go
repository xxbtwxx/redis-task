package consumer

import (
	"context"
	"errors"
	"redis-task/config"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type consumersListManager interface {
	Add(context.Context, string) error
	Remove(context.Context, string) error
}

type consumerManager struct {
	consumers            []*consumer
	wg                   *sync.WaitGroup
	consumersListManager consumersListManager
}

func NewManager(
	cfg *config.Consumers,
	msgProvider messageProvider,
	consumersListManager consumersListManager,
	processor processor,
) (*consumerManager, error) {
	errs := make([]error, 0, cfg.Count)
	consumers := []*consumer{}

	for range cfg.Count {
		consumer := newConsumer(
			uuid.NewString(),
			msgProvider,
			processor,
		)

		err := consumersListManager.Add(context.Background(), consumer.id)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		consumers = append(consumers, consumer)
	}

	if len(errs) != 0 {
		for _, consumer := range consumers {
			err := consumersListManager.Remove(context.Background(), consumer.id)
			if err != nil {
				log.Error().Err(err).Msgf("failed to remove consumer %s", consumer.id)
			}
		}

		return nil, errors.Join(errs...)
	}

	return &consumerManager{
		wg:                   &sync.WaitGroup{},
		consumersListManager: consumersListManager,
		consumers:            consumers,
	}, nil
}

func (cm *consumerManager) Start() {
	for _, consumer := range cm.consumers {
		cm.wg.Add(1)
		consumer.Consume(cm.wg.Done)
	}
}

func (cm *consumerManager) Teardown() {
	cm.wg.Wait()
	for _, consumer := range cm.consumers {
		err := cm.consumersListManager.Remove(context.Background(), consumer.id)
		if err != nil {
			log.Error().Err(err).Msgf("failed to remove consumer %s", consumer.id)
		}
	}
}
