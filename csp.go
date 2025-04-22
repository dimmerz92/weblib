package weblib

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
)

// GenerateNonce generates a random nonce of the given size in bytes and returns it as a hex encoded string.
func GenerateNonce(size uint) (string, error) {
	bytes := make([]byte, size)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

// WithNonce returns a middleware closure that generates random nonces of the given size in bytes and sets the
// Content-Security-Policy response header in addition to setting it in the context for automatice usage by templ.
func WithNonce(size uint) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nonce, err := GenerateNonce(size)
			if err != nil {
				log.Printf("nonce could not be generated: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Security-Policy", fmt.Sprintf("script-src 'nonce-%s'", nonce))
			next.ServeHTTP(w, r.WithContext(templ.WithNonce(r.Context(), nonce)))
		})
	}
}
