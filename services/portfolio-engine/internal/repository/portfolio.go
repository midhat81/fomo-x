package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/midhat81/fomo-x/services/portfolio-engine/internal/domain"
)

// PortfolioRepository handles wallet-level reads/writes: ensuring wallets
// and tokens exist, and recording transactions/trades.
type PortfolioRepository struct {
	pool *pgxpool.Pool
}

// NewPortfolioRepository creates a new portfolio repository.
func NewPortfolioRepository(pool *pgxpool.Pool) *PortfolioRepository {
	return &PortfolioRepository{pool: pool}
}

// EnsureWallet inserts the wallet if it doesn't already exist, and bumps
// last_active_at if it does.
func (r *PortfolioRepository) EnsureWallet(ctx context.Context, address string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO wallets (address, first_seen_at, last_active_at, created_at)
		VALUES ($1, now(), now(), now())
		ON CONFLICT (address)
		DO UPDATE SET last_active_at = now()
	`, address)

	if err != nil {
		return fmt.Errorf("failed to ensure wallet: %w", err)
	}
	return nil
}

// EnsureToken inserts the token if it doesn't already exist. Symbol/name
// are left blank for now — Day 1 ingestion doesn't resolve token metadata
// yet, only the token address itself.
func (r *PortfolioRepository) EnsureToken(ctx context.Context, address string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO tokens (address, decimals, created_at)
		VALUES ($1, 9, now())
		ON CONFLICT (address) DO NOTHING
	`, address)

	if err != nil {
		return fmt.Errorf("failed to ensure token: %w", err)
	}
	return nil
}

// RecordTransaction inserts a transaction row, returning its generated ID.
// If the signature already exists (duplicate event), it returns the
// existing transaction's ID instead of erroring — this is part of our
// idempotency handling.
func (r *PortfolioRepository) RecordTransaction(ctx context.Context, signature, wallet string, success bool) (string, error) {
	var id string
	err := r.pool.QueryRow(ctx, `
		INSERT INTO transactions (signature, wallet_address, success, processed_at, created_at)
		VALUES ($1, $2, $3, now(), now())
		ON CONFLICT (signature) DO UPDATE SET signature = EXCLUDED.signature
		RETURNING id
	`, signature, wallet, success).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("failed to record transaction: %w", err)
	}
	return id, nil
}

// RecordTrade inserts a trade row keyed by event_id for idempotency. If the
// event_id already exists, this is a no-op (ErrDuplicate is returned so the
// caller can log/skip cleanly instead of double-processing).
var ErrDuplicateTrade = fmt.Errorf("trade event already processed")

func (r *PortfolioRepository) RecordTrade(ctx context.Context, eventID, transactionID, wallet, token, side string, quantity, price float64, tradedAt time.Time, signature string) error {
	tag, err := r.pool.Exec(ctx, `
		INSERT INTO trades (event_id, transaction_id, wallet_address, token_address, side, quantity, price, signature, traded_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now())
		ON CONFLICT (event_id) DO NOTHING
	`, eventID, transactionID, wallet, token, side, quantity, price, signature, tradedAt)

	if err != nil {
		return fmt.Errorf("failed to record trade: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrDuplicateTrade
	}
	return nil
}

// GetPortfolio assembles a full Portfolio for a wallet from its positions.
// Cash is left at zero — this is a paper-trading/on-chain-observation
// system, not a custodial wallet, so "cash" isn't tracked separately.
func (r *PortfolioRepository) GetPortfolio(ctx context.Context, wallet string, positions []domain.Position) domain.Portfolio {
	return domain.Portfolio{
		Wallet:    wallet,
		Cash:      0,
		Positions: positions,
	}
}