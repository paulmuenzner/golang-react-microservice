# ==========================================
# Environment variables from .env
# ==========================================
include .env
export

.DEFAULT_GOAL := help


# ==========================================
# Make Imports
# ==========================================
include make/status.mk
include make/db-dev.mk
include make/utility.mk
include make/dev.mk
include make/prod.mk
include make/cleanup.mk
include make/db-prod.mk
include make/frontend.mk


# ==========================================
# Global Targets & Help
# ==========================================

help: ## Show this help message
	@echo "════════════════════════════════════════════════════════════════"
	@echo "Available Commands:"
	@echo "════════════════════════════════════════════════════════════════"
	@echo ""
	@echo "🐳 Development:"
	@grep -E '^dev[^:]*:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "🎨 Frontend:"
	@grep -E '^frontend[^:]*:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[35m%-25s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "🗄️  Database:"
	@grep -E '^db[^:]*:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[34m%-25s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "════════════════════════════════════════════════════════════════"

.PHONY: help



