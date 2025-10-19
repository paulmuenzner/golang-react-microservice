# Microservice Blueprint: Go Monorepo with Hot-Reload

This project serves as a comprehensive and robust blueprint for developing modern microservices using Go, managed within a monorepo structure. It is designed for fast, containerized development using **Podman/Docker** and features **Hot-Reload** capabilities via `air`, ensuring a highly efficient development loop right from the start. It is designed for scalability, maintainability, and clean separation of concerns.

## Features

- Monorepo Structure: Centralized repository for multiple independent services (service-a, service-b) and shared libraries (shared/go).

- Isolated Dependency Management: Uses Go's module replace directives to link services to the internal shared library without requiring external publishing.

- Containerized Development: Utilizes podman-compose (compatible with docker-compose) for environment consistency and simplified dependency handling.

- Instant Hot-Reload (air): Changes in service files or the shared library automatically trigger a recompile and restart within the container, saving valuable development time.

- Clean Build Workflow: All common tasks (setup, dev, test, production builds) are automated through the Makefile.

- Robust Development Setup: Non-destructive Dockerfile commands ensure that local go.mod files are protected from container-level modifications.

## Quick Start

```bash
make init  # Initialize dependencies
make dev   # Start development environment
```

## Testing 🧪

```bash
make test
```

## Deployment 🚢

```bash
make prod
make prod-up
```

## Contributing 🤝

See [CONTRIBUTING.md](CONTRIBUTING.md)

## License 📄

MIT License - see [LICENSE](LICENSE)

## Architecture

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
│ │ ├── logger/              # Service module file
│ │ │   └── logger.go
│ │
│ │
│ │
│ │
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
