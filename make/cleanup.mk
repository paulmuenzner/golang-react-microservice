# ==========================================
# Cleanup & Deletion Targets
# ==========================================

.PHONY: clean clean-all clean-volumes clean-images clean-tmp clean-containers delete delete-vol

clean: ## Remove containers, volumes, and temporary files
	@echo "════════════════════════════════════════════════════════════════"
	@echo "🧹 Cleaning up..."
	@echo "════════════════════════════════════════════════════════════════"
	@echo ""
	@echo "📦 Stopping and removing containers..."
	@podman-compose -f docker-compose.yml down -v 2>/dev/null || true
	@podman-compose -f docker-compose.migrate.yml down -v 2>/dev/null || true
	@podman-compose -f docker-compose.prod.yml down -v 2>/dev/null || true
	@echo ""
	@echo "🗑️  Removing temporary files..."
	@find app/backend/*/tmp -type f -delete 2>/dev/null || true
	@find app/backend/*/tmp -type d -empty -delete 2>/dev/null || true
	@echo ""
	@echo "🖼️  Removing dangling images..."
	@podman images -f "dangling=true" -q | xargs -r podman rmi 2>/dev/null || true
	@echo ""
	@echo "════════════════════════════════════════════════════════════════"
	@echo "✅ Cleanup complete!"
	@echo "════════════════════════════════════════════════════════════════"

clean-all: ## Deep clean (containers, volumes, images, cache)
	@echo "════════════════════════════════════════════════════════════════"
	@echo "🔥 DEEP CLEAN - This will remove EVERYTHING!"
	@echo "════════════════════════════════════════════════════════════════"
	@read -p "Are you sure? Type 'yes' to continue: " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		echo ""; \
		echo "📦 Stopping all containers..."; \
		podman stop $$(podman ps -aq) 2>/dev/null || true; \
		echo ""; \
		echo "🗑️  Removing all containers..."; \
		podman rm $$(podman ps -aq) 2>/dev/null || true; \
		echo ""; \
		echo "💾 Removing all volumes..."; \
		podman volume prune -f 2>/dev/null || true; \
		echo ""; \
		echo "🖼️  Removing all images..."; \
		podman rmi $$(podman images -q) 2>/dev/null || true; \
		echo ""; \
		echo "🧹 Cleaning build cache..."; \
		podman system prune -af 2>/dev/null || true; \
		echo ""; \
		echo "🗑️  Removing temporary files..."; \
		find app/backend/*/tmp -type f -delete 2>/dev/null || true; \
		find app/backend/*/tmp -type d -empty -delete 2>/dev/null || true; \
		echo ""; \
		echo "════════════════════════════════════════════════════════════════"; \
		echo "✅ Deep clean complete!"; \
		echo "════════════════════════════════════════════════════════════════"; \
	else \
		echo "❌ Clean cancelled."; \
	fi

clean-volumes: ## Remove all volumes (keeps containers)
	@echo "💾 Removing all volumes..."
	@podman-compose -f docker-compose.yml down -v 2>/dev/null || true
	@podman-compose -f docker-compose.migrate.yml down -v 2>/dev/null || true
	@podman-compose -f docker-compose.prod.yml down -v 2>/dev/null || true
	@podman volume prune -f 2>/dev/null || true
	@echo "✅ Volumes removed"

clean-images: ## Remove all dangling and unused images
	@echo "🖼️  Removing images..."
	@podman images -f "dangling=true" -q | xargs -r podman rmi 2>/dev/null || true
	@podman image prune -af 2>/dev/null || true
	@echo "✅ Images cleaned"

clean-tmp: ## Remove temporary files from all services
	@echo "🗑️  Removing temporary files..."
	@find app/backend/*/tmp -type f -delete 2>/dev/null || true
	@find app/backend/*/tmp -type d -empty -delete 2>/dev/null || true
	@echo "✅ Temporary files removed"

clean-containers: ## Stop and remove all containers
	@echo "📦 Removing containers..."
	@podman-compose -f docker-compose.yml down 2>/dev/null || true
	@podman-compose -f docker-compose.migrate.yml down 2>/dev/null || true
	@podman-compose -f docker-compose.prod.yml down 2>/dev/null || true
	@echo "✅ Containers removed"

delete: ## Delete ALL containers, images, and system data (DESTRUCTIVE!)
	@bash -c '\
	echo "⚠️  WARNING: This will delete ALL Podman data!"; \
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
	echo "⚠️  WARNING: This will delete ALL Podman volumes!"; \
	read -p "Are you sure? [y/N] " REPLY; \
	if [ "$$REPLY" = "y" ] || [ "$$REPLY" = "Y" ]; then \
		echo "🧹 Deleting all volumes..."; \
		podman volume rm $$(podman volume ls -q) 2>/dev/null || true; \
		echo "✅ Volumes deleted!"; \
	else \
		echo "❌ Deletion cancelled."; \
	fi
