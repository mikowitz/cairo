// ABOUTME: Tests for custom Cairo error types (SurfaceError, ContextError, PatternError).
// ABOUTME: Verifies errors.Is, errors.As, error messages, and type assertions.
package cairo_test

import (
	"errors"
	"testing"

	cairo "github.com/mikowitz/cairo"
	"github.com/mikowitz/cairo/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorTypes(t *testing.T) {
	t.Run("SurfaceError type assertion via errors.As", func(t *testing.T) {
		surfErr := &cairo.SurfaceError{Status: status.InvalidFormat, SurfaceType: "image"}
		var se *cairo.SurfaceError
		require.True(t, errors.As(surfErr, &se))
		assert.Equal(t, status.InvalidFormat, se.Status)
		assert.Equal(t, "image", se.SurfaceType)
	})

	t.Run("ContextError type assertion via errors.As", func(t *testing.T) {
		ctxErr := &cairo.ContextError{Status: status.NoCurrentPoint, Operation: "stroke"}
		var ce *cairo.ContextError
		require.True(t, errors.As(ctxErr, &ce))
		assert.Equal(t, status.NoCurrentPoint, ce.Status)
		assert.Equal(t, "stroke", ce.Operation)
	})

	t.Run("PatternError type assertion via errors.As", func(t *testing.T) {
		patErr := &cairo.PatternError{Status: status.PatternTypeMismatch, PatternType: "linear"}
		var pe *cairo.PatternError
		require.True(t, errors.As(patErr, &pe))
		assert.Equal(t, status.PatternTypeMismatch, pe.Status)
		assert.Equal(t, "linear", pe.PatternType)
	})

	t.Run("errors.As returns false for wrong type", func(t *testing.T) {
		surfErr := &cairo.SurfaceError{Status: status.InvalidFormat}
		var ce *cairo.ContextError
		assert.False(t, errors.As(surfErr, &ce))
	})
}

func TestErrorUnwrapping(t *testing.T) {
	t.Run("SurfaceError unwraps to status", func(t *testing.T) {
		surfErr := &cairo.SurfaceError{Status: status.InvalidFormat}
		assert.True(t, errors.Is(surfErr, status.InvalidFormat))
		assert.False(t, errors.Is(surfErr, status.NoCurrentPoint))
	})

	t.Run("ContextError unwraps to status", func(t *testing.T) {
		ctxErr := &cairo.ContextError{Status: status.NoCurrentPoint}
		assert.True(t, errors.Is(ctxErr, status.NoCurrentPoint))
		assert.False(t, errors.Is(ctxErr, status.InvalidFormat))
	})

	t.Run("PatternError unwraps to status", func(t *testing.T) {
		patErr := &cairo.PatternError{Status: status.PatternTypeMismatch}
		assert.True(t, errors.Is(patErr, status.PatternTypeMismatch))
		assert.False(t, errors.Is(patErr, status.InvalidFormat))
	})

	t.Run("errors.Is with matching SurfaceError", func(t *testing.T) {
		surfErr := &cairo.SurfaceError{Status: status.InvalidFormat, SurfaceType: "image"}
		target := &cairo.SurfaceError{Status: status.InvalidFormat}
		assert.True(t, errors.Is(surfErr, target))
	})

	t.Run("errors.Is with mismatched SurfaceType", func(t *testing.T) {
		surfErr := &cairo.SurfaceError{Status: status.InvalidFormat, SurfaceType: "image"}
		target := &cairo.SurfaceError{Status: status.InvalidFormat, SurfaceType: "pdf"}
		assert.False(t, errors.Is(surfErr, target))
	})

	t.Run("errors.Is with matching ContextError", func(t *testing.T) {
		ctxErr := &cairo.ContextError{Status: status.NoCurrentPoint, Operation: "stroke"}
		target := &cairo.ContextError{Status: status.NoCurrentPoint}
		assert.True(t, errors.Is(ctxErr, target))
	})

	t.Run("errors.Is with mismatched Operation", func(t *testing.T) {
		ctxErr := &cairo.ContextError{Status: status.NoCurrentPoint, Operation: "stroke"}
		target := &cairo.ContextError{Status: status.NoCurrentPoint, Operation: "fill"}
		assert.False(t, errors.Is(ctxErr, target))
	})

	t.Run("errors.Is with matching PatternError", func(t *testing.T) {
		patErr := &cairo.PatternError{Status: status.PatternTypeMismatch, PatternType: "linear"}
		target := &cairo.PatternError{Status: status.PatternTypeMismatch}
		assert.True(t, errors.Is(patErr, target))
	})

	t.Run("errors.Is with mismatched PatternType", func(t *testing.T) {
		patErr := &cairo.PatternError{Status: status.PatternTypeMismatch, PatternType: "linear"}
		target := &cairo.PatternError{Status: status.PatternTypeMismatch, PatternType: "solid"}
		assert.False(t, errors.Is(patErr, target))
	})
}

func TestErrorContext(t *testing.T) {
	t.Run("SurfaceError includes surface type in message", func(t *testing.T) {
		surfErr := &cairo.SurfaceError{Status: status.InvalidFormat, SurfaceType: "image"}
		assert.Contains(t, surfErr.Error(), "image")
	})

	t.Run("SurfaceError without surface type has simpler message", func(t *testing.T) {
		surfErr := &cairo.SurfaceError{Status: status.InvalidFormat}
		assert.Contains(t, surfErr.Error(), "cairo surface error")
		assert.NotContains(t, surfErr.Error(), "()")
	})

	t.Run("ContextError includes operation in message", func(t *testing.T) {
		ctxErr := &cairo.ContextError{Status: status.NoCurrentPoint, Operation: "stroke"}
		assert.Contains(t, ctxErr.Error(), "stroke")
	})

	t.Run("ContextError without operation has simpler message", func(t *testing.T) {
		ctxErr := &cairo.ContextError{Status: status.NoCurrentPoint}
		assert.Contains(t, ctxErr.Error(), "cairo context error")
		assert.NotContains(t, ctxErr.Error(), "()")
	})

	t.Run("PatternError includes pattern type in message", func(t *testing.T) {
		patErr := &cairo.PatternError{Status: status.PatternTypeMismatch, PatternType: "linear"}
		assert.Contains(t, patErr.Error(), "linear")
	})

	t.Run("PatternError without pattern type has simpler message", func(t *testing.T) {
		patErr := &cairo.PatternError{Status: status.PatternTypeMismatch}
		assert.Contains(t, patErr.Error(), "cairo pattern error")
		assert.NotContains(t, patErr.Error(), "()")
	})
}
