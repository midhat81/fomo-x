package domain

import "time"

// Execution represents the result of simulating an order against the
// market: the price it filled at (after slippage), any fees charged, and
// how much of the order quantity was actually filled.
type Execution struct {
	OrderID     string
	FilledQty   float64
	FillPrice   float64
	SlippageBps float64 // slippage in basis points, e.g. 50 = 0.50%
	FeeAmount   float64
	ExecutedAt  time.Time
}

// TotalCost returns the total cost (or proceeds, for a sell) of this
// execution including fees.
func (e Execution) TotalCost() float64 {
	return (e.FilledQty * e.FillPrice) + e.FeeAmount
}

// IsFullyFilled reports whether the execution filled the entire requested
// quantity (vs. a partial fill).
func (e Execution) IsFullyFilled(requestedQty float64) bool {
	return e.FilledQty >= requestedQty
}
