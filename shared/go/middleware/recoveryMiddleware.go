package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/app/shared/go/utils/logger"
)

// ==========================================
// RECOVERY MIDDLEWARE
// ==========================================

// RecoveryMiddleware recovers from panics and logs the error
// This MUST be the outermost middleware to catch ALL panics
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requestID := GetRequestID(r)

				// Log the panic with comprehensive details
				logger.ErrorWithFields("Panic recovered",
					fmt.Errorf("%v", err),
					map[string]interface{}{
						"request_id":   requestID,
						"method":       r.Method,
						"path":         r.URL.Path,
						"ip":           getClientIP(r),
						"user_agent":   r.UserAgent(),
						"headers":      r.Header,
						"query_params": r.URL.Query(),
						"stack":        string(debug.Stack()),
					},
				)

				// Return 500 to client
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
