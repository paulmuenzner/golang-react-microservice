# API Gateway Middleware Architecture

## Current Stack (in order)
1. Recovery → Panic handling
2. RequestID → UUID generation
3. IPExtraction → Client IP extraction (Cloudflare-compatible)
4. CloudflareValidation → Header spoofing detection
5. Timeout → 30s limit
6. MaxBytes → 10MB limit
7. CORS → Whitelist-based
8. SecurityHeaders → OWASP headers
9. RateLimit → 100 req/s per-IP (production)
10. Compression → gzip
11. Logging → Structured JSON to Loki

## File Structure
```

 middleware/
 ├── additionalMiddleware.go           
 ├── cloudflareValidationMiddleware.go  // Header spoofing detection
 ├── compressionMiddleware.go
 ├── corsMiddleware.go                  // Whitelist-based
 ├── healthMiddleware.go      
 ├── ipExtractionMiddleware.go          // Client IP extraction (Cloudflare-compatible)
 ├── loggingMiddleware.go               // Structured JSON to Loki
 ├── maxBytesMiddleware.go
 ├── middlewareBuilder.go               // Stacks all middleware functions
 ├── rateLimitMiddleware.go
 ├── recoveryMiddleware.go              // Panic handling
 ├── requestMiddleware.go               // UUID generation
 ├── reverseProxyMiddleware.go
 ├── securityMiddleware.go              // OWASP headers
 └── timeoutMiddleware.go

```

## Key Design Decisions
- IP extracted once in IPExtractionMiddleware, cached in context
- Request ID used across all middlewares for correlated logging
- Rate limiting per-IP, not global
- Cloudflare header validation for security

## Environment Variables
- `ENVIRONMENT`: "development" | "production"
- `USE_CLOUDFLARE`: "true" | "false"

## Open Questions / TODOs
- [ ] Add metrics middleware (Prometheus)
- [ ] Implement CSRF protection
- [ ] Add circuit breaker for backend services
```

