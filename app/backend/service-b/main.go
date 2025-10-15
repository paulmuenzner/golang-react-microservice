package main

import (
	"fmt"
	"net/http"

	shared "github.com/company/shared/go"
)

func main() {
	logger := shared.NewLogger("SERVICE-B")
	logger.Println("Starting on localhost:8081")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		msg := shared.Greet("Service B")
		logger.Printf("Request from %s", r.RemoteAddr)
		fmt.Fprintf(w, "%s\n", msg)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})

	if err := http.ListenAndServe("localhost:8081", nil); err != nil {
		logger.Fatal(err)
	}
}
