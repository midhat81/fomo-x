package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/midhat81/fomo-x/services/portfolio-engine/internal/domain"
)

// PositionRepository handles persistence of positions to Postgres.
type PositionRepository struct {
	pool *pgxpool.Pool
}

// NewPositionRepository creates a new position repository.
func NewPositionRepository(pool *pgxpool.Pool) *PositionRepository {
	return &PositionRepository{pool: pool}
}

// Upsert inserts or updates a wallet's position in a given token, keyed by
// (wallet_address, token_address). This mirrors the domain layer's
// UpsertPosition behavior but persists it.
func (r *PositionRepository) Upsert(ctx context.Context, pos domain.Position) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO positions (wallet_address, token_address, quantity, average_entry, updated_at)
		VALUES ($1, $2, $3, $4, now())
		ON CONFLICT (wallet_address, token_address)
		DO UPDATE SET quantity = $3, average_entry = $4, updated_at = now()
	`, pos.Wallet, pos.Token, pos.Quantity, pos.AverageEntry)

	if err != nil {
		return fmt.Errorf("failed to upsert position: %w", err)
	}
	return nil
}

// GetByWallet returns all positions held by a wallet. CurrentPrice is left
// at zero here — it's populated separately from live market data by the
// caller before computing unrealized P&L.
func (r *PositionRepository) GetByWallet(ctx context.Context, wallet string) ([]domain.Position, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT wallet_address, token_address, quantity, average_entry
		FROM positions
		WHERE wallet_address = $1 AND quantity > 0
	`, wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to query positions: %w", err)
	}
	defer rows.Close()

	var positions []domain.Position
	for rows.Next() {
		var p domain.Position
		if err := rows.Scan(&p.Wallet, &p.Token, &p.Quantity, &p.AverageEntry); err != nil {
			return nil, fmt.Errorf("failed to scan position: %w", err)
		}
		positions = append(positions, p)
	}
	return positions, rows.Err()
}

// GetOne returns a single wallet+token position, and whether it exists.
func (r *PositionRepository) GetOne(ctx context.Context, wallet, token string) (domain.Position, bool, error) {
	var p domain.Position
	err := r.pool.QueryRow(ctx, `
		SELECT wallet_address, token_address, quantity, average_entry
		FROM positions
		WHERE wallet_address = $1 AND token_address = $2
	`, wallet, token).Scan(&p.Wallet, &p.Token, &p.Quantity, &p.AverageEntry)

	if err != nil {
		return domain.Position{}, false, nil
	}
	return p, true, nil
}