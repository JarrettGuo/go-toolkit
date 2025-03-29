# Go-Toolkit

[![Go Report Card](https://goreportcard.com/badge/github.com/go-toolkit)](https://goreportcard.com/report/github.com/go-toolkit)
[![GoDoc](https://godoc.org/github.com/go-toolkit?status.svg)](https://pkg.go.dev/github.com/go-toolkit)
[![License](https://img.shields.io/github/license/go-toolkit/go-toolkit.svg)](LICENSE)

Go-Toolkit 是一组高性能、生产就绪的 Go 工具集合，用于构建现代应用程序。每个工具都可以独立使用，也可以与工具集中的其他组件一起使用。

## 可用工具

### 限流器 (Rate Limiter)

一个高性能、可扩展的分布式限流库，基于 Redis 实现，支持滑动窗口算法，适用于微服务和 API 网关等场景。

#### 快速开始

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
	// 创建 Redis 客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	
	// 创建限流器，1秒内最多允许10个请求
	limiter := ratelimit.NewRedisSlidingWindowLimiter(redisClient, time.Second, 10)
	
	ctx := context.Background()
	
	// 使用限流器
	limited, err := limiter.Limit(ctx, "user:123")
	if err != nil {
		log.Fatalf("限流器错误: %v", err)
	}
	
	if limited {
		log.Println("请求被限流")
	} else {
		log.Println("请求通过")
	}
}