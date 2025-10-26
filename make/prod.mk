# ==========================================
# Production Commands
# ==========================================

.PHONY: prod-build prod-migrate prod-backup prod-deploy prod-up prod-down prod-logs prod-status prod-restart prod-rebuild

prod-build: ## Build production Docker images (no migration)
	@echo "════════════════════════════════════════════════════════════════"
	@echo "🏗️  Building Production Images"
	@echo "════════════════════════════════════════════════════════════════"
	@echo "📝 Using GO_VERSION=$(GO_VERSION)"
	@podman-compose -f docker-compose.prod.yml build \
		--build-arg GO_VERSION=$(GO_VERSION)
	@echo "✅ Production images built!"

prod-migrate: ## Run production database migrations (with safety prompt)
	@echo "════════════════════════════════════════════════════════════════"
	@echo "⚠️  PRODUCTION DATABASE MIGRATION"
	@echo "════════════════════════════════════════════════════════════════"
	@echo "📝 Configuration:"
	@echo "   DB_USER:     $(POSTGRES_USER)"
	@echo "   DB_NAME:     $(POSTGRES_NAME)"
	@echo "   DB_PORT:     $(PORT_POSTGRES)"
	@echo "   Direction:   UP"
	@echo ""
	@echo "This will modify the production database schema."
	@echo ""
	@read -p "Continue with production migration? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "🚀 Running production migrations..."; \
		MIGRATION_DIRECTION=up GO_VERSION=$(GO_VERSION) podman-compose -f docker-compose.migrate.prod.yml up --abort-on-container-exit; \
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
		echo "⚠️  Cannot connect to production database"

prod-restart: prod-down prod-up ## Restart production services (without rebuild)

prod-rebuild: prod-down prod-build prod-up ## Rebuild and restart production services

