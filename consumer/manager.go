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
) (*consumerManager, error) {
	consumers := []*consumer{}

	errs := make([]error, 0, cfg.Count)
	for range cfg.Count {
		consumer := newConsumer(uuid.NewString(), msgProvider)
		err := consumersListManager.Add(context.Background(), consumer.id)
		if err != nil {
			errs = append(errs, err)
		}

		consumers = append(consumers, consumer)

		log.Debug().Msgf("started consumer: %s", consumer.id)
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
		consumers:            consumers,
		wg:                   &sync.WaitGroup{},
		consumersListManager: consumersListManager,
	}, nil
}

func (cm *consumerManager) Start() {
	for _, consumer := range cm.consumers {
		cm.wg.Add(1)
		consumer.Consume(cm.consumerDoneCallback)
	}
}

func (cm *consumerManager) consumerDoneCallback(id string) error {
	defer log.Debug().Msgf("removed consumer: %s", id)
	defer cm.wg.Done()
	return cm.consumersListManager.Remove(context.Background(), id)
}

func (cm *consumerManager) Teardown() {
	cm.wg.Wait()
}
