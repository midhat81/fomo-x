package config

import "os"

// Config holds all environment-driven settings for the trading engine.
type Config struct {
	PostgresURL  string
	KafkaBrokers string
}

// Load reads configuration from environment variables.
func Load() Config {
	return Config{
		PostgresURL:  getEnv("POSTGRES_URL", "postgres://fomox:fomox@127.0.0.1:5433/fomox?sslmode=disable"),
		KafkaBrokers: getEnv("KAFKA_BROKERS", "127.0.0.1:19092"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
