package weblib

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

// Render writes any number of templ components to the response writer with the given status.
//
// Rendering multiple components can be especially helpful when using HTMX out of band swaps:
// https://htmx.org/attributes/hx-swap-oob/
func Render(w http.ResponseWriter, r *http.Request, status int, tmpls ...templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	for _, tmpl := range tmpls {
		err := tmpl.Render(r.Context(), buf)
		if err != nil {
			return fmt.Errorf("failed to render template: %w", err)
		}
	}

	w.WriteHeader(status)
	_, err := w.Write(buf.Bytes())
	return err
}

// IsHTMX returns true if the request is HTMX, otherwise false.
func IsHTMX(r *http.Request) bool {
	return r.Header.Get("Hx-Request") == "true"
}

// ConditionalRender checks if the request is HTMX and conditionally renders the partial templ component, otherwise
// renders the page templ component.
//
// This function can be helpful in the event that a page is navigated to via the browser address bar instead of a
// HTMX boosted anchor element.
func ConditionalRender(w http.ResponseWriter, r *http.Request, status int, page, partial templ.Component) error {
	if IsHTMX(r) {
		return Render(w, r, status, partial)
	}

	return Render(w, r, status, page)
}

// Redirect checks if the request is HTMX and adds the Hx-Redirect header and reassigns the status to 200 OK.
// Otherwise, a standard HTTP redirect is performed with the given status and route.
//
// The status is changed to 200 OK if the request is HTMX due to the way that HTMX handles redirects.
// HTMX does not see 3xx status redirects and so requires a 2xx status.
// More info: https://github.com/bigskysoftware/htmx/issues/2052#issuecomment-1979805051
func Redirect(w http.ResponseWriter, r *http.Request, status int, route string) {
	if IsHTMX(r) {
		w.Header().Add("Hx-Redirect", route)
		status = http.StatusOK
	}

	http.Redirect(w, r, route, status)
}
