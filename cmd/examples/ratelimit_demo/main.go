package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"

	"github.com/go-toolkit/pkg/ratelimit"
)

// Define Prometheus metrics
var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method", "status"},
	)

	ratelimitTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ratelimit_total",
			Help: "Total number of rate limited requests",
		},
		[]string{"path", "method"},
	)
)

func init() {
	// Register Prometheus metrics
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(ratelimitTotal)
}

// RateLimitMiddleware creates a Redis-based rate limiting middleware
func RateLimitMiddleware(limiter ratelimit.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use client IP as the rate limit key
		key := "ratelimit:" + c.ClientIP()

		// Check if request should be rate limited
		limited, err := limiter.Limit(c, key)
		if err != nil {
			log.Printf("Rate limiter error: %v", err)
			c.String(http.StatusInternalServerError, "Internal Server Error")
			c.Abort()
			return
		}

		if limited {
			// Record rate limited request
			ratelimitTotal.WithLabelValues(c.Request.URL.Path, c.Request.Method).Inc()
			c.String(http.StatusTooManyRequests, "Too Many Requests")
			c.Abort()
			return
		}

		// Continue processing the request
		c.Next()
	}
}

// PrometheusMiddleware creates middleware for collecting Prometheus metrics
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Record request completion time
		duration := time.Since(start).Seconds()

		// Update Prometheus metrics
		status := http.StatusText(c.Writer.Status())
		requestsTotal.WithLabelValues(c.Request.URL.Path, c.Request.Method, status).Inc()
		requestDuration.WithLabelValues(c.Request.URL.Path, c.Request.Method, status).Observe(duration)
	}
}

func main() {
	// Create Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Create rate limiter: allow max 10 requests per second
	limiter := ratelimit.NewRedisSlidingWindowLimiter(redisClient, time.Second, 10)

	// Create Gin router
	r := gin.Default()

	// Add middleware
	r.Use(PrometheusMiddleware())

	// Expose Prometheus metrics
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Create a rate-limited route group
	limited := r.Group("/limited")
	limited.Use(RateLimitMiddleware(limiter))
	limited.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	// Create an unlimited route group for comparison
	unlimited := r.Group("/unlimited")
	unlimited.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World (Unlimited)!")
	})

	// Start server
	log.Println("Server starting, listening on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}
