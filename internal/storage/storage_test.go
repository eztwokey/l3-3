package storage

import (
	"testing"
)

func TestErrNotFound(t *testing.T) {
	if ErrNotFound.Error() != "not found" {
		t.Errorf("expected 'not found', got %q", ErrNotFound.Error())
	}
}
