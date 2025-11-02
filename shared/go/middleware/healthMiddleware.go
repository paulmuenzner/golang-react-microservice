package middleware

import (
	"fmt"
	"net/http"

	"github.com/app/shared/go/utils/logger"
)

// ==========================================
// HEALTH CHECK HANDLER
// ==========================================
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// Extract request ID
	requestID := ""
	if id := r.Context().Value("request_id"); id != nil {
		requestID = id.(string)
	}
	ip := getClientIP(r)

	logger.InfoWithFields("Health Check", map[string]interface{}{
		"request_id": requestID,
		"ip":         ip,
		"path":       r.URL.Path,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"OK"}`)
}
