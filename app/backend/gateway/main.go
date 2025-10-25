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

	// Register service routes
	mux.Handle("/service-a/", middleware.NewProxy("http://service-a:8080", "/service-a"))
	mux.Handle("/service-b/", middleware.NewProxy("http://service-b:8080", "/service-b"))

	// Health check endpoint
	mux.HandleFunc("/health", middleware.HealthHandler)

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
