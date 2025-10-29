package context

import (
	"runtime"
	"testing"

	"github.com/mikowitz/cairo/pattern"
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

// TestContextMoveTo verifies that MoveTo sets the current point.
func TestContextMoveTo(t *testing.T) {
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

	// MoveTo should set current point
	ctx.MoveTo(50.0, 75.0)
	st := ctx.Status()
	assert.Equal(t, status.Success, st, "Status should be Success after MoveTo")

	// Verify current point is set correctly
	x, y, err := ctx.GetCurrentPoint()
	require.NoError(t, err, "GetCurrentPoint should not error after MoveTo")
	assert.InDelta(t, 50.0, x, 0.001, "X coordinate should match")
	assert.InDelta(t, 75.0, y, 0.001, "Y coordinate should match")
}

// TestContextLineTo verifies that LineTo adds a line and updates the current point.
func TestContextLineTo(t *testing.T) {
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

	// Start with MoveTo to establish a current point
	ctx.MoveTo(10.0, 20.0)

	// LineTo should create a line and update current point
	ctx.LineTo(50.0, 60.0)
	st := ctx.Status()
	assert.Equal(t, status.Success, st, "Status should be Success after LineTo")

	// Verify current point is updated
	x, y, err := ctx.GetCurrentPoint()
	require.NoError(t, err, "GetCurrentPoint should not error after LineTo")
	assert.InDelta(t, 50.0, x, 0.001, "X coordinate should match LineTo destination")
	assert.InDelta(t, 60.0, y, 0.001, "Y coordinate should match LineTo destination")

	// Multiple LineTo calls should work
	ctx.LineTo(80.0, 90.0)
	x, y, err = ctx.GetCurrentPoint()
	require.NoError(t, err, "GetCurrentPoint should not error after second LineTo")
	assert.InDelta(t, 80.0, x, 0.001, "X coordinate should match second LineTo")
	assert.InDelta(t, 90.0, y, 0.001, "Y coordinate should match second LineTo")
}

// TestContextRectangle verifies that Rectangle creates a rectangular path.
func TestContextRectangle(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 200, 200)
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

	// Rectangle should create a closed rectangular path
	ctx.Rectangle(10.0, 20.0, 100.0, 50.0)
	st := ctx.Status()
	assert.Equal(t, status.Success, st, "Status should be Success after Rectangle")

	// Test various rectangle dimensions
	testCases := []struct {
		name       string
		x, y, w, h float64
	}{
		{"Square", 0.0, 0.0, 50.0, 50.0},
		{"Wide rectangle", 10.0, 10.0, 100.0, 20.0},
		{"Tall rectangle", 20.0, 30.0, 30.0, 80.0},
		{"Single pixel", 50.0, 50.0, 1.0, 1.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx.Rectangle(tc.x, tc.y, tc.w, tc.h)
			st := ctx.Status()
			assert.Equal(t, status.Success, st, "Status should be Success after Rectangle")
		})
	}
}

// TestContextClosePath verifies that ClosePath closes the current path.
func TestContextClosePath(t *testing.T) {
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

	// Create a path
	ctx.MoveTo(10.0, 10.0)
	ctx.LineTo(50.0, 10.0)
	ctx.LineTo(50.0, 50.0)

	// Close the path
	ctx.ClosePath()
	st := ctx.Status()
	assert.Equal(t, status.Success, st, "Status should be Success after ClosePath")

	// ClosePath on empty path should be safe
	ctx.NewPath()
	ctx.ClosePath()
	st = ctx.Status()
	assert.Equal(t, status.Success, st, "ClosePath on empty path should not error")
}

// TestContextNewPath verifies that NewPath clears the current path.
func TestContextNewPath(t *testing.T) {
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

	// Create a path with a current point
	ctx.MoveTo(50.0, 50.0)
	x, y, err := ctx.GetCurrentPoint()
	require.NoError(t, err, "Should have current point after MoveTo")
	assert.InDelta(t, 50.0, x, 0.001, "X should be 50")
	assert.InDelta(t, 50.0, y, 0.001, "Y should be 50")

	// NewPath should clear the path and current point
	ctx.NewPath()
	st := ctx.Status()
	assert.Equal(t, status.Success, st, "Status should be Success after NewPath")

	// Current point should no longer be defined
	_, _, err = ctx.GetCurrentPoint()
	assert.Error(t, err, "GetCurrentPoint should error after NewPath clears the path")

	// Multiple NewPath calls should be safe
	ctx.NewPath()
	ctx.NewPath()
	st = ctx.Status()
	assert.Equal(t, status.Success, st, "Multiple NewPath calls should succeed")
}

// TestContextGetCurrentPoint verifies getting the current point.
func TestContextGetCurrentPoint(t *testing.T) {
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

	// Set current point with MoveTo
	ctx.MoveTo(25.5, 37.75)
	x, y, err := ctx.GetCurrentPoint()
	require.NoError(t, err, "GetCurrentPoint should not error")
	assert.InDelta(t, 25.5, x, 0.001, "X coordinate should match")
	assert.InDelta(t, 37.75, y, 0.001, "Y coordinate should match")

	// Update with LineTo
	ctx.LineTo(100.0, 200.0)
	x, y, err = ctx.GetCurrentPoint()
	require.NoError(t, err, "GetCurrentPoint should not error after LineTo")
	assert.InDelta(t, 100.0, x, 0.001, "X coordinate should be updated")
	assert.InDelta(t, 200.0, y, 0.001, "Y coordinate should be updated")
}

func TestContextHasCurrentPointNoPoint(t *testing.T) {
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

	// No current point initially
	assert.False(t, ctx.HasCurrentPoint(), "HasCurrentPoint should be false when no current point")

	// Set a current point
	ctx.MoveTo(10.0, 20.0)
	assert.True(t, ctx.HasCurrentPoint(), "HasCurrentPoint should be true when current point is set")

	// NewPath clears the current point
	ctx.NewPath()
	assert.False(t, ctx.HasCurrentPoint(), "HasCurrentPoint should be false when no current point")

	// Rectangle creates a path and Cairo should have a defined current point
	// after a closed subpath like Rectangle
	ctx.Rectangle(10.0, 10.0, 50.0, 50.0)
	assert.True(t, ctx.HasCurrentPoint(), "HasCurrentPoint should be true when current point is set")
}

// TestContextGetCurrentPointNoPoint verifies error when no current point exists.
func TestContextGetCurrentPointNoPoint(t *testing.T) {
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

	// No current point initially
	_, _, err = ctx.GetCurrentPoint()
	assert.Error(t, err, "GetCurrentPoint should error when no current point")

	// Set a current point
	ctx.MoveTo(10.0, 20.0)
	_, _, err = ctx.GetCurrentPoint()
	require.NoError(t, err, "Should have current point after MoveTo")

	// NewPath clears the current point
	ctx.NewPath()
	_, _, err = ctx.GetCurrentPoint()
	assert.Error(t, err, "GetCurrentPoint should error after NewPath")

	// Rectangle creates a path and Cairo should have a defined current point
	// after a closed subpath like Rectangle
	ctx.Rectangle(10.0, 10.0, 50.0, 50.0)
	_, _, err = ctx.GetCurrentPoint()
	require.NoError(t, err, "Should have current point after MoveTo")
}

// TestContextPathOperationsAfterClose verifies path operations are safe after close.
func TestContextPathOperationsAfterClose(t *testing.T) {
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

	// All path operations should be safe no-ops after close
	ctx.MoveTo(10.0, 20.0)
	ctx.LineTo(30.0, 40.0)
	ctx.Rectangle(5.0, 5.0, 20.0, 20.0)
	ctx.ClosePath()
	ctx.NewPath()
	ctx.NewSubPath()

	// GetCurrentPoint should return error indicating closed context
	_, _, err = ctx.GetCurrentPoint()
	assert.Error(t, err, "GetCurrentPoint should error on closed context")

	// Status should indicate closed/null pointer
	st := ctx.Status()
	assert.Equal(t, status.NullPointer, st, "Status after close should be NullPointer")
}

// TestContextFill verifies that Fill renders and consumes the current path.
func TestContextFill(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 200, 200)
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

	// Create a path and fill it
	ctx.Rectangle(50.0, 50.0, 100.0, 100.0)
	ctx.SetSourceRGB(1.0, 0.0, 0.0) // Red
	ctx.Fill()

	// Verify status is still success
	st := ctx.Status()
	assert.Equal(t, status.Success, st, "Status should be Success after Fill")

	// After Fill, there should be no current point (path consumed)
	assert.False(t, ctx.HasCurrentPoint(), "Fill should consume the path, removing current point")
}

// TestContextFillPreserve verifies that FillPreserve renders but keeps the path.
func TestContextFillPreserve(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 200, 200)
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

	// Create a path
	ctx.Rectangle(50.0, 50.0, 100.0, 100.0)
	ctx.SetSourceRGB(0.0, 1.0, 0.0) // Green

	// FillPreserve should keep the path
	ctx.FillPreserve()

	st := ctx.Status()
	assert.Equal(t, status.Success, st, "Status should be Success after FillPreserve")

	// After FillPreserve, current point should still exist (path preserved)
	assert.True(t, ctx.HasCurrentPoint(), "FillPreserve should preserve the path and current point")

	// We should be able to stroke the same path
	ctx.SetSourceRGB(0.0, 0.0, 1.0) // Blue
	ctx.Stroke()
	st = ctx.Status()
	assert.Equal(t, status.Success, st, "Should be able to stroke after FillPreserve")
}

// TestContextStroke verifies that Stroke renders and consumes the current path.
func TestContextStroke(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 200, 200)
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

	// Create a path and stroke it
	ctx.SetLineWidth(2.0)
	ctx.Rectangle(50.0, 50.0, 100.0, 100.0)
	ctx.SetSourceRGB(0.0, 0.0, 1.0) // Blue
	ctx.Stroke()

	st := ctx.Status()
	assert.Equal(t, status.Success, st, "Status should be Success after Stroke")

	// After Stroke, there should be no current point (path consumed)
	assert.False(t, ctx.HasCurrentPoint(), "Stroke should consume the path, removing current point")
}

// TestContextStrokePreserve verifies that StrokePreserve renders but keeps the path.
func TestContextStrokePreserve(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 200, 200)
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

	// Create a path
	ctx.SetLineWidth(3.0)
	ctx.Rectangle(50.0, 50.0, 100.0, 100.0)
	ctx.SetSourceRGB(1.0, 0.0, 1.0) // Magenta

	// StrokePreserve should keep the path
	ctx.StrokePreserve()

	st := ctx.Status()
	assert.Equal(t, status.Success, st, "Status should be Success after StrokePreserve")

	// After StrokePreserve, current point should still exist (path preserved)
	assert.True(t, ctx.HasCurrentPoint(), "StrokePreserve should preserve the path and current point")

	// We should be able to fill the same path
	ctx.SetSourceRGBA(1.0, 1.0, 0.0, 0.5) // Semi-transparent yellow
	ctx.Fill()
	st = ctx.Status()
	assert.Equal(t, status.Success, st, "Should be able to fill after StrokePreserve")
}

// TestContextPaint verifies that Paint paints the current source everywhere.
func TestContextPaint(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 200, 200)
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

	// Paint with a solid color
	ctx.SetSourceRGB(0.5, 0.5, 0.5) // Gray
	ctx.Paint()

	st := ctx.Status()
	assert.Equal(t, status.Success, st, "Status should be Success after Paint")

	// Paint with transparency
	ctx.SetSourceRGBA(1.0, 0.0, 0.0, 0.3) // Semi-transparent red
	ctx.Paint()

	st = ctx.Status()
	assert.Equal(t, status.Success, st, "Status should be Success after Paint with alpha")
}

// TestContextSetLineWidth verifies that SetLineWidth sets the line width for stroking.
func TestContextSetLineWidth(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 200, 200)
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

	// Test various line widths
	testCases := []struct {
		name  string
		width float64
	}{
		{"Thin line", 1.0},
		{"Medium line", 5.0},
		{"Thick line", 10.0},
		{"Very thin", 0.5},
		{"Very thick", 20.0},
		{"Negative width", -5.0}, // Cairo clamps negative to 0
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx.SetLineWidth(tc.width)
			st := ctx.Status()
			assert.Equal(t, status.Success, st, "Status should be Success after SetLineWidth")

			// Draw a line with this width
			ctx.MoveTo(10.0, 10.0)
			ctx.LineTo(100.0, 100.0)
			ctx.Stroke()

			st = ctx.Status()
			assert.Equal(t, status.Success, st, "Status should be Success after Stroke")
		})
	}
}

// TestContextGetLineWidth verifies GetLineWidth returns correct values in different scenarios.
func TestContextGetLineWidth(t *testing.T) {
	t.Run("Default width on new context", func(t *testing.T) {
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

		// Default line width should be 2.0
		width := ctx.GetLineWidth()
		assert.Equal(t, 2.0, width, "Default line width should be 2.0")
	})

	t.Run("Returns set width", func(t *testing.T) {
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

		// Test various widths
		testWidths := []float64{1.0, 5.0, 10.0, 0.5, 100.0}
		for _, expectedWidth := range testWidths {
			ctx.SetLineWidth(expectedWidth)
			actualWidth := ctx.GetLineWidth()
			assert.Equal(t, expectedWidth, actualWidth, "GetLineWidth should return the width that was set")
		}
	})

	t.Run("Negative width clamped to zero", func(t *testing.T) {
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

		// Set negative width - Cairo should clamp to 0
		ctx.SetLineWidth(-5.0)
		width := ctx.GetLineWidth()
		assert.Equal(t, 0.0, width, "Negative line width should be clamped to 0.0")
	})

	t.Run("Returns zero after close", func(t *testing.T) {
		surf, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err, "Failed to create surface")
		defer func() {
			err := surf.Close()
			assert.NoError(t, err, "Failed to close surface")
		}()

		ctx, err := NewContext(surf)
		require.NoError(t, err, "Failed to create context")

		// Set a width before closing
		ctx.SetLineWidth(10.0)

		// Close the context
		err = ctx.Close()
		assert.NoError(t, err, "Closing context should not error")

		// GetLineWidth after close should return 0.0
		width := ctx.GetLineWidth()
		assert.Equal(t, 0.0, width, "GetLineWidth after close should return 0.0")
	})
}

// TestContextRenderAfterClose verifies render operations are safe after close.
func TestContextRenderAfterClose(t *testing.T) {
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

	// All render operations should be safe no-ops after close
	ctx.Fill()
	ctx.FillPreserve()
	ctx.Stroke()
	ctx.StrokePreserve()
	ctx.Paint()
	ctx.SetLineWidth(5.0)

	// Status should indicate closed/null pointer
	st := ctx.Status()
	assert.Equal(t, status.NullPointer, st, "Status after close should be NullPointer")
}

// TestContextIntegrationFillStroke is an integration test combining path operations with rendering.
func TestContextIntegrationFillStroke(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 300, 300)
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

	// Integration test: Create path, set color, and fill
	ctx.NewPath()
	ctx.Rectangle(50.0, 50.0, 100.0, 100.0)
	ctx.SetSourceRGB(1.0, 0.0, 0.0) // Red
	ctx.Fill()

	st := ctx.Status()
	assert.Equal(t, status.Success, st, "Integration test should complete successfully")

	// Draw another shape with stroke
	ctx.NewPath()
	ctx.Rectangle(175.0, 175.0, 100.0, 100.0)
	ctx.SetSourceRGB(0.0, 0.0, 1.0) // Blue
	ctx.SetLineWidth(3.0)
	ctx.Stroke()

	st = ctx.Status()
	assert.Equal(t, status.Success, st, "Second shape should complete successfully")

	// Test FillPreserve + Stroke on same path
	ctx.NewPath()
	ctx.Rectangle(100.0, 175.0, 50.0, 50.0)
	ctx.SetSourceRGBA(0.0, 1.0, 0.0, 0.7) // Semi-transparent green
	ctx.FillPreserve()
	ctx.SetSourceRGB(0.0, 0.0, 0.0) // Black outline
	ctx.SetLineWidth(2.0)
	ctx.Stroke()

	st = ctx.Status()
	assert.Equal(t, status.Success, st, "FillPreserve + Stroke combination should work")
}

// TestContextSetSource verifies that SetSource sets a pattern as the source.
func TestContextSetSource(t *testing.T) {
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

	t.Run("set_solid_pattern_rgb", func(t *testing.T) {
		// Create a solid RGB pattern
		pat, err := pattern.NewSolidPatternRGB(1.0, 0.0, 0.0) // Red
		require.NoError(t, err, "Failed to create solid pattern")
		defer pat.Close()

		// Set it as the source
		ctx.SetSource(pat)

		// Verify status is still success
		st := ctx.Status()
		assert.Equal(t, status.Success, st, "Status should be Success after SetSource")
	})

	t.Run("set_solid_pattern_rgba", func(t *testing.T) {
		// Create a solid RGBA pattern
		pat, err := pattern.NewSolidPatternRGBA(0.0, 1.0, 0.0, 0.5) // Semi-transparent green
		require.NoError(t, err, "Failed to create solid pattern")
		defer pat.Close()

		// Set it as the source
		ctx.SetSource(pat)

		// Verify status is still success
		st := ctx.Status()
		assert.Equal(t, status.Success, st, "Status should be Success after SetSource")
	})

	t.Run("set_multiple_patterns", func(t *testing.T) {
		// Set multiple patterns in sequence
		pat1, err := pattern.NewSolidPatternRGB(1.0, 0.0, 0.0)
		require.NoError(t, err)
		defer pat1.Close()

		ctx.SetSource(pat1)
		assert.Equal(t, status.Success, ctx.Status())

		pat2, err := pattern.NewSolidPatternRGBA(0.0, 0.0, 1.0, 0.8)
		require.NoError(t, err)
		defer pat2.Close()

		ctx.SetSource(pat2)
		assert.Equal(t, status.Success, ctx.Status())
	})
}

// TestContextGetSource verifies that GetSource returns the current source pattern.
func TestContextGetSource(t *testing.T) {
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

	t.Run("get_default_source", func(t *testing.T) {
		// GetSource should work even on a newly created context
		src, err := ctx.GetSource()
		require.NoError(t, err, "GetSource should not return an error")
		require.NotNil(t, src, "GetSource should return a non-nil pattern")

		// The default source should have success status
		st := src.Status()
		assert.Equal(t, status.Success, st, "Default source should have Success status")
	})

	t.Run("get_after_set_pattern", func(t *testing.T) {
		// Create and set a solid pattern
		pat, err := pattern.NewSolidPatternRGB(1.0, 0.0, 0.0)
		require.NoError(t, err)
		defer pat.Close()

		ctx.SetSource(pat)

		// Get the source back
		src, err := ctx.GetSource()
		require.NoError(t, err, "GetSource should not return an error")
		require.NotNil(t, src, "GetSource should return a non-nil pattern")
		assert.Equal(t, status.Success, src.Status(), "Retrieved source should have Success status")
	})
}

// TestContextGetSourceAfterSetSourceRGB verifies that GetSource returns a SolidPattern
// after using SetSourceRGB.
func TestContextGetSourceAfterSetSourceRGB(t *testing.T) {
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

	// Test various colors
	testCases := []struct {
		name    string
		r, g, b float64
	}{
		{"red", 1.0, 0.0, 0.0},
		{"green", 0.0, 1.0, 0.0},
		{"blue", 0.0, 0.0, 1.0},
		{"white", 1.0, 1.0, 1.0},
		{"black", 0.0, 0.0, 0.0},
		{"gray", 0.5, 0.5, 0.5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set source using RGB
			ctx.SetSourceRGB(tc.r, tc.g, tc.b)

			// Get the source - should be a solid pattern
			src, err := ctx.GetSource()
			require.NoError(t, err, "GetSource should not return an error after SetSourceRGB")
			require.NotNil(t, src, "GetSource should return a non-nil pattern after SetSourceRGB")

			// Verify it's a pattern with success status
			st := src.Status()
			assert.Equal(t, status.Success, st, "Source pattern should have Success status")

			// The pattern should be usable (we can't easily verify the color,
			// but we can verify it's a valid pattern by checking its status)
			assert.NotNil(t, src, "Retrieved source should be a valid pattern")
		})
	}
}

// TestContextGetSourceAfterSetSourceRGBA verifies that GetSource returns a SolidPattern
// after using SetSourceRGBA.
func TestContextGetSourceAfterSetSourceRGBA(t *testing.T) {
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

	// Test various colors with alpha
	testCases := []struct {
		name       string
		r, g, b, a float64
	}{
		{"opaque_red", 1.0, 0.0, 0.0, 1.0},
		{"transparent_green", 0.0, 1.0, 0.0, 0.0},
		{"semi_transparent_blue", 0.0, 0.0, 1.0, 0.5},
		{"quarter_opacity_white", 1.0, 1.0, 1.0, 0.25},
		{"three_quarter_opacity_black", 0.0, 0.0, 0.0, 0.75},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set source using RGBA
			ctx.SetSourceRGBA(tc.r, tc.g, tc.b, tc.a)

			// Get the source - should be a solid pattern
			src, err := ctx.GetSource()
			require.NoError(t, err, "GetSource should not return an error after SetSourceRGBA")
			require.NotNil(t, src, "GetSource should return a non-nil pattern after SetSourceRGBA")

			// Verify it's a pattern with success status
			st := src.Status()
			assert.Equal(t, status.Success, st, "Source pattern should have Success status")

			// The pattern should be usable
			assert.NotNil(t, src, "Retrieved source should be a valid pattern")
		})
	}
}

// TestContextSetSourcePatternAfterClose verifies that SetSource is safe after close.
func TestContextSetSourcePatternAfterClose(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
	require.NoError(t, err, "Failed to create surface")
	defer func() {
		err := surf.Close()
		assert.NoError(t, err, "Failed to close surface")
	}()

	ctx, err := NewContext(surf)
	require.NoError(t, err, "Failed to create context")

	// Create a pattern
	pat, err := pattern.NewSolidPatternRGB(1.0, 0.0, 0.0)
	require.NoError(t, err)
	defer pat.Close()

	// Close the context
	err = ctx.Close()
	assert.NoError(t, err, "Closing context should not error")

	// SetSource after close should be safe (no-op)
	ctx.SetSource(pat)

	// Status should indicate closed/null pointer
	st := ctx.Status()
	assert.Equal(t, status.NullPointer, st, "Status after close should be NullPointer")
}

// TestContextGetSourceAfterClose verifies that GetSource returns an error after close.
func TestContextGetSourceAfterClose(t *testing.T) {
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

	// GetSource after close should return an error
	src, err := ctx.GetSource()
	assert.Error(t, err, "GetSource after close should return an error")
	assert.Nil(t, src, "GetSource after close should return nil pattern")
}

// TestContextSourcePatternIntegration verifies integration between SetSource,
// GetSource, and drawing operations.
func TestContextSourcePatternIntegration(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 200, 200)
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

	t.Run("set_source_and_draw", func(t *testing.T) {
		// Create a solid pattern
		pat, err := pattern.NewSolidPatternRGB(1.0, 0.0, 0.0) // Red
		require.NoError(t, err)
		defer pat.Close()

		// Set as source
		ctx.SetSource(pat)

		// Draw with the pattern
		ctx.Rectangle(10.0, 10.0, 50.0, 50.0)
		ctx.Fill()

		// Should complete successfully
		st := ctx.Status()
		assert.Equal(t, status.Success, st, "Drawing with pattern source should succeed")
	})

	t.Run("set_source_rgb_get_and_reuse", func(t *testing.T) {
		// Set source using convenience method
		ctx.SetSourceRGB(0.0, 1.0, 0.0) // Green

		// Get the pattern that was created
		src, err := ctx.GetSource()
		require.NoError(t, err, "GetSource should not return an error")
		require.NotNil(t, src, "Should get a pattern from SetSourceRGB")

		// Draw with it
		ctx.Rectangle(70.0, 10.0, 50.0, 50.0)
		ctx.Fill()

		st := ctx.Status()
		assert.Equal(t, status.Success, st, "Drawing should succeed")

		// The source should still be accessible
		src2, err := ctx.GetSource()
		require.NoError(t, err, "GetSource should not return an error after drawing")
		assert.NotNil(t, src2, "Source should still be accessible after drawing")
	})

	t.Run("set_source_rgba_get_and_reuse", func(t *testing.T) {
		// Set source using convenience method with alpha
		ctx.SetSourceRGBA(0.0, 0.0, 1.0, 0.7) // Semi-transparent blue

		// Get the pattern that was created
		src, err := ctx.GetSource()
		require.NoError(t, err, "GetSource should not return an error")
		require.NotNil(t, src, "Should get a pattern from SetSourceRGBA")

		// Draw with it
		ctx.Rectangle(130.0, 10.0, 50.0, 50.0)
		ctx.Fill()

		st := ctx.Status()
		assert.Equal(t, status.Success, st, "Drawing should succeed")

		// The source should still be accessible
		src2, err := ctx.GetSource()
		require.NoError(t, err, "GetSource should not return an error after drawing")
		assert.NotNil(t, src2, "Source should still be accessible after drawing")
	})

	t.Run("switch_between_patterns_and_rgb", func(t *testing.T) {
		// Create a pattern
		pat, err := pattern.NewSolidPatternRGBA(1.0, 1.0, 0.0, 0.5) // Yellow
		require.NoError(t, err)
		defer pat.Close()

		// Set pattern
		ctx.SetSource(pat)
		ctx.Rectangle(10.0, 70.0, 30.0, 30.0)
		ctx.Fill()

		// Switch to RGB
		ctx.SetSourceRGB(1.0, 0.0, 1.0) // Magenta
		ctx.Rectangle(50.0, 70.0, 30.0, 30.0)
		ctx.Fill()

		// Switch to RGBA
		ctx.SetSourceRGBA(0.0, 1.0, 1.0, 0.8) // Cyan
		ctx.Rectangle(90.0, 70.0, 30.0, 30.0)
		ctx.Fill()

		// All should succeed
		st := ctx.Status()
		assert.Equal(t, status.Success, st, "Switching between source types should work")
	})
}

// TestContextGetSourceBorrowedReference is a regression test for proper reference
// counting when GetSource() returns a borrowed reference from Cairo.
//
// Background:
// The cairo_get_source() C function returns a BORROWED reference - meaning the
// pattern is owned by the context and the caller should NOT destroy it.
// However, our Go wrapper (PatternFromC) sets up a finalizer that WILL destroy
// the pattern when the Go object is garbage collected.
//
// This test verifies that:
// 1. GetSource() properly increments the pattern's reference count (via cairo_pattern_reference)
// 2. The returned Go Pattern owns its own reference
// 3. The Context still owns its original reference
// 4. No double-free occurs when both the Pattern and Context are destroyed
//
// Test scenario:
// - Create a context (which has a default source pattern)
// - Call GetSource() to get the pattern (borrowed reference)
// - Let the returned pattern go out of scope WITHOUT explicit Close()
// - Force GC to run the pattern's finalizer
// - Close the context
// - If reference counting is correct, no crash/double-free occurs
func TestContextGetSourceBorrowedReference(t *testing.T) {
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

	t.Run("get_source_without_explicit_close", func(t *testing.T) {
		// Get the source pattern (borrowed reference from Cairo's perspective)
		src, err := ctx.GetSource()
		require.NoError(t, err, "GetSource should not return an error")
		require.NotNil(t, src, "GetSource should return a non-nil pattern")

		// Verify the pattern is valid
		st := src.Status()
		assert.Equal(t, status.Success, st, "Source pattern should have Success status")

		// Key point: We do NOT call src.Close() here
		// The pattern goes out of scope and its finalizer will run during GC
		// If reference counting is wrong, this will cause a double-free:
		// 1. Pattern finalizer calls cairo_pattern_destroy
		// 2. Context close also calls cairo_pattern_destroy (via cairo_destroy)

		// Force garbage collection to run the pattern's finalizer
		// This simulates what would happen naturally over time
		src = nil
		runtime.GC()
		runtime.GC() // Run twice to ensure finalizers execute

		// If we reach here without crashing, reference counting is working
		// The context should still be valid
		st = ctx.Status()
		assert.Equal(t, status.Success, st, "Context should still be valid after pattern GC")
	})

	t.Run("get_source_multiple_times_without_close", func(t *testing.T) {
		// This tests that we can safely get the source multiple times
		// and let them all be GC'd without explicit Close()

		for i := 0; i < 5; i++ {
			src, err := ctx.GetSource()
			require.NoError(t, err, "GetSource iteration %d should not error", i)
			require.NotNil(t, src, "GetSource iteration %d should return pattern", i)

			// Let it go out of scope - no explicit Close()
			_ = src
		}

		// Force GC to clean up all the patterns
		runtime.GC()
		runtime.GC()

		// Context should still be valid
		st := ctx.Status()
		assert.Equal(t, status.Success, st, "Context should be valid after multiple GetSource calls")
	})

	t.Run("get_source_then_set_new_source", func(t *testing.T) {
		// Get the current source
		oldSrc, err := ctx.GetSource()
		require.NoError(t, err)
		require.NotNil(t, oldSrc)

		// Set a completely new source
		newPat, err := pattern.NewSolidPatternRGB(1.0, 0.0, 0.0)
		require.NoError(t, err)
		defer newPat.Close()

		ctx.SetSource(newPat)

		// The old source is no longer the active source
		// Let it be GC'd without explicit Close()
		oldSrc = nil
		runtime.GC()
		runtime.GC()

		// Everything should still work
		st := ctx.Status()
		assert.Equal(t, status.Success, st, "Context should be valid after source replacement")

		// We should be able to get the new source
		currentSrc, err := ctx.GetSource()
		require.NoError(t, err)
		require.NotNil(t, currentSrc)
	})
}
