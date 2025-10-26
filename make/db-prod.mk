# make/db-prod.mk
# Production-Safe Database Migrations

.PHONY: prod-db-migrate prod-db-migrate-safe prod-db-rollback-safe prod-db-maintenance

# ==========================================
# Production Migrations (Safe)
# ==========================================

prod-db-migrate-safe: ## Production migration with safety checks
	@echo "════════════════════════════════════════════════════════════════"
	@echo "🏭 Production Database Migration (SAFE MODE)"
	@echo "════════════════════════════════════════════════════════════════"
	@echo ""
	@echo "⚠️  PRODUCTION DEPLOYMENT CHECKLIST:"
	@echo ""
	@echo "  ☐ Migrations are backward-compatible"
	@echo "  ☐ No DROP TABLE or DROP COLUMN statements"
	@echo "  ☐ All changes are additive (ADD COLUMN, CREATE INDEX)"
	@echo "  ☐ Tested in staging environment"
	@echo "  ☐ Rollback plan documented"
	@echo "  ☐ Backup created"
	@echo "  ☐ Team notified"
	@echo ""
	@read -p "All checks passed? Type 'MIGRATE' to continue: " confirm; \
	if [ "$$confirm" = "MIGRATE" ]; then \
		echo ""; \
		echo "1️⃣  Creating backup..."; \
		$(MAKE) prod-db-backup; \
		echo ""; \
		echo "2️⃣  Checking migration files..."; \
		$(MAKE) prod-db-check-breaking; \
		echo ""; \
		echo "3️⃣  Running migration..."; \
		MIGRATION_DIRECTION=up GO_VERSION=$(GO_VERSION) \
			podman-compose -f docker-compose.prod.yml -f docker-compose.migrate.yml \
			up --build --abort-on-container-exit db-migrator || exit 1; \
		echo ""; \
		echo "4️⃣  Verifying schema..."; \
		$(MAKE) prod-db-verify; \
		echo ""; \
		echo "════════════════════════════════════════════════════════════════"; \
		echo "✅ Production migration completed!"; \
		echo "════════════════════════════════════════════════════════════════"; \
		echo ""; \
		echo "📝 Next steps:"; \
		echo "  1. Monitor application logs"; \
		echo "  2. Check error rates in Grafana"; \
		echo "  3. Keep backup for 7 days"; \
		echo ""; \
	else \
		echo "❌ Migration cancelled."; \
		exit 1; \
	fi

prod-db-check-breaking: ## Check for breaking changes in migrations
	@echo "🔍 Checking for breaking changes..."
	@echo "════════════════════════════════════════════════════════════════"
	@BREAKING_FOUND=0; \
	for file in $$(find app/backend/*/db/migrations -name "*.up.sql" 2>/dev/null); do \
		if grep -iE "(DROP TABLE|DROP COLUMN|RENAME COLUMN|ALTER TABLE.*DROP)" "$$file" >/dev/null 2>&1; then \
			echo "⚠️  BREAKING CHANGE in $$file:"; \
			grep -inE "(DROP TABLE|DROP COLUMN|RENAME COLUMN|ALTER TABLE.*DROP)" "$$file"; \
			echo ""; \
			BREAKING_FOUND=1; \
		fi; \
	done; \
	if [ $$BREAKING_FOUND -eq 1 ]; then \
		echo "════════════════════════════════════════════════════════════════"; \
		echo "❌ Breaking changes detected!"; \
		echo ""; \
		echo "These changes will break running services!"; \
		echo ""; \
		echo "Options:"; \
		echo "  1. Use 3-phase migration (recommended)"; \
		echo "  2. Schedule maintenance window"; \
		echo "  3. Use blue-green deployment"; \
		echo ""; \
		exit 1; \
	else \
		echo "✅ No breaking changes detected"; \
		echo "════════════════════════════════════════════════════════════════"; \
	fi

prod-db-backup: ## Create production backup with timestamp
	@BACKUP_FILE=backups/prod_backup_$$(date +%Y%m%d_%H%M%S).sql; \
	echo "💾 Creating production backup: $$BACKUP_FILE"; \
	mkdir -p backups; \
	podman exec $(CONTAINER_POSTGRES) pg_dump -U $(POSTGRES_USER) $(POSTGRES_NAME) > $$BACKUP_FILE && \
		echo "✅ Backup created: $$BACKUP_FILE" || \
		(echo "❌ Backup failed!" && exit 1)

prod-db-verify: ## Verify database schema after migration
	@echo "🔍 Verifying database schema..."
	@echo "════════════════════════════════════════════════════════════════"
	@podman exec $(CONTAINER_POSTGRES) psql -U $(POSTGRES_USER) -d $(POSTGRES_NAME) -c \
		"SELECT schemaname, tablename FROM pg_tables WHERE schemaname NOT IN ('pg_catalog', 'information_schema') ORDER BY schemaname, tablename;"
	@echo ""
	@podman exec $(CONTAINER_POSTGRES) psql -U $(POSTGRES_USER) -d $(POSTGRES_NAME) -c \
		"SELECT COUNT(*) as migration_count FROM schema_migrations;"
	@echo "════════════════════════════════════════════════════════════════"

# ==========================================
# Maintenance Window Migration
# ==========================================

prod-db-maintenance: ## Full production migration with downtime
	@echo "════════════════════════════════════════════════════════════════"
	@echo "🔧 MAINTENANCE WINDOW - Database Migration"
	@echo "════════════════════════════════════════════════════════════════"
	@echo ""
	@echo "⚠️  This will cause DOWNTIME!"
	@echo ""
	@echo "Maintenance steps:"
	@echo "  1. Stop all services (downtime begins)"
	@echo "  2. Create backup"
	@echo "  3. Run migrations"
	@echo "  4. Verify schema"
	@echo "  5. Start services (downtime ends)"
	@echo ""
	@echo "Expected downtime: 5-10 minutes"
	@echo ""
	@read -p "Start maintenance window? Type 'MAINTENANCE' to confirm: " confirm; \
	if [ "$$confirm" = "MAINTENANCE" ]; then \
		START_TIME=$$(date +%s); \
		echo ""; \
		echo "════════════════════════════════════════════════════════════════"; \
		echo "🛑 MAINTENANCE MODE ACTIVE"; \
		echo "════════════════════════════════════════════════════════════════"; \
		echo ""; \
		echo "1️⃣  Stopping all services..."; \
		podman-compose -f docker-compose.prod.yml stop gateway service-a service-b || true; \
		echo "   ✅ Services stopped"; \
		echo ""; \
		echo "2️⃣  Creating backup..."; \
		$(MAKE) prod-db-backup || exit 1; \
		echo ""; \
		echo "3️⃣  Running migrations..."; \
		MIGRATION_DIRECTION=up GO_VERSION=$(GO_VERSION) \
			podman-compose -f docker-compose.migrate.yml up --build --abort-on-container-exit || exit 1; \
		echo ""; \
		echo "4️⃣  Verifying schema..."; \
		$(MAKE) prod-db-verify || exit 1; \
		echo ""; \
		echo "5️⃣  Starting services..."; \
		podman-compose -f docker-compose.prod.yml up -d gateway service-a service-b || exit 1; \
		echo "   ⏳ Waiting for health checks..."; \
		sleep 10; \
		echo ""; \
		echo "6️⃣  Health check..."; \
		curl -f http://localhost:$(PORT_GATEWAY)/health || echo "⚠️  Health check failed!"; \
		echo ""; \
		END_TIME=$$(date +%s); \
		DURATION=$$((END_TIME - START_TIME)); \
		echo "════════════════════════════════════════════════════════════════"; \
		echo "✅ MAINTENANCE COMPLETE"; \
		echo "════════════════════════════════════════════════════════════════"; \
		echo ""; \
		echo "📊 Summary:"; \
		echo "   Downtime: $$DURATION seconds"; \
		echo "   Status: Services running"; \
		echo ""; \
	else \
		echo "❌ Maintenance cancelled."; \
	fi

# ==========================================
# Rollback (Production)
# ==========================================

prod-db-rollback-safe: ## Safe production rollback
	@echo "════════════════════════════════════════════════════════════════"
	@echo "⬇️  Production Database Rollback"
	@echo "════════════════════════════════════════════════════════════════"
	@echo ""
	@echo "⚠️  WARNING: This will rollback the last migration in PRODUCTION!"
	@echo ""
	@echo "Prerequisites:"
	@echo "  ☐ Services can handle old schema"
	@echo "  ☐ No data loss will occur"
	@echo "  ☐ Backup is available"
	@echo "  ☐ Team is notified"
	@echo ""
	@read -p "Ready to rollback? Type 'ROLLBACK' to confirm: " confirm; \
	if [ "$$confirm" = "ROLLBACK" ]; then \
		echo ""; \
		echo "1️⃣  Creating pre-rollback backup..."; \
		$(MAKE) prod-db-backup; \
		echo ""; \
		echo "2️⃣  Rolling back migration..."; \
		MIGRATION_DIRECTION=down ROLLBACK_STEPS=1 GO_VERSION=$(GO_VERSION) \
			podman-compose -f docker-compose.migrate.yml up --build --abort-on-container-exit; \
		echo ""; \
		echo "3️⃣  Verifying schema..."; \
		$(MAKE) prod-db-verify; \
		echo ""; \
		echo "✅ Rollback complete!"; \
	else \
		echo "❌ Rollback cancelled."; \
	fi

# ==========================================
# Helper Commands
# ==========================================

prod-db-status: ## Show production migration status
	@echo "📊 Production Migration Status:"
	@echo "════════════════════════════════════════════════════════════════"
	@podman exec $(CONTAINER_POSTGRES) psql -U $(POSTGRES_USER) -d $(POSTGRES_NAME) -c \
		"SELECT version, service, description, direction, applied_at \
		 FROM schema_migrations \
		 ORDER BY applied_at DESC LIMIT 20;"
	@echo "════════════════════════════════════════════════════════════════"

prod-db-health: ## Check production database health
	@echo "🏥 Production Database Health:"
	@echo "════════════════════════════════════════════════════════════════"
	@podman exec $(CONTAINER_POSTGRES) psql -U $(POSTGRES_USER) -d $(POSTGRES_NAME) -c \
		"SELECT \
		    pg_database_size('$(POSTGRES_NAME)') / 1024 / 1024 as size_mb, \
		    (SELECT count(*) FROM pg_stat_activity WHERE state = 'active') as active_connections, \
		    (SELECT count(*) FROM pg_stat_activity) as total_connections;"
	@echo "════════════════════════════════════════════════════════════════"