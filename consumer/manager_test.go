package consumer

import (
	"context"
	"redis-task/config"
	"redis-task/redis"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	extRedis "github.com/redis/go-redis/v9"
)

// TestManager completely tests the manager functionallities
func TestManager(t *testing.T) {
	cfg := &config.Config{
		Redis: &config.Redis{
			Address: "127.0.0.1",
			Port:    "6379",
			Channel: "test-channel",
		},
		Consumers: &config.Consumers{
			Count:    2,
			ListName: "list-name",
		},
	}

	// make sure that the redis client setup is ok
	redisClient, err := redis.New(cfg.Redis)
	assert.NoError(t, err)
	defer redisClient.Teardown()

	pubSub := redisClient.Subscribe(context.Background(), cfg.Redis.Channel)

	list := redisClient.List(cfg.Consumers.ListName)

	ctrl := gomock.NewController(t)
	mockProcessor := NewMockprocessor(ctrl)

	manager, err := NewManager(cfg.Consumers, pubSub, list, mockProcessor)
	assert.NoError(t, err)

	assert.Len(t, manager.consumers, cfg.Consumers.Count)

	rawRedis := extRedis.NewClient(
		&extRedis.Options{
			Addr: "127.0.0.1:6379",
		},
	)
	defer rawRedis.Close()

	// ensure that our redis connection is ok
	pingErr := rawRedis.Ping(context.Background()).Err()
	assert.NoError(t, pingErr)

	// make sure that we have the same number of elements in the list
	// as the number of consumers
	llenCmd := rawRedis.LLen(context.Background(), cfg.Consumers.ListName)
	assert.NoError(t, llenCmd.Err())
	assert.EqualValues(t, cfg.Consumers.Count, llenCmd.Val())

	// get the elements in the list
	lrangeCmd := rawRedis.LRange(context.Background(), cfg.Consumers.ListName, 0, llenCmd.Val())
	assert.NoError(t, lrangeCmd.Err())

	// ensure that the consumer IDs are in the list
	for _, consumer := range manager.consumers {
		assert.True(t, in(consumer.id, lrangeCmd.Val()))
	}

	// we are publishing 3 messages
	// and we expect the mock `Process` to be called 3 times
	// this ensures that all messages are consumed only once
	pubCmd := rawRedis.Publish(context.Background(), cfg.Redis.Channel, []byte(`{"message_id": "1"}`))
	assert.NoError(t, pubCmd.Err())
	pubCmd = rawRedis.Publish(context.Background(), cfg.Redis.Channel, []byte(`{"message_id": "2"}`))
	assert.NoError(t, pubCmd.Err())
	pubCmd = rawRedis.Publish(context.Background(), cfg.Redis.Channel, []byte(`{"message_id": "3"}`))
	assert.NoError(t, pubCmd.Err())
	mockProcessor.EXPECT().Process(gomock.Any(), gomock.Any(), gomock.Any()).Times(3).Return(nil)

	// start the consumers
	manager.Start()
	// give them some time to full consume the messages
	// before we close the iterators they are reading from
	time.Sleep(1 * time.Second)
	pubSub.Teardown()
	manager.Teardown()

	// ensure that after teardown
	// we've removed the existing consumers
	llenCmd = rawRedis.LLen(context.Background(), cfg.Consumers.ListName)
	assert.NoError(t, llenCmd.Err())
	assert.EqualValues(t, 0, llenCmd.Val())
}

func in[T comparable](e T, els []T) bool {
	for _, el := range els {
		if e == el {
			return true
		}
	}

	return false
}
