package execution

// GetMarketPrice returns a simulated market price for a token. This is a
// placeholder for Day 3 — real market data (from a DEX quote API or price
// oracle) is a later-day concern. For now it returns a fixed demo price
// per token so the risk engine and execution simulator have something
// concrete to work against. In production this would call out to a real
// price feed.
func GetMarketPrice(token string) float64 {
	// Fixed demo price. Every token currently prices at $100 for
	// simplicity — swap this out once a real price source exists.
	return 100.0
}
