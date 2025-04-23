package weblib

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Middleware func(http.Handler) http.Handler

type gzipRW struct {
	io.Writer
	http.ResponseWriter
}

// Chain applies the given middlewares to the next handler in the given order and returns it.
func Chain(next http.Handler, middlewares ...Middleware) http.Handler {
	for _, mw := range middlewares {
		next = mw(next)
	}

	return next
}

// Logger logs the method, path, and time taken for the given handler.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"method: %s\troute: %s\ttime_taken: %dms",
			r.Method,
			r.URL.Path,
			time.Since(start).Milliseconds(),
		)
	})
}

// Write provides a custom gzip capable response writer write method.
func (w gzipRW) Write(b []byte) (int, error) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}

	return w.Writer.Write(b)
}

// Gzip applies gzip compression to a response if it is an accepted encoding.
func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gz := gzip.NewWriter(w)
		defer gz.Close()

		gzrw := gzipRW{Writer: gz, ResponseWriter: w}

		next.ServeHTTP(gzrw, r)
	})
}
