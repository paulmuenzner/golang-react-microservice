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
