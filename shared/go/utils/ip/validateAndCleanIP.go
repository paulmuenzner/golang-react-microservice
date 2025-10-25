package ip

import (
	"net"
	"strings"
)

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
