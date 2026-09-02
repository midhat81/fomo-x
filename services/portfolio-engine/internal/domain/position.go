package domain

// Position represents a wallet's current holding in a single token.
type Position struct {
	Wallet       string
	Token        string
	Quantity     float64
	AverageEntry float64
	CurrentPrice float64
}

// MarketValue returns the current market value of the position.
func (p Position) MarketValue() float64 {
	return p.Quantity * p.CurrentPrice
}

// UnrealizedPnL returns the unrealized profit/loss on the position based on
// the difference between current price and average entry price.
func (p Position) UnrealizedPnL() float64 {
	return (p.CurrentPrice - p.AverageEntry) * p.Quantity
}

// UnrealizedPnLPercent returns the unrealized P&L as a percentage of cost basis.
func (p Position) UnrealizedPnLPercent() float64 {
	costBasis := p.AverageEntry * p.Quantity
	if costBasis == 0 {
		return 0
	}
	return (p.UnrealizedPnL() / costBasis) * 100
}

// ApplyBuy updates the position with a new buy, recalculating the
// volume-weighted average entry price.
func (p *Position) ApplyBuy(quantity, price float64) {
	totalCost := (p.AverageEntry * p.Quantity) + (price * quantity)
	p.Quantity += quantity
	if p.Quantity > 0 {
		p.AverageEntry = totalCost / p.Quantity
	}
}

// ApplySell reduces the position's quantity. Average entry price is
// unchanged on a sell — only realized P&L (calculated by the caller,
// see pnl.go) is affected.
func (p *Position) ApplySell(quantity float64) {
	p.Quantity -= quantity
	if p.Quantity < 0 {
		p.Quantity = 0
	}
	if p.Quantity == 0 {
		p.AverageEntry = 0
	}
}
