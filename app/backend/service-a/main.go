package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	shared "github.com/app/shared/go"
	config "github.com/app/shared/go/config"
	logger "github.com/app/shared/go/utils/logger"
	"github.com/google/uuid"
)

// ==========================================
// MIDDLEWARE: Request-ID extraction
// ==========================================
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Request-ID read from header (from Gateway)
		requestID := r.Header.Get("X-Request-ID")

		// If no header (direct call), create new ID
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Store request-ID in context
		ctx := context.WithValue(r.Context(), "request_id", requestID)
		r = r.WithContext(ctx)

		// Request-ID back to Response Header (for Client)
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r)
	})
}

// ==========================================
// MIDDLEWARE 2: HTTP Logging
// ==========================================
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := r.Context().Value("request_id").(string)
		reqLogger := logger.WithRequestID(requestID)

		reqLogger.InfoWithFields("Incoming request to Service-A", map[string]interface{}{
			"method": r.Method,
			"path":   r.URL.Path,
		})

		wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}
		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		reqLogger.HTTP(r.Method, r.URL.Path, wrapped.statusCode, duration, r.RemoteAddr)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func main() {
	logger.Init("SERVICE-A", os.Getenv("ENVIRONMENT"))

	// Load Config
	routeConfig, err := config.LoadRouteConfig()
	if err != nil {
		logger.FatalWithFields("Failed to load route configuration", err, nil)
	}
	listenHost := routeConfig.Backend.ListenHost
	listenPort := routeConfig.Backend.ListenPort
	listenAddress := fmt.Sprintf("%s:%s", listenHost, listenPort)

	// ==========================================
	// Mux with middleware
	// ==========================================
	mux := http.NewServeMux()

	// Define routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// ==========================================
		// Extract request ID from context for logging
		// ==========================================
		requestID := r.Context().Value("request_id").(string)
		reqLogger := logger.WithRequestID(requestID)

		msg := shared.Greet("Service A")
		reqLogger.InfoWithFields(
			"Processing request in Service A",
			map[string]interface{}{
				"remote_addr": r.RemoteAddr,
				"path":        r.URL.Path,
			},
		)
		fmt.Fprintf(w, "%s\n", msg)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	// ==========================================
	// Wrap mux with middleware
	// ==========================================
	handler := requestIDMiddleware(loggingMiddleware(mux))

	// Create server
	srv := &http.Server{
		Addr:    listenAddress,
		Handler: handler, // ← Mit Middleware!
	}

	// Start server
	go func() {
		logger.Info(fmt.Sprintf("Server is ready to handle requests on %s", listenAddress))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.FatalWithFields(
				"Could not start HTTP server listener",
				err,
				map[string]interface{}{"address": listenAddress},
			)
		}
	}()

	// Graceful shutdown...
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info(fmt.Sprintf("Received signal: %v. Shutting down gracefully...", sig))

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(fmt.Sprintf("Server forced to shutdown: %s", err))
	}

	logger.Info("Server stopped gracefully")
}

// func main() {
// 	logger.Init("SERVICE-A", os.Getenv("ENVIRONMENT"))

// 	// Load Config
// 	routeConfig, err := config.LoadRouteConfig()
// 	if err != nil {
// 		logger.FatalWithFields("Failed to load route configuration", err, nil)
// 	}
// 	listenHost := routeConfig.Backend.ListenHost
// 	listenPort := routeConfig.Backend.ListenPort
// 	listenAddress := fmt.Sprintf("%s:%s", listenHost, listenPort)

// 	// Define routes
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		msg := shared.Greet("Service A")
// 		logger.InfoWithFields(
// 			"Incoming request received",
// 			map[string]interface{}{"remote_addr": r.RemoteAddr},
// 		)
// 		fmt.Fprintf(w, "%s\n", msg)
// 	})

// 	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		fmt.Fprint(w, "OK")
// 	})

// 	// Create server with explicit configuration
// 	srv := &http.Server{
// 		Addr:    listenAddress,
// 		Handler: nil, // Uses default ServeMux
// 	}

// 	// Start server in a goroutine
// 	go func() {
// 		logger.Info(fmt.Sprintf("Server is ready to handle requests on %s", listenAddress))
// 		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			logger.FatalWithFields(
// 				"Could not start HTTP server listener",
// 				err,
// 				map[string]interface{}{"address": listenAddress},
// 			)
// 		}
// 	}()

// 	// Setup signal catching
// 	quit := make(chan os.Signal, 1)
// 	// Catch SIGINT (Ctrl+C) and SIGTERM (docker stop)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

// 	// Block until signal is received
// 	sig := <-quit
// 	logger.Info(fmt.Sprintf("Received signal: %v. Shutting down gracefully...", sig))

// 	// Create a deadline for shutdown
// 	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
// 	defer cancel()

// 	// Attempt graceful shutdown
// 	if err := srv.Shutdown(ctx); err != nil {
// 		logger.Fatal(fmt.Sprintf("Server forced to shutdown: %s", err))
// 	}

// 	logger.Info("Server stopped gracefully")
// }
