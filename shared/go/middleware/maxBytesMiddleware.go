package middleware

import (
	"net/http"

	"github.com/app/shared/go/utils/logger"
)

// ==========================================
// MAX BYTES MIDDLEWARE
// ==========================================

// MaxBytesMiddleware limits the size of request bodies
// maxBytes: maximum allowed request body size in bytes
func MaxBytesMiddleware(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Limit request body size
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

			// Check if content length exceeds limit
			if r.ContentLength > maxBytes {
				// ✅ Event-based logging mit Helper
				logger.LogMiddlewareEvent(
					logger.EventRequestTooLarge,
					GetRequestID(r),
					GetClientIPFromContext(r),
					r.Method,
					r.URL.Path,
					r.UserAgent(),
					map[string]interface{}{
						"content_length": r.ContentLength,
						"max_bytes":      maxBytes,
						"exceeded_by":    r.ContentLength - maxBytes,
					},
				)

				http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
