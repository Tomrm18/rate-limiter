package algos_test

import (
	"testing"
	"time"

	"github.com/Tomrm18/rate-limiter/mocks"
	"github.com/Tomrm18/rate-limiter/pkg/algos"
	"github.com/stretchr/testify/assert"
)

func TestFixedWindowAllow(t *testing.T) {
	mockClock := mocks.NewMockClock()
	window := algos.NewFixedWindowRateLimiter(2, time.Second*5, mockClock)

	// first request
	result, err := window.Allow(TEST_KEY)
	assert.True(t, result.Success)
	assert.Nil(t, err)

	// second request
	result, err = window.Allow(TEST_KEY)
	assert.True(t, result.Success)
	assert.Nil(t, err)

	// third request, rejected
	result, err = window.Allow(TEST_KEY)
	assert.False(t, result.Success)
	assert.Nil(t, err)

	// advance time to allow the window to roll over
	mockClock.Advance(time.Second * 6)

	// first request
	result, err = window.Allow(TEST_KEY)
	assert.True(t, result.Success)
	assert.Nil(t, err)

	// second request
	result, err = window.Allow(TEST_KEY)
	assert.True(t, result.Success)
	assert.Nil(t, err)

	// third request, rejected
	result, err = window.Allow(TEST_KEY)
	assert.False(t, result.Success)
	assert.Nil(t, err)
}

func TestFixedWindowAllowN(t *testing.T) {
	mockClock := mocks.NewMockClock()
	window := algos.NewFixedWindowRateLimiter(2, time.Second*5, mockClock)

	// first request
	result, err := window.AllowN(TEST_KEY, 1)
	assert.True(t, result.Success)
	assert.Nil(t, err)

	// second request
	result, err = window.AllowN(TEST_KEY, 1)
	assert.True(t, result.Success)
	assert.Nil(t, err)

	// third request, rejected
	result, err = window.AllowN(TEST_KEY, 1)
	assert.False(t, result.Success)
	assert.Nil(t, err)

	// advance time to allow the window to roll over
	mockClock.Advance(time.Second * 6)

	// first request
	result, err = window.AllowN(TEST_KEY, 2)
	assert.True(t, result.Success)
	assert.Nil(t, err)

	// second request
	result, err = window.AllowN(TEST_KEY, 1)
	assert.False(t, result.Success)
	assert.Nil(t, err)

	// advance time to allow the window to roll over
	mockClock.Advance(time.Second * 6)

	// first request, rejected
	result, err = window.AllowN(TEST_KEY, 3)
	assert.Nil(t, result)
	assert.Equal(t, algos.ErrNGreaterThanCapacity, err)

	// second request
	result, err = window.AllowN(TEST_KEY, 2)
	assert.True(t, result.Success)
	assert.Nil(t, err)
}
