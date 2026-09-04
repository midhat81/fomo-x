package risk

import "github.com/midhat81/fomo-x/services/trading-engine/internal/domain"

// AccountState is the snapshot of a wallet's current trading activity that
// the risk engine needs to evaluate an order against. In a real system
// this would be fetched live from the portfolio engine; for Day 3 it's
// passed in directly by the caller (wired to real data in a later step
// once the trading engine talks to portfolio-engine's tables).
type AccountState struct {
	Wallet               string
	CurrentPositionValue float64 // USD value of existing position in the traded token, before this order
	RealizedLossToday    float64 // positive number representing today's realized losses so far
}

// CheckOrder evaluates an order against the given risk limits and account
// state, returning whether it passes.
func CheckOrder(order domain.Order, price float64, limits domain.RiskLimits, account AccountState) domain.RiskCheckResult {
	orderValue := order.Quantity * price

	if orderValue > limits.MaxOrderValue {
		return domain.RiskFailed("order value exceeds max order value limit")
	}

	if order.Side == domain.OrderSideBuy {
		projectedPositionValue := account.CurrentPositionValue + orderValue
		if projectedPositionValue > limits.MaxPositionValue {
			return domain.RiskFailed("projected position value exceeds max position limit")
		}
	}

	if account.RealizedLossToday >= limits.MaxDailyLoss {
		return domain.RiskFailed("daily loss limit already reached")
	}

	return domain.RiskPassed()
}
