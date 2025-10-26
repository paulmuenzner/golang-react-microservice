## Purpose

Short, actionable instructions for AI coding agents working on this monorepo. Focus on the project's conventions, build/test flows, integration points, and concrete examples.

## Big picture (what this repo is)
- Monorepo for multiple Go microservices: `app/backend/gateway`, `service-a`, `service-b` and shared libraries under `shared/go`.
- Gateway routes to internal services and applies a heavy middleware stack (rate-limiting, security headers, logging, request-id, etc.).
- Centralized logging stack (Loki, Promtail, Grafana) and a dedicated DB migrator built from `shared/go/utils/db/migrator`.

## Key files & directories to inspect
- `Makefile` — primary entry for developer workflows (init, dev, lint, test, prod, db-migrate). Use it for recommended commands.
- `docker-compose.yml` — dev composition (gateway, service-a, service-b, postgres, loki, promtail, grafana). Uses `podman-compose` in Makefile but is Docker-compatible.
- `docker-compose.migrate.yml` — used by migration tasks; builds `db-migrator` from `shared/go/utils/db/migrator/Dockerfile.migrator` and copies migrations from `app/backend/*/db/migrations`.
- `shared/go/config/load_config.go` — configuration loader uses runtime.Caller and resolves JSON config files relative to package location. Always use the provided helper functions (e.g., `LoadRouteConfig`) instead of hardcoding paths.
- `shared/go/utils/logger/` — structured JSON logging and event code conventions. Look at `events_*.go` for event code patterns used across middleware and services.
- `app/backend/*/main.go` and `app/backend/*/Dockerfile` — service entrypoints and how middleware is composed; useful for adding new routes or middleware behavior.

## Developer workflows (concrete commands)
- Initialize modules/tools: `make init` (runs `go mod tidy` across services). Install tools with `make install-tools`.
- Start full dev environment (with hot-reload): `make dev` (invokes `podman-compose up --build`). The Makefile uses `GO_VERSION` and `podman-compose` by default.
- Start services individually: `make dev-a`, `make dev-b`, `make dev-g`.
- Run linters/tests: `make lint`, `make test`, `make check` (formats, lints, tests).
- Build production images: `make prod` and `make prod-up` / `make prod-stop` to manage prod containers.
- Run DB migrations: `make db-migrate` (this uses `docker-compose.migrate.yml` and the `db-migrator` image). For a dry-run/test: `make db-test-migrations`.
- Connect to DB shell: `make db-connect` (exec into `postgres-dev`).

Notes:
- The repository's Makefile expects `podman-compose`, though `docker-compose` often works interchangeably in CI or local environments; prefer `podman-compose` to match author conventions.
- Environment variables are read from `.env` (used by compose / migrator). Ensure `.env` or `.env.example` is populated.

## Project-specific conventions & patterns
- Shared module import path: code imports the shared library as `github.com/app/shared/go`. Respect that module path when adding packages.
- Config loading: packages use `loadRelativeConfig` (runtime.Caller) to find JSON config files under `shared/data/config`. When writing tests or moving files, keep relative locations stable.
- Migrations: place SQL files under `app/backend/<service>/db/migrations`. The migrator Dockerfile copies `app/backend/*/db/migrations` into `/app/migrations` at build time. Naming convention and ordering is handled by the migrator (see migrator code).
- Hot-reload & volumes: dev Dockerfiles mount source into `/build/...` and rely on `air` for live reloading. Keep `/tmp` or `tmp` directories for runtime artifacts out of version control.
- Request ID propagation: the gateway sets `X-Request-ID` header; services read it from headers and store it in request context (key: `request_id`). When adding instrumentation or correlation, use the same header and context key.

## Integration points & external deps
- PostgreSQL: defined in `docker-compose.yml` (container `postgres-dev`). DB credentials come from env variables in `.env`.
- Logging: Grafana/Loki/Promtail stack configured in `docker-compose.yml` and `config/` files. Services log structured JSON; use Loki queries referencing `container`/`service` labels.
- Cloudflare/CDN: middleware contains Cloudflare-aware IP extraction and validation (see `middleware/cloudflareValidationMiddleware.go`). Keep this when deploying behind Cloudflare.

## How to add common changes (examples)
- Add a DB migration for service-a:
  1. Create SQL file `app/backend/service-a/db/migrations/002_add_table_xyz.up.sql`.
  2. Run `make db-migrate` (it will build the migrator, copy migrations, and apply them against the `postgres` container in `docker-compose.migrate.yml`).

- Add a new middleware to the gateway:
  1. Implement middleware in `shared/go/middleware/` (follow existing middleware signatures that accept and return `http.Handler`).
  2. Register it in `middleware/middlewareBuilder.go` or `middleware/middlewareBuilder` chain so it appears in the global stack used by the gateway.
  3. Write unit tests focused on middleware behavior and an integration test using `make dev` to validate runtime behavior.

## Pitfalls & gotchas
- Config file paths are resolved relative to the package; moving JSON files or packages can break `LoadRouteConfig` and similar helpers.
- The project favors Podman (`podman-compose`) in Makefile scripts — CI or contributors using Docker may need minor command adjustments.
- Do not manually publish `shared/go` — the monorepo uses module replacement locally. Use `make init` to regenerate tidy modules.

## Where to look for more context
- `Makefile` (workflow & targets)
- `docker-compose.yml` and `docker-compose.migrate.yml` (run-time topology)
- `shared/go/config/load_config.go` (config resolution pattern)
- `shared/go/utils/logger/` (event codes and logging conventions)
- `shared/go/utils/db/migrator/` (migrator implementation and Dockerfile)
- `app/backend/*/db/migrations` (migration SQL files)

If anything here is unclear or you want extra examples (e.g., unit test snippets, a sample migration, or a small middleware PR), tell me which area to expand and I will iterate.
