package middleware

import (
	"context"
	"net/http"

	"github.com/app/shared/go/utils/security"
)

// ==========================================
// REQUEST ID MIDDLEWARE
// ==========================================
// RequestIDMiddleware generates a unique ID for each request
// This MUST be one of the first middlewares so other middlewares can use the ID
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if Request ID already exists (e.g., from load balancer)
		requestID := r.Header.Get("X-Request-ID")

		// Generate new ID if not present
		if requestID == "" {
			requestID = security.GenerateID()
		}

		// Add Request ID to context for downstream use
		ctx := context.WithValue(r.Context(), "request_id", requestID)
		r = r.WithContext(ctx)

		// Add Request ID to response headers for client tracing
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r)
	})
}
