package algos

import (
	"math"
	"rate-limiter/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const TEST_KEY = "A"

func TestBucketAllow(t *testing.T) {
	mockClock := mocks.NewMockClock()
	bucket := NewBucket(10, 1, time.Second, mockClock)

	// for the purposes of this test, manually reduce tokens to 1
	bucket.tokens = 1

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
	assert.EqualValues(t, 2, bucket.tokens)
	assert.True(t, result.Success)
	assert.Nil(t, err)

	result, err = bucket.Allow(TEST_KEY)
	assert.EqualValues(t, 1, bucket.tokens)
	assert.True(t, result.Success)
	assert.Nil(t, err)

	result, err = bucket.Allow(TEST_KEY)
	assert.EqualValues(t, 0, bucket.tokens)
	assert.True(t, result.Success)
	assert.Nil(t, err)

	// no tokens are left, should reject
	result, err = bucket.Allow(TEST_KEY)
	assert.False(t, result.Success)
	assert.Nil(t, err)
}

func TestBucketAllowN(t *testing.T) {
	mockClock := mocks.NewMockClock()
	bucket := NewBucket(10, 2, time.Second, mockClock)

	// for the purposes of this test, manually reduce tokens to 1
	bucket.tokens = 1

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
	assert.EqualValues(t, 0, bucket.tokens)

	// make a request for 0 requests to be allowed, should receive error
	_, err = bucket.AllowN(TEST_KEY, 0)
	assert.Equal(t, ErrInvalidN, err)

	// make a request greater than the capacity of the bucket, should receive error
	_, err = bucket.AllowN(TEST_KEY, math.MaxUint)
	assert.Equal(t, ErrNGreaterThanCapacity, err)
}
