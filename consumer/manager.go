package consumer

import (
	"redis-task/config"
	"sync"

	"github.com/google/uuid"
)

type consumerManager struct {
	consumers []*consumer
	wg        *sync.WaitGroup
}

func NewManager(cfg *config.Consumers, msgProvider messageProvider) *consumerManager {
	consumers := []*consumer{}
	for range cfg.Count {
		consumers = append(consumers, newConsumer(uuid.NewString(), msgProvider))
	}

	return &consumerManager{
		consumers: consumers,
		wg:        &sync.WaitGroup{},
	}
}

func (cm *consumerManager) Start() {
	for _, consumer := range cm.consumers {
		cm.wg.Add(1)
		consumer.Consume(cm.wg.Done)
	}
}

func (cm *consumerManager) Teardown() {
	cm.wg.Wait()
}
