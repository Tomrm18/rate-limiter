package algos

import (
	"errors"
	"sync"
	"time"

	"github.com/Tomrm18/rate-limiter/internal/clock"
	"github.com/Tomrm18/rate-limiter/internal/common"
	"github.com/Tomrm18/rate-limiter/internal/util"
	"github.com/Tomrm18/rate-limiter/pkg/store"
)

// private struct for defining the bucket state, this struct is serialsed as a string and used with the store
type bucketState struct {
	Tokens     uint      `json:"tokens"`
	LastRefill time.Time `json:"last_refill"`
}

// the Bucket will contain a number of tokens, that will refresh at a fixed rate
// if the number of tokens is positive, we allow a request
// if it is zero, reject the request
type BucketRateLimiter struct {
	mu sync.Mutex
	// the amount of tokens to add back to the bucket at each refill
	refillAmount uint
	// seconds each refill is executed
	refillSeconds time.Duration
	// clock
	clock clock.Clock
	// the maximum amount of tokens the bucket can hold
	capacity uint
	// the store used for getting, setting and deleting state
	store store.Store
}

func NewBucketRateLimiter(capacity, refillAmount uint, refillSeconds time.Duration, clock clock.Clock, store store.Store) *BucketRateLimiter {
	return &BucketRateLimiter{
		capacity:      capacity,
		refillAmount:  refillAmount,
		refillSeconds: refillSeconds,
		clock:         clock,
		store:         store,
	}
}

func (b *BucketRateLimiter) Allow(key string) (*common.Result, error) {
	return b.AllowN(key, 1)
}

func (b *BucketRateLimiter) AllowN(key string, n uint) (*common.Result, error) {
	allowed := false

	if n == 0 {
		return nil, ErrInvalidN
	}
	if n > b.capacity {
		return nil, ErrNGreaterThanCapacity
	}

	// lock the resources until the function is done, prevents race conditions
	b.mu.Lock()
	defer b.mu.Unlock()

	state, err := b.loadState(key)
	if err != nil {
		return nil, err
	}

	elapsed := b.clock.Since(state.LastRefill)

	if elapsed >= b.refillSeconds {
		b.refillTokens(elapsed, &state)
		state.LastRefill = b.clock.Now()
	}

	if state.Tokens >= n {
		state.Tokens -= n
		allowed = true
	}

	err = b.saveState(key, state)
	if err != nil {
		return nil, err
	}

	return b.buildResult(allowed, state), nil
}

func (b *BucketRateLimiter) loadState(key string) (bucketState, error) {
	entry, err := b.store.Get(key)
	if err != nil {
		if errors.Is(err, store.ErrKeyNotFound) {
			return bucketState{
				Tokens:     b.capacity,
				LastRefill: b.clock.Now(),
			}, nil
		}
		return bucketState{}, err
	}
	return store.Unmarshal[bucketState](entry)
}

func (b *BucketRateLimiter) saveState(key string, state bucketState) error {
	serialised, err := store.Marshal(state)
	if err != nil {
		return err
	}
	return b.store.Set(key, serialised)
}

func (b *BucketRateLimiter) refillTokens(elapsed time.Duration, state *bucketState) {
	refillMult := uint(elapsed / b.refillSeconds)
	state.Tokens += b.refillAmount * refillMult
	state.Tokens = util.Min(state.Tokens, b.capacity)
}

func (b *BucketRateLimiter) buildResult(res bool, state bucketState) *common.Result {
	var timeUntilRefill time.Duration
	if res {
		timeUntilRefill = 0
	} else {
		timeUntilRefill = state.LastRefill.Add(b.refillSeconds).Sub(b.clock.Now())
	}
	return common.NewResult(res, state.Tokens, timeUntilRefill, b.capacity)
}
