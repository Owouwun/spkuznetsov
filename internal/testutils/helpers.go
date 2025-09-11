package testutils

import (
	"errors"
	"testing"
)

func AssertError(t *testing.T, expected, actual error) {
	if !errors.Is(expected, actual) {
		t.Errorf("expected error: '%v', got: '%v'", expected, actual)
		return
	}
}
