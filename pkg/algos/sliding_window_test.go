package algos_test

import (
	"testing"
	"time"

	"github.com/Tomrm18/rate-limiter/mocks"
	"github.com/Tomrm18/rate-limiter/pkg/algos"
	"github.com/stretchr/testify/assert"
)

func TestSlidingWindowAllow(t *testing.T) {
	mockClock := mocks.NewMockClock()
	window := algos.NewSlidingWindowRateLimiter(2, time.Second*5, mockClock)

	// first request
	// 0 / 5 = 0
	// previous weight = 1 - 0, = 1
	// count = 0 * 1 + 0
	// count = 0
	// 0 is under 2, should pass
	result, err := window.Allow(TEST_KEY)
	assert.True(t, result.Success)
	assert.Nil(t, err)

	// second request
	// previous 1 request should partially count
	// 0 / 5 = 0
	// previous weight = 1 - 0, = 1
	// count = 0 * 1 + 1
	// count = 1
	// 1 is under 2, should pass
	result, err = window.Allow(TEST_KEY)
	assert.True(t, result.Success)
	assert.Nil(t, err)

	// third request
	// 0 / 5 = 0
	// previous weight = 1 - 0, = 1
	// count = 0 * 1 + 2
	// count = 2
	// 2 is not under 2, should reject
	result, err = window.Allow(TEST_KEY)
	assert.False(t, result.Success)
	assert.Nil(t, err)

	// advance time to allow the window to roll over
	mockClock.Advance(time.Second * 6)

	// first request in new window
	// 1 / 5 = 0.2
	// previous weight = 1 - 0.2, = 0.8
	// count = 2 * (0.8) + 0
	// count = 1.6
	// 1.6 is under 2, should pass
	result, err = window.Allow(TEST_KEY)
	assert.True(t, result.Success)
	assert.Nil(t, err)

	// second request in new window
	// 1 / 5 = 0.2
	// previous weight = 1 - 0.2, = 0.8
	// count = 2 * (0.8) + 1
	// count = 2.6
	// 2.6 is over 2, should reject
	result, err = window.Allow(TEST_KEY)
	assert.False(t, result.Success)
	assert.Nil(t, err)

	// advance by 2 seconds
	mockClock.Advance(time.Second * 2)

	// third request in new window
	// 3 / 5 = 0.6
	// previous weight = 1 - 0.6, = 0.4
	// count = 2 * (0.4) + 1
	// count = 1.8
	// 1.8 is under 2, should pass
	result, err = window.Allow(TEST_KEY)
	assert.True(t, result.Success)
	assert.Nil(t, err)

	// advance double the window size, previous requests should become 0
	mockClock.Advance(time.Second * 9)
	// 0 / 5 = 0
	// previous weight = 1 - 0, = 1
	// count = 0 * 1 + 0
	// count = 0
	// 0 is under 2, should pass
	result, err = window.Allow(TEST_KEY)
	assert.True(t, result.Success)
	assert.Nil(t, err)
}
