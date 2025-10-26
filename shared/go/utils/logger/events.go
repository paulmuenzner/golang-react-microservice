// shared/go/utils/logger/events.go

package logger

// ==========================================
// LOG EVENT BASE TYPES
// ==========================================
// Core definitions for the event-based logging system

type ILogEvent struct {
	Code        string   // Unique identifier (e.g., "MW-RL-001")
	Component   string   // Module/category (e.g., "middleware.ratelimit")
	Message     string   // Human-readable message template
	Level       LogLevel // INFO, WARN, ERROR, DEBUG
	Description string   // Detailed description for documentation
}

type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

// ==========================================
// EVENT CATEGORIES (Prefixes)
// ==========================================

const (
	CategoryMiddleware  = "MW"   // Middleware events
	CategoryAuth        = "AUTH" // Authentication/Authorization
	CategoryService     = "SVC"  // Service lifecycle
	CategoryDatabase    = "DB"   // Database operations
	CategoryAPI         = "API"  // API/Business logic
	CategorySecurity    = "SEC"  // Security-specific
	CategoryPerformance = "PERF" // Performance monitoring
	CategoryIntegration = "INT"  // External integrations
)

// ==========================================
// COMPONENT NAMESPACES
// ==========================================

const (
	ComponentMiddlewareRateLimit   = "middleware.ratelimit"
	ComponentMiddlewareTimeout     = "middleware.timeout"
	ComponentMiddlewareMaxBytes    = "middleware.maxbytes"
	ComponentMiddlewareCORS        = "middleware.cors"
	ComponentMiddlewareSecurity    = "middleware.security"
	ComponentMiddlewareCloudflare  = "middleware.cloudflare"
	ComponentMiddlewareRecovery    = "middleware.recovery"
	ComponentMiddlewareLogging     = "middleware.logging"
	ComponentMiddlewareCompression = "middleware.compression"

	ComponentAuthJWT      = "auth.jwt"
	ComponentAuthSession  = "auth.session"
	ComponentAuthPassword = "auth.password"

	ComponentServiceLifecycle = "service.lifecycle"
	ComponentServiceHealth    = "service.health"

	ComponentDatabaseConnection  = "database.connection"
	ComponentDatabaseQuery       = "database.query"
	ComponentDatabaseMigration   = "database.migration"
	ComponentDatabasePerformance = "database.performance"
)
