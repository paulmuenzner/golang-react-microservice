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


# ==========================================
# Global Targets & Help
# ==========================================

.PHONY: help
help: ## Help information
	@echo "════════════════════════════════════════════════════════════════"
	@echo "Verfügbare Kommandos (Modularisiert):"
	@echo "════════════════════════════════════════════════════════════════"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2}'
	@echo ""



