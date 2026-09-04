package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/midhat81/fomo-x/services/trading-engine/internal/domain"
	"github.com/midhat81/fomo-x/services/trading-engine/internal/repository"
)

const testPostgresURL = "postgres://fomox:fomox@127.0.0.1:5433/fomox?sslmode=disable"

func TestIdempotency_DuplicateOrderCreatesOnlyOneRow(t *testing.T) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, testPostgresURL)
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}
	defer pool.Close()

	wallet := "IdemTestWallet1111111111111111111111111111"
	token := "IdemTestToken11111111111111111111111111111"

	_, err = pool.Exec(ctx, `INSERT INTO wallets (address) VALUES ($1) ON CONFLICT (address) DO NOTHING`, wallet)
	if err != nil {
		t.Fatalf("failed to ensure test wallet: %v", err)
	}

	_, err = pool.Exec(ctx, `INSERT INTO tokens (address) VALUES ($1) ON CONFLICT (address) DO NOTHING`, token)
	if err != nil {
		t.Fatalf("failed to ensure test token: %v", err)
	}

	idempotencyKey := "idem-test-" + uuid.NewString()

	repo := repository.NewOrderRepository(pool)

	order := domain.NewOrder(idempotencyKey, wallet, token, domain.OrderSideBuy, 5)

	first, err := repo.Create(ctx, order)
	if err != nil {
		t.Fatalf("expected first order creation to succeed, got error: %v", err)
	}
	if first.ID == "" {
		t.Fatal("expected first order to have a generated ID")
	}

	second, err := repo.Create(ctx, order)
	if err != repository.ErrDuplicateOrder {
		t.Fatalf("expected ErrDuplicateOrder on second submission, got: %v", err)
	}
	if second.ID != first.ID {
		t.Fatalf("expected duplicate submission to return the same order ID: first=%s second=%s", first.ID, second.ID)
	}

	var count int
	err = pool.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE idempotency_key = $1`, idempotencyKey).Scan(&count)
	if err != nil {
		t.Fatalf("failed to count orders: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected exactly 1 order row for key, found %d", count)
	}

	t.Logf("idempotency test passed: %d order row exists for 2 submissions of the same idempotency key", count)
}
