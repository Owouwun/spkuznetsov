package testutils

import "testing"

func AssertError(t *testing.T, expected, actual error) {
	if actual != expected {
		t.Errorf("Expected error: '%v', got: '%v'", expected, actual)
	}
}
