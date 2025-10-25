package middleware

import (
	"net/http"
	"time"

	"github.com/app/shared/go/utils/logger"
)

// ==========================================
// LOGGING MIDDLEWARE
// ==========================================
// LoggingMiddleware logs all HTTP requests and responses
// IMPORTANT: This relies on RequestIDMiddleware being applied first!
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Get Request ID from context (set by RequestIDMiddleware)
		requestID := GetRequestID(r)

		// Create logger with Request ID
		reqLogger := logger.WithRequestID(requestID)

		// Log incoming request
		reqLogger.InfoWithFields("Incoming request", map[string]interface{}{
			"method":      r.Method,
			"path":        r.URL.Path,
			"ip":          getClientIP(r),
			"user_agent":  r.UserAgent(),
			"content_len": r.ContentLength,
		})

		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Process request
		next.ServeHTTP(wrapped, r)

		// Log completed request
		duration := time.Since(start)
		reqLogger.HTTP(
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			duration,
			getClientIP(r),
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
