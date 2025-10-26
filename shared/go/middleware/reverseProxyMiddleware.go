package middleware

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// ==========================================
// REVERSE PROXY HANDLER
// ==========================================
func NewProxy(target string, prefix string) http.Handler {
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

		// Forward Request-ID header for distributed tracing
		if requestID := req.Context().Value("request_id"); requestID != nil {
			req.Header.Set("X-Request-ID", requestID.(string))
		}
	}

	// Configure transport with reasonable timeouts
	proxy.Transport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return proxy
}
