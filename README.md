# Go-Toolkit

[![Go Report Card](https://goreportcard.com/badge/github.com/go-toolkit)](https://goreportcard.com/report/github.com/go-toolkit)
[![GoDoc](https://godoc.org/github.com/go-toolkit?status.svg)](https://pkg.go.dev/github.com/go-toolkit)
[![License](https://img.shields.io/github/license/go-toolkit/go-toolkit.svg)](LICENSE)

A collection of high-performance, production-ready Go utilities for building modern applications.

## Tools

### Rate Limiter

Distributed rate limiting solution using Redis with sliding window algorithm.

```go
// Quick start example
import (
    "github.com/redis/go-redis/v9"
    "github.com/go-toolkit/pkg/ratelimit"
)

// Create a limiter: 10 requests per second
limiter := ratelimit.NewRedisSlidingWindowLimiter(
    redis.NewClient(&redis.Options{Addr: "localhost:6379"}), 
    time.Second, 
    10,
)

// Check if request should be limited
limited, err := limiter.Limit(ctx, "user:123")