package mocks

import (
	"time"
)

// mockClock implements the Clock interface
type mockClock struct {
	now time.Time
}

func NewMockClock() *mockClock {
	return &mockClock{
		now: time.Unix(0, 0),
	}
}

func (c *mockClock) Now() time.Time {
	return c.now
}

func (c *mockClock) Since(t time.Time) time.Duration {
	return c.now.Sub(t)
}

func (c *mockClock) Advance(d time.Duration) {
	c.now = c.now.Add(d)
}
