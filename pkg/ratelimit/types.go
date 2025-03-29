package ratelimit

import "context"

type Limiter interface {
	// Limit checks if the given key should be rate limited.
	// Returns true if the request should be limited, false otherwise.
	// An error is returned if the rate limiting operation fails.
	Limit(ctx context.Context, key string) (bool, error)
}
