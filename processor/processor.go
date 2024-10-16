package processor

import (
	"context"
	"encoding/json"
	"math/rand/v2"
	"redis-task/metrics"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type streamManager interface {
	Add(ctx context.Context, value interface{}) error
}

type processor struct {
	streamManager streamManager
}

func New(streamManager streamManager) *processor {
	return &processor{
		streamManager: streamManager,
	}
}

func (p *processor) Process(ctx context.Context, consumerID, rawMsg string) error {
	var err error
	var processDuration time.Duration
	defer func() {
		state := "PASS"
		if err != nil {
			state = "FAIL"
		}

		metrics.ObserveProcessingTimes(consumerID, state, processDuration.Seconds())
	}()

	msg := &message{}
	err = json.Unmarshal([]byte(rawMsg), msg)
	if err != nil {
		return err
	}

	start := time.Now()
	data := p.process(consumerID, msg)
	processDuration = time.Since(start)

	jsonData, _ := json.Marshal(data)
	err = p.streamManager.Add(ctx, jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (p *processor) process(consumerID string, message *message) *processedData {
	log.Debug().Str("consumer id", consumerID).Any("message", *message).Msg("processing message")
	processingTime := rand.Int64N(500)
	time.Sleep(time.Duration(processingTime+100) * time.Millisecond)

	return &processedData{
		MessageID:       message.MessageID,
		ConsumerID:      consumerID,
		ProcessedResult: uuid.NewString(),
	}
}
