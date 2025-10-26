# Logger Package - Event-Based Structured Logging

## Overview
This package provides structured logging with event codes for tracking, analytics, and monitoring.

## Event Structure
Every log event has:
- **Code**: Unique identifier (e.g., `MW-RL-001`)
- **Component**: Module/category (e.g., `middleware.ratelimit`)
- **Message**: Human-readable message
- **Level**: DEBUG, INFO, WARN, ERROR
- **Description**: Detailed explanation for documentation

## Event Categories

| Prefix | Category | File |
|--------|----------|------|
| `MW-*` | Middleware | `events_middleware.go` |
| `AUTH-*` | Authentication | `events_auth.go` |
| `SVC-*` | Service Lifecycle | `events_service.go` |
| `DB-*` | Database | `events_database.go` |
| `API-*` | API/Business Logic | `events_api.go` |

## Usage Examples

### Basic Event Logging
```go
logger.LogEvent(logger.EventRateLimitExceeded, map[string]interface{}{
    "request_id": requestID,
    "ip": clientIP,
    "rate_limit": 100,
})
```

### Middleware Event
```go
logger.LogMiddlewareEvent(
    logger.EventRequestTooLarge,
    requestID, ip, method, path, userAgent,
    map[string]interface{}{
        "content_length": contentLength,
        "max_bytes": maxBytes,
    },
)
```

### Authentication Event
```go
logger.LogAuthEvent(
    logger.EventAuthFailed,
    userID, requestID, ip,
    map[string]interface{}{
        "reason": "invalid_password",
    },
)
```

### Database Event
```go
logger.LogDatabaseQuerySlow(
    query,
    duration,
    200*time.Millisecond, // threshold
    map[string]interface{}{
        "table": "users",
    },
)
```

## Grafana/Loki Queries
```logql
// All rate limit events
{job="gateway"} | json | event_code="MW-RL-001"

// All middleware warnings
{job="gateway"} | json | component=~"middleware.*" | level="warn"

// Top IPs triggering rate limits (24h)
topk(10, sum by (ip) (count_over_time({job="gateway"} | json | event_code="MW-RL-001" [24h])))

// All security events
{job="gateway"} | json | event_code=~"MW-SEC-.*"

// Slow database queries
{job="gateway"} | json | event_code="DB-QRY-001"
```

## Adding New Events

1. Choose appropriate file based on category
2. Define event constant with unique code
3. Use helper function for logging
4. Update this README with the new event code

## Event Code Convention

Format: `CATEGORY-SUBCATEGORY-NUMBER`

Examples:
- `MW-RL-001` = Middleware > Rate Limiting > Event #1
- `AUTH-JWT-003` = Authentication > JWT > Event #3
- `DB-CON-002` = Database > Connection > Event #2