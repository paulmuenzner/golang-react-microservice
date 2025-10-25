package ip

import (
	"net/http"

	"github.com/app/shared/go/utils/conversion"
	"github.com/app/shared/go/utils/logger"
	"github.com/app/shared/go/utils/misc"
)

// GetClientIPWithContext returns client IP with additional context
// Useful for detailed logging and debugging
func GetClientIPWithContext(r *http.Request) map[string]string {
	context := map[string]string{
		"ip":               getClientIP(r),
		"remote_addr":      r.RemoteAddr,
		"x_forwarded_for":  r.Header.Get("X-Forwarded-For"),
		"x_real_ip":        r.Header.Get("X-Real-IP"),
		"cf_connecting_ip": r.Header.Get("CF-Connecting-IP"),
		"true_client_ip":   r.Header.Get("True-Client-IP"),
		"via_cloudflare":   conversion.BoolToString(IsCloudflareRequest(r)),
	}

	// Log if IP extraction seems suspicious
	if context["ip"] == "unknown" || isPrivateIP(context["ip"]) {
		requestID := misc.GetRequestID(r)
		logger.WarnWithFields("Suspicious IP extraction", map[string]interface{}{
			"request_id": requestID,
			"context":    context,
		})
	}

	return context
}
