package post_test

import (
	"errors"
	"testing"

	"github.com/jared-wallace/website-go/internal/repository/post"
)

// TestErrNotFound_Sentinel verifies the ErrNotFound sentinel is a distinct, unwrappable error.
// This is a pure unit test — no database required.
func TestErrNotFound_Sentinel(t *testing.T) {
	if post.ErrNotFound == nil {
		t.Fatal("ErrNotFound must not be nil")
	}

	// Wrapping should still match via errors.Is
	wrapped := errors.Join(post.ErrNotFound, errors.New("context"))
	if !errors.Is(wrapped, post.ErrNotFound) {
		t.Errorf("errors.Is(wrapped, ErrNotFound) = false, want true")
	}
}

// Integration tests below — guarded by build tag.
// Run with: go test -tags=integration ./internal/repository/post/...
