// middleware/ipextraction.go

package middleware

import (
	"context"
	"net/http"
	"os"

	"github.com/app/shared/go/utils/logger"
)

// ==========================================
// IP EXTRACTION MIDDLEWARE (Optional)
// ==========================================

// IPExtractionMiddleware adds client IP to request context
// Useful for consistent IP access and debugging
func IPExtractionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := getClientIP(r)

		// Add to context for easy access by other handlers
		ctx := context.WithValue(r.Context(), "client_ip", clientIP)
		r = r.WithContext(ctx)

		// Log IP extraction in development for debugging
		if os.Getenv("ENVIRONMENT") == "development" {
			requestID := GetRequestID(r)
			ipContext := GetClientIPWithContext(r)

			logger.DebugWithFields("IP extracted", map[string]interface{}{
				"request_id": requestID,
				"ip_context": ipContext,
			})
		}

		next.ServeHTTP(w, r)
	})
}

// GetClientIPFromContext retrieves client IP from request context
func GetClientIPFromContext(r *http.Request) string {
	if ip := r.Context().Value("client_ip"); ip != nil {
		return ip.(string)
	}
	return getClientIP(r) // Fallback
}
