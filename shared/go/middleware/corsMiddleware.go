package middleware

import (
	"net/http"

	"github.com/app/shared/go/utils/logger"
)

// ==========================================
// CORS MIDDLEWARE
// ==========================================

// CORSMiddleware handles Cross-Origin Resource Sharing
// allowedOrigins: whitelist of allowed origins
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			if origin != "" {
				if isOriginAllowed(origin, allowedOrigins) {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				} else {
					// Log blocked CORS request
					requestID := GetRequestID(r)
					logger.WarnWithFields("CORS request blocked", map[string]interface{}{
						"request_id":      requestID,
						"method":          r.Method,
						"path":            r.URL.Path,
						"origin":          origin,
						"ip":              getClientIP(r),
						"user_agent":      r.UserAgent(),
						"allowed_origins": allowedOrigins,
					})
				}
			}

			// Set allowed methods
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")

			// Set allowed headers
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Authorization, Content-Type, X-CSRF-Token, X-Request-ID, X-API-Key")

			// Set exposed headers
			w.Header().Set("Access-Control-Expose-Headers",
				"Link, X-Request-ID, X-RateLimit-Limit, X-RateLimit-Remaining")

			// Set max age for preflight cache (24 hours)
			w.Header().Set("Access-Control-Max-Age", "86400")

			// Handle preflight OPTIONS request
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
