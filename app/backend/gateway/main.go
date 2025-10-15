// gateway/main.go

package main

import (
	"fmt"
	"net/http"

	shared "github.com/company/shared/go"
)

func main() {
	logger := shared.NewLogger("GATEWAY")
	logger.Println("Starting on 0.0.0.0:8082")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		msg := shared.Greet("Gateway")
		logger.Printf("Request from %s", r.RemoteAddr)
		fmt.Fprintf(w, "%s\n", msg)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Gateway OK!!")
	})

	if err := http.ListenAndServe("0.0.0.0:8082", nil); err != nil {
		logger.Fatal(err)
	}
}
