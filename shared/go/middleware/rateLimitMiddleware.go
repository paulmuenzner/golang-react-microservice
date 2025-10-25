// shared/go/middleware/rateLimitMiddleware.go

package middleware

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/app/shared/go/utils/logger"
	"golang.org/x/time/rate"
)

// ==========================================
// RATE LIMITING MIDDLEWARE (PER IP)
// ==========================================

type ipRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

func newIPRateLimiter(r rate.Limit, b int) *ipRateLimiter {
	return &ipRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

func (i *ipRateLimiter) getLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

// Global rate limiter instance
var limiter *ipRateLimiter

// RateLimitMiddleware limits requests per IP address
// requestsPerSecond: maximum requests allowed per second per IP
// burst: maximum burst size (allows temporary spikes)
func RateLimitMiddleware(requestsPerSecond float64, burst int) func(http.Handler) http.Handler {
	limiter = newIPRateLimiter(rate.Limit(requestsPerSecond), burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := GetClientIPFromContext(r)
			ipLimiter := limiter.getLimiter(ip)

			if !ipLimiter.Allow() {
				// ✅ Mit IP-Context für Security-Monitoring
				logger.LogMiddlewareEventWithIPContext(
					logger.EventRateLimitExceeded,
					GetRequestID(r),
					GetIPContextFromContext(r),
					r.Method,
					r.URL.Path,
					r.UserAgent(),
					map[string]interface{}{
						"rate_limit": requestsPerSecond,
						"burst":      burst,
					},
				)

				w.Header().Set("Retry-After", "1")
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%.0f", requestsPerSecond))
				w.Header().Set("X-RateLimit-Remaining", "0")
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
