package weblib

// Coalesce returns the first non nil value and true if it exists, otherwise nil and false.
func Coalesce[T any](vals ...*T) (*T, bool) {
	for _, v := range vals {
		if v != nil {
			return v, true
		}
	}

	return nil, false
}

// IIF is an inline if statement. if the condition is true, v1 is returned, otherwise v2 is returned.
func IIF[T any](condition bool, v1, v2 T) T {
	if condition {
		return v1
	}

	return v2
}

// Default returns the first non-zero value and true if it exists, otherwise a zero value and false.
func Default[T comparable](vals ...T) (T, bool) {
	var zero T
	for _, v := range vals {
		if v != zero {
			return v, true
		}
	}

	return zero, false
}
