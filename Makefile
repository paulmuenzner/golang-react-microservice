
# Central version management
GO_VERSION ?= 1.25

.PHONY: help init dev dev-a dev-b dev-g test lint prod prod-up prod-stop stop clean delete stat info storage

.DEFAULT_GOAL := help

help: ## Show this help message
	@echo "📖 Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: install-tools

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


dev-local:
	@echo "💻 Starting local development..."
	@echo ""
	@echo "Open 2 terminals and run:"
	@echo " Terminal 1: make dev-local-a"
	@echo " Terminal 2: make dev-local-b"
	@echo " Terminal 2: make dev-local-g"
	@echo ""

dev: ## Start all services in development mode with hot-reload
	@echo "🐳 Starting services (Go ${GO_VERSION})..."
	GO_VERSION=${GO_VERSION} podman-compose up --build

dev-a: ## Start only service-a in development mode
	@echo "🐳 Starting service-a only (Go ${GO_VERSION})..."
	podman-compose up service-a

dev-b: ## Start only service-b in development mode
	@echo "🐳 Starting service-b only (Go ${GO_VERSION})..."
	GO_VERSION=${GO_VERSION} podman-compose up service-b

dev-g: ## Start only gateway in development mode
	@echo "🐳 Starting gateway only (Go ${GO_VERSION})..."
	GO_VERSION=${GO_VERSION} podman-compose up gateway

lint: ## Run golangci-lint on all Go code
	@echo "🔍 Running linters..."
	@cd shared/go && golangci-lint run
	@cd app/backend/service-a && golangci-lint run
	@cd app/backend/service-b && golangci-lint run
	@cd app/backend/gateway && golangci-lint run
	@echo "✅ Linting complete!"

fmt: ## Format all Go code
	@echo "🎨 Formatting code..."
	@cd shared/go && go fmt ./...
	@cd app/backend/service-a && go fmt ./...
	@cd app/backend/service-b && go fmt ./...
	@cd app/backend/gateway && go fmt ./...
	@echo "✅ Formatting complete!"

prod: ## Build production Docker images
	@echo "🏗️  Building production images..."
	podman-compose -f docker-compose.prod.yml build \
		--build-arg GO_VERSION=${GO_VERSION}
	@echo "✅ Production images built!"

prod-up: prod ## Start production containers
	@echo "🚀 Starting production services..."
	podman-compose -f docker-compose.prod.yml up -d
	@echo "✅ Production services running!"
	@echo ""
	@echo "📍 Services available at:"
	@echo "Service A: http://localhost:8080"
	@echo "Service B: http://localhost:8081"
	@echo "Service B: http://localhost:8082"

prod-stop: ## Stop production containers
	@echo "🛑 Stopping production services..."
	podman-compose -f docker-compose.prod.yml stop
	@echo "✅ Production services stopped!"

stop: ## Stop development containers
	@echo "🛑 Stopping development services..."
	podman-compose stop
	@echo "✅ Development services stopped!"

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