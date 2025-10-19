# Golang React Microservice Template

A production-ready monorepo template for building microservices with Go, featuring centralized logging and monitoring.

## Features

### Architecture
- **Monorepo Structure**: Centralized repository for multiple independent services (service-a, service-b) and shared libraries (shared/go)
- **Isolated Dependency Management**: Uses Go's module replace directives to link services to the internal shared library without requiring external publishing
- **API Gateway Pattern**: Single entry point (gateway) routing to internal microservices

### Development Experience
- **Containerized Development**: Utilizes podman-compose (compatible with docker-compose) for environment consistency and simplified dependency handling
- **Instant Hot-Reload (Air)**: Changes in service files or the shared library automatically trigger a recompile and restart within the container, saving valuable development time
- **Clean Build Workflow**: All common tasks (setup, dev, test, production builds) are automated through the Makefile
- **Robust Development Setup**: Non-destructive Dockerfile commands ensure that local go.mod files are protected from container-level modifications

### Logging & Monitoring
- **Centralized Logging (Loki)**: All service logs automatically collected and stored in a time-series database
- **Log Visualization (Grafana)**: Beautiful web UI for searching, filtering, and analyzing logs across all services
- **Automatic Log Collection (Promtail)**: Zero-config log aggregation from all Docker/Podman containers
- **Structured JSON Logs**: All services output structured logs for easy parsing and analysis
- **Log Retention**: Configurable retention periods (7 days dev, 90 days prod)


## Prerequisites

- Go 1.25+
- Podman
- Make


## Quick Start

### 1. Initial Setup
```bash
# Clone repository
git clone 
cd golang-react-microservice

# Copy environment file
cp .env.example .env

# Initialize dependencies
make init
```

### 2. Development
```bash
# Start all services with hot-reload
make dev

# Or start services individually
make dev-a  # service-a only
make dev-b  # service-b only
make dev-g  # gateway only
```

**Services available at:**
- Gateway: http://localhost:8080
- Grafana (Logs): http://localhost:3000 (admin/admin)

## Testing 🧪

```bash
make test
```

## Production Deployment 🚢
```bash
# Build production images
make prod

# Start production environment
make prod-up

# Stop production
make prod-stop
```


## Project Structure

```
/project/
├── .gitignore               # Files ignored by Git (Node, Go, Python artifacts)
├── .dockerignore            # Files excluded from the Docker build context (speeds up builds)
├── docker-compose.yml       # Development environment definition (with Hot-Reload)
├── docker-compose.prod.yml  # Production environment definition
├── Makefile                 # Automation script for setup, dev, test, and cleanup
├── README.md                # This file
│
├── shared/
│ ├── go/
│ │ ├── go.mod               # Shared module definition
│ │ ├── shared.go            # Example shared logic (e.g., logging, helpers)
│ │ ├── utils/              # Utils
│ │     └── logger/
│ │
│ ├── react/
│ │
│ ├── data/
│
└── app/
  └── backend/
    ├── gateway/
    │   ├── go.mod             # Service module file
    │   ├── main.go            # Service entry point
    │   ├── Dockerfile         # Service-specific build instructions
    │   └── .air.toml          # Air Hot-Reload configuration
    │
    ├── service-a/
    │   ├── go.mod             # Service module file
    │   ├── main.go            # Service entry point
    │   ├── Dockerfile         # Service-specific build instructions
    │   └── .air.toml          # Air Hot-Reload configuration
    │
    └── service-b/
        ├── go.mod           # Service module file
        ├── main.go          # Service entry point
        ├── Dockerfile       # Service-specific build instructions
        └── .air.toml        # Air Hot-Reload configuration

```

## Troubleshooting

### Grafana shows no data
1. Check Loki is running: `podman ps | grep loki`
2. Check Promtail logs: `podman logs promtail`
3. Verify datasource: Grafana → Configuration → Data Sources


### Grafana Query Examples

1. All gateway logs: `{container="gateway-dev"}`
2. Only errors: `{container="gateway-dev"} | json | level="error"`
3. Certain user: `{container="gateway-dev"} | json | user_id="12345"`
4. Tracing request ID in all services: `{service=~".*"} | json | request_id="550e8400..."`
5. Slow Requests: `{service=~".*"} | json | duration_ms > 1000`
6. HTTP Errors (4xx, 5xx): `{service=~".*"} | json | status >= 400`
7. Certain operations: `{service=~".*"} | json | component="email_worker"`

### Services won't start
```bash
make clean   # Remove all containers
make dev     # Restart
```

### Hot-reload not working
- Ensure volumes are correctly mounted
- Check Air logs in service output


## Contributing 

See [CONTRIBUTING.md](CONTRIBUTING.md)

## License 

MIT License - see [LICENSE](LICENSE)