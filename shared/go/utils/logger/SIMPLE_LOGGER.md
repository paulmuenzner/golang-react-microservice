# Logger

Simple structured logging for microservices with request tracking support.

---

## Quick Start

### 1. Initialize

```go
package main

import (
    "os"
    "github.com/app/shared/go/utils/logger"
)

func main() {
    // Initialize once at startup
    logger.Init("SERVICE-A", os.Getenv("ENVIRONMENT"))
    
    logger.Info("Server started")
}
```

---

## Basic Logging

### Simple Messages

```go
logger.Debug("Debug information")              // Only in development
logger.Info("General information")              
logger.Warn("Warning message")                  
logger.Error("Error occurred", err)             
logger.Fatal("Critical error", err)             // Exits app!
```

### With Additional Fields

```go
logger.InfoWithFields("User logged in", map[string]interface{}{
    "user_id": "123",
    "ip":      "192.168.1.1",
})

logger.ErrorWithFields("Database error", err, map[string]interface{}{
    "query": "SELECT * FROM users",
    "duration_ms": 1500,
})
```

**Output:**
```json
{"level":"info","service":"SERVICE-A","user_id":"123","ip":"192.168.1.1","msg":"User logged in"}
```

---

## Request Tracking

### Setup

**Helper Function:**
```go
// GetRequestID extracts request ID from context
func GetRequestID(r *http.Request) string {
    if id := r.Context().Value("request_id"); id != nil {
        return id.(string)
    }
    return ""
}
```

**Middleware to add Request ID:**
```go
import (
    "context"
    "github.com/google/uuid"
)

func requestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get from header or generate
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        // Store in context
        ctx := context.WithValue(r.Context(), "request_id", requestID)
        r = r.WithContext(ctx)
        
        // Add to response
        w.Header().Set("X-Request-ID", requestID)
        
        next.ServeHTTP(w, r)
    })
}
```

### Usage in Handlers

```go
func usersHandler(w http.ResponseWriter, r *http.Request) {
    // Get request ID
    requestID := GetRequestID(r)
    
    // Create logger with request ID
    reqLog := logger.WithRequestID(requestID)
    
    // All logs now have request_id automatically!
    reqLog.Info("Fetching users from database")
    reqLog.InfoWithFields("Query completed", map[string]interface{}{
        "count": 10,
        "duration_ms": 45,
    })
    
    // Pass to other functions
    users := fetchUsers(reqLog)
}

func fetchUsers(log *logger.Logger) []User {
    log.Debug("Connecting to database")
    // ... database query ...
    log.Debug("Query executed")
    return users
}
```

**Output:**
```json
{"level":"info","request_id":"abc-123","msg":"Fetching users from database"}
{"level":"debug","request_id":"abc-123","msg":"Connecting to database"}
{"level":"debug","request_id":"abc-123","msg":"Query executed"}
{"level":"info","request_id":"abc-123","count":10,"duration_ms":45,"msg":"Query completed"}
```

---

## Contextual Logging

Create loggers with additional context:

```go
// With request ID
reqLog := logger.WithRequestID("abc-123")

// With user ID
userLog := logger.WithUserID("user-456")

// With single field
txLog := logger.WithField("transaction_id", "tx-789")

// With multiple fields
log := logger.WithFields(map[string]interface{}{
    "component": "auth",
    "version": "1.0",
})

// Chain them
log := logger.WithRequestID(requestID).WithField("user_id", userID)
log.Info("User authenticated")
```

---

## HTTP Request Logging

```go
import "time"

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        requestID := GetRequestID(r)
        reqLog := logger.WithRequestID(requestID)
        
        reqLog.InfoWithFields("Request received", map[string]interface{}{
            "method": r.Method,
            "path": r.URL.Path,
        })
        
        next.ServeHTTP(w, r)
        
        duration := time.Since(start)
        reqLog.HTTP(r.Method, r.URL.Path, 200, duration, r.RemoteAddr)
    })
}
```

---

## Complete Example

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "time"
    
    "github.com/app/shared/go/utils/logger"
    "github.com/google/uuid"
)

// Helper to get request ID from context
func GetRequestID(r *http.Request) string {
    if id := r.Context().Value("request_id"); id != nil {
        return id.(string)
    }
    return ""
}

// Middleware to add request ID
func requestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        ctx := context.WithValue(r.Context(), "request_id", requestID)
        r = r.WithContext(ctx)
        w.Header().Set("X-Request-ID", requestID)
        
        next.ServeHTTP(w, r)
    })
}

// Middleware to log requests
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        requestID := GetRequestID(r)
        reqLog := logger.WithRequestID(requestID)
        
        reqLog.InfoWithFields("Incoming request", map[string]interface{}{
            "method": r.Method,
            "path": r.URL.Path,
        })
        
        next.ServeHTTP(w, r)
        
        duration := time.Since(start)
        reqLog.HTTP(r.Method, r.URL.Path, 200, duration, r.RemoteAddr)
    })
}

// Example handler
func usersHandler(w http.ResponseWriter, r *http.Request) {
    requestID := GetRequestID(r)
    reqLog := logger.WithRequestID(requestID)
    
    reqLog.Info("Fetching users")
    
    // Simulate work
    time.Sleep(10 * time.Millisecond)
    
    reqLog.InfoWithFields("Users fetched", map[string]interface{}{
        "count": 5,
    })
    
    fmt.Fprint(w, "Users list")
}

func main() {
    // Initialize logger
    logger.Init("SERVICE-A", os.Getenv("ENVIRONMENT"))
    
    logger.Info("Starting server...")
    
    // Setup routes
    mux := http.NewServeMux()
    mux.HandleFunc("/users", usersHandler)
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, "OK")
    })
    
    // Apply middleware
    handler := requestIDMiddleware(loggingMiddleware(mux))
    
    // Start server
    logger.Info("Server listening on :8080")
    if err := http.ListenAndServe(":8080", handler); err != nil {
        logger.Fatal("Server failed", err)
    }
}
```

---

## API Reference

### Initialization

| Function | Description |
|----------|-------------|
| `Init(serviceName, environment string)` | Initialize logger (call once at startup) |

### Simple Logging

| Function | Description |
|----------|-------------|
| `Debug(msg string)` | Debug message (dev only) |
| `Info(msg string)` | Info message |
| `Warn(msg string)` | Warning message |
| `Error(msg string, err error)` | Error with error object |
| `Fatal(msg string, err error)` | Fatal error (exits app!) |

### Logging with Fields

| Function | Description |
|----------|-------------|
| `InfoWithFields(msg string, fields map[string]interface{})` | Info with structured data |
| `DebugWithFields(...)` | Debug with fields |
| `WarnWithFields(...)` | Warning with fields |
| `ErrorWithFields(msg string, err error, fields map[string]interface{})` | Error with fields |
| `FatalWithFields(...)` | Fatal with fields |

### Contextual Logging

| Function | Description |
|----------|-------------|
| `WithRequestID(id string) *Logger` | Create logger with request ID |
| `WithUserID(id string) *Logger` | Create logger with user ID |
| `WithField(key string, value interface{}) *Logger` | Create logger with single field |
| `WithFields(fields map[string]interface{}) *Logger` | Create logger with multiple fields |
| `WithContext(ctx context.Context) *Logger` | Create logger from context |

### HTTP Logging

| Function | Description |
|----------|-------------|
| `HTTP(method, path string, status int, duration time.Duration, clientIP string)` | Log HTTP request |
| `HTTPWithFields(method, path string, status int, duration time.Duration, fields map[string]interface{})` | HTTP request with custom fields |

---

## Best Practices

### ✅ DO

```go
// Initialize once at startup
logger.Init("SERVICE-A", os.Getenv("ENVIRONMENT"))

// Use contextual loggers for requests
reqLog := logger.WithRequestID(GetRequestID(r))
reqLog.Info("Processing request")

// Pass loggers to functions
func processData(log *logger.Logger) {
    log.Info("Processing started")
}

// Use appropriate log levels
logger.Debug("Cache miss")           // Development
logger.Info("User logged in")        // Important events
logger.Warn("Slow query")            // Warnings
logger.Error("DB failed", err)       // Errors
logger.Fatal("Config missing", err)  // Unrecoverable
```

### ❌ DON'T

```go
// Don't call Init multiple times
logger.Init("SERVICE-A", "dev")
logger.Init("SERVICE-A", "dev")  // ❌

// Don't log sensitive data
logger.InfoWithFields("User data", map[string]interface{}{
    "password": "secret",  // ❌ Never!
})

// Don't overuse Fatal (it exits the app!)
if err != nil {
    logger.Fatal("Minor error", err)  // ❌
}
```

---

## Log Levels

| Level | When to Use | Production |
|-------|-------------|------------|
| `Debug` | Detailed debugging info | ❌ Not shown |
| `Info` | Important events | ✅ Shown |
| `Warn` | Warning conditions | ✅ Shown |
| `Error` | Error conditions | ✅ Shown |
| `Fatal` | Unrecoverable errors | ✅ Shown + exits |

---

## Environment Configuration

```bash
# .env
ENVIRONMENT=development  # or "production"
```

**Development:**
- Debug level enabled
- Pretty console output (optional)

**Production:**
- Info level and above
- JSON output for log aggregation