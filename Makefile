# ==========================================
# Environment Variables (aus .env laden)
# ==========================================
include .env
export

.PHONY: help init dev dev-a dev-b dev-g test lint prod prod-up prod-stop stop clean delete stat info storage

.DEFAULT_GOAL := help

# ==========================================
# Utility Commands
# ==========================================

help: ## Show this help message
	@echo "════════════════════════════════════════════════════════════════"
	@echo "Available Commands:"
	@echo "════════════════════════════════════════════════════════════════"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2}'
	@echo ""

clean: ## Clean all containers, volumes, and images
	@echo "🧹 Cleaning up..."
	@podman-compose -f docker-compose.yml down -v
	@podman-compose -f docker-compose.migrate.yml down -v
	@podman-compose -f docker-compose.prod.yml down -v
	@echo "✅ Cleanup complete"

install-tools: ## Install development tools (golangci-lint, air)
	@echo "🔧 Installing development tools..."
	@which golangci-lint > /dev/null || { \
		echo "  → Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	}
	@which air > /dev/null || { \
		echo "  → Installing air..."; \
		go install github.com/air-verse/air@latest; \
	}
	@echo "✅ All tools installed!"


init: ## Initialize Go modules for all services
	@echo "🔧 Initializing (Go ${GO_VERSION})..."
	@cd shared/go && go mod tidy
	@cd app/backend/service-a && go mod tidy
	@cd app/backend/service-b && go mod tidy
	@cd app/backend/gateway && go mod tidy
	@echo "✅ Initialization complete!"

lint: ## Run golangci-lint on all Go code
	@echo "🔍 Running linters..."
	@cd shared/go && golangci-lint run
	@cd app/backend/service-a && golangci-lint run
	@cd app/backend/service-b && golangci-lint run
	@cd app/backend/gateway && golangci-lint run
	@echo "✅ Linting complete!"
	

.PHONY: dev dev-up dev-down dev-logs dev-logs-gateway dev-restart \
        db-migrate db-status db-tables db-connect db-migrations-list db-test db-clean \
        prod-build prod-migrate prod-backup prod-deploy prod-up prod-down prod-logs prod-status prod-restart prod-rebuild \
        help clean install-tools init lint

.DEFAULT_GOAL := help





# ==========================================
# Development Commands
# ==========================================

dev: db-migrate dev-up ## Start development environment (migration first, then services)
	@echo ""
	@echo "════════════════════════════════════════════════════════════════"
	@echo "✅ Development environment ready!"
	@echo "════════════════════════════════════════════════════════════════"
	@echo "📍 Gateway:  http://localhost:$(PORT_GATEWAY)"
	@echo "📍 Grafana:  http://localhost:$(PORT_GRAFANA) (admin/admin)"
	@echo ""

dev-up: ## Start development services (without migrations)
	@echo "🐳 Starting development services..."
	@podman-compose -f docker-compose.yml up -d
	@echo "⏳ Waiting for services to be healthy..."
	@sleep 5
	@echo "✅ Services running"

dev-down: ## Stop development services
	@echo "🛑 Stopping development services..."
	@podman-compose -f docker-compose.yml down
	@echo "✅ Services stopped"

dev-logs: ## Show logs from all services
	@podman-compose -f docker-compose.yml logs -f

dev-logs-gateway: ## Show gateway logs only
	@podman-compose -f docker-compose.yml logs -f gateway

dev-restart: ## Restart development environment
	@$(MAKE) dev-down
	@$(MAKE) dev

dev-quiet: ## Start services without logging stack logs
	@echo "🐳 Starting logging stack in background..."
	podman-compose up -d loki promtail grafana
	@sleep 5
	@echo "🐳 Starting services with logs..."
	podman-compose up gateway service-a service-b

stop: ## Stop development containers
	@echo "🛑 Stopping development services..."
	podman-compose stop
	@echo "✅ Development services stopped!"


# ==========================================
# Production Commands
# ==========================================

prod-build: ## Build production Docker images (no migration)
	@echo "════════════════════════════════════════════════════════════════"
	@echo "🏗️  Building Production Images"
	@echo "════════════════════════════════════════════════════════════════"
	@echo "📝 Using GO_VERSION=$(GO_VERSION)"
	@podman-compose -f docker-compose.prod.yml build \
		--build-arg GO_VERSION=$(GO_VERSION)
	@echo "✅ Production images built!"

prod-migrate: ## Run production database migrations (with safety prompt)
	@echo "════════════════════════════════════════════════════════════════"
	@echo "⚠️  PRODUCTION DATABASE MIGRATION"
	@echo "════════════════════════════════════════════════════════════════"
	@echo "This will modify the production database schema."
	@echo ""
	@read -p "Continue with production migration? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "🚀 Running production migrations..."; \
		GO_VERSION=$(GO_VERSION) podman-compose -f docker-compose.migrate.prod.yml up --abort-on-container-exit; \
		podman-compose -f docker-compose.migrate.prod.yml down; \
		echo "✅ Production migration complete!"; \
	else \
		echo "❌ Migration cancelled."; \
		exit 1; \
	fi

prod-backup: ## Create production database backup
	@echo "💾 Creating production database backup..."
	@mkdir -p backups
	@BACKUP_FILE=backups/backup_$(shell date +%Y%m%d_%H%M%S).sql; \
	podman exec $$(podman ps -q -f name=postgres) pg_dump -U $(POSTGRES_USER) $(POSTGRES_NAME) > $$BACKUP_FILE; \
	echo "✅ Backup created: $$BACKUP_FILE"

prod-deploy: prod-backup prod-migrate prod-build prod-up ## Full production deployment (backup → migrate → build → deploy)
	@echo ""
	@echo "════════════════════════════════════════════════════════════════"
	@echo "✅ Production Deployment Complete!"
	@echo "════════════════════════════════════════════════════════════════"

prod-up: ## Start production services
	@echo "🚀 Starting production services..."
	@podman-compose -f docker-compose.prod.yml up -d
	@echo "⏳ Waiting for services to be healthy..."
	@sleep 10
	@echo "✅ Production services running"

prod-down: ## Stop production services
	@echo "🛑 Stopping production services..."
	@podman-compose -f docker-compose.prod.yml down
	@echo "✅ Production services stopped"

prod-logs: ## Show production logs
	@podman-compose -f docker-compose.prod.yml logs -f

prod-status: ## Show production migration status
	@echo "📊 Production Migration Status:"
	@echo "════════════════════════════════════════════════════════════════"
	@podman exec $$(podman ps -q -f name=postgres) psql -U $(POSTGRES_USER) -d $(POSTGRES_NAME) -c \
		"SELECT version, service, description, to_char(applied_at, 'YYYY-MM-DD HH24:MI:SS') as applied \
		 FROM schema_migrations ORDER BY applied_at DESC LIMIT 10;" 2>/dev/null || \
		echo "⚠️  Cannot connect to production database"

prod-restart: prod-down prod-up ## Restart production services (without rebuild)

prod-rebuild: prod-down prod-build prod-up ## Rebuild and restart production services



# ==========================================
# Cleanup & Deletion
# ==========================================

clean: ## Remove containers, volumes, and temporary files
	@echo "🧹 Cleaning up..."
	@echo "  → Stopping and removing dev containers..."
	podman-compose down -v 2>/dev/null || true
	@echo "  → Stopping and removing prod containers..."
	podman-compose -f docker-compose.prod.yml down -v 2>/dev/null || true
	@echo "  → Removing temporary files..."
	rm -rf app/backend/service-a/tmp app/backend/service-b/tmp app/backend/gateway/tmp
	@echo "  → Removing dangling images..."
	podman rmi $$(podman images -f "dangling=true" -q) 2>/dev/null || true
	@echo "✅ Cleanup complete!"

delete: ## Delete ALL containers, images, and system data (DESTRUCTIVE!)
	@bash -c '\
	echo "⚠️  WARNING: This will delete ALL Podman data!"; \
	read -p "Are you sure? [y/N] " REPLY; \
	if [ "$$REPLY" = "y" ] || [ "$$REPLY" = "Y" ]; then \
		echo "🧹 Deleting all images..."; \
		echo "Stopping all containers..."; \
		podman stop -a; \
		echo "Deleting all stopped containers..."; \
		podman rm -a; \
		echo "Deleting all images..."; \
		podman image prune -a --force; \
		echo "System clean up (networks, volumes, build cache)..."; \
		podman system prune -a --force; \
		echo "✅ Deleted!"; \
	else \
		echo "❌ Deletion cancelled."; \
	fi'

delete-vol: ## Delete ALL Podman volumes (DESTRUCTIVE!)
	@bash -c '\
	echo "⚠️  WARNING: This will delete ALL Podman volumes!"; \
	read -p "Are you sure? [y/N] " REPLY; \
	if [ "$$REPLY" = "y" ] || [ "$$REPLY" = "Y" ]; then \
		echo "🧹 Deleting all volumes..."; \
		podman volume rm $$(podman volume ls -q) 2>/dev/null || true; \
		echo "✅ Volumes deleted!"; \
	else \
		echo "❌ Deletion cancelled."; \
	fi'	



# ==========================================
# Database Migration & Status
# ==========================================

db-migrate: ## Run all database migrations
	@echo "════════════════════════════════════════════════════════════════"
	@echo "🚀 Starting Database Migration"
	@echo "════════════════════════════════════════════════════════════════"
	@echo "📝 Configuration:"
	@echo "   DB_USER:     $(POSTGRES_USER)"
	@echo "   DB_NAME:     $(POSTGRES_NAME)"
	@echo "   DB_PORT:     $(PORT_POSTGRES)"
	@echo ""
	@echo "🐳 Starting PostgreSQL and running migrations..."
	@GO_VERSION=$(GO_VERSION) podman-compose -f docker-compose.migrate.yml up --build --abort-on-container-exit
	@echo ""
	@echo "🧹 Cleaning up migration containers..."
	@podman-compose -f docker-compose.migrate.yml down
	@echo ""
	@echo "════════════════════════════════════════════════════════════════"
	@echo "✅ Migration complete!"
	@echo "════════════════════════════════════════════════════════════════"

db-status: ## Show current migration status (requires running DB)
	@echo "📊 Current Migration Status:"
	@echo "════════════════════════════════════════════════════════════════"
	@podman exec $(CONTAINER_POSTGRES) psql -U $(POSTGRES_USER) -d $(POSTGRES_NAME) -c \
		"SELECT version, service, description, direction, to_char(applied_at, 'YYYY-MM-DD HH24:MI:SS') as applied \
		 FROM schema_migrations ORDER BY service, version;" 2>/dev/null || \
		echo "⚠️  Database container not running. Start it with 'make db-migrate' or 'make dev'"

db-test: ## Full migration test cycle
	@echo "════════════════════════════════════════════════════════════════"
	@echo "🔍 Full Migration Test Cycle"
	@echo "════════════════════════════════════════════════════════════════"
	@echo ""
	@echo "1️⃣  Starting PostgreSQL..."
	@GO_VERSION=$(GO_VERSION) podman-compose -f docker-compose.migrate.yml up -d postgres
	@echo "⏳ Waiting for PostgreSQL to be ready..."
	@sleep 8
	@echo ""
	@echo "2️⃣  Running migrations..."
	@GO_VERSION=$(GO_VERSION) podman-compose -f docker-compose.migrate.yml up --build db-migrator
	@echo ""
	@echo "3️⃣  Validation Results:"
	@echo "────────────────────────────────────────────────────────────────"
	@echo ""
	@echo "📊 Applied Migrations:"
	@podman exec $(CONTAINER_POSTGRES) psql -U $(POSTGRES_USER) -d $(POSTGRES_NAME) -c \
		"SELECT version, service, description, to_char(applied_at, 'YYYY-MM-DD HH24:MI:SS') as applied \
		 FROM schema_migrations ORDER BY service, version;" 2>/dev/null || echo "Failed to query migrations"
	@echo ""
	@echo "📚 Created Tables:"
	@podman exec $(CONTAINER_POSTGRES) psql -U $(POSTGRES_USER) -d $(POSTGRES_NAME) -c "\dt" 2>/dev/null || echo "Failed to list tables"
	@echo ""
	@echo "4️⃣  Cleaning up..."
	@podman-compose -f docker-compose.migrate.yml down
	@echo ""
	@echo "════════════════════════════════════════════════════════════════"
	@echo "✅ Test complete!"
	@echo "════════════════════════════════════════════════════════════════"

db-tables: ## List all database tables (requires running DB)
	@echo "📚 Database Tables:"
	@echo "════════════════════════════════════════════════════════════════"
	@podman exec $(CONTAINER_POSTGRES) psql -U $(POSTGRES_USER) -d $(POSTGRES_NAME) -c "\dt" 2>/dev/null || \
		echo "⚠️  Database container not running."

db-connect: ## Connect to PostgreSQL shell (requires running DB)
	@echo "🔌 Connecting to PostgreSQL..."
	@podman exec -it $(CONTAINER_POSTGRES) psql -U $(POSTGRES_USER) -d $(POSTGRES_NAME)

db-migrations-list: ## List all available migration files
	@echo "📁 Available Migration Files:"
	@echo "════════════════════════════════════════════════════════════════"
	@find app/backend/*/db/migrations -name "*.sql" 2>/dev/null | sort || echo "No migrations found"

db-clean: ## ⚠️  DANGER: Remove database volume (deletes all data!)
	@echo "════════════════════════════════════════════════════════════════"
	@echo "⚠️  WARNING: This will DELETE all database data!"
	@echo "════════════════════════════════════════════════════════════════"
	@read -p "Are you absolutely sure? Type 'YES' to confirm: " confirm; \
	if [ "$$confirm" = "YES" ]; then \
		echo "🗑️  Removing database volume..."; \
		podman-compose -f docker-compose.migrate.yml down -v; \
		podman-compose -f docker-compose.yml down -v; \
		echo "✅ Database volumes deleted"; \
	else \
		echo "❌ Operation cancelled."; \
	fi


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
		echo "   DB_PORT:     5432"; \
		echo "   Direction:   DOWN"; \
		echo "   Steps:       1"; \
		echo ""; \
		echo "🐳 Starting PostgreSQL and running rollback..."; \
		MIGRATION_DIRECTION=down ROLLBACK_STEPS=1 GO_VERSION=$(GO_VERSION) podman-compose -f docker-compose.migrate.yml up --build --abort-on-container-exit; \
		echo ""; \
		echo "🧹 Cleaning up migration containers..."; \
		podman-compose -f docker-compose.migrate.yml down; \
		echo ""; \
		echo "════════════════════════════════════════════════════════════════"; \
		echo "✅ Rollback complete!"; \
		echo "════════════════════════════════════════════════════════════════"; \
	else \
		echo "❌ Rollback cancelled."; \
	fi

db-rollback-all: ## Rollback ALL migrations (DANGEROUS!)
	@echo "════════════════════════════════════════════════════════════════"
	@echo "⚠️  DANGER: Rollback ALL Migrations"
	@echo "════════════════════════════════════════════════════════════════"
	@echo "⚠️  WARNING: This will rollback ALL migrations!"
	@echo "⚠️  This will DELETE all your database schema!"
	@echo ""
	@read -p "Are you ABSOLUTELY SURE? Type 'ROLLBACK ALL' to continue: " confirm; \
	if [ "$$confirm" = "ROLLBACK ALL" ]; then \
		echo "📝 Configuration:"; \
		echo "   DB_USER:     $(POSTGRES_USER)"; \
		echo "   DB_NAME:     $(POSTGRES_NAME)"; \
		echo "   Direction:   DOWN"; \
		echo "   Mode:        ALL"; \
		echo ""; \
		echo "🐳 Starting PostgreSQL and running full rollback..."; \
		MIGRATION_DIRECTION=down ROLLBACK_ALL=true GO_VERSION=$(GO_VERSION) podman-compose -f docker-compose.migrate.yml up --build --abort-on-container-exit; \
		echo ""; \
		echo "🧹 Cleaning up migration containers..."; \
		podman-compose -f docker-compose.migrate.yml down; \
		echo ""; \
		echo "════════════════════════════════════════════════════════════════"; \
		echo "✅ Full rollback complete!"; \
		echo "════════════════════════════════════════════════════════════════"; \
	else \
		echo "❌ Rollback cancelled."; \
	fi



.PHONY: db-migrate db-rollback db-rollback-all db-status db-test db-connect db-tables db-migrations-list db-clean db-rollback





# ==========================================
# Metrics & Logs
# ==========================================

stat: ## Show Podman images and running containers
	@echo "🔍 Podman Status"
	@echo ""
	@echo "📦 Images:"
	@podman images
	@echo ""
	@echo "🏃 Running Containers:"
	@podman ps
	@echo ""
	@echo "💤 All Containers:"
	@podman ps -a

info: ## Show Podman system information
	@echo "🔍 Podman System Information"
	@podman info

storage: ## Show Podman storage usage
	@echo "💾 Storage Usage"
	@podman system df

logs: ## Show logs from all development services
	@echo "📋 Showing logs (Ctrl+C to exit)..."
	podman-compose logs -f

logs-a: ## Show logs from service-a only
	@echo "📋 Showing service-a logs (Ctrl+C to exit)..."
	podman-compose logs -f service-a

logs-b: ## Show logs from service-b only
	@echo "📋 Showing service-b logs (Ctrl+C to exit)..."
	podman-compose logs -f service-b

logs-g: ## Show logs from gateway only
	@echo "📋 Showing gateway logs (Ctrl+C to exit)..."
	podman-compose logs -f gateway

check: fmt lint test ## Run all checks (format, lint, test)
	@echo "✅ All checks passed!"

ci: init check ## Run CI pipeline locally
	@echo "✅ CI pipeline complete!"

logs-grafana: ## Open Grafana in browser
	@echo "🔍 Opening Grafana..."
	@echo "URL: http://localhost:3000"
	@echo "Login: admin / admin"
	@xdg-open http://localhost:3000 2>/dev/null || open http://localhost:3000 2>/dev/null || echo "Open manually: http://localhost:3000"

logs-loki: ## Show Loki logs
	@echo "📋 Showing Loki logs..."
	podman-compose logs -f loki