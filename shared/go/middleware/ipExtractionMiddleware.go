// middleware/ipextraction.go

package middleware

import (
	"context"
	"net/http"
	"os"

	ip "github.com/app/shared/go/utils/ip"
	"github.com/app/shared/go/utils/logger"
)

// Output example of ip.GetClientIPWithContext(r):
// {
//     "ip":                 "203.0.113.45",           // The final resolved IP
//     "remote_addr":        "104.16.0.1:443",         // Connection IP (Cloudflare)
//     "x_forwarded_for":    "203.0.113.45, 10.0.0.1", // Proxy Chain
//     "x_real_ip":          "203.0.113.45",           // Nginx header
//     "cf_connecting_ip":   "203.0.113.45",           // Cloudflare header
//     "true_client_ip":     "",                       // Enterprise CDN header
//     "via_cloudflare":     "true",                   // Cloudflare detected?
// }

// ==========================================
// IP EXTRACTION MIDDLEWARE
// ==========================================

// IPExtractionMiddleware extracts and caches client IP + context
func IPExtractionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract IP
		clientIP := getClientIP(r)

		// Get detailed context for debugging (cached for reuse)
		ipContext := ip.GetClientIPWithContext(r)

		// Add both to context
		ctx := context.WithValue(r.Context(), "client_ip", clientIP)
		ctx = context.WithValue(ctx, "ip_context", ipContext)
		r = r.WithContext(ctx)

		// Add IP to response headers for client tracing
		w.Header().Set("X-Client-IP", clientIP)

		// Log IP extraction in development for debugging
		if os.Getenv("ENVIRONMENT") == "development" {
			requestID := GetRequestID(r)
			logger.DebugWithFields("IP extracted", map[string]interface{}{
				"request_id": requestID,
				"ip":         clientIP,
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

// GetIPContextFromContext retrieves detailed IP context from request context
func GetIPContextFromContext(r *http.Request) map[string]string {
	if ctx := r.Context().Value("ip_context"); ctx != nil {
		return ctx.(map[string]string)
	}
	return ip.GetClientIPWithContext(r) // Fallback
}
