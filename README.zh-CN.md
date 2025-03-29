# README.zh-CN.md (中文版)

```markdown
# Go-Toolkit

[![Go Report Card](https://goreportcard.com/badge/github.com/go-toolkit)](https://goreportcard.com/report/github.com/go-toolkit)
[![GoDoc](https://godoc.org/github.com/go-toolkit?status.svg)](https://pkg.go.dev/github.com/go-toolkit)
[![License](https://img.shields.io/github/license/go-toolkit/go-toolkit.svg)](LICENSE)

高性能、生产级 Go 开发工具集，用于构建现代应用程序。

## 工具列表

### 限流器 (Rate Limiter)

基于 Redis 的分布式限流解决方案，采用滑动窗口算法。

```go
// 快速开始示例
import (
    "github.com/redis/go-redis/v9"
    "github.com/go-toolkit/pkg/ratelimit"
)

// 创建限流器：每秒最多10个请求
limiter := ratelimit.NewRedisSlidingWindowLimiter(
    redis.NewClient(&redis.Options{Addr: "localhost:6379"}), 
    time.Second, 
    10,
)

// 检查请求是否应该被限流
limited, err := limiter.Limit(ctx, "user:123")