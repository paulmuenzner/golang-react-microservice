# make/frontend.mk

# ==========================================
# Frontend Configuration
# ==========================================

FRONTEND_DIR := app/frontend/public
NPM := npm

# ==========================================
# Frontend Development
# ==========================================

frontend-dev: ## Start frontend development server
	@echo "🚀 Starting Next.js development server..."
	@cd $(FRONTEND_DIR) && $(NPM) run dev

frontend-dev-clean: ## Start frontend with cache cleanup
	@echo "🧹 Cleaning caches and starting Next.js..."
	@cd $(FRONTEND_DIR) && rm -rf .next node_modules/.cache .turbo
	@cd $(FRONTEND_DIR) && $(NPM) run dev

frontend-dev-fresh: ## Fresh install and start
	@echo "🔄 Fresh install and starting Next.js..."
	@cd $(FRONTEND_DIR) && rm -rf node_modules package-lock.json .next .turbo node_modules/.cache
	@cd $(FRONTEND_DIR) && $(NPM) install
	@cd $(FRONTEND_DIR) && $(NPM) run dev

# ==========================================
# Frontend Build & Start
# ==========================================

frontend-build: ## Build frontend for production
	@echo "🔨 Building Next.js for production..."
	@cd $(FRONTEND_DIR) && $(NPM) run build

frontend-start: ## Start production server
	@echo "🚀 Starting Next.js production server..."
	@cd $(FRONTEND_DIR) && $(NPM) start

frontend-build-start: frontend-build frontend-start ## Build and start production

# ==========================================
# Frontend Dependencies
# ==========================================

frontend-install: ## Install frontend dependencies
	@echo "📦 Installing frontend dependencies..."
	@cd $(FRONTEND_DIR) && $(NPM) install

frontend-update: ## Update frontend dependencies
	@echo "⬆️  Updating frontend dependencies..."
	@cd $(FRONTEND_DIR) && $(NPM) update

frontend-audit: ## Audit frontend dependencies
	@echo "🔍 Auditing frontend dependencies..."
	@cd $(FRONTEND_DIR) && $(NPM) audit

frontend-audit-fix: ## Fix frontend vulnerabilities
	@echo "🔧 Fixing frontend vulnerabilities..."
	@cd $(FRONTEND_DIR) && $(NPM) audit fix

# ==========================================
# Frontend Cleanup
# ==========================================

frontend-clean: ## Clean frontend caches
	@echo "🧹 Cleaning frontend caches..."
	@cd $(FRONTEND_DIR) && rm -rf .next node_modules/.cache .turbo

frontend-clean-all: ## Clean everything including node_modules
	@echo "🧹 Cleaning everything..."
	@cd $(FRONTEND_DIR) && rm -rf .next node_modules/.cache .turbo node_modules package-lock.json

# ==========================================
# Frontend Linting & Type Checking
# ==========================================

frontend-lint: ## Lint frontend code
	@echo "🔍 Linting frontend code..."
	@cd $(FRONTEND_DIR) && $(NPM) run lint

frontend-type-check: ## Type check frontend code
	@echo "🔍 Type checking frontend code..."
	@cd $(FRONTEND_DIR) && npx tsc --noEmit

# ==========================================
# Frontend Testing (wenn du Tests hast)
# ==========================================

frontend-test: ## Run frontend tests
	@echo "🧪 Running frontend tests..."
	@cd $(FRONTEND_DIR) && $(NPM) test

frontend-test-watch: ## Run frontend tests in watch mode
	@echo "🧪 Running frontend tests in watch mode..."
	@cd $(FRONTEND_DIR) && $(NPM) test -- --watch

# ==========================================
# Frontend Utilities
# ==========================================

frontend-shell: ## Open shell in frontend directory
	@cd $(FRONTEND_DIR) && bash

frontend-upgrade-next: ## Upgrade Next.js to latest
	@echo "⬆️  Upgrading Next.js..."
	@cd $(FRONTEND_DIR) && $(NPM) install next@latest react@latest react-dom@latest

frontend-downgrade-next: ## Downgrade Next.js to 15.1.3 (stable)
	@echo "⬇️  Downgrading Next.js to 15.1.3..."
	@cd $(FRONTEND_DIR) && $(NPM) install next@15.1.3
	@cd $(FRONTEND_DIR) && rm -rf .next
	@echo "✅ Downgraded to Next.js 15.1.3"

# ==========================================
# Docker Frontend Commands
# ==========================================

frontend-docker-build: ## Build frontend Docker image
	@echo "🐳 Building frontend Docker image..."
	@podman build -f $(FRONTEND_DIR)/Dockerfile \
		--target production \
		-t frontend:latest \
		$(FRONTEND_DIR)

frontend-docker-dev: ## Start frontend in Docker
	@echo "🐳 Starting frontend in Docker..."
	@podman-compose up -d frontend
	@echo "✅ Frontend running in Docker"
	@echo "📍 http://localhost:$(PORT_FRONTEND)"

frontend-docker-logs: ## Show frontend Docker logs
	@podman-compose logs -f frontend

frontend-docker-restart: ## Restart frontend container
	@podman-compose restart frontend

# ==========================================
# Combined Commands
# ==========================================

frontend-reset: frontend-clean-all frontend-install ## Complete reset
	@echo "✅ Frontend reset complete"

frontend-fix-turbopack: frontend-clean frontend-downgrade-next ## Fix Turbopack crashes
	@echo "✅ Turbopack fix applied (downgraded to Next.js 15)"