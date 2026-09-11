package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
)

type clientRate struct {
	tokens     float64
	lastRefill time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	clients  map[string]*clientRate
	rate     float64 // tokens per second
	capacity float64 // max burst capacity
}

func NewRateLimiter(rate, capacity float64) *RateLimiter {
	return &RateLimiter{
		clients:  make(map[string]*clientRate),
		rate:     rate,
		capacity: capacity,
	}
}

func RateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		limiter.mu.Lock()
		client, exists := limiter.clients[ip]
		now := time.Now()

		if !exists {
			client = &clientRate{
				tokens:     limiter.capacity - 1,
				lastRefill: now,
			}
			limiter.clients[ip] = client
			limiter.mu.Unlock()
			c.Next()
			return
		}

		// Calculate tokens to add based on elapsed time
		elapsed := now.Sub(client.lastRefill).Seconds()
		client.tokens += elapsed * limiter.rate
		if client.tokens > limiter.capacity {
			client.tokens = limiter.capacity
		}
		client.lastRefill = now

		if client.tokens < 1.0 {
			limiter.mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, dto.ErrorAPIResponse("Too many requests, please slow down"))
			return
		}

		client.tokens -= 1.0
		limiter.mu.Unlock()

		c.Next()
	}
}
