.PHONY: init dev dev-local dev-local-g dev-local-a dev-local-b prod stop clean test stat info storage delete


init:
	@echo "🔧 Initializing Go modules..."
	@cd shared/go && go mod tidy
	@cd app/backend/service-a && go mod tidy
	@cd app/backend/service-b && go mod tidy
	@cd app/backend/gateway && go mod tidy
	@echo "✅ Done!"

dev-local:
	@echo "💻 Starting local development..."
	@echo ""
	@echo "Open 2 terminals and run:"
	@echo " Terminal 1: make dev-local-a"
	@echo " Terminal 2: make dev-local-b"
	@echo " Terminal 2: make dev-local-g"
	@echo ""

dev-local-g:
	@echo "🚀 Starting gateway locally with Air..."
	@cd app/backend/gateway && air	

dev-local-a:
	@echo "🚀 Starting service-a locally with Air..."
	@cd app/backend/service-a && air

dev-local-b:
	@echo "🚀 Starting service-b locally with Air..."
	@cd app/backend/service-b && air

dev:
	@echo "🐳 Starting with Podman + Air hot-reload..."
	podman-compose up --build

prod:
	@echo "🏗️ Building production images..."
	podman-compose build service-a
	podman-compose build service-b
	podman-compose build gateway
	@echo "✅ Production images built!"
	@echo ""
	@echo "Run with: podman-compose -f docker-compose.prod.yml up"

prod-up:
	@echo "🚀 Starting production containers..."
	podman-compose -f docker-compose.prod.yml up -d
	@echo "✅ Production services running!"
	@echo ""
	@echo "Service A: http://localhost:8080"
	@echo "Service B: http://localhost:8081"
	@echo "Service B: http://localhost:8082"

stop:
	@echo "🛑 Stopping Docker services..."
	podman-compose -f docker-compose.yml stop

prod-stop:
	@echo "🛑 Stopping Docker services..."
	podman-compose -f docker-compose.prod.yml stop

clean:
	@echo "🧹 Cleaning up. Removing containers, volumes, and images..."
	@echo "Dev container..."
	podman-compose -f docker-compose.yml down -v
	rm -rf app/backend/service-a/tmp app/backend/service-b/tmp
	@echo "Prod container..."
	podman-compose -f docker-compose.prod.yml down -v
	@echo "🧹 Delete unused images..."
	podman rmi $(podman images -f "dangling=true" -q) 2>/dev/null || true
	@echo "🧹 Delete temp air folder..."
	rm -rf app/backend/service-a/tmp app/backend/service-b/tmp
	@echo "✅ Clean up complete!"

delete:
	@echo "🧹 Delete all images..."
	@echo " Stop all containers..."
	podman stop -a
	@echo "Delete all stopped containers..."
	podman rm -a
	@echo "Delete all images..."
	podman image prune -a --force
	@echo "System clean up (networks, volumes, build cache)..."
	podman system prune -a --force
	@echo "✅ Deleted!"

test:
	@echo "🧪 Running tests..."
	@cd app/backend/service-a && go test -v ./...
	@cd app/backend/service-b && go test -v ./...
	@echo "✅ Tests passed!"

stat:
	@echo "🔍 Podman images and running containers..."
	@echo "Images (podman images) ---"
	podman images
	@echo "Running Containers (podman ps) ---"
	podman ps

info:
	@echo "🔍 Podman information..."
	podman info

storage:
	@echo "💾 Storage Usage..."
	podman system df