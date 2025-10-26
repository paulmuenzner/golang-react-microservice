// shared/go/utils/logger/events_database.go

package logger

// ==========================================
// DATABASE EVENTS (DB prefix)
// ==========================================

// Connection Events (DB-CON-xxx)
var (
	EventDBConnected = ILogEvent{
		Code:        "DB-CON-001",
		Component:   ComponentDatabaseConnection,
		Message:     "Database connected",
		Level:       LevelInfo,
		Description: "Successfully connected to database",
	}

	EventDBConnectionFailed = ILogEvent{
		Code:        "DB-CON-002",
		Component:   ComponentDatabaseConnection,
		Message:     "Database connection failed",
		Level:       LevelError,
		Description: "Failed to establish database connection",
	}

	EventDBDisconnected = ILogEvent{
		Code:        "DB-CON-003",
		Component:   ComponentDatabaseConnection,
		Message:     "Database disconnected",
		Level:       LevelWarn,
		Description: "Database connection lost",
	}

	EventDBReconnecting = ILogEvent{
		Code:        "DB-CON-004",
		Component:   ComponentDatabaseConnection,
		Message:     "Database reconnecting",
		Level:       LevelInfo,
		Description: "Attempting to reconnect to database",
	}

	EventDBConnectionPoolExhausted = ILogEvent{
		Code:        "DB-CON-005",
		Component:   ComponentDatabaseConnection,
		Message:     "Connection pool exhausted",
		Level:       LevelError,
		Description: "Database connection pool has no available connections",
	}
)

// Query Events (DB-QRY-xxx)
var (
	EventDBQuerySlow = ILogEvent{
		Code:        "DB-QRY-001",
		Component:   ComponentDatabasePerformance,
		Message:     "Slow query detected",
		Level:       LevelWarn,
		Description: "Database query exceeded performance threshold",
	}

	EventDBQueryFailed = ILogEvent{
		Code:        "DB-QRY-002",
		Component:   ComponentDatabaseQuery,
		Message:     "Query failed",
		Level:       LevelError,
		Description: "Database query execution failed",
	}

	EventDBQueryTimeout = ILogEvent{
		Code:        "DB-QRY-003",
		Component:   ComponentDatabaseQuery,
		Message:     "Query timeout",
		Level:       LevelError,
		Description: "Database query exceeded timeout limit",
	}

	EventDBDeadlockDetected = ILogEvent{
		Code:        "DB-QRY-004",
		Component:   ComponentDatabaseQuery,
		Message:     "Deadlock detected",
		Level:       LevelError,
		Description: "Database deadlock encountered",
	}
)

// Migration Events (DB-MIG-xxx)
var (
	EventDBMigrationStarted = ILogEvent{
		Code:        "DB-MIG-001",
		Component:   ComponentDatabaseMigration,
		Message:     "Migration started",
		Level:       LevelInfo,
		Description: "Database migration process started",
	}

	EventDBMigrationCompleted = ILogEvent{
		Code:        "DB-MIG-002",
		Component:   ComponentDatabaseMigration,
		Message:     "Migration completed",
		Level:       LevelInfo,
		Description: "Database migration successfully completed",
	}

	EventDBMigrationFailed = ILogEvent{
		Code:        "DB-MIG-003",
		Component:   ComponentDatabaseMigration,
		Message:     "Migration failed",
		Level:       LevelError,
		Description: "Database migration failed",
	}

	EventDBMigrationRolledBack = ILogEvent{
		Code:        "DB-MIG-004",
		Component:   ComponentDatabaseMigration,
		Message:     "Migration rolled back",
		Level:       LevelWarn,
		Description: "Database migration was rolled back",
	}
)
