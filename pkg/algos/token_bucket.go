package algos

import (
	"sync"
	"time"

	"github.com/Tomrm18/rate-limiter/internal/clock"
	"github.com/Tomrm18/rate-limiter/internal/common"
	"github.com/Tomrm18/rate-limiter/internal/util"
)

// the Bucket will contain a number of tokens, that will refresh at a fixed rate
// if the number of tokens is positive, we allow a request
// if it is zero, reject the request
type BucketRateLimiter struct {
	mu sync.Mutex
	// the amount of tokens in the bucket
	tokens uint
	// the amount of tokens to add back to the bucket at each refill
	refillAmount uint
	// seconds each refill is executed
	refillSeconds time.Duration
	// time the last refill occured
	lastRefill time.Time
	// clock
	clock clock.Clock
	// the maximum amount of tokens the bucket can hold
	capacity uint
}

func NewBucketRateLimiter(tokens, capacity, refillAmount uint, refillSeconds time.Duration, clock clock.Clock) *BucketRateLimiter {
	return &BucketRateLimiter{
		tokens:        tokens,
		capacity:      capacity,
		refillAmount:  refillAmount,
		refillSeconds: refillSeconds,
		clock:         clock,
		lastRefill:    clock.Now(),
	}
}

func (b *BucketRateLimiter) Allow(key string) (*common.Result, error) {
	return b.AllowN(key, 1)
}

func (b *BucketRateLimiter) AllowN(_ string, n uint) (*common.Result, error) {
	// lock the resources until the function is done, prevents race conditions
	b.mu.Lock()
	defer b.mu.Unlock()

	if n == 0 {
		return nil, ErrInvalidN
	}
	if n > b.capacity {
		return nil, ErrNGreaterThanCapacity
	}

	elapsed := b.clock.Since(b.lastRefill)

	if elapsed >= b.refillSeconds {
		b.refillTokens(elapsed)
		b.lastRefill = b.clock.Now()
	}

	if b.tokens >= n {
		b.tokens -= n
		return b.buildResult(true), nil
	}

	return b.buildResult(false), nil
}

func (b *BucketRateLimiter) refillTokens(elapsed time.Duration) {
	refillMult := uint(elapsed / b.refillSeconds)
	b.tokens += b.refillAmount * refillMult
	b.tokens = util.Min(b.tokens, b.capacity)
}

func (b *BucketRateLimiter) buildResult(res bool) *common.Result {
	var timeUntilRefill time.Duration
	if res {
		timeUntilRefill = 0
	} else {
		timeUntilRefill = b.lastRefill.Add(b.refillSeconds).Sub(b.clock.Now())
	}
	return common.NewResult(res, b.tokens, timeUntilRefill, b.capacity)
}
