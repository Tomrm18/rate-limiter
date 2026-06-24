package algos

import (
	"rate-limiter/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const TEST_KEY = "A"

func TestBucketAllow(t *testing.T) {
	mockClock := mocks.NewMockClock()
	bucket := NewBucket(1, 1, 1, mockClock)

	// first request, use 1 token
	allowed, err := bucket.Allow(TEST_KEY)
	assert.True(t, allowed)
	assert.Nil(t, err)

	// no tokens are left, should reject
	allowed, err = bucket.Allow(TEST_KEY)
	assert.False(t, allowed)
	assert.Nil(t, err)

	// advance 1 second, one token should of been added
	mockClock.Advance(time.Second)
	allowed, err = bucket.Allow(TEST_KEY)
	assert.True(t, allowed)
	assert.Nil(t, err)

	// no tokens are left, should reject
	allowed, err = bucket.Allow(TEST_KEY)
	assert.False(t, allowed)
	assert.Nil(t, err)

	// advance 3 seconds, three tokens should of been added
	mockClock.Advance(time.Second * 3)

	allowed, err = bucket.Allow(TEST_KEY)
	assert.EqualValues(t, 2, bucket.tokens)
	assert.True(t, allowed)
	assert.Nil(t, err)

	allowed, err = bucket.Allow(TEST_KEY)
	assert.EqualValues(t, 1, bucket.tokens)
	assert.True(t, allowed)
	assert.Nil(t, err)

	allowed, err = bucket.Allow(TEST_KEY)
	assert.EqualValues(t, 0, bucket.tokens)
	assert.True(t, allowed)
	assert.Nil(t, err)

	// no tokens are left, should reject
	allowed, err = bucket.Allow(TEST_KEY)
	assert.False(t, allowed)
	assert.Nil(t, err)
}
