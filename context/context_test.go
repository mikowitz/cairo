package context

import (
	"testing"

	"github.com/mikowitz/cairo/status"
	"github.com/mikowitz/cairo/surface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewContext verifies that a Context can be created from an ImageSurface.
func TestNewContext(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
	require.NoError(t, err, "Failed to create surface")
	defer func() {
		err := surf.Close()
		assert.NoError(t, err, "failed to clause the surface")
	}()

	ctx, err := NewContext(surf)
	require.NoError(t, err, "Failed to create context")
	require.NotNil(t, ctx, "Context should not be nil")
	defer func() {
		err := ctx.Close()
		assert.NoError(t, err, "failed to close context")
	}()

	// Verify the context has a valid status
	st := ctx.Status()
	assert.Equal(t, status.Success, st, "Context should have Success status")
}

// TestNewContextNilSurface verifies that creating a Context with a nil surface returns an error.
func TestNewContextNilSurface(t *testing.T) {
	ctx, err := NewContext(nil)
	assert.Error(t, err, "Creating context with nil surface should return error")
	assert.Nil(t, ctx, "Context should be nil when creation fails")
	assert.Equal(t, status.NullPointer, err, "Error should be NullPointer status")
}

// TestContextClose verifies that Close works correctly and double-close is safe.
func TestContextClose(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
	require.NoError(t, err, "Failed to create surface")
	defer func() {
		err := surf.Close()
		assert.NoError(t, err, "failed to clause the surface")
	}()

	ctx, err := NewContext(surf)
	require.NoError(t, err, "Failed to create context")
	require.NotNil(t, ctx, "Context should not be nil")

	// First close should succeed
	err = ctx.Close()
	assert.NoError(t, err, "First close should not return error")

	// Second close should be safe (no-op)
	err = ctx.Close()
	assert.NoError(t, err, "Second close should not return error")

	// Operations after close should be safe (no-ops)
	ctx.Save()    // Should not panic
	ctx.Restore() // Should not panic
	st := ctx.Status()
	assert.Equal(t, status.NullPointer, st, "Status after close should be NullPointer")
}

// TestContextStatus verifies that Status returns the correct status.
func TestContextStatus(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
	require.NoError(t, err, "Failed to create surface")
	defer func() {
		err := surf.Close()
		assert.NoError(t, err, "failed to clause the surface")
	}()

	ctx, err := NewContext(surf)
	require.NoError(t, err, "Failed to create context")
	defer func() {
		err := ctx.Close()
		assert.NoError(t, err, "failed to close context")
	}()

	// New context should have Success status
	st := ctx.Status()
	assert.Equal(t, status.Success, st, "New context should have Success status")
}

// TestContextSaveRestore verifies the save/restore stack works correctly.
func TestContextSaveRestore(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
	require.NoError(t, err, "Failed to create surface")
	defer func() {
		err := surf.Close()
		assert.NoError(t, err, "failed to clause the surface")
	}()

	ctx, err := NewContext(surf)
	require.NoError(t, err, "Failed to create context")
	defer func() {
		err := ctx.Close()
		assert.NoError(t, err, "failed to close context")
	}()

	// Save should not cause errors
	ctx.Save()
	st := ctx.Status()
	assert.Equal(t, status.Success, st, "Status after Save should be Success")

	// Restore should not cause errors
	ctx.Restore()
	st = ctx.Status()
	assert.Equal(t, status.Success, st, "Status after Restore should be Success")

	// Multiple nested saves and restores
	ctx.Save()
	ctx.Save()
	ctx.Save()
	st = ctx.Status()
	assert.Equal(t, status.Success, st, "Status after multiple Saves should be Success")

	ctx.Restore()
	ctx.Restore()
	ctx.Restore()
	st = ctx.Status()
	assert.Equal(t, status.Success, st, "Status after multiple Restores should be Success")
}

// TestContextSaveRestoreImbalance verifies that restoring without a matching save causes an error.
func TestContextSaveRestoreImbalance(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
	require.NoError(t, err, "Failed to create surface")
	defer func() {
		err := surf.Close()
		assert.NoError(t, err, "failed to clause the surface")
	}()

	ctx, err := NewContext(surf)
	require.NoError(t, err, "Failed to create context")
	defer func() {
		err := ctx.Close()
		assert.NoError(t, err, "failed to close context")
	}()

	// Restore without Save should cause InvalidRestore status
	ctx.Restore()
	st := ctx.Status()
	assert.Equal(t, status.InvalidRestore, st, "Restore without Save should set InvalidRestore status")
}

// TestContextCloseIndependentOfSurface verifies that closing the context
// doesn't affect the surface and vice versa.
func TestContextCloseIndependentOfSurface(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
	require.NoError(t, err, "Failed to create surface")
	defer func() {
		err := surf.Close()
		assert.NoError(t, err, "failed to clause the surface")
	}()

	ctx, err := NewContext(surf)
	require.NoError(t, err, "Failed to create context")

	// Close context first
	err = ctx.Close()
	assert.NoError(t, err, "Closing context should not error")

	// Surface should still be valid
	st := surf.Status()
	assert.Equal(t, status.Success, st, "Surface should still have Success status after context close")

	// Can create another context from the same surface
	ctx2, err := NewContext(surf)
	require.NoError(t, err, "Should be able to create another context from same surface")
	defer func() {
		err := ctx2.Close()
		assert.NoError(t, err, "failed to close second context")
	}()

	st = ctx2.Status()
	assert.Equal(t, status.Success, st, "New context should have Success status")
}

// TestContextMultipleContextsOnSameSurface verifies that multiple contexts
// can be created on the same surface.
func TestContextMultipleContextsOnSameSurface(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
	require.NoError(t, err, "Failed to create surface")
	defer func() {
		err := surf.Close()
		assert.NoError(t, err, "failed to clause the surface")
	}()

	// Create first context
	ctx1, err := NewContext(surf)
	require.NoError(t, err, "Failed to create first context")
	defer func() {
		err := ctx1.Close()
		assert.NoError(t, err, "failed to close first context")
	}()

	// Create second context on same surface
	ctx2, err := NewContext(surf)
	require.NoError(t, err, "Failed to create second context")
	defer func() {
		err := ctx2.Close()
		assert.NoError(t, err, "failed to close second context")
	}()

	// Both contexts should have Success status
	st1 := ctx1.Status()
	assert.Equal(t, status.Success, st1, "First context should have Success status")

	st2 := ctx2.Status()
	assert.Equal(t, status.Success, st2, "Second context should have Success status")
}

// TestContextCreationWithDifferentSurfaceFormats verifies that contexts
// can be created with different surface formats.
func TestContextCreationWithDifferentSurfaceFormats(t *testing.T) {
	formats := []struct {
		name   string
		format surface.Format
	}{
		{"ARGB32", surface.FormatARGB32},
		{"RGB24", surface.FormatRGB24},
		{"A8", surface.FormatA8},
		{"A1", surface.FormatA1},
	}

	for _, tc := range formats {
		t.Run(tc.name, func(t *testing.T) {
			surf, err := surface.NewImageSurface(tc.format, 100, 100)
			require.NoError(t, err, "Failed to create %s surface", tc.name)
			defer func() {
				err := surf.Close()
				assert.NoError(t, err, "failed to clause the surface")
			}()

			ctx, err := NewContext(surf)
			require.NoError(t, err, "Failed to create context for %s surface", tc.name)
			defer func() {
				err := ctx.Close()
				assert.NoError(t, err, "failed to close context")
			}()

			st := ctx.Status()
			assert.Equal(t, status.Success, st, "Context for %s surface should have Success status", tc.name)
		})
	}
}

// TestContextSetSourceRGB verifies that SetSourceRGB sets a color and maintains success status.
func TestContextSetSourceRGB(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
	require.NoError(t, err, "Failed to create surface")
	defer func() {
		err := surf.Close()
		assert.NoError(t, err, "Failed to close surface")
	}()

	ctx, err := NewContext(surf)
	require.NoError(t, err, "Failed to create context")
	defer func() {
		err := ctx.Close()
		assert.NoError(t, err, "Failed to close context")
	}()

	// Test various color combinations
	testCases := []struct {
		name    string
		r, g, b float64
	}{
		{"Red", 1.0, 0.0, 0.0},
		{"Green", 0.0, 1.0, 0.0},
		{"Blue", 0.0, 0.0, 1.0},
		{"White", 1.0, 1.0, 1.0},
		{"Black", 0.0, 0.0, 0.0},
		{"Half intensity gray", 0.5, 0.5, 0.5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx.SetSourceRGB(tc.r, tc.g, tc.b)
			st := ctx.Status()
			assert.Equal(t, status.Success, st, "Status should be Success after SetSourceRGB")
		})
	}
}

// TestContextSetSourceRGBA verifies that SetSourceRGBA sets a color with alpha and maintains success status.
func TestContextSetSourceRGBA(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
	require.NoError(t, err, "Failed to create surface")
	defer func() {
		err := surf.Close()
		assert.NoError(t, err, "Failed to close surface")
	}()

	ctx, err := NewContext(surf)
	require.NoError(t, err, "Failed to create context")
	defer func() {
		err := ctx.Close()
		assert.NoError(t, err, "failed to close context")
	}()

	// Test various color and alpha combinations
	testCases := []struct {
		name       string
		r, g, b, a float64
	}{
		{"Opaque red", 1.0, 0.0, 0.0, 1.0},
		{"Semi-transparent blue", 0.0, 0.0, 1.0, 0.5},
		{"Fully transparent green", 0.0, 1.0, 0.0, 0.0},
		{"Quarter opacity white", 1.0, 1.0, 1.0, 0.25},
		{"Three-quarter opacity black", 0.0, 0.0, 0.0, 0.75},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx.SetSourceRGBA(tc.r, tc.g, tc.b, tc.a)
			st := ctx.Status()
			assert.Equal(t, status.Success, st, "Status should be Success after SetSourceRGBA")
		})
	}
}

// TestContextSetSourceAfterClose verifies that setting source after close is safe.
func TestContextSetSourceAfterClose(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
	require.NoError(t, err, "Failed to create surface")
	defer func() {
		err := surf.Close()
		assert.NoError(t, err, "Failed to close surface")
	}()

	ctx, err := NewContext(surf)
	require.NoError(t, err, "Failed to create context")

	// Close the context
	err = ctx.Close()
	assert.NoError(t, err, "Closing context should not error")

	// Setting source after close should be safe (no-op)
	ctx.SetSourceRGB(1.0, 0.0, 0.0)
	ctx.SetSourceRGBA(0.0, 1.0, 0.0, 0.5)

	// Status should indicate closed/null pointer
	st := ctx.Status()
	assert.Equal(t, status.NullPointer, st, "Status after close should be NullPointer")
}
