package common

import "time"

// Result is the result of a call to Allow or AllowN
type Result struct {
	Success         bool
	TokensRemaining uint
	TimeUntilRefill time.Duration
	Limit           uint
}

func NewResult(res bool, remaining uint, refillTime time.Duration, limit uint) *Result {
	return &Result{
		Success:         res,
		TokensRemaining: remaining,
		TimeUntilRefill: refillTime,
		Limit:           limit,
	}
}
