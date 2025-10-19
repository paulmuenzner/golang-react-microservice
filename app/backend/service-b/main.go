package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	shared "github.com/app/shared/go"
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

func main() {
	logger.Init("SERVICE-B", os.Getenv("ENVIRONMENT"))
	logger.Info("Starting on 0.0.0.0:8080")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		msg := shared.Greet("Service B")
		logger.InfoWithFields(
			"Incoming request received",
			map[string]interface{}{"remote_addr": r.RemoteAddr},
		)

		fmt.Fprintf(w, "%s\n", msg)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		logger.Error("Starting on 0.0.0.0:8080")
	}
}
