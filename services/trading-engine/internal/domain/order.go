package domain

import "time"

// OrderSide indicates whether an order buys or sells.
type OrderSide string

const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

// OrderType indicates how an order should be executed. Day 3 only
// implements market orders (immediate execution at simulated market
// price) — limit orders are a possible later addition, not required
// by the spec.
type OrderType string

const (
	OrderTypeMarket OrderType = "MARKET"
)

// OrderStatus tracks an order through its lifecycle.
type OrderStatus string

const (
	StatusCreated         OrderStatus = "CREATED"
	StatusRiskCheck       OrderStatus = "RISK_CHECK"
	StatusAccepted        OrderStatus = "ACCEPTED"
	StatusExecuting       OrderStatus = "EXECUTING"
	StatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	StatusFilled          OrderStatus = "FILLED"
	StatusCancelled       OrderStatus = "CANCELLED"
	StatusRejected        OrderStatus = "REJECTED"
)

// Order represents a single paper-trading order.
type Order struct {
	ID             string
	IdempotencyKey string // client-supplied key; duplicate keys must not create duplicate orders
	Wallet         string
	Token          string
	Side           OrderSide
	Type           OrderType
	Quantity       float64
	Status         OrderStatus
	RejectReason   string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewOrder constructs a new order in the CREATED state.
func NewOrder(idempotencyKey, wallet, token string, side OrderSide, quantity float64) Order {
	now := time.Now().UTC()
	return Order{
		IdempotencyKey: idempotencyKey,
		Wallet:         wallet,
		Token:          token,
		Side:           side,
		Type:           OrderTypeMarket,
		Quantity:       quantity,
		Status:         StatusCreated,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// Reject transitions the order to REJECTED with a reason.
func (o *Order) Reject(reason string) {
	o.Status = StatusRejected
	o.RejectReason = reason
	o.UpdatedAt = time.Now().UTC()
}

// Accept transitions the order to ACCEPTED (passed risk check).
func (o *Order) Accept() {
	o.Status = StatusAccepted
	o.UpdatedAt = time.Now().UTC()
}

// Fill transitions the order to FILLED.
func (o *Order) Fill() {
	o.Status = StatusFilled
	o.UpdatedAt = time.Now().UTC()
}
