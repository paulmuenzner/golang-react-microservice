package main

// For production, you should add features such as timeouts, circuit breakers (e.g. goresilience), retries, logging, authentication, and metrics.

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	logger "github.com/app/shared/go/utils/logger"
)

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
	}

	// Optional: adjust timeouts via Transport
	proxy.Transport = &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return proxy
}

func main() {
	logger.Init("gateway", os.Getenv("ENVIRONMENT"))
	logger.Info("Gateway starting on localhost:8080")

	// Proxy /service-a -> http://service-a:8080
	http.Handle("/service-a/", newProxy("http://service-a:8080", "/service-a"))

	// Proxy /service-b -> http://service-b:8080
	http.Handle("/service-b/", newProxy("http://service-b:8080", "/service-b"))

	// Health
	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("gateway OK"))
	})

	logger.Info(fmt.Sprintf("Gateway listening on %s", listenAddr))

	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		logger.FatalWithFields(
			"Failed to start HTTP server",
			err,
			map[string]interface{}{
				"address": listenAddr,
			},
		)
	}
}
