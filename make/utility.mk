# ==========================================
# Utility Commands (init, lint, tools)
# ==========================================

.PHONY: install-tools init lint check ci fmt test

install-tools: ## Install development tools (golangci-lint, air)
	@echo "🔧 Installing development tools..."
	@which golangci-lint > /dev/null || { \
		echo "  → Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	}
	@which air > /dev/null || { \
		echo "  → Installing air..."; \
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

test: ## Run go tests
	@echo "🧪 Running tests..."
	@go test ./... -v -short
