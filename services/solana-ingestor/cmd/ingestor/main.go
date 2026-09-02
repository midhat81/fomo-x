package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/midhat81/fomo-x/services/solana-ingestor/internal/checkpoint"
	"github.com/midhat81/fomo-x/services/solana-ingestor/internal/config"
	"github.com/midhat81/fomo-x/services/solana-ingestor/internal/kafka"
	"github.com/midhat81/fomo-x/services/solana-ingestor/internal/solana"
	"github.com/midhat81/fomo-x/services/solana-ingestor/internal/worker"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on system environment variables")
	}

	cfg := config.Load()

	log.Println("FOMO-X Solana Ingestor starting...")
	log.Printf("Solana WS URL: %s\n", cfg.SolanaWSURL)
	log.Printf("Kafka Brokers: %s\n", cfg.KafkaBrokers)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	wsClient, err := solana.NewWSClient(ctx, cfg.SolanaWSURL)
	if err != nil {
		log.Fatalf("failed to connect to solana websocket: %v", err)
	}
	defer wsClient.Close()

	producer := kafka.NewProducer(cfg.KafkaBrokers)
	defer producer.Close()

	cp := checkpoint.NewStore("checkpoint.json")

	w := worker.New(wsClient, producer, cp)

	log.Println("Worker started. Listening for trade activity... (Ctrl+C to stop)")

	if err := w.Run(ctx); err != nil {
		log.Fatalf("worker stopped with error: %v", err)
	}

	log.Println("Ingestor shut down cleanly.")
}