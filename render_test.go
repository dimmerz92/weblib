package weblib

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockComponent struct {
	content string
}

func (m mockComponent) Render(ctx context.Context, w io.Writer) error {
	_, err := w.Write([]byte(m.content))
	return err
}

type errorComponent struct{}

func (e errorComponent) Render(ctx context.Context, w io.Writer) error {
	return fmt.Errorf("error")
}

var (
	page    = mockComponent{content: "page"}
	partial = mockComponent{content: "partial"}
	errComp = errorComponent{}
)

func TestIsHTMX(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Hx-Request", "true")
	if !IsHTMX(req) {
		t.Error("expected request to be HTMX")
	}

	req.Header.Set("Hx-Request", "false")
	if IsHTMX(req) {
		t.Error("expected request to not be HTMX")
	}
}

func TestRender(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	err := Render(w, r, http.StatusOK, page)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != page.content {
		t.Errorf("expected body to be '%s', got %s", page.content, w.Body.String())
	}
}

func TestRender_MultipleComponents(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	err := Render(w, r, http.StatusOK, page, partial)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := page.content + partial.content
	if w.Body.String() != expected {
		t.Errorf("expected combined content %q, got %q", expected, w.Body.String())
	}
}

func TestRender_WithError(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	err := Render(w, r, http.StatusOK, errComp)
	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestConditionalRender_HTMX(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Hx-Request", "true")

	err := ConditionalRender(w, r, http.StatusOK, page, partial)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if w.Body.String() != partial.content {
		t.Errorf("expected partial content, got %q", w.Body.String())
	}
}

func TestConditionalRender_NonHTMX(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	err := ConditionalRender(w, r, http.StatusOK, page, partial)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if w.Body.String() != page.content {
		t.Errorf("expected page content, got %q", w.Body.String())
	}
}

func TestRedirect_HTMX(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Hx-Request", "true")

	Redirect(w, r, http.StatusFound, "/new-location")

	if w.Header().Get("Hx-Redirect") != "/new-location" {
		t.Errorf("expected Hx-Redirect header, got %q", w.Header().Get("Hx-Redirect"))
	}
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 OK for HTMX, got %d", w.Code)
	}
}

func TestRedirect_NonHTMX(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	Redirect(w, r, http.StatusFound, "/new-location")

	if w.Header().Get("Location") != "/new-location" {
		t.Errorf("expected Location header, got %q", w.Header().Get("Location"))
	}
	if w.Code != http.StatusFound {
		t.Errorf("expected status 302 Found, got %d", w.Code)
	}
}
