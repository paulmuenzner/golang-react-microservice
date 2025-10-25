// shared/go/utils/logger/event_helpers.go

package logger

import "time"

// ==========================================
// EVENT LOGGING HELPERS
// ==========================================

// LogEvent logs a structured event with standardized fields
func LogEvent(event ILogEvent, fields map[string]interface{}) {
	if AppLogger == nil {
		return
	}

	// Add standard event fields
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["event_code"] = event.Code
	fields["component"] = event.Component

	// Log at appropriate level
	switch event.Level {
	case LevelDebug:
		AppLogger.DebugWithFields(event.Message, fields)
	case LevelInfo:
		AppLogger.InfoWithFields(event.Message, fields)
	case LevelWarn:
		AppLogger.WarnWithFields(event.Message, fields)
	case LevelError:
		if err, ok := fields["error"].(error); ok {
			AppLogger.ErrorWithFields(event.Message, err, fields)
		} else {
			AppLogger.ErrorWithFields(event.Message, nil, fields)
		}
	}
}

// LogEventWithContext logs an event with request context
func LogEventWithContext(event ILogEvent, requestID, ip string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}

	fields["request_id"] = requestID
	fields["ip"] = ip

	LogEvent(event, fields)
}

// ==========================================
// MIDDLEWARE-SPECIFIC HELPERS
// ==========================================

// LogMiddlewareEvent logs middleware events with standard HTTP fields
func LogMiddlewareEvent(event ILogEvent, requestID, ip, method, path, userAgent string, additionalFields map[string]interface{}) {
	fields := map[string]interface{}{
		"request_id": requestID,
		"ip":         ip,
		"method":     method,
		"path":       path,
		"user_agent": userAgent,
	}

	// Merge additional fields
	for k, v := range additionalFields {
		fields[k] = v
	}

	LogEvent(event, fields)
}

// LogMiddlewareEventWithIPContext logs middleware events with detailed IP context
func LogMiddlewareEventWithIPContext(event ILogEvent, requestID string, ipContext map[string]string, method, path, userAgent string, additionalFields map[string]interface{}) {
	fields := map[string]interface{}{
		"request_id": requestID,
		"ip":         ipContext["ip"],
		"ip_context": ipContext,
		"method":     method,
		"path":       path,
		"user_agent": userAgent,
	}

	for k, v := range additionalFields {
		fields[k] = v
	}

	LogEvent(event, fields)
}

// ==========================================
// AUTH-SPECIFIC HELPERS
// ==========================================

// LogAuthEvent logs authentication/authorization events
func LogAuthEvent(event ILogEvent, userID, requestID, ip string, additionalFields map[string]interface{}) {
	fields := map[string]interface{}{
		"user_id":    userID,
		"request_id": requestID,
		"ip":         ip,
	}

	for k, v := range additionalFields {
		fields[k] = v
	}

	LogEvent(event, fields)
}

// ==========================================
// DATABASE-SPECIFIC HELPERS
// ==========================================

// LogDatabaseEvent logs database events with query context
func LogDatabaseEvent(event ILogEvent, query string, duration time.Duration, additionalFields map[string]interface{}) {
	fields := map[string]interface{}{
		"query":       query,
		"duration_ms": duration.Milliseconds(),
	}

	for k, v := range additionalFields {
		fields[k] = v
	}

	LogEvent(event, fields)
}

// LogDatabaseQuerySlow logs slow query with automatic threshold
func LogDatabaseQuerySlow(query string, duration time.Duration, threshold time.Duration, additionalFields map[string]interface{}) {
	if duration > threshold {
		fields := map[string]interface{}{
			"query":        query,
			"duration_ms":  duration.Milliseconds(),
			"threshold_ms": threshold.Milliseconds(),
			"exceeded_by":  (duration - threshold).Milliseconds(),
		}

		for k, v := range additionalFields {
			fields[k] = v
		}

		LogEvent(EventDBQuerySlow, fields)
	}
}

// ==========================================
// SERVICE-SPECIFIC HELPERS
// ==========================================

// LogServiceEvent logs service lifecycle events
func LogServiceEvent(event ILogEvent, serviceName, version string, additionalFields map[string]interface{}) {
	fields := map[string]interface{}{
		"service_name": serviceName,
		"version":      version,
	}

	for k, v := range additionalFields {
		fields[k] = v
	}

	LogEvent(event, fields)
}

// ==========================================
// API-SPECIFIC HELPERS
// ==========================================

// LogAPIEvent logs API/business logic events
func LogAPIEvent(event ILogEvent, requestID, endpoint string, additionalFields map[string]interface{}) {
	fields := map[string]interface{}{
		"request_id": requestID,
		"endpoint":   endpoint,
	}

	for k, v := range additionalFields {
		fields[k] = v
	}

	LogEvent(event, fields)
}

// LogExternalAPICall logs external API integration events
func LogExternalAPICall(event ILogEvent, externalService, endpoint string, duration time.Duration, statusCode int, additionalFields map[string]interface{}) {
	fields := map[string]interface{}{
		"external_service": externalService,
		"endpoint":         endpoint,
		"duration_ms":      duration.Milliseconds(),
		"status_code":      statusCode,
	}

	for k, v := range additionalFields {
		fields[k] = v
	}

	LogEvent(event, fields)
}
