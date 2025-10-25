package middleware

// import (
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/app/shared/go/utils/logger"
// )

// // ==========================================
// // ADDITIONAL SECURITY MIDDLEWARES
// // ==========================================

// // 1. CSRF Protection Middleware
// func CSRFMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Skip CSRF for safe methods
// 		if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
// 			next.ServeHTTP(w, r)
// 			return
// 		}

// 		// Verify CSRF token
// 		token := r.Header.Get("X-CSRF-Token")
// 		if token == "" {
// 			http.Error(w, "CSRF token missing", http.StatusForbidden)
// 			return
// 		}

// 		// TODO: Validate token against session
// 		// if !validateCSRFToken(token, session) {
// 		//     http.Error(w, "Invalid CSRF token", http.StatusForbidden)
// 		//     return
// 		// }

// 		next.ServeHTTP(w, r)
// 	})
// }

// // 2. IP Whitelist/Blacklist Middleware
// func IPFilterMiddleware(whitelist []string, blacklist []string) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			ip := getClientIP(r)

// 			// Check blacklist first
// 			for _, blocked := range blacklist {
// 				if ip == blocked || strings.HasPrefix(ip, blocked) {
// 					http.Error(w, "Forbidden", http.StatusForbidden)
// 					return
// 				}
// 			}

// 			// If whitelist exists, check it
// 			if len(whitelist) > 0 {
// 				allowed := false
// 				for _, allowedIP := range whitelist {
// 					if ip == allowedIP || strings.HasPrefix(ip, allowedIP) {
// 						allowed = true
// 						break
// 					}
// 				}
// 				if !allowed {
// 					http.Error(w, "Forbidden", http.StatusForbidden)
// 					return
// 				}
// 			}

// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

// // 3. Slow Loris Attack Protection
// func SlowLorisMiddleware(readTimeout time.Duration) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			// Already handled by http.Server ReadTimeout, but can add additional logic
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

// // 4. Request Method Validation
// func MethodFilterMiddleware(allowedMethods []string) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			allowed := false
// 			for _, method := range allowedMethods {
// 				if r.Method == method {
// 					allowed = true
// 					break
// 				}
// 			}

// 			if !allowed {
// 				w.Header().Set("Allow", strings.Join(allowedMethods, ", "))
// 				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 				return
// 			}

// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

// // 5. API Key Authentication Middleware
// func APIKeyMiddleware(validKeys map[string]bool) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			apiKey := r.Header.Get("X-API-Key")
// 			if apiKey == "" {
// 				http.Error(w, "API key required", http.StatusUnauthorized)
// 				return
// 			}

// 			if !validKeys[apiKey] {
// 				logger.WarnWithFields("Invalid API key", map[string]interface{}{
// 					"ip":   getClientIP(r),
// 					"path": r.URL.Path,
// 				})
// 				http.Error(w, "Invalid API key", http.StatusUnauthorized)
// 				return
// 			}

// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

// // 6. Circuit Breaker Middleware (for backend services)
// // Prevents cascading failures by temporarily blocking requests to failing services
// func CircuitBreakerMiddleware(threshold int, timeout time.Duration) func(http.Handler) http.Handler {
// 	// TODO: Implement full circuit breaker pattern
// 	// Consider using: github.com/sony/gobreaker
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

// // 7. Metrics/Monitoring Middleware
// func MetricsMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()

// 		// Wrap response writer to capture status code
// 		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

// 		next.ServeHTTP(wrapped, r)

// 		duration := time.Since(start)

// 		// Log metrics
// 		logger.InfoWithFields("Request metrics", map[string]interface{}{
// 			"method":      r.Method,
// 			"path":        r.URL.Path,
// 			"status":      wrapped.statusCode,
// 			"duration_ms": duration.Milliseconds(),
// 			"ip":          getClientIP(r),
// 		})

// 		// TODO: Send to monitoring system (Prometheus, Datadog, etc.)
// 	})
// }

// type responseWriter struct {
// 	http.ResponseWriter
// 	statusCode int
// }

// func (rw *responseWriter) WriteHeader(code int) {
// 	rw.statusCode = code
// 	rw.ResponseWriter.WriteHeader(code)
// }
// ```

// ---

// ## 📊 Vollständige Middleware-Kategorien

// | Kategorie | Middleware | Zweck |
// |-----------|-----------|-------|
// | **Resilience** | Recovery | Panic handling |
// | | Timeout | Request time limits |
// | | Circuit Breaker | Prevent cascading failures |
// | **Security** | Security Headers | OWASP headers |
// | | CORS | Cross-origin control |
// | | CSRF | Cross-site request forgery |
// | | IP Filter | Whitelist/Blacklist |
// | | API Key | Authentication |
// | **Traffic Control** | Rate Limiting | Per-IP throttling |
// | | Max Bytes | Request size limits |
// | **Performance** | Compression | Gzip responses |
// | | Caching | Response caching |
// | **Observability** | Logging | Request/response logs |
// | | Metrics | Performance tracking |
// | | Request ID | Distributed tracing |

// ---

// ## ✅ Optimale Reihenfolge
// ```
// Request kommt rein
// ↓
// 1. Recovery (fängt ALLES)
// ↓
// 2. Timeout (Gesamt-Zeit-Limit)
// ↓
// 3. Max Bytes (frühe Ablehnung großer Requests)
// ↓
// 4. CORS (Preflight schnell behandeln)
// ↓
// 5. Security Headers (immer setzen)
// ↓
// 6. Rate Limiting (nach Security, vor Business Logic)
// ↓
// 7. Compression (vor Logging)
// ↓
// 8. Request ID (für Tracing)
// ↓
// 9. Logging (sieht alles)
// ↓
// 10. Router/Proxy (Business Logic)
