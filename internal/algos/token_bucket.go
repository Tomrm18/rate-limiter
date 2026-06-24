package algos

import (
	"rate-limiter/internal/clock"
	"time"
)

// the Bucket will contain a number of tokens, that will refresh at a fixed rate
// if the number of tokens is positive, we allow a request
// if it is zero, reject the request
type Bucket struct {
	// the amount of tokens in the bucket
	tokens uint
	// the amount of tokens to add back to the bucket at each refill
	refillAmount uint
	// seconds each refill is executed
	refillSeconds uint
	// time the last refill occured
	lastRefill time.Time
	// clock
	clock clock.Clock
}

func NewBucket(tokens, refillAmount, refillSeconds uint, clock clock.Clock) *Bucket {
	return &Bucket{
		tokens:        tokens,
		refillAmount:  refillAmount,
		refillSeconds: refillSeconds,
		clock:         clock,
		lastRefill:    clock.Now(),
	}
}

func (b *Bucket) Allow(key string) (bool, error) {
	// update the amount of tokens in the bucket based on the amount of time passed
	// this is known as a 'lazy refill', and allows us to avoid manually running a loop to refill the bucket
	elapsed := b.clock.Since(b.lastRefill)

	if elapsed >= (clock.Duration(b.refillSeconds) * clock.Second) {
		refillMult := uint(elapsed / (time.Duration(b.refillSeconds) * time.Second))
		b.tokens += b.refillAmount * refillMult
		b.lastRefill = b.clock.Now()
	}

	if b.tokens == 0 {
		return false, nil
	}
	b.tokens -= 1
	return true, nil
}
