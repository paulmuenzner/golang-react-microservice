package middleware

import (
	"net/http"
	"os"

	"github.com/app/shared/go/utils/logger"
)

// ==========================================
// SECURITY HEADERS MIDDLEWARE
// ==========================================

// SecurityHeadersMiddleware sets comprehensive security headers
// Implements OWASP recommendations for secure HTTP headers
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Enable XSS protection (legacy browsers)
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Prevent clickjacking attacks
		w.Header().Set("X-Frame-Options", "DENY")

		// Strict Transport Security (HSTS)
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

		// Content Security Policy (CSP)
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https:; " +
			"font-src 'self' data:; " +
			"connect-src 'self'; " +
			"frame-ancestors 'none'; " +
			"base-uri 'self'; " +
			"form-action 'self'"
		w.Header().Set("Content-Security-Policy", csp)

		// Referrer Policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions Policy
		permissions := "accelerometer=(), " +
			"camera=(), " +
			"geolocation=(), " +
			"gyroscope=(), " +
			"magnetometer=(), " +
			"microphone=(), " +
			"payment=(), " +
			"usb=()"
		w.Header().Set("Permissions-Policy", permissions)

		// Remove server information
		w.Header().Set("Server", "")
		w.Header().Del("X-Powered-By")

		// Log if insecure request (HTTP instead of HTTPS) in production
		if r.TLS == nil && r.Header.Get("X-Forwarded-Proto") != "https" {
			if os.Getenv("ENVIRONMENT") == "production" {
				requestID := GetRequestID(r)
				logger.WarnWithFields("Insecure HTTP request in production", map[string]interface{}{
					"request_id": requestID,
					"method":     r.Method,
					"path":       r.URL.Path,
					"ip":         getClientIP(r),
					"user_agent": r.UserAgent(),
				})
			}
		}

		next.ServeHTTP(w, r)
	})
}
