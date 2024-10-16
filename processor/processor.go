package processor

import (
	"context"
	"encoding/json"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
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
	msg := &message{}
	err := json.Unmarshal([]byte(rawMsg), msg)
	if err != nil {
		return err
	}

	start := time.Now()
	data := p.process(consumerID, msg)
	processDuration := time.Since(start)

	jsonData, _ := json.Marshal(data)
	err = p.streamManager.Add(ctx, jsonData)
	if err != nil {
		return err
	}

	_ = processDuration

	return nil
}

func (p *processor) process(consumerID string, message *message) *processedData {
	processingTime := rand.Int64N(500)
	time.Sleep(time.Duration(processingTime+100) * time.Millisecond)

	return &processedData{
		MessageID:       message.MessageID,
		ConsumerID:      consumerID,
		ProcessedResult: uuid.NewString(),
	}
}
