package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/app/shared/go/utils/logger"
)

// ==========================================
// TIMEOUT MIDDLEWARE
// ==========================================
// TimeoutMiddleware ensures requests don't run indefinitely
// Returns 504 Gateway Timeout if request takes longer than specified duration
func TimeoutMiddleware(next http.Handler, timeout time.Duration) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create context with timeout
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		// Replace request context with timeout context
		r = r.WithContext(ctx)

		// Channel to signal completion
		done := make(chan bool, 1)

		// Track start time for duration calculation
		start := time.Now()

		// Run handler in goroutine
		go func() {
			next.ServeHTTP(w, r)
			done <- true
		}()

		// Wait for completion or timeout
		select {
		case <-done:
			// Request completed successfully
			return
		case <-ctx.Done():
			// Timeout occurred
			requestID := GetRequestID(r)
			duration := time.Since(start)

			logger.WarnWithFields("Request timeout", map[string]interface{}{
				"request_id":  requestID,
				"method":      r.Method,
				"path":        r.URL.Path,
				"ip":          getClientIP(r),
				"user_agent":  r.UserAgent(),
				"timeout":     timeout.String(),
				"duration":    duration.String(),
				"duration_ms": duration.Milliseconds(),
			})

			http.Error(w, "Gateway Timeout", http.StatusGatewayTimeout)
		}
	})
}
