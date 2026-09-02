package config

import "os"

// Config holds all environment-driven settings for the portfolio engine.
type Config struct {
	PostgresURL  string
	RedisURL     string
	KafkaBrokers string
}

// Load reads configuration from environment variables.
func Load() Config {
	return Config{
		PostgresURL:  getEnv("POSTGRES_URL", "postgres://fomox:fomox@localhost:5432/fomox?sslmode=disable"),
		RedisURL:     getEnv("REDIS_URL", "localhost:6379"),
		KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost:19092"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
