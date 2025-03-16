package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
	rate     int
	window   time.Duration
}

func (mw *MiddlewareManager) RateLimit() gin.HandlerFunc {
	limiter := &RateLimiter{
		requests: make(map[string][]time.Time),
		rate:     mw.cfg.RateLimit.Rate,
		window:   time.Minute,
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.Allow(ip) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			return
		}

		c.Next()
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Remove old requests
	var recent []time.Time
	for _, t := range rl.requests[key] {
		if t.After(windowStart) {
			recent = append(recent, t)
		}
	}

	rl.requests[key] = recent

	// Check rate limit
	if len(recent) >= rl.rate {
		return false
	}

	// Add new request
	rl.requests[key] = append(rl.requests[key], now)
	return true
}
