package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/midhat81/fomo-x/services/portfolio-engine/internal/domain"
)

// defaultTTL controls how long cached portfolio/position data stays fresh
// before a reader falls back to Postgres. Short enough to avoid stale P&L,
// long enough to absorb bursts of repeated reads.
const defaultTTL = 10 * time.Second

// Cache wraps a Redis client for caching portfolio and position data.
type Cache struct {
	client *redis.Client
}

// NewCache creates a Redis cache client from an address like "localhost:6379".
func NewCache(addr string) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &Cache{client: client}
}

// SetPortfolio caches a wallet's full portfolio.
func (c *Cache) SetPortfolio(ctx context.Context, portfolio domain.Portfolio) error {
	key := fmt.Sprintf("portfolio:%s", portfolio.Wallet)
	data, err := json.Marshal(portfolio)
	if err != nil {
		return fmt.Errorf("failed to marshal portfolio: %w", err)
	}
	return c.client.Set(ctx, key, data, defaultTTL).Err()
}

// GetPortfolio reads a cached portfolio. Returns (portfolio, true) on hit,
// (zero, false) on miss — caller should fall back to Postgres on miss.
func (c *Cache) GetPortfolio(ctx context.Context, wallet string) (domain.Portfolio, bool) {
	key := fmt.Sprintf("portfolio:%s", wallet)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.Portfolio{}, false
	}

	var portfolio domain.Portfolio
	if err := json.Unmarshal(data, &portfolio); err != nil {
		return domain.Portfolio{}, false
	}
	return portfolio, true
}

// InvalidatePortfolio removes a wallet's cached portfolio, e.g. after a new
// trade is processed and the cached value is now stale.
func (c *Cache) InvalidatePortfolio(ctx context.Context, wallet string) error {
	key := fmt.Sprintf("portfolio:%s", wallet)
	return c.client.Del(ctx, key).Err()
}

// SetPosition caches a single wallet+token position.
func (c *Cache) SetPosition(ctx context.Context, pos domain.Position) error {
	key := fmt.Sprintf("position:%s:%s", pos.Wallet, pos.Token)
	data, err := json.Marshal(pos)
	if err != nil {
		return fmt.Errorf("failed to marshal position: %w", err)
	}
	return c.client.Set(ctx, key, data, defaultTTL).Err()
}

// Close closes the underlying Redis client.
func (c *Cache) Close() error {
	return c.client.Close()
}
