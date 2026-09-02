package config

import "os"

// Config holds all environment-driven settings for the Solana ingestor service.
type Config struct {
	SolanaRPCURL string
	SolanaWSURL  string
	KafkaBrokers string
}

// Load reads configuration from environment variables.
func Load() Config {
	return Config{
		SolanaRPCURL: getEnv("SOLANA_RPC_URL", ""),
		SolanaWSURL:  getEnv("SOLANA_WS_URL", ""),
		KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost:9092"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}