package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-toolkit/pkg/ratelimit"
)

// TestRedisSlidingWindowLimiter_Basic tests basic rate limiting functionality using miniredis
func TestRedisSlidingWindowLimiter_Basic(t *testing.T) {
	// Create a mock Redis server
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	// Connect to mock Redis
	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	defer func() { _ = client.Close() }()

	// Create limiter allowing max 3 requests per second
	limiter := ratelimit.NewRedisSlidingWindowLimiter(client, time.Second, 3)
	ctx := context.Background()
	key := "test-limiter:user:123"

	// First request should not be limited
	limited, err := limiter.Limit(ctx, key)
	require.NoError(t, err)
	assert.False(t, limited, "First request should not be rate limited")

	// Second request should not be limited
	limited, err = limiter.Limit(ctx, key)
	require.NoError(t, err)
	assert.False(t, limited, "Second request should not be rate limited")

	// Third request should not be limited
	limited, err = limiter.Limit(ctx, key)
	require.NoError(t, err)
	assert.False(t, limited, "Third request should not be rate limited")

	// Fourth request should be limited
	limited, err = limiter.Limit(ctx, key)
	require.NoError(t, err)
	assert.True(t, limited, "Fourth request should be rate limited")
}

// TestRedisSlidingWindowLimiter_Sliding tests window sliding behavior
func TestRedisSlidingWindowLimiter_Sliding(t *testing.T) {
	// Create a mock Redis server
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	// Connect to mock Redis
	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	defer func() { _ = client.Close() }()

	// Create limiter allowing max 2 requests per 500ms
	interval := 500 * time.Millisecond
	limiter := ratelimit.NewRedisSlidingWindowLimiter(client, interval, 2)
	ctx := context.Background()
	key := "test-limiter:user:456"

	// First request should not be limited
	limited, err := limiter.Limit(ctx, key)
	require.NoError(t, err)
	assert.False(t, limited, "First request should not be rate limited")

	// Second request should not be limited
	limited, err = limiter.Limit(ctx, key)
	require.NoError(t, err)
	assert.False(t, limited, "Second request should not be rate limited")

	// Third request should be limited
	limited, err = limiter.Limit(ctx, key)
	require.NoError(t, err)
	assert.True(t, limited, "Third request should be rate limited")

	// Wait for window to pass
	time.Sleep(interval + 50*time.Millisecond)

	// After window passes, new requests should be allowed
	limited, err = limiter.Limit(ctx, key)
	require.NoError(t, err)
	assert.False(t, limited, "Request after window passes should not be rate limited")
}

// TestRedisSlidingWindowLimiter_DifferentKeys tests that different keys don't affect each other
func TestRedisSlidingWindowLimiter_DifferentKeys(t *testing.T) {
	// Create a mock Redis server
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	// Connect to mock Redis
	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	defer func() { _ = client.Close() }()

	// Create limiter allowing max 2 requests per second
	limiter := ratelimit.NewRedisSlidingWindowLimiter(client, time.Second, 2)
	ctx := context.Background()
	key1 := "test-limiter:user:789"
	key2 := "test-limiter:user:101112"

	// First key's requests
	limited, err := limiter.Limit(ctx, key1)
	require.NoError(t, err)
	assert.False(t, limited, "First request for key1 should not be rate limited")

	limited, err = limiter.Limit(ctx, key1)
	require.NoError(t, err)
	assert.False(t, limited, "Second request for key1 should not be rate limited")

	limited, err = limiter.Limit(ctx, key1)
	require.NoError(t, err)
	assert.True(t, limited, "Third request for key1 should be rate limited")

	// Second key's requests should not be affected by first key
	limited, err = limiter.Limit(ctx, key2)
	require.NoError(t, err)
	assert.False(t, limited, "First request for key2 should not be rate limited")

	limited, err = limiter.Limit(ctx, key2)
	require.NoError(t, err)
	assert.False(t, limited, "Second request for key2 should not be rate limited")

	limited, err = limiter.Limit(ctx, key2)
	require.NoError(t, err)
	assert.True(t, limited, "Third request for key2 should be rate limited")
}

// TestRedisSlidingWindowLimiter_Error tests error handling
func TestRedisSlidingWindowLimiter_Error(t *testing.T) {
	// Create a mock Redis server
	s, err := miniredis.Run()
	require.NoError(t, err)

	// Connect to mock Redis
	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	// Create limiter
	limiter := ratelimit.NewRedisSlidingWindowLimiter(client, time.Second, 3)
	ctx := context.Background()
	key := "test-limiter:user:error"

	// Normal request
	limited, err := limiter.Limit(ctx, key)
	require.NoError(t, err)
	assert.False(t, limited)

	// Close Redis connection to simulate error
	s.Close()
	_ = client.Close()

	// Should return error when Redis connection is closed
	_, err = limiter.Limit(ctx, key)
	assert.Error(t, err)
}
