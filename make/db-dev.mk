# make/db.mk
# Database Management Commands
# Include this in main Makefile with: include make/db.mk

.PHONY: db-migrate db-rollback db-rollback-all db-status db-tables db-schemas \
        db-connect db-logs db-clean db-reset db-backup db-restore \
        db-tablespaces db-storage-report

# ==========================================
# Migration Commands
# ==========================================

db-migrate: ## Run all database migrations (UP)
	@echo "════════════════════════════════════════════════════════════════"
	@echo "⬆️  Starting Database Migration (UP)"
	@echo "════════════════════════════════════════════════════════════════"
	@echo "📝 Configuration:"
	@echo "   DB_USER:     $(POSTGRES_USER)"
	@echo "   DB_NAME:     $(POSTGRES_NAME)"
	@echo "   DB_PORT:     $(PORT_POSTGRES)"
	@echo "   Direction:   UP"
	@echo ""
	@echo "🔍 Checking for port conflicts..."
	@if lsof -Pi :$(PORT_POSTGRES) -sTCP:LISTEN -t >/dev/null 2>&1 || \
	   ss -tuln 2>/dev/null | grep -q ":$(PORT_POSTGRES) " || \
	   netstat -tuln 2>/dev/null | grep -q ":$(PORT_POSTGRES) "; then \
		if podman ps --filter "name=$(CONTAINER_POSTGRES)" --format "{{.Names}}" | grep -q "$(CONTAINER_POSTGRES)"; then \
			echo "⚠️  Development database is running on port $(PORT_POSTGRES)"; \
			echo "🔌 Connecting to running database..."; \
			echo ""; \
			echo "🔨 Building migrator image..."; \
			podman build -f shared/go/utils/db/migrator/Dockerfile.migrator \
				--build-arg GO_VERSION=$(GO_VERSION) \
				-t golang-react-microservice_db-migrator:latest . 2>&1 | grep -E "(STEP|COMMIT|Successfully)" || true; \
			echo ""; \
			echo "🐳 Starting migrator container..."; \
			podman rm -f $(CONTAINER_MIGRATOR) 2>/dev/null || true; \
			podman run --rm \
				--name $(CONTAINER_MIGRATOR) \
				--network $(NETWORK_BACKEND) \
				-e DB_HOST=postgres \
				-e DB_PORT=5432 \
				-e DB_USER=$(POSTGRES_USER) \
				-e DB_PASSWORD=$(POSTGRES_PASSWORD) \
				-e DB_NAME=$(POSTGRES_NAME) \
				-e DB_SSLMODE=$(DB_SSLMODE_DEV) \
				-e DB_MAX_OPEN_CONNS=$(DB_MAX_OPEN_CONNS) \
				-e DB_MAX_IDLE_CONNS=$(DB_MAX_IDLE_CONNS) \
				-e MIGRATION_DIRECTION=up \
				-e MIGRATIONS_ROOT=/build/migrations \
				localhost/golang-react-microservice_db-migrator:latest; \
			MIGRATE_EXIT=$$?; \
			if [ $$MIGRATE_EXIT -ne 0 ]; then \
				echo ""; \
				echo "════════════════════════════════════════════════════════════════"; \
				echo "❌ Migration FAILED (exit code: $$MIGRATE_EXIT)"; \
				echo "════════════════════════════════════════════════════════════════"; \
				echo ""; \
				echo "Check logs above for errors."; \
				exit 1; \
			fi; \
		else \
			echo "❌ Port $(PORT_POSTGRES) is in use by another process!"; \
			echo ""; \
			echo "Check what's using the port:"; \
			echo "  lsof -i :$(PORT_POSTGRES)"; \
			echo "  ss -tuln | grep $(PORT_POSTGRES)"; \
			echo ""; \
			exit 1; \
		fi; \
	else \
		echo "✅ Port $(PORT_POSTGRES) is available"; \
		echo ""; \
		echo "🐳 Starting PostgreSQL and running migrations..."; \
		MIGRATION_DIRECTION=up GO_VERSION=$(GO_VERSION) \
			podman-compose -f docker-compose.migrate.yml up --build --abort-on-container-exit 2>&1; \
		MIGRATE_EXIT=$$?; \
		echo ""; \
		echo "🧹 Cleaning up..."; \
		podman-compose -f docker-compose.migrate.yml down >/dev/null 2>&1 || true; \
		echo ""; \
		if [ $$MIGRATE_EXIT -ne 0 ]; then \
			echo "════════════════════════════════════════════════════════════════"; \
			echo "❌ Migration FAILED (exit code: $$MIGRATE_EXIT)"; \
			echo "════════════════════════════════════════════════════════════════"; \
			echo ""; \
			echo "Common issues:"; \
			echo "  • Migration file has syntax errors"; \
			echo "  • Database connection failed"; \
			echo "  • Migration already applied"; \
			echo ""; \
			echo "Check logs above for details."; \
			exit 1; \
		fi; \
	fi
	@echo "════════════════════════════════════════════════════════════════"
	@echo "✅ Migration complete!"
	@echo "════════════════════════════════════════════════════════════════"

db-rollback: ## Rollback last migration (DOWN)
	@echo "════════════════════════════════════════════════════════════════"
	@echo "⬇️  Starting Database Rollback (DOWN)"
	@echo "════════════════════════════════════════════════════════════════"
	@echo "⚠️  WARNING: This will rollback the last migration!"
	@echo ""
	@read -p "Are you sure? Type 'yes' to continue: " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		echo "📝 Configuration:"; \
		echo "   DB_USER:     $(POSTGRES_USER)"; \
		echo "   DB_NAME:     $(POSTGRES_NAME)"; \
		echo "   DB_PORT:     $(PORT_POSTGRES)"; \
		echo "   Direction:   DOWN"; \
		echo "   Steps:       1"; \
		echo ""; \
		echo "🔍 Checking for port conflicts..."; \
		if lsof -Pi :$(PORT_POSTGRES) -sTCP:LISTEN -t >/dev/null 2>&1 || \
		   ss -tuln 2>/dev/null | grep -q ":$(PORT_POSTGRES) " || \
		   netstat -tuln 2>/dev/null | grep -q ":$(PORT_POSTGRES) "; then \
			if podman ps --filter "name=$(CONTAINER_POSTGRES)" --format "{{.Names}}" | grep -q "$(CONTAINER_POSTGRES)"; then \
				echo "⚠️  Development database is running on port $(PORT_POSTGRES)"; \
				echo ""; \
				echo "Options:"; \
				echo "  1. Stop development environment first: make dev-down"; \
				echo "  2. Use the running database (will connect to it)"; \
				echo ""; \
				read -p "Connect to running database? (y/n): " use_running; \
				if [ "$$use_running" != "y" ]; then \
					echo "❌ Rollback cancelled. Please stop dev environment first."; \
					exit 1; \
				fi; \
				echo ""; \
				echo "🔌 Connecting to running database..."; \
				echo "🔨 Building migrator image..."; \
				podman build -f shared/go/utils/db/migrator/Dockerfile.migrator \
					--build-arg GO_VERSION=$(GO_VERSION) \
					-t golang-react-microservice_db-migrator:latest . 2>&1 | grep -E "(STEP|COMMIT|Successfully)" || true; \
				echo ""; \
				echo "🐳 Starting migrator container..."; \
				podman rm -f $(CONTAINER_MIGRATOR) 2>/dev/null || true; \
				podman run --rm \
					--name $(CONTAINER_MIGRATOR) \
					--network $(NETWORK_BACKEND) \
					-e DB_HOST=postgres \
					-e DB_PORT=$(PORT_POSTGRES) \
					-e DB_USER=$(POSTGRES_USER) \
					-e DB_PASSWORD=$(POSTGRES_PASSWORD) \
					-e DB_NAME=$(POSTGRES_NAME) \
					-e DB_SSLMODE=$(DB_SSLMODE_DEV) \
					-e DB_MAX_OPEN_CONNS=$(DB_MAX_OPEN_CONNS) \
					-e DB_MAX_IDLE_CONNS=$(DB_MAX_IDLE_CONNS) \
					-e MIGRATION_DIRECTION=down \
					-e ROLLBACK_STEPS=1 \
					-e MIGRATIONS_ROOT=/build/migrations \
					localhost/golang-react-microservice_db-migrator:latest; \
				ROLLBACK_EXIT=$$?; \
				if [ $$ROLLBACK_EXIT -ne 0 ]; then \
					echo ""; \
					echo "════════════════════════════════════════════════════════════════"; \
					echo "❌ Rollback FAILED (exit code: $$ROLLBACK_EXIT)"; \
					echo "════════════════════════════════════════════════════════════════"; \
					echo ""; \
					echo "Check logs above for errors."; \
					exit 1; \
				fi; \
			else \
				echo "❌ Port $(PORT_POSTGRES) is in use by another process!"; \
				echo ""; \
				echo "Check what's using the port:"; \
				echo "  lsof -i :$(PORT_POSTGRES)"; \
				echo "  ss -tuln | grep $(PORT_POSTGRES)"; \
				echo ""; \
				exit 1; \
			fi; \
		else \
			echo "✅ Port $(PORT_POSTGRES) is available"; \
			echo ""; \
			echo "🧹 Ensuring clean state..."; \
			podman-compose -f docker-compose.migrate.yml down 2>/dev/null || true; \
			echo ""; \
			echo "🐳 Starting PostgreSQL and running rollback..."; \
			MIGRATION_DIRECTION=down ROLLBACK_STEPS=1 GO_VERSION=$(GO_VERSION) \
				podman-compose -f docker-compose.migrate.yml up --build --abort-on-container-exit; \
			ROLLBACK_EXIT=$$?; \
			echo ""; \
			echo "🧹 Cleaning up migration containers..."; \
			podman-compose -f docker-compose.migrate.yml down 2>/dev/null || true; \
			echo ""; \
			if [ $$ROLLBACK_EXIT -ne 0 ]; then \
				echo "════════════════════════════════════════════════════════════════"; \
				echo "❌ Rollback FAILED (exit code: $$ROLLBACK_EXIT)"; \
				echo "════════════════════════════════════════════════════════════════"; \
				echo ""; \
				echo "Common issues:"; \
				echo "  • Migration file has syntax errors"; \
				echo "  • Database connection failed"; \
				echo "  • No migrations to rollback"; \
				echo ""; \
				echo "Check logs above for details."; \
				exit 1; \
			fi; \
		fi; \
		echo "════════════════════════════════════════════════════════════════"; \
		echo "✅ Rollback complete!"; \
		echo "════════════════════════════════════════════════════════════════"; \
	else \
		echo "❌ Rollback cancelled."; \
	fi

db-rollback-all: ## Rollback ALL migrations (DANGEROUS!)
	@echo "════════════════════════════════════════════════════════════════"
	@echo "🔥 DANGER: Rollback ALL Migrations"
	@echo "════════════════════════════════════════════════════════════════"
	@echo "⚠️  WARNING: This will rollback ALL migrations!"
	@echo "⚠️  This will DELETE all your database schema!"
	@echo ""
	@read -p "Type 'ROLLBACK ALL' to confirm: " confirm; \
	if [ "$$confirm" = "ROLLBACK ALL" ]; then \
		echo "📝 Configuration:"; \
		echo "   DB_USER:     $(POSTGRES_USER)"; \
		echo "   DB_NAME:     $(POSTGRES_NAME)"; \
		echo "   Direction:   DOWN"; \
		echo "   Mode:        ALL"; \
		echo ""; \
		echo "🧹 Ensuring clean state..."; \
		podman-compose -f docker-compose.migrate.yml down || true; \
		echo ""; \
		echo "🐳 Starting PostgreSQL and running full rollback..."; \
		MIGRATION_DIRECTION=down ROLLBACK_ALL=true GO_VERSION=$(GO_VERSION) \
			podman-compose -f docker-compose.migrate.yml up --build --abort-on-container-exit || true; \
		echo ""; \
		echo "🧹 Cleaning up migration containers..."; \
		podman-compose -f docker-compose.migrate.yml down || true; \
		echo ""; \
		echo "════════════════════════════════════════════════════════════════"; \
		echo "✅ Full rollback complete!"; \
		echo "════════════════════════════════════════════════════════════════"; \
	else \
		echo "❌ Rollback cancelled."; \
	fi

# ==========================================
# Status & Info Commands
# ==========================================

db-status: ## Show current migration status
	@echo "📊 Migration Status:"
	@echo "════════════════════════════════════════════════════════════════"
	@if ! podman ps --filter "name=$(CONTAINER_POSTGRES)" --format "{{.Names}}" | grep -q "$(CONTAINER_POSTGRES)"; then \
		echo "⚠️  Database container not running. Starting it..."; \
		podman-compose up -d postgres 2>/dev/null || true; \
		sleep 3; \
	fi
	@podman exec $(CONTAINER_POSTGRES) psql -U $(POSTGRES_USER) -d $(POSTGRES_NAME) -c \
		"SELECT version, service, description, direction, to_char(applied_at, 'YYYY-MM-DD HH24:MI:SS') as applied \
		 FROM schema_migrations \
		 ORDER BY service, version, applied_at DESC;" 2>/dev/null || \
		echo "❌ Could not connect. Try 'make dev' first."
	@echo "════════════════════════════════════════════════════════════════"

db-tables: ## Show all user database tables (excludes system tables)
	@echo "📚 Database Tables:"
	@echo "════════════════════════════════════════════════════════════════"
	@if ! podman ps --filter "name=$(CONTAINER_POSTGRES)" --format "{{.Names}}" | grep -q "$(CONTAINER_POSTGRES)"; then \
		echo "⚠️  Database container not running. Starting it..."; \
		podman-compose up -d postgres 2>/dev/null || true; \
		sleep 3; \
	fi
	@podman exec $(CONTAINER_POSTGRES) psql -U $(POSTGRES_USER) -d $(POSTGRES_NAME) -c \
		"SELECT schemaname as \"Schema\", tablename as \"Table\", \
		pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as \"Size\" \
		FROM pg_tables \
		WHERE schemaname NOT IN ('pg_catalog', 'information_schema', 'pg_toast') \
		ORDER BY schemaname, tablename;" 2>/dev/null || \
		echo "❌ Could not connect. Try 'make dev' first."
	@echo "════════════════════════════════════════════════════════════════"

db-tables-all: ## Show ALL tables including system tables
	@echo "📚 All Database Tables (including system):"
	@echo "════════════════════════════════════════════════════════════════"
	@if ! podman ps --filter "name=$(CONTAINER_POSTGRES)" --format "{{.Names}}" | grep -q "$(CONTAINER_POSTGRES)"; then \
		echo "⚠️  Database container not running. Starting it..."; \
		podman-compose up -d postgres 2>/dev/null || true; \
		sleep 3; \
	fi
	@podman exec $(CONTAINER_POSTGRES) psql -U $(POSTGRES_USER) -d $(POSTGRES_NAME) -c "\dt *.*" 2>/dev/null || \
		echo "❌ Could not connect."
	@echo "════════════════════════════════════════════════════════════════"

db-schemas: ## Show all schemas in database
	@echo "📂 Database Schemas:"
	@echo "════════════════════════════════════════════════════════════════"
	@if ! podman ps --filter "name=$(CONTAINER_POSTGRES)" --format "{{.Names}}" | grep -q "$(CONTAINER_POSTGRES)"; then \
		echo "⚠️  Database container not running. Starting it..."; \
		podman-compose up -d postgres 2>/dev/null || true; \
		sleep 3; \
	fi
	@podman exec $(CONTAINER_POSTGRES) psql -U $(POSTGRES_USER) -d $(POSTGRES_NAME) -c "\dn+" 2>/dev/null || \
		echo "❌ Could not connect."
	@echo "════════════════════════════════════════════════════════════════"

# ==========================================
# Connection Commands
# ==========================================

db-connect: ## Connect to PostgreSQL via psql
	@echo "🔌 Connecting to PostgreSQL..."
	@if ! podman ps --filter "name=$(CONTAINER_POSTGRES)" --format "{{.Names}}" | grep -q "$(CONTAINER_POSTGRES)"; then \
		echo "⚠️  Database container not running. Starting it..."; \
		podman-compose up -d postgres 2>/dev/null || true; \
		sleep 3; \
	fi
	@echo "════════════════════════════════════════════════════════════════"
	@echo "Connected to: $(POSTGRES_NAME)@$(CONTAINER_POSTGRES)"
	@echo "Useful commands:"
	@echo "  \\dt          - List tables"
	@echo "  \\dt+         - List tables with sizes"
	@echo "  \\dn+         - List schemas"
	@echo "  \\l           - List databases"
	@echo "  \\q           - Quit"
	@echo "════════════════════════════════════════════════════════════════"
	@podman exec -it $(CONTAINER_POSTGRES) psql -U $(POSTGRES_USER) -d $(POSTGRES_NAME) || \
		echo "❌ Could not connect. Try 'make dev' first."

db-logs: ## Show PostgreSQL logs
	@echo "📋 PostgreSQL Logs:"
	@echo "════════════════════════════════════════════════════════════════"
	@podman logs $(CONTAINER_POSTGRES) --tail 50 2>/dev/null || \
		echo "❌ Container not running. Try 'make dev' first."

db-logs-follow: ## Follow PostgreSQL logs
	@echo "📋 Following PostgreSQL Logs (Ctrl+C to stop):"
	@echo "════════════════════════════════════════════════════════════════"
	@podman logs -f $(CONTAINER_POSTGRES) 2>/dev/null || \
		echo "❌ Container not running. Try 'make dev' first."

# ==========================================
# Cleanup Commands
# ==========================================

db-clean: ## Remove database containers and volumes
	@echo "🧹 Cleaning database containers and volumes..."
	@echo "════════════════════════════════════════════════════════════════"
	@echo "⚠️  WARNING: This will delete all database data!"
	@echo ""
	@read -p "Type 'yes' to confirm: " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		echo "🗑️  Stopping containers..."; \
		podman-compose -f docker-compose.yml down -v 2>/dev/null || true; \
		podman-compose -f docker-compose.migrate.yml down -v 2>/dev/null || true; \
		echo "🗑️  Removing volume..."; \
		podman volume rm $(VOLUME_POSTGRES) 2>/dev/null || true; \
		echo "✅ Database cleaned!"; \
	else \
		echo "❌ Cleanup cancelled."; \
	fi
	@echo "════════════════════════════════════════════════════════════════"

db-reset: ## Reset database (clean + migrate)
	@echo "🔄 Resetting Database..."
	@echo "════════════════════════════════════════════════════════════════"
	@echo "⚠️  WARNING: This will delete all data and re-run migrations!"
	@echo ""
	@read -p "Type 'RESET' to confirm: " confirm; \
	if [ "$$confirm" = "RESET" ]; then \
		echo "🗑️  Cleaning database..."; \
		podman-compose -f docker-compose.yml down -v 2>/dev/null || true; \
		podman-compose -f docker-compose.migrate.yml down -v 2>/dev/null || true; \
		podman volume rm $(VOLUME_POSTGRES) 2>/dev/null || true; \
		echo ""; \
		echo "🚀 Running migrations..."; \
		$(MAKE) db-migrate; \
		echo ""; \
		echo "════════════════════════════════════════════════════════════════"; \
		echo "✅ Database reset complete!"; \
		echo "════════════════════════════════════════════════════════════════"; \
	else \
		echo "❌ Reset cancelled."; \
	fi

# ==========================================
# Backup & Restore Commands
# ==========================================

db-backup: ## Backup database to file (Usage: make db-backup [FILE=backup.sql])
	@BACKUP_FILE=$${FILE:-backups/db_backup_$$(date +%Y%m%d_%H%M%S).sql}; \
	echo "💾 Backing up database to $$BACKUP_FILE..."; \
	echo "════════════════════════════════════════════════════════════════"; \
	if ! podman ps --filter "name=$(CONTAINER_POSTGRES)" --format "{{.Names}}" | grep -q "$(CONTAINER_POSTGRES)"; then \
		echo "❌ Database container not running. Try 'make dev' first."; \
		exit 1; \
	fi; \
	mkdir -p $$(dirname $$BACKUP_FILE); \
	podman exec $(CONTAINER_POSTGRES) pg_dump -U $(POSTGRES_USER) $(POSTGRES_NAME) > $$BACKUP_FILE 2>/dev/null && \
		echo "✅ Backup saved to: $$BACKUP_FILE" || \
		echo "❌ Backup failed!"; \
	echo "════════════════════════════════════════════════════════════════"

db-restore: ## Restore database from file (Usage: make db-restore FILE=backup.sql)
	@if [ -z "$(FILE)" ]; then \
		echo "❌ Error: FILE parameter required"; \
		echo "Usage: make db-restore FILE=backup.sql"; \
		exit 1; \
	fi; \
	if [ ! -f "$(FILE)" ]; then \
		echo "❌ Error: File $(FILE) not found!"; \
		exit 1; \
	fi; \
	echo "📥 Restoring database from $(FILE)..."; \
	echo "════════════════════════════════════════════════════════════════"; \
	echo "⚠️  WARNING: This will overwrite current database!"; \
	echo ""; \
	read -p "Type 'yes' to confirm: " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		if ! podman ps --filter "name=$(CONTAINER_POSTGRES)" --format "{{.Names}}" | grep -q "$(CONTAINER_POSTGRES)"; then \
			echo "⚠️  Starting database container..."; \
			podman-compose up -d postgres 2>/dev/null || true; \
			sleep 3; \
		fi; \
		cat $(FILE) | podman exec -i $(CONTAINER_POSTGRES) psql -U $(POSTGRES_USER) $(POSTGRES_NAME) && \
			echo "✅ Database restored successfully!" || \
			echo "❌ Restore failed!"; \
	else \
		echo "❌ Restore cancelled."; \
	fi; \
	echo "════════════════════════════════════════════════════════════════"

# ==========================================
# Tablespace Commands (Multi-Storage)
# ==========================================

db-tablespaces: ## Show all tablespaces and their sizes
	@echo "🗄️  Database Tablespaces:"
	@echo "════════════════════════════════════════════════════════════════"
	@if ! podman ps --filter "name=$(CONTAINER_POSTGRES)" --format "{{.Names}}" | grep -q "$(CONTAINER_POSTGRES)"; then \
		echo "⚠️  Database container not running. Starting it..."; \
		podman-compose up -d postgres 2>/dev/null || true; \
		sleep 3; \
	fi
	@podman exec $(CONTAINER_POSTGRES) psql -U $(POSTGRES_USER) -d $(POSTGRES_NAME) -c \
		"SELECT spcname as \"Tablespace\", \
		pg_tablespace_location(oid) as \"Location\", \
		pg_size_pretty(pg_tablespace_size(spcname)) as \"Size\" \
		FROM pg_tablespace \
		ORDER BY spcname;" 2>/dev/null || \
		echo "❌ Could not connect."
	@echo "════════════════════════════════════════════════════════════════"

db-storage-report: ## Detailed storage usage report by tablespace
	@echo "📊 Storage Usage Report:"
	@echo "════════════════════════════════════════════════════════════════"
	@if ! podman ps --filter "name=$(CONTAINER_POSTGRES)" --format "{{.Names}}" | grep -q "$(CONTAINER_POSTGRES)"; then \
		echo "⚠️  Database container not running. Starting it..."; \
		podman-compose up -d postgres 2>/dev/null || true; \
		sleep 3; \
	fi
	@podman exec $(CONTAINER_POSTGRES) psql -U $(POSTGRES_USER) -d $(POSTGRES_NAME) -c \
		"SELECT \
		    schemaname as \"Schema\", \
		    tablename as \"Table\", \
		    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as \"Total Size\", \
		    pg_size_pretty(pg_relation_size(schemaname||'.'||tablename)) as \"Table Size\", \
		    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename) - pg_relation_size(schemaname||'.'||tablename)) as \"Indexes\", \
		    CASE \
		        WHEN reltablespace = 0 THEN 'pg_default (SSD)' \
		        ELSE (SELECT spcname FROM pg_tablespace WHERE oid = reltablespace) \
		    END as \"Tablespace\" \
		FROM pg_tables t \
		JOIN pg_class c ON t.tablename = c.relname \
		WHERE schemaname NOT IN ('pg_catalog', 'information_schema') \
		ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;" 2>/dev/null || \
		echo "❌ Could not connect."
	@echo ""
	@echo "📈 Total Size by Tablespace:"
	@podman exec $(CONTAINER_POSTGRES) psql -U $(POSTGRES_USER) -d $(POSTGRES_NAME) -c \
		"SELECT \
		    COALESCE( \
		        (SELECT spcname FROM pg_tablespace WHERE oid = reltablespace), \
		        'pg_default (SSD)' \
		    ) as \"Tablespace\", \
		    COUNT(*) as \"Tables\", \
		    pg_size_pretty(SUM(pg_total_relation_size(schemaname||'.'||tablename))) as \"Total Size\" \
		FROM pg_tables t \
		JOIN pg_class c ON t.tablename = c.relname \
		WHERE schemaname NOT IN ('pg_catalog', 'information_schema') \
		GROUP BY reltablespace \
		ORDER BY SUM(pg_total_relation_size(schemaname||'.'||tablename)) DESC;" 2>/dev/null
	@echo "════════════════════════════════════════════════════════════════"

db-move-to-hdd: ## Move a table to HDD (Usage: make db-move-to-hdd TABLE=schema.table)
	@if [ -z "$(TABLE)" ]; then \
		echo "❌ Error: TABLE parameter required"; \
		echo "Usage: make db-move-to-hdd TABLE=schema.table"; \
		exit 1; \
	fi; \
	echo "📦 Moving table $(TABLE) to HDD (cold_storage)..."; \
	podman exec $(CONTAINER_POSTGRES) psql -U $(POSTGRES_USER) -d $(POSTGRES_NAME) -c \
		"ALTER TABLE $(TABLE) SET TABLESPACE cold_storage;" 2>/dev/null && \
		echo "✅ Table $(TABLE) moved to HDD" || \
		echo "❌ Failed to move table. Does cold_storage tablespace exist?"

db-move-to-ssd: ## Move a table to SSD (Usage: make db-move-to-ssd TABLE=schema.table)
	@if [ -z "$(TABLE)" ]; then \
		echo "❌ Error: TABLE parameter required"; \
		echo "Usage: make db-move-to-ssd TABLE=schema.table"; \
		exit 1; \
	fi; \
	echo "⚡ Moving table $(TABLE) to SSD (pg_default)..."; \
	podman exec $(CONTAINER_POSTGRES) psql -U $(POSTGRES_USER) -d $(POSTGRES_NAME) -c \
		"ALTER TABLE $(TABLE) SET TABLESPACE pg_default;" 2>/dev/null && \
		echo "✅ Table $(TABLE) moved to SSD" || \
		echo "❌ Failed to move table"

# ==========================================
# Help
# ==========================================

db-help: ## Show all database commands
	@echo "════════════════════════════════════════════════════════════════"
	@echo "📚 Database Commands"
	@echo "════════════════════════════════════════════════════════════════"
	@echo ""
	@echo "🔧 Migration:"
	@echo "  make db-migrate          - Run all migrations (UP)"
	@echo "  make db-rollback         - Rollback last migration (DOWN)"
	@echo "  make db-rollback-all     - Rollback ALL migrations (DANGEROUS)"
	@echo ""
	@echo "📊 Status & Info:"
	@echo "  make db-status           - Show migration status"
	@echo "  make db-tables           - Show user tables"
	@echo "  make db-tables-all       - Show all tables (incl. system)"
	@echo "  make db-schemas          - Show all schemas"
	@echo ""
	@echo "🔌 Connection:"
	@echo "  make db-connect          - Connect to PostgreSQL via psql"
	@echo "  make db-logs             - Show PostgreSQL logs"
	@echo "  make db-logs-follow      - Follow PostgreSQL logs"
	@echo ""
	@echo "🧹 Cleanup:"
	@echo "  make db-clean            - Remove containers and volumes"
	@echo "  make db-reset            - Clean + migrate from scratch"
	@echo ""
	@echo "💾 Backup & Restore:"
	@echo "  make db-backup           - Backup database"
	@echo "  make db-backup FILE=x    - Backup to specific file"
	@echo "  make db-restore FILE=x   - Restore from file"
	@echo ""
	@echo "🗄️  Multi-Storage (Optional):"
	@echo "  make db-tablespaces      - Show tablespaces"
	@echo "  make db-storage-report   - Storage usage by tablespace"
	@echo "  make db-move-to-hdd      - Move table to HDD"
	@echo "  make db-move-to-ssd      - Move table to SSD"
	@echo ""
	@echo "════════════════════════════════════════════════════════════════"