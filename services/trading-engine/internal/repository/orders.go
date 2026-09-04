package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/midhat81/fomo-x/services/trading-engine/internal/domain"
)

// ErrDuplicateOrder is returned when an order with the same idempotency
// key has already been created — the caller should treat this as "no-op,
// return the existing order" rather than an error.
var ErrDuplicateOrder = fmt.Errorf("order with this idempotency key already exists")

// OrderRepository handles persistence of orders to Postgres.
type OrderRepository struct {
	pool *pgxpool.Pool
}

// NewOrderRepository creates a new order repository.
func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

// Create inserts a new order. If an order with the same idempotency_key
// already exists, it returns (existingOrder, ErrDuplicateOrder) instead
// of creating a duplicate — this is the core idempotency guarantee for
// Day 3.
func (r *OrderRepository) Create(ctx context.Context, order domain.Order) (domain.Order, error) {
	var id string
	err := r.pool.QueryRow(ctx, `
		INSERT INTO orders (idempotency_key, wallet_address, token_address, side, order_type, quantity, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`, order.IdempotencyKey, order.Wallet, order.Token, order.Side, order.Type, order.Quantity, order.Status, order.CreatedAt, order.UpdatedAt).Scan(&id)

	if err != nil {
		// Postgres unique_violation error code is 23505.
		if isUniqueViolation(err) {
			existing, getErr := r.GetByIdempotencyKey(ctx, order.IdempotencyKey)
			if getErr != nil {
				return domain.Order{}, fmt.Errorf("failed to fetch existing order after duplicate: %w", getErr)
			}
			return existing, ErrDuplicateOrder
		}
		return domain.Order{}, fmt.Errorf("failed to create order: %w", err)
	}

	order.ID = id
	return order, nil
}

// GetByIdempotencyKey fetches an order by its idempotency key.
func (r *OrderRepository) GetByIdempotencyKey(ctx context.Context, key string) (domain.Order, error) {
	var o domain.Order
	err := r.pool.QueryRow(ctx, `
		SELECT id, idempotency_key, wallet_address, token_address, side, order_type, quantity, status, created_at, updated_at
		FROM orders
		WHERE idempotency_key = $1
	`, key).Scan(&o.ID, &o.IdempotencyKey, &o.Wallet, &o.Token, &o.Side, &o.Type, &o.Quantity, &o.Status, &o.CreatedAt, &o.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Order{}, fmt.Errorf("order not found: %w", err)
		}
		return domain.Order{}, fmt.Errorf("failed to fetch order: %w", err)
	}
	return o, nil
}

// UpdateStatus updates an order's status and optionally its reject reason.
func (r *OrderRepository) UpdateStatus(ctx context.Context, orderID string, status domain.OrderStatus, rejectReason string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE orders
		SET status = $1, reject_reason = NULLIF($2, ''), updated_at = now()
		WHERE id = $3
	`, status, rejectReason, orderID)

	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	return nil
}

// RecordExecution stores the execution result on the order row.
func (r *OrderRepository) RecordExecution(ctx context.Context, orderID string, exec domain.Execution) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE orders
		SET filled_quantity = $1, fill_price = $2, slippage_bps = $3, fee_amount = $4, updated_at = now()
		WHERE id = $5
	`, exec.FilledQty, exec.FillPrice, exec.SlippageBps, exec.FeeAmount, orderID)

	if err != nil {
		return fmt.Errorf("failed to record execution: %w", err)
	}
	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr interface{ SQLState() string }
	if errors.As(err, &pgErr) {
		return pgErr.SQLState() == "23505"
	}
	return false
}
