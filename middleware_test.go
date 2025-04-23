package weblib

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

func TestChain(t *testing.T) {
	testMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test", "true")
			next.ServeHTTP(w, r)
		})
	}

	handler := Chain(http.HandlerFunc(handler), testMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Header().Get("X-Test") != "true" {
		t.Errorf("expected X-Test header to be set")
	} else if rec.Body.String() != "hello world" {
		t.Errorf("expected body to be 'hello world', got %q", rec.Body.String())
	}
}

func TestLogger(t *testing.T) {
	handler := Logger(http.HandlerFunc(handler))

	req := httptest.NewRequest(http.MethodGet, "/logtest", nil)
	rec := httptest.NewRecorder()

	// redirect logs to buffer
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(io.Discard)

	handler.ServeHTTP(rec, req)

	if !strings.Contains(buf.String(), "method: GET") {
		t.Error("expected log to contain method GET")
	} else if rec.Body.String() != "hello world" {
		t.Errorf("expected body to be 'hello world', got %q", rec.Body.String())
	}
}

func TestGzip(t *testing.T) {
	handler := Gzip(http.HandlerFunc(handler))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Errorf("expected Content-Encoding to be gzip")
	}

	gr, err := gzip.NewReader(rec.Body)
	if err != nil {
		t.Fatalf("failed to create gzip reader: %v", err)
	}
	defer gr.Close()

	body, err := io.ReadAll(gr)
	if err != nil {
		t.Fatalf("failed to read gzip body: %v", err)
	}

	if string(body) != "hello world" {
		t.Errorf("expected body to be 'hello world', got %q", string(body))
	}
}

func TestGzipWithoutHeader(t *testing.T) {
	handler := Gzip(http.HandlerFunc(handler))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Body.String() != "hello world" {
		t.Errorf("expected body to be 'hello world', got %q", rec.Body.String())
	} else if rec.Header().Get("Content-Encoding") == "gzip" {
		t.Errorf("did not expect gzip encoding")
	}
}
