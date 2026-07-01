package algos

import (
	"sync"
	"time"

	"github.com/Tomrm18/rate-limiter/internal/clock"
	"github.com/Tomrm18/rate-limiter/internal/common"
)

type FixedWindowRateLimiter struct {
	mu              sync.Mutex
	limit           uint
	windowSize      time.Duration
	windowStartTime time.Time
	windowRequests  uint
	clock           clock.Clock
}

func NewFixedWindowRateLimiter(limit uint, windowSize time.Duration, clock clock.Clock) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		limit:           limit,
		windowSize:      windowSize,
		windowStartTime: clock.Now(),
		windowRequests:  0,
		clock:           clock,
	}
}

func (f *FixedWindowRateLimiter) Allow(key string) (*common.Result, error) {
	return f.AllowN(key, 1)
}

func (f *FixedWindowRateLimiter) AllowN(_ string, n uint) (*common.Result, error) {
	// lock the resources until the function is done, prevents race conditions
	f.mu.Lock()
	defer f.mu.Unlock()

	if n > f.limit {
		return nil, ErrNGreaterThanCapacity
	}

	elapsed := f.clock.Since(f.windowStartTime)
	// need to reset to the next window
	if elapsed >= f.windowSize {
		f.windowStartTime = f.windowStartTime.Add(f.windowSize)
		f.windowRequests = 0
	}

	if f.windowRequests >= f.limit {
		return f.buildResult(false, f.limit-f.windowRequests), nil
	}

	f.windowRequests += n

	return f.buildResult(true, f.limit-f.windowRequests), nil
}

func (f *FixedWindowRateLimiter) buildResult(res bool, remaining uint) *common.Result {
	return common.NewResult(res, remaining, 0, f.limit)
}
