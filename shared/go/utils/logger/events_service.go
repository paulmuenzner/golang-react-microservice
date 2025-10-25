// shared/go/utils/logger/events_service.go

package logger

// ==========================================
// SERVICE EVENTS (SVC prefix)
// ==========================================

// Lifecycle Events (SVC-LC-xxx)
var (
	EventServiceStarted = ILogEvent{
		Code:        "SVC-LC-001",
		Component:   ComponentServiceLifecycle,
		Message:     "Service started",
		Level:       LevelInfo,
		Description: "Service successfully started and listening",
	}

	EventServiceStopping = ILogEvent{
		Code:        "SVC-LC-002",
		Component:   ComponentServiceLifecycle,
		Message:     "Service stopping",
		Level:       LevelInfo,
		Description: "Service shutdown initiated",
	}

	EventServiceStopped = ILogEvent{
		Code:        "SVC-LC-003",
		Component:   ComponentServiceLifecycle,
		Message:     "Service stopped",
		Level:       LevelInfo,
		Description: "Service gracefully stopped",
	}

	EventServiceStartupFailed = ILogEvent{
		Code:        "SVC-LC-004",
		Component:   ComponentServiceLifecycle,
		Message:     "Service startup failed",
		Level:       LevelError,
		Description: "Service failed to start",
	}

	EventServiceRestarting = ILogEvent{
		Code:        "SVC-LC-005",
		Component:   ComponentServiceLifecycle,
		Message:     "Service restarting",
		Level:       LevelInfo,
		Description: "Service is restarting",
	}
)

// Health Check Events (SVC-HC-xxx)
var (
	EventHealthCheckOK = ILogEvent{
		Code:        "SVC-HC-001",
		Component:   ComponentServiceHealth,
		Message:     "Health check passed",
		Level:       LevelDebug,
		Description: "Service health check successful",
	}

	EventHealthCheckFailed = ILogEvent{
		Code:        "SVC-HC-002",
		Component:   ComponentServiceHealth,
		Message:     "Health check failed",
		Level:       LevelError,
		Description: "Service health check indicated unhealthy state",
	}

	EventHealthCheckDegraded = ILogEvent{
		Code:        "SVC-HC-003",
		Component:   ComponentServiceHealth,
		Message:     "Health check degraded",
		Level:       LevelWarn,
		Description: "Service partially healthy but degraded",
	}

	EventReadinessCheckFailed = ILogEvent{
		Code:        "SVC-HC-004",
		Component:   ComponentServiceHealth,
		Message:     "Readiness check failed",
		Level:       LevelWarn,
		Description: "Service not ready to accept traffic",
	}
)

// Configuration Events (SVC-CFG-xxx)
var (
	EventConfigLoaded = ILogEvent{
		Code:        "SVC-CFG-001",
		Component:   ComponentServiceLifecycle,
		Message:     "Configuration loaded",
		Level:       LevelInfo,
		Description: "Service configuration successfully loaded",
	}

	EventConfigReloaded = ILogEvent{
		Code:        "SVC-CFG-002",
		Component:   ComponentServiceLifecycle,
		Message:     "Configuration reloaded",
		Level:       LevelInfo,
		Description: "Service configuration reloaded",
	}

	EventConfigInvalid = ILogEvent{
		Code:        "SVC-CFG-003",
		Component:   ComponentServiceLifecycle,
		Message:     "Invalid configuration",
		Level:       LevelError,
		Description: "Configuration validation failed",
	}
)
