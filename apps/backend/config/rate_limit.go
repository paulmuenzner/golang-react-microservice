// config/rate_limit.go
package config

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	mu               sync.Mutex
	requests         map[string]int
	blockedUntil     map[string]time.Time // Stores the time until the IP is blocked
	maxRequests      int
	interval         time.Duration
	cooldownDuration time.Duration // Duration for cooldown (10 minutes)
}

func NewRateLimiter(maxRequests int, interval, cooldownDuration time.Duration) *RateLimiter {
	return &RateLimiter{
		requests:         make(map[string]int),
		blockedUntil:     make(map[string]time.Time),
		maxRequests:      maxRequests,
		interval:         interval,
		cooldownDuration: cooldownDuration,
	}
}

// LimitIP checks if the IP is within rate limits and handles cooldown logic
func (rl *RateLimiter) LimitIP(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Check if the IP is currently blocked
	if blockedUntil, ok := rl.blockedUntil[ip]; ok {
		if time.Now().Before(blockedUntil) {
			// If the IP is still in cooldown period, reject the request
			return false
		}
		// If cooldown is over, reset the request count and unblock the IP
		delete(rl.blockedUntil, ip)
		rl.requests[ip] = 0
	}

	// Reset the request count periodically
	if _, ok := rl.requests[ip]; !ok {
		rl.requests[ip] = 0
	}

	// Check if the limit is exceeded
	if rl.requests[ip] >= rl.maxRequests {
		// Block the IP for the cooldown period
		rl.blockedUntil[ip] = time.Now().Add(rl.cooldownDuration)
		return false
	}

	// Increment the request count
	rl.requests[ip]++

	return true
}

// ResetRequests resets the request count for all IPs periodically
func (rl *RateLimiter) ResetRequests() {
	for {
		time.Sleep(rl.interval)
		rl.mu.Lock()
		for ip := range rl.requests {
			rl.requests[ip] = 0
		}
		rl.mu.Unlock()
	}
}

func RateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !rl.LimitIP(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": "Rate limit exceeded. Please try again later."})
			c.Abort()
			return
		}

		// Continue with the next middleware/handler
		c.Next()
	}
}
