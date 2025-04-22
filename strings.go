package weblib

// TrimQuotes trims single, double, or backtick quotes from a string and returns it.
func TrimQuotes(s string) string {
	if len(s) >= 2 {
		lastChar := s[len(s)-1]
		if lastChar == s[0] && (lastChar == '"' || lastChar == '\'' || lastChar == '`') {
			return s[1 : len(s)-1]
		}
	}

	return s
}
