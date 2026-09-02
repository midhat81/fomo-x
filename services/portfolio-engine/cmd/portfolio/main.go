package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/midhat81/fomo-x/services/portfolio-engine/internal/cache"
	"github.com/midhat81/fomo-x/services/portfolio-engine/internal/config"
	"github.com/midhat81/fomo-x/services/portfolio-engine/internal/kafka"
	"github.com/midhat81/fomo-x/services/portfolio-engine/internal/repository"
	"github.com/midhat81/fomo-x/services/portfolio-engine/internal/worker"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on system environment variables")
	}

	cfg := config.Load()

	log.Println("FOMO-X Portfolio Engine starting...")
	log.Printf("Postgres URL: %s\n", cfg.PostgresURL)
	log.Printf("Redis URL: %s\n", cfg.RedisURL)
	log.Printf("Kafka Brokers: %s\n", cfg.KafkaBrokers)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("postgres ping failed: %v", err)
	}
	log.Println("Connected to Postgres.")

	c := cache.NewCache(cfg.RedisURL)
	defer c.Close()

	portfolioRepo := repository.NewPortfolioRepository(pool)
	positionRepo := repository.NewPositionRepository(pool)

	consumer := kafka.NewConsumer(cfg.KafkaBrokers, "portfolio-engine")
	defer consumer.Close()

	p := worker.New(consumer, portfolioRepo, positionRepo, c)

	log.Println("Portfolio processor started. Consuming trade events... (Ctrl+C to stop)")

	if err := p.Run(ctx); err != nil {
		log.Fatalf("processor stopped with error: %v", err)
	}

	log.Println("Portfolio engine shut down cleanly.")
}
