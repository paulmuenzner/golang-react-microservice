package main

// For production, you should add features such as timeouts, circuit breakers (e.g. goresilience), retries, logging, authentication, and metrics.

import (
	"net/http"
	"os"
	"time"

	middleware "github.com/app/shared/go/middleware"
	logger "github.com/app/shared/go/utils/logger"
)

const listenAddr = ":8080"

// ==========================================
// MAIN
// ==========================================
func main() {
	// Initialize logger
	logger.Init("gateway", os.Getenv("ENVIRONMENT"))
	logger.Info("Gateway starting")

	// Create router
	mux := http.NewServeMux()

	// Root route - Gateway info
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"service": "gateway",
			"version": "1.0.0",
			"status": "running",
			"endpoints": {
				"health": "/health",
				"service-a": "/service-a/*",
				"service-b": "/service-b/*"
			}
		}`))
	})

	// Register service routes
	mux.Handle("/api/service-a/", middleware.NewProxy("http://service-a:8080", "/api/service-a"))
	mux.Handle("/api/service-b/", middleware.NewProxy("http://service-b:8080", "/api/service-b"))

	// Health check endpoint
	mux.HandleFunc("/api/health", middleware.HealthHandler)

	// Build complete middleware stack
	handler := middleware.BuildMiddlewareStack(mux)

	// Configure server with timeouts
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  15 * time.Second, // Time to read request
		WriteTimeout: 15 * time.Second, // Time to write response
		IdleTimeout:  60 * time.Second, // Keep-alive timeout
	}

	// Start server
	logger.Info("Gateway listening on :8080")
	if err := srv.ListenAndServe(); err != nil {
		logger.FatalWithFields(
			"Failed to start HTTP server",
			err,
			map[string]interface{}{
				"address": srv.Addr,
			},
		)
	}
}
