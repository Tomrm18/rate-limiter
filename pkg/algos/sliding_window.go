package algos

import (
	"errors"
	"math"
	"sync"
	"time"

	"github.com/Tomrm18/rate-limiter/internal/clock"
	"github.com/Tomrm18/rate-limiter/internal/common"
	"github.com/Tomrm18/rate-limiter/pkg/store"
)

type slidingWindowState struct {
	PreviousRequests uint `json:"previous_requests"`
	CurrentRequests  uint `json:"current_requests"`
}

type SlidingWindowRateLimiter struct {
	mu              sync.Mutex
	limit           uint
	windowSize      time.Duration
	windowStartTime time.Time
	clock           clock.Clock
	store           store.Store
}

func NewSlidingWindowRateLimiter(limit uint, windowSize time.Duration, clock clock.Clock, store store.Store) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		limit:           limit,
		windowSize:      windowSize,
		windowStartTime: clock.Now(),
		clock:           clock,
		store:           store,
	}
}

func (s *SlidingWindowRateLimiter) Allow(key string) (*common.Result, error) {
	return s.AllowN(key, 1)
}

func (s *SlidingWindowRateLimiter) AllowN(key string, n uint) (*common.Result, error) {
	if n <= 0 {
		return nil, ErrInvalidN
	}
	if n > s.limit {
		return nil, ErrNGreaterThanCapacity
	}

	allowed := false

	// lock the resources until the function is done, prevents race conditions
	s.mu.Lock()
	defer s.mu.Unlock()

	state, err := s.loadState(key)
	if err != nil {
		return nil, err
	}

	// check if the current window has expired
	elapsed := s.clock.Since(s.windowStartTime)
	expired := elapsed >= s.windowSize

	if expired {
		// check if we have rolled over more than once
		if elapsed >= (2 * s.windowSize) {
			state.PreviousRequests = 0
		} else {
			state.PreviousRequests = state.CurrentRequests
		}
		state.CurrentRequests = 0
		s.windowStartTime = s.windowStartTime.Add(s.windowSize)
	}

	count := s.getCount(state)

	if math.Floor(count)+float64(n) <= float64(s.limit) {
		allowed = true
		state.CurrentRequests += n
	}

	err = s.saveState(key, state)
	if err != nil {
		return nil, err
	}

	return s.buildResult(allowed, uint(float64(s.limit)-math.Floor(count)-float64(n))), nil
}

func (s *SlidingWindowRateLimiter) loadState(key string) (slidingWindowState, error) {
	entry, err := s.store.Get(key)
	if err != nil {
		if errors.Is(err, store.ErrKeyNotFound) {
			return slidingWindowState{
				PreviousRequests: 0,
				CurrentRequests:  0,
			}, nil
		}
		return slidingWindowState{}, err
	}
	return store.Unmarshal[slidingWindowState](entry)
}

func (s *SlidingWindowRateLimiter) saveState(key string, state slidingWindowState) error {
	serialised, err := store.Marshal(state)
	if err != nil {
		return err
	}
	return s.store.Set(key, serialised)
}

func (s *SlidingWindowRateLimiter) getCount(state slidingWindowState) float64 {
	elapsedFraction := float64(s.clock.Since(s.windowStartTime)) / float64(s.windowSize)
	return float64(state.PreviousRequests)*(1-elapsedFraction) + float64(state.CurrentRequests)
}

func (s *SlidingWindowRateLimiter) buildResult(res bool, remaining uint) *common.Result {
	timeUntilNextWindow := s.windowStartTime.Add(s.windowSize).Sub(s.clock.Now())
	return common.NewResult(res, remaining, timeUntilNextWindow, s.limit)
}
