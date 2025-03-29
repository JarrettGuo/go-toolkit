# Go-Toolkit

[![Go Report Card](https://goreportcard.com/badge/github.com/go-toolkit)](https://goreportcard.com/report/github.com/go-toolkit)
[![GoDoc](https://godoc.org/github.com/go-toolkit?status.svg)](https://pkg.go.dev/github.com/go-toolkit)
[![License](https://img.shields.io/github/license/go-toolkit/go-toolkit.svg)](LICENSE)

Go-Toolkit 是一组高性能、生产就绪的 Go 工具集合，用于构建现代应用程序。每个工具都可以独立使用，也可以与工具集中的其他组件一起使用。

## 可用工具

### 限流器 (Rate Limiter)

一个高性能、可扩展的分布式限流库，基于 Redis 实现，支持滑动窗口算法，适用于微服务和 API 网关等场景。

[查看详细文档](pkg/ratelimit/README.md)

## 开发

### 环境要求

- Go 1.18+
- Redis (用于限流器)
- Docker & Docker Compose (用于集成测试)

### 构建和测试

```bash
# 运行所有检查和测试
make all

# 仅构建
make build

# 运行测试
make test

# 运行集成测试（需要 Redis）
make integration-test

# 启动开发环境
make docker-up

我会将快速开始部分放入各个README文件中。以下是包含了快速开始的完整README文件：
主项目README.md（中文版）
markdownCopy# Go-Toolkit

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