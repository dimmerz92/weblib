package weblib

import "testing"

func TestTrimQuotes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Test case 1: String with double quotes
		{"\"quoted\"", "quoted"},
		// Test case 2: String with single quotes
		{"'quoted'", "quoted"},
		// Test case 3: String with backtick quotes
		{"`quoted`", "quoted"},
		// Test case 4: String without any quotes
		{"noquotes", "noquotes"},
		// Test case 5: String with different characters on both ends
		{"\"wrongquote'", "\"wrongquote'"},
		// Test case 6: String with one quote on the end
		{"wrongquote'", "wrongquote'"},
		// Test case 7: Single character quotes (no trimming)
		{"\"", "\""},
		{"'", "'"},
		{"`", "`"},
		// Test case 8: Empty string
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := TrimQuotes(tt.input)
			if result != tt.expected {
				t.Errorf("TrimQuotes(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
