package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	kafkago "github.com/segmentio/kafka-go"

	"github.com/midhat81/fomo-x/services/solana-ingestor/internal/solana"
)

const TradesTopic = "solana.trades"

// Producer publishes trade events to Kafka.
type Producer struct {
	writer *kafkago.Writer
}

// NewProducer creates a Kafka producer for the given broker list
// (comma-separated, e.g. "localhost:9092").
func NewProducer(brokers string) *Producer {
	brokerList := strings.Split(brokers, ",")

	writer := &kafkago.Writer{
		Addr:                   kafkago.TCP(brokerList...),
		Topic:                  TradesTopic,
		Balancer:               &kafkago.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	return &Producer{writer: writer}
}

// PublishTrade serializes and publishes a TradeEvent to Kafka, keyed by
// event ID for idempotency downstream.
func (p *Producer) PublishTrade(ctx context.Context, evt solana.TradeEvent) error {
	payload, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("failed to marshal trade event: %w", err)
	}

	msg := kafkago.Message{
		Key:   []byte(evt.EventID),
		Value: payload,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to publish trade event: %w", err)
	}

	return nil
}

// Close closes the underlying Kafka writer.
func (p *Producer) Close() error {
	return p.writer.Close()
}