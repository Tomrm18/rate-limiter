package clock

import "time"

type Duration = time.Duration

const (
	Nanosecond  Duration = 1
	Microsecond          = 1000 * Nanosecond
	Millisecond          = 1000 * Microsecond
	Second               = 1000 * Millisecond
)

// Clock is an interface that wraps the time functions, used for testing
type Clock interface {
	Now() time.Time
	Since(time.Time) time.Duration
}

// clock implements the Clock interface
type clock struct{}

func New() Clock {
	return &clock{}
}

func (c *clock) Now() time.Time {
	return time.Now()
}

func (c *clock) Since(t time.Time) time.Duration {
	return time.Since(t)
}
