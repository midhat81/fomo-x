package domain

import "time"

// TradeSide mirrors the side of a trade event.
type TradeSide string

const (
	SideBuy  TradeSide = "BUY"
	SideSell TradeSide = "SELL"
)

// RealizedPnLEntry records the profit/loss realized from a single sell.
type RealizedPnLEntry struct {
	Wallet    string
	Token     string
	Quantity  float64
	EntryCost float64
	ExitValue float64
	PnL       float64
	Timestamp time.Time
}

// CalculateRealizedPnL computes the realized P&L for selling `quantity`
// units out of a position with the given average entry price, at the
// given exit price.
func CalculateRealizedPnL(wallet, token string, quantity, averageEntry, exitPrice float64) RealizedPnLEntry {
	entryCost := averageEntry * quantity
	exitValue := exitPrice * quantity

	return RealizedPnLEntry{
		Wallet:    wallet,
		Token:     token,
		Quantity:  quantity,
		EntryCost: entryCost,
		ExitValue: exitValue,
		PnL:       exitValue - entryCost,
		Timestamp: time.Now().UTC(),
	}
}

// DailyPnLSummary aggregates realized and unrealized P&L for a wallet
// over a day.
type DailyPnLSummary struct {
	Wallet        string
	Date          time.Time
	RealizedPnL   float64
	UnrealizedPnL float64
	TotalPnL      float64
}

// NewDailyPnLSummary builds a summary from realized and unrealized totals.
func NewDailyPnLSummary(wallet string, date time.Time, realized, unrealized float64) DailyPnLSummary {
	return DailyPnLSummary{
		Wallet:        wallet,
		Date:          date,
		RealizedPnL:   realized,
		UnrealizedPnL: unrealized,
		TotalPnL:      realized + unrealized,
	}
}
