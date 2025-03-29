# Go-Toolkit

Go-Toolkit 是一个实用工具集合，提供各种常用功能的高效实现。

## 组件

- **ratelimit**: 基于Redis的滑动窗口限流器
- 更多组件开发中...

## 安装

```bash
go get github.com/yourusername/go-toolkit
```

## 使用示例

### Redis滑动窗口限流器

```go
package main

import (
    "context"
    "time"
    
    "github.com/redis/go-redis/v9"
    "github.com/yourusername/go-toolkit/pkg/ratelimit"
)

func main() {
    // 创建Redis客户端
    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    
    // 创建限流器(1秒内最多5个请求)
    limiter := ratelimit.NewRedisSlidingWindowLimiter(client, time.Second, 5)
    
    // 检查限流
    limited, err := limiter.Limit(context.Background(), "user:123")
    if err != nil {
        panic(err)
    }
    
    if limited {
        // 请求被限流
    } else {
        // 请求正常处理
    }
}
```

### Gin中间件集成

```go
import "github.com/yourusername/go-toolkit/pkg/ratelimit"

func main() {
    // 创建限流器
    limiter := ratelimit.NewRedisSlidingWindowLimiter(redisClient, time.Second, 10)
    
    // 在Gin中使用
    r := gin.Default()
    r.Use(func(c *gin.Context) {
        limited, _ := limiter.Limit(c, "rate:"+c.ClientIP())
        if limited {
            c.AbortWithStatus(429) // Too Many Requests
            return
        }
        c.Next()
    })
}
```

## 测试

```bash
# 运行单元测试
make unit-test

# 运行集成测试(需要Redis)
make integration-test

# 启动Redis
make run-redis