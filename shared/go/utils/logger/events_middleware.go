// shared/go/utils/logger/events_middleware.go

package logger

// ==========================================
// MIDDLEWARE EVENTS (MW prefix)
// ==========================================

// Rate Limiting Events (MW-RL-xxx)
var (
	EventRateLimitExceeded = ILogEvent{
		Code:        "MW-RL-001",
		Component:   ComponentMiddlewareRateLimit,
		Message:     "Rate limit exceeded",
		Level:       LevelWarn,
		Description: "Client exceeded allowed request rate per IP",
	}

	EventRateLimitIPBlocked = ILogEvent{
		Code:        "MW-RL-002",
		Component:   ComponentMiddlewareRateLimit,
		Message:     "IP temporarily blocked",
		Level:       LevelError,
		Description: "IP address blocked after multiple rate limit violations",
	}

	EventRateLimitConfigChanged = ILogEvent{
		Code:        "MW-RL-003",
		Component:   ComponentMiddlewareRateLimit,
		Message:     "Rate limit configuration changed",
		Level:       LevelInfo,
		Description: "Rate limit settings were updated",
	}
)

// Request Size Events (MW-SZ-xxx)
var (
	EventRequestTooLarge = ILogEvent{
		Code:        "MW-SZ-001",
		Component:   ComponentMiddlewareMaxBytes,
		Message:     "Request body too large",
		Level:       LevelWarn,
		Description: "Request body exceeded maximum allowed size",
	}

	EventRequestSizeWarning = ILogEvent{
		Code:        "MW-SZ-002",
		Component:   ComponentMiddlewareMaxBytes,
		Message:     "Request approaching size limit",
		Level:       LevelWarn,
		Description: "Request body size is close to maximum (>80%)",
	}
)

// Timeout Events (MW-TO-xxx)
var (
	EventRequestTimeout = ILogEvent{
		Code:        "MW-TO-001",
		Component:   ComponentMiddlewareTimeout,
		Message:     "Request timeout",
		Level:       LevelWarn,
		Description: "Request exceeded maximum processing time",
	}

	EventSlowRequest = ILogEvent{
		Code:        "MW-TO-002",
		Component:   ComponentMiddlewareTimeout,
		Message:     "Slow request detected",
		Level:       LevelWarn,
		Description: "Request took longer than expected but didn't timeout",
	}
)

// Security Events (MW-SEC-xxx)
var (
	EventCloudflareSpoof = ILogEvent{
		Code:        "MW-SEC-001",
		Component:   ComponentMiddlewareCloudflare,
		Message:     "Cloudflare header spoofing detected",
		Level:       LevelError,
		Description: "Request claims to be from Cloudflare but IP validation failed",
	}

	EventCORSBlocked = ILogEvent{
		Code:        "MW-SEC-002",
		Component:   ComponentMiddlewareCORS,
		Message:     "CORS request blocked",
		Level:       LevelWarn,
		Description: "Request origin not in whitelist",
	}

	EventInsecureHTTP = ILogEvent{
		Code:        "MW-SEC-003",
		Component:   ComponentMiddlewareSecurity,
		Message:     "Insecure HTTP request in production",
		Level:       LevelWarn,
		Description: "HTTP request received when HTTPS is required",
	}

	EventSuspiciousHeaders = ILogEvent{
		Code:        "MW-SEC-004",
		Component:   ComponentMiddlewareSecurity,
		Message:     "Suspicious request headers detected",
		Level:       LevelWarn,
		Description: "Request contains potentially malicious headers",
	}

	EventIPSpoofingAttempt = ILogEvent{
		Code:        "MW-SEC-005",
		Component:   ComponentMiddlewareSecurity,
		Message:     "IP spoofing attempt detected",
		Level:       LevelError,
		Description: "Mismatch between connection IP and forwarded headers",
	}
)

// Recovery Events (MW-REC-xxx)
var (
	EventPanicRecovered = ILogEvent{
		Code:        "MW-REC-001",
		Component:   ComponentMiddlewareRecovery,
		Message:     "Panic recovered",
		Level:       LevelError,
		Description: "Application panic caught and recovered",
	}
)

// Compression Events (MW-CMP-xxx)
var (
	EventCompressionFailed = ILogEvent{
		Code:        "MW-CMP-001",
		Component:   ComponentMiddlewareCompression,
		Message:     "Compression failed",
		Level:       LevelWarn,
		Description: "Failed to compress response",
	}

	EventCompressionSkipped = ILogEvent{
		Code:        "MW-CMP-002",
		Component:   ComponentMiddlewareCompression,
		Message:     "Compression skipped",
		Level:       LevelDebug,
		Description: "Response too small or wrong content-type for compression",
	}
)

// Request Logging Events (MW-LOG-xxx)
var (
	EventRequestIncoming = ILogEvent{
		Code:        "MW-LOG-001",
		Component:   ComponentMiddlewareLogging,
		Message:     "Incoming request",
		Level:       LevelInfo,
		Description: "New HTTP request received",
	}

	EventRequestCompleted = ILogEvent{
		Code:        "MW-LOG-002",
		Component:   ComponentMiddlewareLogging,
		Message:     "Request completed",
		Level:       LevelInfo,
		Description: "HTTP request processing completed",
	}

	EventRequestFailed = ILogEvent{
		Code:        "MW-LOG-003",
		Component:   ComponentMiddlewareLogging,
		Message:     "Request failed",
		Level:       LevelError,
		Description: "HTTP request processing failed",
	}
)
