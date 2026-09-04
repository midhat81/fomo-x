package domain

// RiskCheckResult is the outcome of running an order through the risk
// engine before it's allowed to execute.
type RiskCheckResult struct {
	Passed bool
	Reason string // populated when Passed is false
}

// Passed is a convenience constructor for a successful risk check.
func RiskPassed() RiskCheckResult {
	return RiskCheckResult{Passed: true}
}

// RiskFailed is a convenience constructor for a failed risk check.
func RiskFailed(reason string) RiskCheckResult {
	return RiskCheckResult{Passed: false, Reason: reason}
}

// RiskLimits defines the configurable thresholds the risk engine enforces
// for a given account. These mirror the example numbers from the spec
// (max trade $1,000, max position $2,000, max daily loss $500 on a
// $10,000 account) but are expressed generically so they can be tuned.
type RiskLimits struct {
	MaxOrderValue    float64 // max USD value of a single order
	MaxPositionValue float64 // max USD value of a single position after this order
	MaxDailyLoss     float64 // max USD realized loss allowed in a day before new orders are blocked
}

// DefaultRiskLimits returns the example limits from the spec, scaled to a
// $10,000 paper account.
func DefaultRiskLimits() RiskLimits {
	return RiskLimits{
		MaxOrderValue:    1000,
		MaxPositionValue: 2000,
		MaxDailyLoss:     500,
	}
}
