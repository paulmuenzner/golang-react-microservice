// shared/go/utils/db/migrator/main.go

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	interfaces "github.com/app/shared/go/interfaces/config"
	db "github.com/app/shared/go/utils/db"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// DatabaseType repräsentiert verschiedene Datenbank-Typen
type DatabaseType string

const (
	PostgreSQL DatabaseType = "postgresql"
	VectorDB   DatabaseType = "vectordb"
	GraphDB    DatabaseType = "graphdb"
)

// MigrationConfig hält die Konfiguration für Migrationen
type MigrationConfig struct {
	MigrationsRoot string
	EnabledDBs     []DatabaseType
	Direction      string // "up" oder "down"
}

func main() {
	logger := log.Default()
	logger.Println("🚀 Starting Database Migrator...")
	logger.Println(strings.Repeat("=", 60))

	// Migration-Richtung aus ENV (default: up)
	direction := strings.ToLower(getEnv("MIGRATION_DIRECTION", "up"))
	if direction != "up" && direction != "down" {
		logger.Fatalf("❌ Invalid MIGRATION_DIRECTION: %s (must be 'up' or 'down')", direction)
	}

	// Migration-Konfiguration
	config := MigrationConfig{
		MigrationsRoot: getEnv("MIGRATIONS_ROOT", "/build/migrations"),
		EnabledDBs:     []DatabaseType{PostgreSQL},
		Direction:      direction,
	}

	// PostgreSQL Migration
	if contains(config.EnabledDBs, PostgreSQL) {
		logger.Printf("\n📊 PostgreSQL Migration (%s)", strings.ToUpper(config.Direction))
		logger.Println(strings.Repeat("=", 60))

		if err := migratePostgreSQL(config, logger); err != nil {
			logger.Fatalf("❌ PostgreSQL migration failed: %v", err)
		}
		logger.Println("\n✅ PostgreSQL migration completed successfully")
	}

	logger.Println("\n" + strings.Repeat("=", 60))
	logger.Println("🎉 ALL migrations completed successfully!")
}

// migratePostgreSQL führt PostgreSQL-Migrationen durch
func migratePostgreSQL(config MigrationConfig, logger *log.Logger) error {
	// DB-Konfiguration aus ENV
	port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	maxOpen, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	maxIdle, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "5"))

	cfg := interfaces.IPostgreConnectionConfig{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            port,
		User:            getEnv("DB_USER", "postgres"),
		Password:        getEnv("DB_PASSWORD", ""),
		DBName:          getEnv("DB_NAME", "postgres"),
		SSLMode:         getEnv("DB_SSLMODE", "disable"),
		Timeout:         30 * time.Second,
		MaxOpenConns:    maxOpen,
		MaxIdleConns:    maxIdle,
		ConnMaxIdleTime: 5 * time.Minute,
		ConnMaxLifetime: 30 * time.Minute,
	}

	logger.Printf("📝 Configuration:")
	logger.Printf("   Host:      %s", cfg.Host)
	logger.Printf("   Port:      %d", cfg.Port)
	logger.Printf("   User:      %s", cfg.User)
	logger.Printf("   Database:  %s", cfg.DBName)
	logger.Printf("   SSL Mode:  %s", cfg.SSLMode)
	logger.Printf("   Direction: %s", strings.ToUpper(config.Direction))
	logger.Printf("")

	// Verbindung mit Retry-Logik
	dbConn, err := retryConnect(cfg, 15, 3*time.Second, logger)
	if err != nil {
		return fmt.Errorf("DB connection failed: %w", err)
	}
	defer dbConn.Close()

	// Erstelle Migrations-Tracking-Tabelle
	if err := createMigrationsTable(dbConn, logger); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Finde automatisch alle Service-Migrations
	migrationPaths, err := discoverMigrationPaths(config.MigrationsRoot)
	if err != nil {
		return fmt.Errorf("failed to discover migration paths: %w", err)
	}

	if len(migrationPaths) == 0 {
		logger.Println("⚠️  No migration directories found")
		logger.Printf("   Searched in: %s", config.MigrationsRoot)
		return nil
	}

	logger.Printf("\n📁 Found %d migration directory(s):\n", len(migrationPaths))
	for _, path := range migrationPaths {
		logger.Printf("   - %s", path)
	}
	logger.Println()

	// Führe Migrationen für jeden Service aus
	for i, path := range migrationPaths {
		serviceName := extractServiceName(path)

		if config.Direction == "up" {
			logger.Printf("[%d/%d] ⬆️  Migrating UP: %s", i+1, len(migrationPaths), serviceName)
		} else {
			logger.Printf("[%d/%d] ⬇️  Migrating DOWN: %s", i+1, len(migrationPaths), serviceName)
		}
		logger.Println(strings.Repeat("-", 60))

		if config.Direction == "up" {
			if err := runMigrationsUp(dbConn, path, logger); err != nil {
				return fmt.Errorf("migration failed for %s: %w", serviceName, err)
			}
		} else {
			if err := runMigrationsDown(dbConn, path, logger); err != nil {
				return fmt.Errorf("rollback failed for %s: %w", serviceName, err)
			}
		}

		logger.Printf("✅ Service '%s' migrated successfully\n", serviceName)
	}

	return nil
}

// discoverMigrationPaths findet automatisch alle Migrations-Verzeichnisse
func discoverMigrationPaths(root string) ([]string, error) {
	var paths []string

	if _, err := os.Stat(root); os.IsNotExist(err) {
		return paths, fmt.Errorf("migrations root does not exist: %s", root)
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			upFiles, _ := filepath.Glob(filepath.Join(path, "*.up.sql"))
			downFiles, _ := filepath.Glob(filepath.Join(path, "*.down.sql"))
			sqlFiles, _ := filepath.Glob(filepath.Join(path, "*.sql"))

			if len(upFiles) > 0 || len(downFiles) > 0 || len(sqlFiles) > 0 {
				paths = append(paths, path)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Strings(paths)
	return paths, nil
}

// extractServiceName extrahiert den Service-Namen aus dem Pfad
func extractServiceName(path string) string {
	parts := strings.Split(filepath.Clean(path), string(filepath.Separator))

	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == "migrations" && i > 0 {
			return parts[i-1]
		}
	}

	return filepath.Base(path)
}

// createMigrationsTable erstellt die Tracking-Tabelle
func createMigrationsTable(db *sql.DB, logger *log.Logger) error {
	logger.Println("📋 Creating schema_migrations table if not exists...")

	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			version BIGINT NOT NULL,
			service VARCHAR(100) NOT NULL,
			description TEXT,
			direction VARCHAR(10) NOT NULL DEFAULT 'up',
			dirty BOOLEAN NOT NULL DEFAULT FALSE,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(version, service, direction)
		);
		
		CREATE INDEX IF NOT EXISTS idx_migrations_service ON schema_migrations(service);
		CREATE INDEX IF NOT EXISTS idx_migrations_version ON schema_migrations(version);
		CREATE INDEX IF NOT EXISTS idx_migrations_applied ON schema_migrations(applied_at);
	`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	logger.Println("✅ schema_migrations table ready")
	return nil
}

// runMigrationsUp führt alle ausstehenden UP-Migrationen aus
func runMigrationsUp(db *sql.DB, migrationsPath string, logger *log.Logger) error {
	serviceName := extractServiceName(migrationsPath)

	// Lese Migration-Dateien
	upFiles, err := filepath.Glob(filepath.Join(migrationsPath, "*.up.sql"))
	if err != nil {
		return fmt.Errorf("failed to read .up.sql files: %w", err)
	}

	var files []string
	if len(upFiles) > 0 {
		files = upFiles
		logger.Printf("   Found %d .up.sql migration(s)", len(files))
	} else {
		// Fallback: normale .sql Dateien (legacy)
		files, err = filepath.Glob(filepath.Join(migrationsPath, "*.sql"))
		if err != nil {
			return fmt.Errorf("failed to read .sql files: %w", err)
		}
		// Filter out .down.sql files
		var filtered []string
		for _, f := range files {
			if !strings.HasSuffix(f, ".down.sql") {
				filtered = append(filtered, f)
			}
		}
		files = filtered
		logger.Printf("   Found %d .sql migration(s)", len(files))
	}

	if len(files) == 0 {
		logger.Printf("   ℹ️  No migration files in %s", migrationsPath)
		return nil
	}

	sort.Strings(files)

	appliedCount := 0
	skippedCount := 0

	for _, file := range files {
		version, description, _, err := parseMigrationFilename(filepath.Base(file))
		if err != nil {
			logger.Printf("   ⚠️  Skipping invalid filename: %s (%v)", filepath.Base(file), err)
			continue
		}

		// Prüfe ob bereits ausgeführt
		var exists bool
		err = db.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1 AND service = $2 AND direction = 'up')",
			version, serviceName,
		).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if exists {
			logger.Printf("   ⏭️  [v%03d] %s (already applied)", version, description)
			skippedCount++
			continue
		}

		// Lese SQL-Datei
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file: %w", err)
		}

		// Führe Migration in Transaction aus
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", filepath.Base(file), err)
		}

		_, err = tx.Exec(
			"INSERT INTO schema_migrations (version, service, description, direction) VALUES ($1, $2, $3, 'up')",
			version, serviceName, description,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		logger.Printf("   ✅ [v%03d] %s", version, description)
		appliedCount++
	}

	if appliedCount > 0 {
		logger.Printf("   📊 Applied: %d, Skipped: %d", appliedCount, skippedCount)
	} else {
		logger.Printf("   📊 All migrations already applied (skipped: %d)", skippedCount)
	}

	return nil
}

// runMigrationsDown führt Rollback durch (DOWN-Migrationen)
func runMigrationsDown(db *sql.DB, migrationsPath string, logger *log.Logger) error {
	serviceName := extractServiceName(migrationsPath)

	// Lese DOWN-Migration-Dateien
	downFiles, err := filepath.Glob(filepath.Join(migrationsPath, "*.down.sql"))
	if err != nil {
		return fmt.Errorf("failed to read .down.sql files: %w", err)
	}

	if len(downFiles) == 0 {
		logger.Printf("   ⚠️  No .down.sql files found - cannot rollback")
		return nil
	}

	logger.Printf("   Found %d .down.sql migration(s)", len(downFiles))

	// Sortiere ABSTEIGEND (neueste zuerst für Rollback)
	sort.Sort(sort.Reverse(sort.StringSlice(downFiles)))

	// Hole alle angewendeten UP-Migrationen für diesen Service
	rows, err := db.Query(
		"SELECT version FROM schema_migrations WHERE service = $1 AND direction = 'up' ORDER BY version DESC",
		serviceName,
	)
	if err != nil {
		return fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	var appliedVersions []int64
	for rows.Next() {
		var version int64
		if err := rows.Scan(&version); err != nil {
			return fmt.Errorf("failed to scan version: %w", err)
		}
		appliedVersions = append(appliedVersions, version)
	}

	if len(appliedVersions) == 0 {
		logger.Printf("   ℹ️  No migrations to rollback for %s", serviceName)
		return nil
	}

	logger.Printf("   Found %d applied migration(s) to potentially rollback", len(appliedVersions))

	rolledBackCount := 0
	skippedCount := 0

	// Rollback nur die neueste Migration (oder alle mit ENV-Flag)
	rollbackAll := strings.ToLower(getEnv("ROLLBACK_ALL", "false")) == "true"
	rollbackSteps, _ := strconv.Atoi(getEnv("ROLLBACK_STEPS", "1"))

	for i, file := range downFiles {
		if !rollbackAll && i >= rollbackSteps {
			break
		}

		version, description, _, err := parseMigrationFilename(filepath.Base(file))
		if err != nil {
			logger.Printf("   ⚠️  Skipping invalid filename: %s (%v)", filepath.Base(file), err)
			continue
		}

		// Prüfe ob diese Version überhaupt angewendet wurde
		found := false
		for _, v := range appliedVersions {
			if v == version {
				found = true
				break
			}
		}

		if !found {
			logger.Printf("   ⏭️  [v%03d] %s (not applied, skipping)", version, description)
			skippedCount++
			continue
		}

		// Prüfe ob DOWN bereits ausgeführt wurde
		var downExists bool
		err = db.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1 AND service = $2 AND direction = 'down')",
			version, serviceName,
		).Scan(&downExists)
		if err != nil {
			return fmt.Errorf("failed to check down migration status: %w", err)
		}

		if downExists {
			logger.Printf("   ⏭️  [v%03d] %s (already rolled back)", version, description)
			skippedCount++
			continue
		}

		// Lese DOWN-SQL-Datei
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read down migration file: %w", err)
		}

		// Führe Rollback in Transaction aus
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		// Führe DOWN-SQL aus
		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute down migration %s: %w", filepath.Base(file), err)
		}

		// Markiere DOWN als ausgeführt UND entferne UP
		_, err = tx.Exec(
			"DELETE FROM schema_migrations WHERE version = $1 AND service = $2 AND direction = 'up'",
			version, serviceName,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to remove up migration record: %w", err)
		}

		_, err = tx.Exec(
			"INSERT INTO schema_migrations (version, service, description, direction) VALUES ($1, $2, $3, 'down')",
			version, serviceName, description,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record down migration: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		logger.Printf("   ⬇️  [v%03d] %s (rolled back)", version, description)
		rolledBackCount++
	}

	if rolledBackCount > 0 {
		logger.Printf("   📊 Rolled back: %d, Skipped: %d", rolledBackCount, skippedCount)
	} else {
		logger.Printf("   📊 No migrations rolled back (skipped: %d)", skippedCount)
	}

	return nil
}

// parseMigrationFilename extrahiert Version, Beschreibung und Richtung
func parseMigrationFilename(filename string) (int64, string, string, error) {
	name := strings.TrimSuffix(filename, ".sql")

	var direction string
	if strings.HasSuffix(name, ".up") {
		direction = "up"
		name = strings.TrimSuffix(name, ".up")
	} else if strings.HasSuffix(name, ".down") {
		direction = "down"
		name = strings.TrimSuffix(name, ".down")
	} else {
		direction = "up"
	}

	parts := strings.SplitN(name, "_", 2)
	if len(parts) < 2 {
		return 0, "", "", fmt.Errorf("invalid format, expected: NNN_description[.up|.down].sql")
	}

	version, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, "", "", fmt.Errorf("invalid version number: %w", err)
	}

	description := strings.ReplaceAll(parts[1], "_", " ")

	return version, description, direction, nil
}

// retryConnect versucht mehrfach eine Datenbankverbindung herzustellen
func retryConnect(cfg interfaces.IPostgreConnectionConfig, attempts int, delay time.Duration, logger *log.Logger) (*sql.DB, error) {
	logger.Printf("🔌 Attempting to connect (max %d attempts)...", attempts)
	logger.Printf("   Connection string: postgresql://%s:***@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	for i := 0; i < attempts; i++ {
		dbConn, err := db.NewPostgresDB(cfg)
		if err == nil {
			if err := dbConn.Ping(); err == nil {
				logger.Printf("✅ Connected successfully (attempt %d/%d)", i+1, attempts)
				return dbConn, nil
			}
			dbConn.Close()
		}

		if i < attempts-1 {
			logger.Printf("⚠️  Connection failed (attempt %d/%d). Retrying in %v...", i+1, attempts, delay)
			if err != nil {
				logger.Printf("   Error: %v", err)
			}
			time.Sleep(delay)
		}
	}

	return nil, fmt.Errorf("failed to connect after %d attempts", attempts)
}

// Helper functions

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func contains(slice []DatabaseType, item DatabaseType) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
