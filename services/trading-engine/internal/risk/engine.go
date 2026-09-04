package risk

import "github.com/midhat81/fomo-x/services/trading-engine/internal/domain"

// Engine wraps a fixed set of risk limits and exposes order evaluation.
// Keeping this as its own type (rather than just calling CheckOrder
// directly everywhere) means limits could later be made per-wallet or
// per-tier without changing every call site.
type Engine struct {
	Limits domain.RiskLimits
}

// NewEngine creates a risk engine with the given limits.
func NewEngine(limits domain.RiskLimits) *Engine {
	return &Engine{Limits: limits}
}

// Evaluate runs an order through this engine's risk limits against the
// given account state and market price.
func (e *Engine) Evaluate(order domain.Order, price float64, account AccountState) domain.RiskCheckResult {
	return CheckOrder(order, price, e.Limits, account)
}
