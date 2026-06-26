package algos

import (
	"math"
	"sync"
	"time"

	"github.com/Tomrm18/rate-limiter/internal/clock"
	"github.com/Tomrm18/rate-limiter/internal/common"
)

type SlidingWindowRateLimiter struct {
	mu               sync.Mutex
	limit            uint
	windowSize       time.Duration
	previousRequests int
	windowStartTime  time.Time
	windowRequests   int
	clock            clock.Clock
}

func NewSlidingWindowRateLimiter(limit uint, windowSize time.Duration, clock clock.Clock) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		limit:            limit,
		windowSize:       windowSize,
		previousRequests: 0,
		windowStartTime:  clock.Now(),
		windowRequests:   0,
		clock:            clock,
	}
}

func (s *SlidingWindowRateLimiter) Allow(key string) (*common.Result, error) {
	return s.AllowN(key, 1)
}

func (s *SlidingWindowRateLimiter) AllowN(_ string, n int) (*common.Result, error) {
	// lock the resources until the function is done, prevents race conditions
	s.mu.Lock()
	defer s.mu.Unlock()

	if n <= 0 {
		return nil, ErrInvalidN
	}
	if n > int(s.limit) {
		return nil, ErrNGreaterThanCapacity
	}

	// check if the current window has expired
	elapsed := s.clock.Since(s.windowStartTime)
	expired := elapsed >= s.windowSize

	if expired {
		// check if we have rolled over more than once
		if elapsed >= (2 * s.windowSize) {
			s.previousRequests = 0
		} else {
			s.previousRequests = s.windowRequests
		}
		s.windowRequests = 0
		s.windowStartTime = s.windowStartTime.Add(s.windowSize)
	}

	count := s.getCount()

	if math.Floor(count)+float64(n) <= float64(s.limit) {
		s.windowRequests += n
		return s.buildResult(true, uint(float64(s.limit)-math.Floor(count)-float64(n))), nil
	}
	return s.buildResult(false, 0), nil
}

func (s *SlidingWindowRateLimiter) getCount() float64 {
	elapsedFraction := float64(s.clock.Since(s.windowStartTime)) / float64(s.windowSize)
	return float64(s.previousRequests)*(1-elapsedFraction) + float64(s.windowRequests)
}

func (s *SlidingWindowRateLimiter) buildResult(res bool, remaining uint) *common.Result {
	timeUntilNextWindow := s.windowStartTime.Add(s.windowSize).Sub(s.clock.Now())
	return common.NewResult(res, remaining, timeUntilNextWindow, s.limit)
}
