package algos_test

import (
	"math"
	"testing"
	"time"

	"github.com/Tomrm18/rate-limiter/mocks"
	"github.com/Tomrm18/rate-limiter/pkg/algos"

	"github.com/stretchr/testify/assert"
)

const TEST_KEY = "A"

func TestBucketAllow(t *testing.T) {
	mockClock := mocks.NewMockClock()
	bucket := algos.NewBucketRateLimiter(1, 3, 1, time.Second, mockClock)

	// first request, use 1 token
	result, err := bucket.Allow(TEST_KEY)
	assert.True(t, result.Success)
	assert.Nil(t, err)

	// no tokens are left, should reject
	result, err = bucket.Allow(TEST_KEY)
	assert.False(t, result.Success)
	assert.Nil(t, err)

	// advance 1 second, one token should of been added
	mockClock.Advance(time.Second)
	result, err = bucket.Allow(TEST_KEY)
	assert.True(t, result.Success)
	assert.Nil(t, err)

	// no tokens are left, should reject
	result, err = bucket.Allow(TEST_KEY)
	assert.False(t, result.Success)
	assert.Nil(t, err)

	// advance 3 seconds, three tokens should of been added
	mockClock.Advance(time.Second * 3)

	result, err = bucket.Allow(TEST_KEY)
	assert.True(t, result.Success)
	assert.Nil(t, err)

	result, err = bucket.Allow(TEST_KEY)
	assert.True(t, result.Success)
	assert.Nil(t, err)

	result, err = bucket.Allow(TEST_KEY)
	assert.True(t, result.Success)
	assert.Nil(t, err)

	// no tokens are left, should reject
	result, err = bucket.Allow(TEST_KEY)
	assert.False(t, result.Success)
	assert.Nil(t, err)
}

func TestBucketAllowN(t *testing.T) {
	mockClock := mocks.NewMockClock()
	bucket := algos.NewBucketRateLimiter(1, 10, 2, time.Second, mockClock)

	// first request, use 1 token
	result, err := bucket.AllowN(TEST_KEY, 1)
	assert.True(t, result.Success)
	assert.Nil(t, err)

	// no tokens are left, should reject
	result, err = bucket.AllowN(TEST_KEY, 1)
	assert.False(t, result.Success)
	assert.Nil(t, err)

	// advance 5 seconds, 10 tokens should of been added
	mockClock.Advance(time.Second * 5)
	result, err = bucket.AllowN(TEST_KEY, 10)
	assert.True(t, result.Success)
	assert.Nil(t, err)

	// make a request for 0 requests to be allowed, should receive error
	_, err = bucket.AllowN(TEST_KEY, 0)
	assert.Equal(t, algos.ErrInvalidN, err)

	// make a request greater than the capacity of the bucket, should receive error
	_, err = bucket.AllowN(TEST_KEY, math.MaxUint)
	assert.Equal(t, algos.ErrNGreaterThanCapacity, err)
}
