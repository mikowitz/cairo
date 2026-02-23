// ABOUTME: Shared test helpers for context package tests.
// ABOUTME: Provides newTestContext for creating a Context and ImageSurface with automatic cleanup.

package context

import (
	"testing"

	"github.com/mikowitz/cairo/surface"
	"github.com/stretchr/testify/require"
)

// newTestContext creates a Context backed by an ImageSurface of the given dimensions.
// Cleanup is registered via t.Cleanup; callers need not close ctx or surf manually.
func newTestContext(t *testing.T, width, height int) *Context {
	t.Helper()
	surf, err := surface.NewImageSurface(surface.FormatARGB32, width, height)
	require.NoError(t, err)
	ctx, err := NewContext(surf)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = ctx.Close()
		_ = surf.Close()
	})
	return ctx
}
