package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// ==========================================
// COMPRESSION MIDDLEWARE (GZIP)
// ==========================================

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
	wroteHeader bool
}

func (w *gzipResponseWriter) WriteHeader(code int) {
	if !w.wroteHeader {
		// Entferne Content-Length - wird neu berechnet!
		w.Header().Del("Content-Length")
		w.ResponseWriter.WriteHeader(code)
		w.wroteHeader = true
	}
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.Writer.Write(b)
}

func CompressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Nur komprimieren wenn Client es unterstützt
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Setze Header
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length") // ← WICHTIG!

		gz := gzip.NewWriter(w)
		defer gz.Close()

		gzw := &gzipResponseWriter{
			Writer:         gz,
			ResponseWriter: w,
		}
		next.ServeHTTP(gzw, r)
	})
}
