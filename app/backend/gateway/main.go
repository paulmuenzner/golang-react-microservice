package main

// For production, you should add features such as timeouts, circuit breakers (e.g. goresilience), retries, logging, authentication, and metrics.

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	logger "github.com/app/shared/go/utils/logger"
	"github.com/google/uuid"
)

// ==========================================
// LOGGING MIDDLEWARE
// ==========================================
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Generate unique Request ID
		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "request_id", requestID)
		r = r.WithContext(ctx)

		// Create logger with Request ID
		reqLogger := logger.WithRequestID(requestID)

		// Log incoming request
		reqLogger.InfoWithFields("Incoming request", map[string]interface{}{
			"method": r.Method,
			"path":   r.URL.Path,
			"ip":     r.RemoteAddr,
		})

		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}

		// Process request
		next.ServeHTTP(wrapped, r)

		// Log completed request
		duration := time.Since(start)
		reqLogger.HTTP(
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			duration,
			r.RemoteAddr,
		)
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

// ==========================================
// REVERSE PROXY HANDLER
// ==========================================
const listenAddr = ":8080"

func newProxy(target string, prefix string) http.Handler {
	u, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(u)

	origDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		origDirector(req)
		// Strip the prefix from the path (e.g. /service-a/foo -> /foo)
		req.URL.Path = strings.TrimPrefix(req.URL.Path, prefix)
		if req.URL.Path == "" {
			req.URL.Path = "/"
		}

		// Add request-id to header for tracing
		if requestID := req.Context().Value("request_id"); requestID != nil {
			req.Header.Set("X-Request-ID", requestID.(string))
		}
	}

	// Optional: adjust timeouts via Transport
	proxy.Transport = &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return proxy
}

// ==========================================
// MAIN
// ==========================================
func main() {
	// Initialize logger
	logger.Init("gateway", os.Getenv("ENVIRONMENT"))
	logger.Info("Gateway starting")

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register routes
	mux.Handle("/service-a/", newProxy("http://service-a:8080", "/service-a"))
	mux.Handle("/service-b/", newProxy("http://service-b:8080", "/service-b"))

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		logger.Debug("Health check")
		w.Write([]byte("gateway OK"))
	})

	// Wrap entire mux with logging middleware
	handler := loggingMiddleware(mux)

	// Start server
	logger.Info("Gateway listening on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		logger.FatalWithFields(
			"Failed to start HTTP server",
			err,
			map[string]interface{}{
				"address": listenAddr,
			},
		)
	}
}

// func main() {
// 	logger.Init("gateway", os.Getenv("ENVIRONMENT"))
// 	logger.Info("Gateway starting on localhost:8080")

// 	// Proxy /service-a -> http://service-a:8080
// 	http.Handle("/service-a/", newProxy("http://service-a:8080", "/service-a"))

// 	// Proxy /service-b -> http://service-b:8080
// 	http.Handle("/service-b/", newProxy("http://service-b:8080", "/service-b"))

// 	// Health
// 	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
// 		w.Write([]byte("gateway OK"))
// 	})

// 	logger.Info(fmt.Sprintf("Gateway listening on %s", listenAddr))

// 	if err := http.ListenAndServe(listenAddr, nil); err != nil {
// 		logger.FatalWithFields(
// 			"Failed to start HTTP server",
// 			err,
// 			map[string]interface{}{
// 				"address": listenAddr,
// 			},
// 		)
// 	}
// }
