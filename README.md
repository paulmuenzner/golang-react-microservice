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
- **Event-Based Logging System**: Unique event codes (e.g., `MW-RL-001`, `AUTH-JWT-003`) for precise tracking, filtering, and analytics across all components
- **Smart Log Categorization**: Events organized by category (Middleware, Auth, Service, Database, API) with dedicated helper functions for consistent structured logging
- **Log Retention**: Configurable retention periods (7 days dev, 90 days prod)

### API Gateway Middleware Stack
1. **Recovery** - Panic handling with stack traces
2. **Request ID** - UUID generation for tracing
3. **Timeout** - 30s request limit (504 Gateway Timeout)
4. **Size Limit** - 10MB max payload (413 Entity Too Large)
5. **CORS** - Whitelist-based cross-origin control
6. **Security Headers** - OWASP headers (HSTS, CSP, XSS, clickjacking)
7. **Rate Limiting** - Per-IP: 100 req/s, Cloudflare-compatible
8. **Compression** - Automatic gzip compression
9. **Logging** - Structured JSON logs to Loki/Grafana

#### Middleware Configuration Reference

| Middleware | Default Setting | Production | Development | Configurable |
|------------|----------------|------------|-------------|--------------|
| **Timeout** | 30s | 30s | 30s | ✅ Via code |
| **Max Request Size** | 10MB | 10MB | 10MB | ✅ Via code |
| **Rate Limit** | - | 100 req/s (burst 200) | 1000 req/s (burst 2000) | ✅ Via code |
| **CORS Origins** | Whitelist | Strict whitelist | `*` (all) | ✅ Via code |
| **Compression** | Enabled | ✅ | ✅ | ❌ Always on |
| **Security Headers** | Full suite | ✅ | ✅ | ❌ Always on |
| **Request ID** | UUID v4 | ✅ | ✅ | ❌ Always on |
| **Recovery** | Enabled | ✅ | ✅ | ❌ Always on |
| **Logging** | Structured JSON | ✅ | ✅ | ❌ Always on |

### Security Headers Applied

| Header | Value | Purpose |
|--------|-------|---------|
| `X-Content-Type-Options` | `nosniff` | Prevent MIME type sniffing |
| `X-Frame-Options` | `DENY` | Prevent clickjacking attacks |
| `X-XSS-Protection` | `1; mode=block` | Enable XSS filter (legacy) |
| `Strict-Transport-Security` | `max-age=63072000; includeSubDomains; preload` | Force HTTPS for 2 years |
| `Content-Security-Policy` | Restrictive policy | Control resource loading |
| `Referrer-Policy` | `strict-origin-when-cross-origin` | Limit referrer leakage |
| `Permissions-Policy` | Disabled features | Block dangerous browser APIs |
| `Server` | (removed) | Hide server information |


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

[PROJECT_ARCHITECTURE](/PROJECT_ARCHITECTURE.md)

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

## More Feature Details

### Smart Client IP Extraction
The gateway includes intelligent IP detection that works across different deployment scenarios:
- **Cloudflare**: `CF-Connecting-IP` header (highest priority)
- **Enterprise CDNs**: `True-Client-IP` header (Akamai, Cloudflare Enterprise)
- **Nginx/Load Balancers**: `X-Real-IP` header
- **Standard Proxies**: `X-Forwarded-For` header (takes first public IP)
- **Direct Connection**: `RemoteAddr` (fallback)
- **Security**: Filters private/internal IPs and validates against spoofing

### Production Readiness
- **Environment-Aware Configuration**: Different settings for development vs. production
  - Development: Permissive CORS (*), higher rate limits (1000 req/s), verbose logging
  - Production: Strict CORS whitelist, lower rate limits (100 req/s), optimized performance
- **Graceful Error Handling**: All panics caught and logged without service crashes
- **Defense in Depth**: Multiple security layers (headers, rate limiting, size limits, timeouts)
- **Observability**: Comprehensive logging with request IDs for debugging and monitoring
- **Scalability**: Per-IP rate limiting allows horizontal scaling without shared state
- **CDN Compatible**: Works seamlessly with Cloudflare, AWS CloudFront, Google Cloud CDN

### Optional/Future Middlewares
- **CSRF Protection**: Cross-Site Request Forgery prevention with token validation
- **IP Filtering**: Whitelist/blacklist specific IP addresses or ranges
- **API Key Authentication**: Validate `X-API-Key` headers for protected endpoints
- **Circuit Breaker**: Prevent cascading failures by temporarily blocking failing services
- **Metrics Middleware**: Export Prometheus metrics for monitoring
- **Cloudflare Validation**: Verify requests truly originate from Cloudflare IPs

## Contributing 

See [CONTRIBUTING.md](CONTRIBUTING.md)

## License 

MIT License - see [LICENSE](LICENSE)