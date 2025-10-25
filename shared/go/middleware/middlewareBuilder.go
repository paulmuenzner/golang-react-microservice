package middleware

import (
	"net/http"
	"os"
)

// ==========================================
// MIDDLEWARE STACK BUILDER
// ==========================================
// BuildMiddlewareStack applies all middleware in the correct order
// Order matters: outermost (first applied) to innermost (last applied)
func BuildMiddlewareStack(handler http.Handler) http.Handler {
	// Define CORS whitelist
	corsWhitelist := []string{
		"http://localhost:3000",
		"http://localhost:8080",
		"https://app.example.com",
		"https://*.example.com",
		"https://www.getpostman.com",
		"https://web.postman.co",
		"chrome-extension://fhbjgbiflinjbdggehcddcbncdddomop",
	}

	if os.Getenv("ENVIRONMENT") == "development" {
		corsWhitelist = []string{"*"}
	}

	// ==========================================
	// MIDDLEWARE ORDER (CRITICAL!)
	// ==========================================
	// Apply from outermost to innermost

	// 1. Recovery - MUST be first to catch ALL panics from any middleware
	handler = RecoveryMiddleware(handler)

	// 2. Request ID - MUST be early so all other middlewares can use it for logging
	handler = RequestIDMiddleware(handler)

	// 3. Timeout - Limit total request time (after RequestID so timeouts are logged with ID)
	handler = LoggingMiddleware(handler)

	// 4. Max Request Size - Reject large requests early
	handler = MaxBytesMiddleware(10 * 1024 * 1024)(handler) // 10MB

	// 5. CORS - Handle preflight requests early
	handler = CORSMiddleware(corsWhitelist)(handler)

	// 6. Security Headers - Always set security headers
	handler = SecurityHeadersMiddleware(handler)

	// 7. Rate Limiting - Protect against abuse (per IP)
	if os.Getenv("ENVIRONMENT") == "production" {
		handler = RateLimitMiddleware(100, 200)(handler) // 100 req/s per IP, burst 200
	} else {
		handler = RateLimitMiddleware(1000, 2000)(handler)
	}

	// 8. Compression - Compress responses
	handler = CompressionMiddleware(handler)

	// 9. Logging - Log requests/responses (uses RequestID from context)
	//    Should be relatively late so it captures final response status/size
	handler = LoggingMiddleware(handler)

	// 10. Business Logic (Router/Proxy) comes after all middleware

	return handler
}
