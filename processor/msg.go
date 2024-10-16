package processor

type message struct {
	MessageID string `json:"message_id"`
}

type processedData struct {
	MessageID       string `json:"message_id"`
	ConsumerID      string `json:"consumer_id"`
	ProcessedResult string `json:"processed_result"`
}
