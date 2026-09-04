package execution

// CalculateSlippage returns a simulated slippage amount, in basis points,
// based on order size relative to a fixed "liquidity" assumption. This is
// a simplified model for Day 3 — larger orders slip more, capped at a
// reasonable maximum. Real order-book depth data isn't available in this
// paper-trading system, so this approximates the effect.
func CalculateSlippage(orderValue float64) float64 {
	const baseSlippageBps = 5.0       // 0.05% baseline slippage for any trade
	const sizeImpactDivisor = 10000.0 // every $10,000 of order value adds ~1bps
	const maxSlippageBps = 200.0      // cap at 2% to avoid unrealistic blowups

	slippage := baseSlippageBps + (orderValue / sizeImpactDivisor)
	if slippage > maxSlippageBps {
		return maxSlippageBps
	}
	return slippage
}

// ApplySlippage adjusts a market price by the given slippage (in bps),
// worsening the price for the trader: buys pay more, sells receive less.
func ApplySlippage(price float64, slippageBps float64, isBuy bool) float64 {
	slippageFactor := slippageBps / 10000.0
	if isBuy {
		return price * (1 + slippageFactor)
	}
	return price * (1 - slippageFactor)
}
