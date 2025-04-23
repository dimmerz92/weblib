package weblib

import (
	"net/http"
	"strings"
)

// NoBrowse prevents a browser from being able to browse a file server.
func NoBrowse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ServeFiles serves files from the given root. The prefix should match the route used in your handler.
func ServeFiles(prefix, root string, browsable bool) http.Handler {
	fs := http.StripPrefix(prefix, http.FileServer(http.Dir(root)))

	if !browsable {
		return NoBrowse(fs)
	}

	return fs
}
