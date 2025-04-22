package weblib

import (
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestGenerateNonce(t *testing.T) {
	tests := []struct {
		name    string
		size    uint
		wantLen int // Expected length of hex string (2 * size)
	}{
		{
			name:    "size 16",
			size:    16,
			wantLen: 32, // 16 bytes = 32 hex chars
		},
		{
			name:    "size 8",
			size:    8,
			wantLen: 16, // 8 bytes = 16 hex chars
		},
		{
			name:    "zero size",
			size:    0,
			wantLen: 0, // 0 bytes = 0 hex chars
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nonce, err := GenerateNonce(tt.size)
			if err != nil {
				t.Errorf("GenerateNonce returned unexpected error: %v", err)
			}
			if len(nonce) != tt.wantLen {
				t.Errorf("GenerateNonce returned nonce of length %d, want %d", len(nonce), tt.wantLen)
			}

			if tt.wantLen > 0 {
				if _, err := hex.DecodeString(nonce); err != nil {
					t.Errorf("GenerateNonce returned invalid hex: %v", err)
				}
			}
		})
	}
}

func TestWithNonce(t *testing.T) {
	newTestHandler := func(t *testing.T) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nonce := templ.GetNonce(r.Context())
			if nonce == "" {
				t.Error("Nonce not set in context")
			}
			w.WriteHeader(http.StatusOK)
		})
	}

	t.Run("middleware sets CSP and context", func(t *testing.T) {
		middleware := WithNonce(16)
		handler := middleware(newTestHandler(t))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Handler returned status %d, want %d", rr.Code, http.StatusOK)
		}

		csp := rr.Header().Get("Content-Security-Policy")
		if csp == "" {
			t.Error("Content-Security-Policy header not set")
		} else if !regexp.MustCompile(`^script-src 'nonce-[0-9a-f]+'$`).MatchString(csp) {
			t.Errorf("Content-Security-Policy header %q does not match expected format", csp)
		}

		prefix := "script-src 'nonce-"
		suffix := "'"
		if !strings.HasPrefix(csp, prefix) || !strings.HasSuffix(csp, suffix) {
			t.Errorf("Content-Security-Policy header %q has invalid format", csp)
		}

		nonce := strings.TrimPrefix(strings.TrimSuffix(csp, suffix), prefix)
		if len(nonce) != 32 {
			t.Errorf("Nonce length is %d, want 32", len(nonce))
		}
	})
}
