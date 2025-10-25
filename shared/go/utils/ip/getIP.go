package ip

import (
	"net"
	"net/http"
	"strings"

	"github.com/app/shared/go/utils/conversion"
	"github.com/app/shared/go/utils/logger"
	misc "github.com/app/shared/go/utils/misc"
)

// ==========================================
// CLIENT IP EXTRACTION
// ==========================================

// getClientIP extracts the real client IP from request headers
// Handles various proxy scenarios including Cloudflare, AWS, GCP, etc.
// Returns the most trustworthy IP address found
func getClientIP(r *http.Request) string {
	// Priority order of IP sources (most trustworthy first):
	// 1. CF-Connecting-IP (Cloudflare - most reliable)
	// 2. True-Client-IP (Cloudflare Enterprise & Akamai)
	// 3. X-Real-IP (Nginx proxy)
	// 4. X-Forwarded-For (Standard, but can be spoofed)
	// 5. RemoteAddr (Direct connection fallback)

	// 1. Cloudflare: CF-Connecting-IP
	// This is the most reliable when using Cloudflare
	if ip := r.Header.Get("CF-Connecting-IP"); ip != "" {
		if validIP := validateAndCleanIP(ip); validIP != "" {
			return validIP
		}
	}

	// 2. Cloudflare Enterprise / Akamai: True-Client-IP
	if ip := r.Header.Get("True-Client-IP"); ip != "" {
		if validIP := validateAndCleanIP(ip); validIP != "" {
			return validIP
		}
	}

	// 3. Nginx / General: X-Real-IP
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		if validIP := validateAndCleanIP(ip); validIP != "" {
			return validIP
		}
	}

	// 4. Standard proxy chain: X-Forwarded-For
	// Format: "client, proxy1, proxy2"
	// We want the leftmost (original client) IP that is NOT a private/internal IP
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Split by comma and take the first public IP
		ips := strings.Split(xff, ",")
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			if validIP := validateAndCleanIP(ip); validIP != "" {
				// Skip private/internal IPs (they are proxies, not the real client)
				if !isPrivateIP(validIP) {
					return validIP
				}
			}
		}
	}

	// 5. Fallback: Direct connection RemoteAddr
	if ip := r.RemoteAddr; ip != "" {
		// Remove port if present (format: "IP:port")
		if host, _, err := net.SplitHostPort(ip); err == nil {
			ip = host
		}
		if validIP := validateAndCleanIP(ip); validIP != "" {
			return validIP
		}
	}

	// Ultimate fallback
	return "unknown"
}

// validateAndCleanIP validates and cleans an IP address string
// Returns empty string if invalid
func validateAndCleanIP(ip string) string {
	ip = strings.TrimSpace(ip)

	// Parse and validate IP
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return ""
	}

	return parsedIP.String()
}

// isPrivateIP checks if an IP is private/internal
// Private ranges: 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16, 127.0.0.0/8
func isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	// Check for private IPv4 ranges
	privateRanges := []string{
		"10.0.0.0/8",     // Private network
		"172.16.0.0/12",  // Private network
		"192.168.0.0/16", // Private network
		"127.0.0.0/8",    // Loopback
		"169.254.0.0/16", // Link-local
		"::1/128",        // IPv6 loopback
		"fc00::/7",       // IPv6 private
		"fe80::/10",      // IPv6 link-local
	}

	for _, cidr := range privateRanges {
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

// isCloudflareRequest checks if the request is coming through Cloudflare
// Useful for additional validation or specific handling
func IsCloudflareRequest(r *http.Request) bool {
	// Cloudflare adds specific headers
	return r.Header.Get("CF-Ray") != "" ||
		r.Header.Get("CF-Connecting-IP") != "" ||
		r.Header.Get("CF-Visitor") != ""
}

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
