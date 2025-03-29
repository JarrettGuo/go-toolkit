```markdown
# Go-Toolkit

[![Go Report Card](https://goreportcard.com/badge/github.com/go-toolkit)](https://goreportcard.com/report/github.com/go-toolkit)
[![GoDoc](https://godoc.org/github.com/go-toolkit?status.svg)](https://pkg.go.dev/github.com/go-toolkit)
[![License](https://img.shields.io/github/license/go-toolkit/go-toolkit.svg)](LICENSE)

Go-Toolkit is a collection of high-performance, production-ready Go tools for building modern applications. Each tool is designed to be used independently or together with other toolkit components.

## Available Tools

### Rate Limiter

A high-performance, scalable distributed rate limiting library based on Redis, supporting sliding window algorithm, suitable for microservices and API gateway scenarios.

[View detailed documentation](pkg/ratelimit/README_EN.md)

## Development

### Requirements

- Go 1.18+
- Redis (for rate limiter)
- Docker & Docker Compose (for integration testing)

### Build and Test

```bash
# Run all checks and tests
make all

# Build only
make build

# Run tests
make test

# Run integration tests (requires Redis)
make integration-test

# Start development environment
make docker-up

#### Quick Start

```go
package main

import (
	"log"
	"time"
	"context"
	
	"github.com/redis/go-redis/v9"
	"github.com/go-toolkit/pkg/ratelimit"
)

func main() {
	// Create Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	
	// Create rate limiter: maximum 10 requests per second
	limiter := ratelimit.NewRedisSlidingWindowLimiter(redisClient, time.Second, 10)
	
	ctx := context.Background()
	
	// Use the rate limiter
	limited, err := limiter.Limit(ctx, "user:123")
	if err != nil {
		log.Fatalf("Rate limiter error: %v", err)
	}
	
	if limited {
		log.Println("Request rate limited")
	} else {
		log.Println("Request allowed")
	}
}