package main

import (
	"fmt"
	"net/http"
	"os"

	shared "github.com/app/shared/go"
	logger "github.com/app/shared/go/utils/logger"
)

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
