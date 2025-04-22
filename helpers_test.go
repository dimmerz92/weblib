package weblib

import "testing"

var empty *int
var a, b, c = 1, 0, 3

// TestCoalesce tests the Coalesce function
func TestCoalesce(t *testing.T) {
	// Test with non-nil values
	result, ok := Coalesce(&a, empty, empty)
	if !ok || *result != a {
		t.Errorf("Coalesce failed, expected %d, got %v", a, *result)
	}

	// Test with first value being nil
	result, ok = Coalesce(empty, &b, &c)
	if !ok || *result != b {
		t.Errorf("Coalesce failed, expected %d, got %v", b, *result)
	}

	// Test with all values being nil
	result, ok = Coalesce(empty, empty, empty)
	if ok || result != nil {
		t.Errorf("Coalesce failed, expected nil, got %v", result)
	}
}

// TestIIF tests the IIF function
func TestIIF(t *testing.T) {
	// Test when condition is true
	result := IIF(true, a, b)
	if result != a {
		t.Errorf("IIF failed, expected %d, got %v", a, result)
	}

	// Test when condition is false
	result = IIF(false, a, b)
	if result != b {
		t.Errorf("IIF failed, expected %d, got %v", b, result)
	}
}

// TestDefault tests the Default function
func TestDefault(t *testing.T) {
	// Test with non-zero values
	result, ok := Default(a, b, c)
	if !ok || result != a {
		t.Errorf("Default failed, expected %d, got %v", a, result)
	}

	// Test with all zero values
	result, ok = Default(b, b, b)
	if ok || result != 0 {
		t.Errorf("Default failed, expected %d, got %v", b, result)
	}
}
