package middleware

import (
	"net/http"

	"github.com/app/shared/go/utils/logger"
)

// ==========================================
// HEALTH CHECK HANDLER
// ==========================================
func HealthHandler(w http.ResponseWriter, r *http.Request) {

	requestID := ""
	if id := r.Context().Value("request_id"); id != nil {
		requestID = id.(string)
	}

	logger.InfoWithFields("Health Check", map[string]interface{}{
		"request_id": requestID,
		"ip":         r.RemoteAddr,
		"path":       r.URL.Path,
	})
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("gateway OK"))
}
