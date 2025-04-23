package weblib

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeFiles_RootProject(t *testing.T) {
	tests := []struct {
		name         string
		browsable    bool
		url          string
		expectedCode int
	}{
		{
			name:         "Browsable - file access allowed",
			browsable:    true,
			url:          "/static/go.mod",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Non-browsable - file access allowed",
			browsable:    false,
			url:          "/static/go.mod",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Browsable - directory access",
			browsable:    true,
			url:          "/static/",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Non-browsable - directory access blocked",
			browsable:    false,
			url:          "/static/",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := ServeFiles("/static/", ".", tt.browsable)

			req := httptest.NewRequest("GET", tt.url, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("Expected HTTP %d, got %d", tt.expectedCode, rr.Code)
			}
		})
	}
}
