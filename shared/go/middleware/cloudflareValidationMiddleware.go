// middleware/cloudflare.go

package middleware

import (
	"net"
	"net/http"

	ip "github.com/app/shared/go/utils/ip"
	"github.com/app/shared/go/utils/logger"
)

// ==========================================
// CLOUDFLARE VALIDATION
// ==========================================

// CloudflareIPRanges contains Cloudflare's IP ranges
// Update from: https://www.cloudflare.com/ips/
var CloudflareIPRanges = []string{
	// IPv4
	"173.245.48.0/20",
	"103.21.244.0/22",
	"103.22.200.0/22",
	"103.31.4.0/22",
	"141.101.64.0/18",
	"108.162.192.0/18",
	"190.93.240.0/20",
	"188.114.96.0/20",
	"197.234.240.0/22",
	"198.41.128.0/17",
	"162.158.0.0/15",
	"104.16.0.0/13",
	"104.24.0.0/14",
	"172.64.0.0/13",
	"131.0.72.0/22",
	// Add more as needed
}

// ValidateCloudflareRequest checks if request is genuinely from Cloudflare
// Validates the connecting IP against Cloudflare's published IP ranges
func ValidateCloudflareRequest(r *http.Request) bool {
	// If no CF headers, not from Cloudflare
	if !ip.IsCloudflareRequest(r) {
		return false
	}

	// Get the actual connecting IP (not the client IP)
	connectingIP := r.RemoteAddr
	if host, _, err := net.SplitHostPort(connectingIP); err == nil {
		connectingIP = host
	}

	ip := net.ParseIP(connectingIP)
	if ip == nil {
		return false
	}

	// Check if IP is in Cloudflare's ranges
	for _, cidr := range CloudflareIPRanges {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if network.Contains(ip) {
			return true
		}
	}

	return false
}

// CloudflareValidationMiddleware validates requests claiming to be from Cloudflare
func CloudflareValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only validate if request claims to be from Cloudflare
		if ip.IsCloudflareRequest(r) && !ValidateCloudflareRequest(r) {
			requestID := GetRequestID(r)

			logger.WarnWithFields("Spoofed Cloudflare headers detected", map[string]interface{}{
				"request_id":  requestID,
				"remote_addr": r.RemoteAddr,
				"cf_ray":      r.Header.Get("CF-Ray"),
				"path":        r.URL.Path,
			})

			// Optional: block the request or just log
			// http.Error(w, "Forbidden", http.StatusForbidden)
			// return
		}

		next.ServeHTTP(w, r)
	})
}
