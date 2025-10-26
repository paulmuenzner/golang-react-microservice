package ip

import (
	"net"
	"net/http"
	"strings"
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

// isCloudflareRequest checks if the request is coming through Cloudflare
// Useful for additional validation or specific handling
func IsCloudflareRequest(r *http.Request) bool {
	// Cloudflare adds specific headers
	return r.Header.Get("CF-Ray") != "" ||
		r.Header.Get("CF-Connecting-IP") != "" ||
		r.Header.Get("CF-Visitor") != ""
}
