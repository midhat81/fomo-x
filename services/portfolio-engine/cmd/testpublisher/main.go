package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	kafkago "github.com/segmentio/kafka-go"
)

// TradeEvent mirrors the shared event contract.
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

// This is a throwaway Day 2 testing tool. It publishes a sequence of
// synthetic but realistic trade events for one wallet/token so we can
// verify position math, Postgres writes, and cache invalidation without
// waiting on real instruction-level decoding (which comes later).
func main() {
	writer := &kafkago.Writer{
		Addr:                   kafkago.TCP("127.0.0.1:19092"),
		Topic:                  "solana.trades",
		Balancer:               &kafkago.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
	defer writer.Close()

	wallet := "TestWa11et1111111111111111111111111111111"
	token := "TestToken111111111111111111111111111111111"

	events := []TradeEvent{
		{EventID: uuid.NewString(), Wallet: wallet, Signature: "test-sig-" + uuid.NewString(), Token: token, Side: "BUY", Quantity: 10, Price: 100, Timestamp: time.Now().UTC()},
		{EventID: uuid.NewString(), Wallet: wallet, Signature: "test-sig-" + uuid.NewString(), Token: token, Side: "BUY", Quantity: 5, Price: 120, Timestamp: time.Now().UTC()},
		{EventID: uuid.NewString(), Wallet: wallet, Signature: "test-sig-" + uuid.NewString(), Token: token, Side: "SELL", Quantity: 8, Price: 150, Timestamp: time.Now().UTC()},
	}

	ctx := context.Background()

	for _, evt := range events {
		payload, err := json.Marshal(evt)
		if err != nil {
			log.Fatalf("failed to marshal event: %v", err)
		}

		msg := kafkago.Message{
			Key:   []byte(evt.EventID),
			Value: payload,
		}

		if err := writer.WriteMessages(ctx, msg); err != nil {
			log.Fatalf("failed to publish event: %v", err)
		}

		log.Printf("Published test event: %s (%s %v @ %v)\n", evt.EventID, evt.Side, evt.Quantity, evt.Price)
		time.Sleep(500 * time.Millisecond)
	}

	log.Println("Done publishing test events.")
}
