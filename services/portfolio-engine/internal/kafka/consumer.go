package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

const TradesTopic = "solana.trades"

// TradeEvent mirrors the event contract published by the solana-ingestor
// service (packages/contracts/events/trade-event.json).
type TradeEvent struct {
	EventID   string    `json:"event_id"`
	Wallet    string    `json:"wallet"`
	Signature string    `json:"signature"`
	Token     string    `json:"token"`
	Side      string    `json:"side"`
	Quantity  float64   `json:"quantity"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

// Consumer reads trade events from Kafka.
type Consumer struct {
	reader *kafkago.Reader
}

// NewConsumer creates a Kafka consumer for the given broker list
// (comma-separated) and consumer group ID.
func NewConsumer(brokers, groupID string) *Consumer {
	brokerList := strings.Split(brokers, ",")

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: brokerList,
		Topic:   TradesTopic,
		GroupID: groupID,
	})

	return &Consumer{reader: reader}
}

// ReadTrade blocks until the next trade event is available, or ctx is
// cancelled. Returns the decoded event and whether the read succeeded.
func (c *Consumer) ReadTrade(ctx context.Context) (TradeEvent, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return TradeEvent{}, fmt.Errorf("failed to read message: %w", err)
	}

	var evt TradeEvent
	if err := json.Unmarshal(msg.Value, &evt); err != nil {
		return TradeEvent{}, fmt.Errorf("failed to unmarshal trade event: %w", err)
	}

	return evt, nil
}

// Close closes the underlying Kafka reader.
func (c *Consumer) Close() error {
	return c.reader.Close()
}
