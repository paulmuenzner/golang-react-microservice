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
)

func main() {
	logger.Init("service-a", os.Getenv("ENVIRONMENT"))
	logger.Println("Starting on 0.0.0.0:8080")

	// Define routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		msg := shared.Greet("Service A")
		logger.Printf("Request from %s", r.RemoteAddr)
		fmt.Fprintf(w, "%s\n", msg)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	// Create server with explicit configuration
	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: nil, // Uses default ServeMux
	}

	// Start server in a goroutine
	go func() {
		logger.Println("Server is ready to handle requests")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Could not listen on 0.0.0.0:8080: %v\n", err)
		}
	}()

	// Setup signal catching
	quit := make(chan os.Signal, 1)
	// Catch SIGINT (Ctrl+C) and SIGTERM (docker stop)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until signal is received
	sig := <-quit
	logger.Printf("Received signal: %v. Shutting down gracefully...", sig)

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server stopped gracefully")
}
