package ratelimit

import (
	"context"
	_ "embed"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:embed slide_window.lua
var luaSlideWindow string

// RedisSlidingWindowLimiter implements a sliding window rate limiter using Redis.
type RedisSlidingWindowLimiter struct {
	cmd redis.Cmdable
	// interval represents the size of the time window
	interval time.Duration
	// rate represents the maximum number of allowed requests within the time window
	rate int
}

// NewRedisSlidingWindowLimiter creates a new sliding window rate limiter.
// cmd: Redis client interface
// interval: the time window duration
// rate: maximum number of requests allowed in the time window
func NewRedisSlidingWindowLimiter(cmd redis.Cmdable, interval time.Duration, rate int) Limiter {
	return &RedisSlidingWindowLimiter{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

// Limit implements the Limiter interface.
// It checks if the given key should be rate limited using Redis.
func (r *RedisSlidingWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	// Execute rate limiting logic via Redis Lua script
	return r.cmd.Eval(ctx, luaSlideWindow, []string{key}, r.interval.Milliseconds(), r.rate, time.Now().UnixMilli()).Bool()
}
