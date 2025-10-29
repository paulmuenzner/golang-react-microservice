package main

import (
	"fmt"
	"net/http"
	"os"

	shared "github.com/app/shared/go"
	config "github.com/app/shared/go/config"
	logger "github.com/app/shared/go/utils/logger"
)

// ==========================================
// SERVICE A
// ==========================================
func main() {
	logger.Init("SERVICE-B", os.Getenv("ENVIRONMENT"))

	// Load Config
	routeConfig, err := config.LoadRouteConfig()
	if err != nil {
		logger.FatalWithFields("Failed to load route configuration", err, nil)
	}
	listenAddress := fmt.Sprintf("%s:%s", routeConfig.Backend.ListenHost, routeConfig.Backend.ListenPort)

	// Simple Mux - NO middleware needed!
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Request-ID kommt vom Gateway!
		requestID := r.Header.Get("X-Request-ID")
		reqLogger := logger.WithRequestID(requestID)

		msg := shared.Greet("Service B")
		reqLogger.InfoWithFields("Processing request", map[string]interface{}{
			"path": r.URL.Path,
		})

		fmt.Fprintf(w, "%s\n", msg)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	// Simple Server - NO middleware wrapper!
	logger.Info(fmt.Sprintf("Service B ready on %s", listenAddress))
	if err := http.ListenAndServe(listenAddress, mux); err != nil {
		logger.FatalWithFields("Failed to start server", err, nil)
	}
}
