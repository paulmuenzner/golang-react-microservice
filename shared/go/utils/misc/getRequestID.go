package misc

import (
	"net/http"
)

// GetRequestID is a helper to extract Request ID from context
func GetRequestID(r *http.Request) string {
	if id := r.Context().Value("request_id"); id != nil {
		return id.(string)
	}
	return ""
}
