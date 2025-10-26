// shared/go/utils/logger/events_auth.go

package logger

// ==========================================
// AUTHENTICATION EVENTS (AUTH prefix)
// ==========================================

// JWT Events (AUTH-JWT-xxx)
var (
	EventAuthSuccess = ILogEvent{
		Code:        "AUTH-JWT-001",
		Component:   ComponentAuthJWT,
		Message:     "Authentication successful",
		Level:       LevelInfo,
		Description: "User successfully authenticated",
	}

	EventAuthFailed = ILogEvent{
		Code:        "AUTH-JWT-002",
		Component:   ComponentAuthJWT,
		Message:     "Authentication failed",
		Level:       LevelWarn,
		Description: "Authentication attempt failed",
	}

	EventTokenExpired = ILogEvent{
		Code:        "AUTH-JWT-003",
		Component:   ComponentAuthJWT,
		Message:     "Token expired",
		Level:       LevelWarn,
		Description: "JWT token has expired",
	}

	EventTokenInvalid = ILogEvent{
		Code:        "AUTH-JWT-004",
		Component:   ComponentAuthJWT,
		Message:     "Invalid token",
		Level:       LevelWarn,
		Description: "JWT token validation failed",
	}

	EventTokenRefreshed = ILogEvent{
		Code:        "AUTH-JWT-005",
		Component:   ComponentAuthJWT,
		Message:     "Token refreshed",
		Level:       LevelInfo,
		Description: "JWT token successfully refreshed",
	}

	EventTokenRevoked = ILogEvent{
		Code:        "AUTH-JWT-006",
		Component:   ComponentAuthJWT,
		Message:     "Token revoked",
		Level:       LevelWarn,
		Description: "JWT token was revoked",
	}
)

// Session Events (AUTH-SES-xxx)
var (
	EventSessionCreated = ILogEvent{
		Code:        "AUTH-SES-001",
		Component:   ComponentAuthSession,
		Message:     "Session created",
		Level:       LevelInfo,
		Description: "New user session created",
	}

	EventSessionExpired = ILogEvent{
		Code:        "AUTH-SES-002",
		Component:   ComponentAuthSession,
		Message:     "Session expired",
		Level:       LevelInfo,
		Description: "User session expired",
	}

	EventSessionInvalid = ILogEvent{
		Code:        "AUTH-SES-003",
		Component:   ComponentAuthSession,
		Message:     "Invalid session",
		Level:       LevelWarn,
		Description: "Session validation failed",
	}
)

// Password Events (AUTH-PWD-xxx)
var (
	EventPasswordChanged = ILogEvent{
		Code:        "AUTH-PWD-001",
		Component:   ComponentAuthPassword,
		Message:     "Password changed",
		Level:       LevelInfo,
		Description: "User password successfully changed",
	}

	EventPasswordResetRequested = ILogEvent{
		Code:        "AUTH-PWD-002",
		Component:   ComponentAuthPassword,
		Message:     "Password reset requested",
		Level:       LevelInfo,
		Description: "User requested password reset",
	}

	EventPasswordResetCompleted = ILogEvent{
		Code:        "AUTH-PWD-003",
		Component:   ComponentAuthPassword,
		Message:     "Password reset completed",
		Level:       LevelInfo,
		Description: "Password reset successfully completed",
	}

	EventPasswordResetFailed = ILogEvent{
		Code:        "AUTH-PWD-004",
		Component:   ComponentAuthPassword,
		Message:     "Password reset failed",
		Level:       LevelWarn,
		Description: "Password reset attempt failed",
	}

	EventWeakPasswordDetected = ILogEvent{
		Code:        "AUTH-PWD-005",
		Component:   ComponentAuthPassword,
		Message:     "Weak password detected",
		Level:       LevelWarn,
		Description: "User attempted to set weak password",
	}
)

// Authorization Events (AUTH-AZ-xxx)
var (
	EventUnauthorizedAccess = ILogEvent{
		Code:        "AUTH-AZ-001",
		Component:   ComponentAuthJWT,
		Message:     "Unauthorized access attempt",
		Level:       LevelWarn,
		Description: "User attempted to access unauthorized resource",
	}

	EventPermissionDenied = ILogEvent{
		Code:        "AUTH-AZ-002",
		Component:   ComponentAuthJWT,
		Message:     "Permission denied",
		Level:       LevelWarn,
		Description: "User lacks required permissions",
	}

	EventRoleChanged = ILogEvent{
		Code:        "AUTH-AZ-003",
		Component:   ComponentAuthJWT,
		Message:     "User role changed",
		Level:       LevelInfo,
		Description: "User role/permissions updated",
	}
)
