package main

import (
	"log"

	"github.com/midhat81/fomo-x/services/solana-ingestor/internal/config"
)

func main() {
	cfg := config.Load()

	log.Println("FOMO-X Solana Ingestor starting...")
	log.Printf("Solana RPC URL: %s\n", cfg.SolanaRPCURL)
	log.Printf("Solana WS URL: %s\n", cfg.SolanaWSURL)
	log.Printf("Kafka Brokers: %s\n", cfg.KafkaBrokers)

	// Worker startup will be wired in once solana client, parser,
	// decoder, and kafka producer are implemented.
}