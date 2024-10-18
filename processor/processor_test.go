package processor

import (
	"context"
	"encoding/json"
	"redis-task/config"
	"redis-task/metrics"
	"redis-task/redis"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	extRedis "github.com/redis/go-redis/v9"
)

func TestProcessor_WriteProcessedMsg(t *testing.T) {
	// metrics are a global obj
	// and we need them to be initialized
	// for this test to run
	metrics.Expose()

	redisClient, err := redis.New(&config.Redis{
		Address: "127.0.0.1",
		Port:    "6379",
	})
	assert.NoError(t, err)
	defer redisClient.Teardown()

	rawRedis := extRedis.NewClient(
		&extRedis.Options{
			Addr: "127.0.0.1:6379",
		},
	)
	pingCmd := rawRedis.Ping(context.Background())
	assert.NoError(t, pingCmd.Err())

	defer rawRedis.Close()

	stream := redisClient.Stream("some-stream")

	processor := New(stream)
	err = processor.Process(context.Background(), "consumer_id", `{"message_id":"msg_id"}`)
	assert.NoError(t, err)

	xreadCmd := rawRedis.XRead(context.Background(), &extRedis.XReadArgs{
		Count:   1,
		Streams: []string{"some-stream", "0"},
	})
	assert.NoError(t, xreadCmd.Err())
	require.Len(t, xreadCmd.Val(), 1)

	xreadValue := xreadCmd.Val()[0]
	require.Len(t, xreadValue.Messages, 1)

	msg := xreadValue.Messages[0]
	data, ok := msg.Values["data"]
	require.True(t, ok)

	dataString, ok := data.(string)
	require.True(t, ok)

	processedData := &processedData{}

	err = json.Unmarshal([]byte(dataString), processedData)
	assert.NoError(t, err)

	assert.Equal(t, processedData.ConsumerID, "consumer_id")
	assert.Equal(t, processedData.MessageID, "msg_id")
}
