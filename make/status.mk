# ==========================================
# Metrics & Logs Targets
# ==========================================

.PHONY: stat info storage logs logs-a logs-b logs-g logs-grafana logs-loki

stat: ## stat: Show Podman images and running containers
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

info: ## info: Show Podman system information
	@echo "🔍 Podman System Information"
	@podman info

storage: ## storage: Show Podman storage usage
	@echo "💾 Storage Usage"
	@podman system df

logs: ## logs: Show logs from all development services
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

logs-grafana: ## logs-grafana: Open Grafana in browser
	@echo "🔍 Opening Grafana..."
	@echo "URL: http://localhost:$(PORT_GRAFANA)"
	@echo "Login: admin / admin"
	@xdg-open http://localhost:$(PORT_GRAFANA) 2>/dev/null || open http://localhost:$(PORT_GRAFANA) 2>/dev/null || echo "Open manually: http://localhost:$(PORT_GRAFANA)"

logs-loki: ## logs-loki: Show Loki logs
	@echo "📋 Showing Loki logs..."
	podman-compose logs -f loki
