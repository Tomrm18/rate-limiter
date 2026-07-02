package algos

import (
	"errors"
	"sync"
	"time"

	"github.com/Tomrm18/rate-limiter/internal/clock"
	"github.com/Tomrm18/rate-limiter/internal/common"
	"github.com/Tomrm18/rate-limiter/pkg/store"
)

type fixedWindowState struct {
	Requests uint `json:"requests"`
}

type FixedWindowRateLimiter struct {
	mu              sync.Mutex
	limit           uint
	windowSize      time.Duration
	windowStartTime time.Time
	clock           clock.Clock
	store           store.Store
}

func NewFixedWindowRateLimiter(limit uint, windowSize time.Duration, clock clock.Clock, store store.Store) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		limit:           limit,
		windowSize:      windowSize,
		windowStartTime: clock.Now(),
		clock:           clock,
		store:           store,
	}
}

func (f *FixedWindowRateLimiter) Allow(key string) (*common.Result, error) {
	return f.AllowN(key, 1)
}

func (f *FixedWindowRateLimiter) AllowN(key string, n uint) (*common.Result, error) {
	if n > f.limit {
		return nil, ErrNGreaterThanCapacity
	}

	allowed := false

	// lock the resources until the function is done, prevents race conditions
	f.mu.Lock()
	defer f.mu.Unlock()

	state, err := f.loadState(key)
	if err != nil {
		return nil, err
	}

	elapsed := f.clock.Since(f.windowStartTime)
	// need to reset to the next window
	if elapsed >= f.windowSize {
		f.windowStartTime = f.windowStartTime.Add(f.windowSize)
		state.Requests = 0
	}

	if state.Requests < f.limit {
		allowed = true
		state.Requests += n
	}

	err = f.saveState(key, state)
	if err != nil {
		return nil, err
	}

	return f.buildResult(allowed, f.limit-state.Requests), nil
}

func (f *FixedWindowRateLimiter) loadState(key string) (fixedWindowState, error) {
	entry, err := f.store.Get(key)
	if err != nil {
		if errors.Is(err, store.ErrKeyNotFound) {
			return fixedWindowState{
				Requests: 0,
			}, nil
		}
		return fixedWindowState{}, err
	}
	return store.Unmarshal[fixedWindowState](entry)
}

func (f *FixedWindowRateLimiter) saveState(key string, state fixedWindowState) error {
	serialised, err := store.Marshal(state)
	if err != nil {
		return err
	}
	return f.store.Set(key, serialised)
}

func (f *FixedWindowRateLimiter) buildResult(res bool, remaining uint) *common.Result {
	return common.NewResult(res, remaining, 0, f.limit)
}
