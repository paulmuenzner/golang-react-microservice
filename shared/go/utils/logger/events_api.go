// shared/go/utils/logger/events_api.go

package logger

// ==========================================
// API EVENTS (API prefix)
// ==========================================

// Request Processing Events (API-REQ-xxx)
var (
	EventAPIRequestValidationFailed = ILogEvent{
		Code:        "API-REQ-001",
		Component:   "api.validation",
		Message:     "Request validation failed",
		Level:       LevelWarn,
		Description: "API request failed validation",
	}

	EventAPIRequestProcessed = ILogEvent{
		Code:        "API-REQ-002",
		Component:   "api.handler",
		Message:     "Request processed successfully",
		Level:       LevelInfo,
		Description: "API request processed successfully",
	}

	EventAPIRequestFailed = ILogEvent{
		Code:        "API-REQ-003",
		Component:   "api.handler",
		Message:     "Request processing failed",
		Level:       LevelError,
		Description: "API request processing encountered error",
	}
)

// Business Logic Events (API-BIZ-xxx)
var (
	EventBusinessRuleViolation = ILogEvent{
		Code:        "API-BIZ-001",
		Component:   "api.business",
		Message:     "Business rule violation",
		Level:       LevelWarn,
		Description: "Request violated business rules",
	}

	EventResourceNotFound = ILogEvent{
		Code:        "API-BIZ-002",
		Component:   "api.business",
		Message:     "Resource not found",
		Level:       LevelWarn,
		Description: "Requested resource does not exist",
	}

	EventResourceCreated = ILogEvent{
		Code:        "API-BIZ-003",
		Component:   "api.business",
		Message:     "Resource created",
		Level:       LevelInfo,
		Description: "New resource successfully created",
	}

	EventResourceUpdated = ILogEvent{
		Code:        "API-BIZ-004",
		Component:   "api.business",
		Message:     "Resource updated",
		Level:       LevelInfo,
		Description: "Resource successfully updated",
	}

	EventResourceDeleted = ILogEvent{
		Code:        "API-BIZ-005",
		Component:   "api.business",
		Message:     "Resource deleted",
		Level:       LevelInfo,
		Description: "Resource successfully deleted",
	}
)

// External Integration Events (API-INT-xxx)
var (
	EventExternalAPICallStarted = ILogEvent{
		Code:        "API-INT-001",
		Component:   "api.integration",
		Message:     "External API call started",
		Level:       LevelDebug,
		Description: "Initiated call to external API",
	}

	EventExternalAPICallSucceeded = ILogEvent{
		Code:        "API-INT-002",
		Component:   "api.integration",
		Message:     "External API call succeeded",
		Level:       LevelInfo,
		Description: "External API call completed successfully",
	}

	EventExternalAPICallFailed = ILogEvent{
		Code:        "API-INT-003",
		Component:   "api.integration",
		Message:     "External API call failed",
		Level:       LevelError,
		Description: "External API call failed",
	}

	EventExternalAPITimeout = ILogEvent{
		Code:        "API-INT-004",
		Component:   "api.integration",
		Message:     "External API timeout",
		Level:       LevelError,
		Description: "External API call timed out",
	}
)
