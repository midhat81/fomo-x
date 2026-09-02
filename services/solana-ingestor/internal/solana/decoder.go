package solana

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// TradeEvent matches the packages/contracts/events/trade-event.json contract.
type TradeEvent struct {
	EventID   string    `json:"event_id"`
	Wallet    string    `json:"wallet"`
	Signature string    `json:"signature"`
	Token     string    `json:"token"`
	Side      string    `json:"side"`
	Quantity  float64   `json:"quantity"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

// swapKeywords are log substrings that suggest a swap/trade occurred.
// This is a heuristic for Day 1 — real amount/token decoding comes later
// once we parse instruction data directly instead of log text.
var swapKeywords = []string{"Swap", "swap", "Instruction: Swap"}

// DecodeTrade attempts to build a TradeEvent from a parsed transaction.
// Returns (event, true) if the transaction looks like a trade, (zero, false)
// otherwise, so callers can skip irrelevant transactions.
func DecodeTrade(tx ParsedTransaction) (TradeEvent, bool) {
	if !tx.Success {
		return TradeEvent{}, false
	}

	if !containsSwap(tx.RawLogs) {
		return TradeEvent{}, false
	}

	return TradeEvent{
		EventID:   uuid.NewString(),
		Signature: tx.Signature,
		Side:      "UNKNOWN", // refined later once instruction decoding is added
		Timestamp: time.Now().UTC(),
	}, true
}

func containsSwap(logs []string) bool {
	for _, line := range logs {
		for _, kw := range swapKeywords {
			if strings.Contains(line, kw) {
				return true
			}
		}
	}
	return false
}