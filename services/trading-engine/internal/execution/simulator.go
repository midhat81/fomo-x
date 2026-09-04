package execution

import (
	"time"

	"github.com/midhat81/fomo-x/services/trading-engine/internal/domain"
)

// feeBps is the simulated trading fee, in basis points, applied to every
// execution. 30 bps (0.30%) approximates typical Solana DEX swap fees.
const feeBps = 30.0

// partialFillThreshold is the order value (in USD) above which we
// simulate a partial fill instead of a full fill, to exercise the
// PARTIALLY_FILLED status path. This is a deliberate simulation rule for
// Day 3 demo purposes, not a reflection of real market depth.
const partialFillThreshold = 800.0

// Simulate executes an order against a given market price, applying
// slippage and fees, and simulating partial fills for large orders.
func Simulate(order domain.Order, marketPrice float64) domain.Execution {
	orderValue := order.Quantity * marketPrice
	isBuy := order.Side == domain.OrderSideBuy

	slippageBps := CalculateSlippage(orderValue)
	fillPrice := ApplySlippage(marketPrice, slippageBps, isBuy)

	filledQty := order.Quantity
	if orderValue > partialFillThreshold {
		// Simulate an 80% fill on large orders to exercise the
		// PARTIALLY_FILLED status in the demo.
		filledQty = order.Quantity * 0.8
	}

	feeAmount := (filledQty * fillPrice) * (feeBps / 10000.0)

	return domain.Execution{
		OrderID:     order.ID,
		FilledQty:   filledQty,
		FillPrice:   fillPrice,
		SlippageBps: slippageBps,
		FeeAmount:   feeAmount,
		ExecutedAt:  time.Now().UTC(),
	}
}
