package main

import (
	"fmt"
	"net/http"
	"os"

	shared "github.com/app/shared/go"
)

func main() {
	logger.Init("SERVICE-B", os.Getenv("ENVIRONMENT"))
	logger.Println("Starting on 0.0.0.0:8080")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		msg := shared.Greet("Service B")
		logger.Printf("Request from %s", r.RemoteAddr)
		fmt.Fprintf(w, "%s\n", msg)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		logger.Fatal(err)
	}
}
