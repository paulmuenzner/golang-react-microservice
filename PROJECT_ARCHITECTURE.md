# Microservice Architecture Overview

This document outlines the architecture and project structure of the Go/React Monorepo, focusing on clear separation of concerns, shared resources, and efficient deployment.



```
/project/
├── .gitignore                      # Files ignored by Git (Node, Go, Python artifacts)
├── .dockerignore                   # Files excluded from the Docker build context (speeds up builds)
├── .golangci.yml                   # 
├── docker-compose.yml              # Development environment definition (with Hot-Reload)
├── docker-compose.prod.yml         # Production environment definition
├── docker-compose.migrate.yml      # Database migration definition
├── docker-compose.migrate.prod.yml
├── Makefile                        # Automation script for setup, dev, test, migrate, cleanup, ...
├── README.md                  
│
├── shared/
│ ├── go/
│ │ ├── go.mod                      # Shared module definition
│ │ ├── shared.go                   # Example shared logic (e.g., logging, helpers)
│ │ ├── utils/                      # Utils Golang
│ │ │   ├── cache/
│ │ │   ├── conversion/
│ │ │   ├── date/
│ │ │   ├── db/                     # Database connection, query, migration, ...
│ │ │   │   ├── migrator/           # Migration config. Dockerfile.migrator and migration function main.go
│ │ │   │   └── postgres_connection.go
│ │ │   │
│ │ │   ├── files/
│ │ │   ├── ip/
│ │ │   ├── logger/
│ │ │   ├── metrics/
│ │ │   ├── misc/
│ │ │   ├── security/
│ │ │   └── validation/
│ │ │ 
│ │ ├── config/  
│ │ ├── interfaces/  
│ │ ├── middleware/                             // Middleware functions für api gateway
│ │ │   ├── healthMiddleware.go
│ │ │   ├── cloudflareValidationMiddleware.go
│ │ │   ├── loggingMiddleware.go
│ │ │   ├── compressionMiddleware.go
│ │ │   ├── corsMiddleware.go
│ │ │   ├── healthMiddleware.go
│ │ │   ├── ipExtractionMiddleware.go
│ │ │   ├── loggingMiddleware.go
│ │ │   ├── maxBytesMiddleware.go
│ │ │   ├── middlewareBuilder.go
│ │ │   ├── rateLimitMiddleware.go
│ │ │   ├── recoveryMiddleware.go
│ │ │   ├── requestMiddleware.go
│ │ │   ├── reverseProxyMiddleware.go
│ │ │   ├── securityMiddleware.go
│ │ │   └── timeoutMiddleware.go
│ │ │   
│ │ │ 
│ │
│ ├── react/
│ │
│ ├── data/
│ │   └── config/ 
│ │       ├── authConfig.json
│ │       ├── baseConfig.json
│ │       ├── regexConfig.json
│ │       └── routeConfig.json
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
    │   ├── .air.toml          # Air Hot-Reload configuration
    │   └── db/
    │       └── migrations/
    │           └── 001_create_schema.up.sql
    │  
    └── service-b/
        ├── go.mod           # Service module file
        ├── main.go          # Service entry point
        ├── Dockerfile       # Service-specific build instructions
        └── .air.toml        # Air Hot-Reload configuration

```