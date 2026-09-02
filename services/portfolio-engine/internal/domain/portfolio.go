package domain

// Portfolio represents a wallet's full set of holdings.
type Portfolio struct {
	Wallet    string
	Cash      float64
	Positions []Position
}

// TotalPositionsValue sums the market value of all open positions.
func (p Portfolio) TotalPositionsValue() float64 {
	var total float64
	for _, pos := range p.Positions {
		total += pos.MarketValue()
	}
	return total
}

// TotalValue returns cash plus the market value of all positions.
func (p Portfolio) TotalValue() float64 {
	return p.Cash + p.TotalPositionsValue()
}

// TotalUnrealizedPnL sums unrealized P&L across all open positions.
func (p Portfolio) TotalUnrealizedPnL() float64 {
	var total float64
	for _, pos := range p.Positions {
		total += pos.UnrealizedPnL()
	}
	return total
}

// FindPosition returns the position for a given token, and whether it exists.
func (p Portfolio) FindPosition(token string) (Position, bool) {
	for _, pos := range p.Positions {
		if pos.Token == token {
			return pos, true
		}
	}
	return Position{}, false
}

// UpsertPosition replaces an existing position for the same token, or
// appends it if none exists yet.
func (p *Portfolio) UpsertPosition(updated Position) {
	for i, pos := range p.Positions {
		if pos.Token == updated.Token {
			p.Positions[i] = updated
			return
		}
	}
	p.Positions = append(p.Positions, updated)
}
