package ratelimit_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-toolkit/pkg/ratelimit"
)

// Redis address for integration tests
// Can be configured via environment variables or config files
const (
	redisAddr = "localhost:6379" // or use Docker container address
)

// TestRedisSlidingWindowLimiter_Integration tests the limiter with a real Redis instance
func TestRedisSlidingWindowLimiter_Integration(t *testing.T) {
	// Uncomment to skip integration tests
	// t.Skip("Skipping integration test")

	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	defer func() { _ = client.Close() }()

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	require.NoError(t, err, "Failed to connect to Redis")

	// Clean up any existing keys before testing
	client.Del(ctx, "integration-test:limiter:user:123")

	// Create limiter allowing max 5 requests per second
	limiter := ratelimit.NewRedisSlidingWindowLimiter(client, time.Second, 5)
	key := "integration-test:limiter:user:123"

	// First 5 requests should not be limited
	for i := 0; i < 5; i++ {
		limited, err := limiter.Limit(ctx, key)
		require.NoError(t, err)
		assert.False(t, limited, "Request %d should not be rate limited", i+1)
	}

	// 6th request should be rate limited
	limited, err := limiter.Limit(ctx, key)
	require.NoError(t, err)
	assert.True(t, limited, "6th request should be rate limited")

	// Wait for window to slide
	time.Sleep(1100 * time.Millisecond)

	// After window slides, new requests should be allowed
	limited, err = limiter.Limit(ctx, key)
	require.NoError(t, err)
	assert.False(t, limited, "Request after window slide should not be rate limited")
}

// TestRedisSlidingWindowLimiter_Concurrent tests concurrent request handling
func TestRedisSlidingWindowLimiter_Concurrent(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	defer func() { _ = client.Close() }()

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	require.NoError(t, err, "Failed to connect to Redis")

	client.Del(ctx, "integration-test:limiter:concurrent")

	// Create limiter allowing max 10 requests per second
	limiter := ratelimit.NewRedisSlidingWindowLimiter(client, time.Second, 10)
	key := "integration-test:limiter:concurrent"

	// Number of concurrent requests
	concurrency := 20

	var wg sync.WaitGroup
	limitedCount := 0
	var mu sync.Mutex

	// Simulate concurrent requests
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			limited, err := limiter.Limit(ctx, key)
			require.NoError(t, err)

			if limited {
				mu.Lock()
				limitedCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// Should have concurrency - 10 = 10 requests rate limited
	assert.Equal(t, concurrency-10, limitedCount, "Should have correct number of rate limited requests")
}

// TestRedisSlidingWindowLimiter_PreciseSliding tests precise window sliding behavior
func TestRedisSlidingWindowLimiter_PreciseSliding(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	defer func() { _ = client.Close() }()

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	require.NoError(t, err, "Failed to connect to Redis")

	client.Del(ctx, "integration-test:limiter:precise")

	// Create limiter with small window to test precise sliding
	window := 200 * time.Millisecond
	limiter := ratelimit.NewRedisSlidingWindowLimiter(client, window, 2)
	key := "integration-test:limiter:precise"

	// First request
	limited, err := limiter.Limit(ctx, key)
	require.NoError(t, err)
	assert.False(t, limited, "First request should not be rate limited")

	// Second request
	limited, err = limiter.Limit(ctx, key)
	require.NoError(t, err)
	assert.False(t, limited, "Second request should not be rate limited")

	// Third request should be rate limited
	limited, err = limiter.Limit(ctx, key)
	require.NoError(t, err)
	assert.True(t, limited, "Third request should be rate limited")

	// Wait slightly less than window duration to test boundary condition
	time.Sleep(window - 50*time.Millisecond)

	// Window has not completely passed, should still be limited
	limited, err = limiter.Limit(ctx, key)
	require.NoError(t, err)
	assert.True(t, limited, "Request before window passes should be rate limited")

	// Wait longer to ensure window has passed
	time.Sleep(60 * time.Millisecond)

	// After window passes completely, new requests should be allowed
	limited, err = limiter.Limit(ctx, key)
	require.NoError(t, err)
	assert.False(t, limited, "Request after window passes should not be rate limited")
}

// TestRedisSlidingWindowLimiter_TTL tests TTL settings for Redis keys
func TestRedisSlidingWindowLimiter_TTL(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	defer func() { _ = client.Close() }()

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	require.NoError(t, err, "Failed to connect to Redis")

	key := "integration-test:limiter:ttl"
	client.Del(ctx, key)

	// Create limiter with 2 second window
	window := 2 * time.Second
	limiter := ratelimit.NewRedisSlidingWindowLimiter(client, window, 1)

	// Send request to trigger rate limiter
	limited, err := limiter.Limit(ctx, key)
	require.NoError(t, err)
	assert.False(t, limited)

	ttl1, err := client.TTL(ctx, key).Result()
	require.NoError(t, err)
	assert.True(t, ttl1 > 0, "Key should have a positive TTL")

	time.Sleep(500 * time.Millisecond)

	ttl2, err := client.TTL(ctx, key).Result()
	require.NoError(t, err)

	assert.True(t, ttl2 < ttl1, "TTL should decrease over time")
	assert.True(t, ttl2 > 0, "TTL should still be positive")
}
