package context

import (
	"math"
	"testing"

	"github.com/mikowitz/cairo/matrix"
	"github.com/mikowitz/cairo/status"
	"github.com/mikowitz/cairo/surface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestContextGetSetMatrix verifies getting and setting the transformation matrix.
func TestContextGetSetMatrix(t *testing.T) {
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

	t.Run("get_identity_matrix", func(t *testing.T) {
		ctx.IdentityMatrix()

		m, err := ctx.GetMatrix()
		require.NoError(t, err, "GetMatrix should succeed")
		require.NotNil(t, m, "Matrix should not be nil")

		// Verify identity matrix values
		assert.InDelta(t, 1.0, m.XX, 0.001)
		assert.InDelta(t, 0.0, m.YX, 0.001)
		assert.InDelta(t, 0.0, m.XY, 0.001)
		assert.InDelta(t, 1.0, m.YY, 0.001)
		assert.InDelta(t, 0.0, m.X0, 0.001)
		assert.InDelta(t, 0.0, m.Y0, 0.001)
	})

	t.Run("get_matrix_after_transformations", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Translate(10.0, 20.0)
		ctx.Scale(2.0, 3.0)

		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		assert.InDelta(t, 2.0, m.XX, 0.001)
		assert.InDelta(t, 3.0, m.YY, 0.001)
		assert.InDelta(t, 10.0, m.X0, 0.001)
		assert.InDelta(t, 20.0, m.Y0, 0.001)
	})

	t.Run("set_matrix", func(t *testing.T) {
		ctx.IdentityMatrix()

		// Create a custom matrix
		customMatrix := matrix.NewMatrix(2.0, 0.0, 0.0, 3.0, 10.0, 20.0)
		defer customMatrix.Close()

		ctx.SetMatrix(customMatrix)

		assert.Equal(t, status.Success, ctx.Status())

		// Verify matrix was set
		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		assert.InDelta(t, 2.0, m.XX, 0.001)
		assert.InDelta(t, 3.0, m.YY, 0.001)
		assert.InDelta(t, 10.0, m.X0, 0.001)
		assert.InDelta(t, 20.0, m.Y0, 0.001)
	})

	t.Run("round_trip_get_set_matrix", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Translate(15.0, 25.0)
		ctx.Scale(1.5, 2.5)

		// Get the matrix
		m, err := ctx.GetMatrix()
		require.NoError(t, err)

		// Reset to identity
		ctx.IdentityMatrix()

		// Set the matrix back
		ctx.SetMatrix(m)

		// Verify it matches original
		m2, err := ctx.GetMatrix()
		require.NoError(t, err)
		assert.InDelta(t, m.XX, m2.XX, 0.001)
		assert.InDelta(t, m.YY, m2.YY, 0.001)
		assert.InDelta(t, m.X0, m2.X0, 0.001)
		assert.InDelta(t, m.Y0, m2.Y0, 0.001)
	})

	t.Run("get_matrix_after_close", func(t *testing.T) {
		surf2, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surf2.Close()

		ctx2, err := NewContext(surf2)
		require.NoError(t, err)

		// Close the context
		err = ctx2.Close()
		require.NoError(t, err)

		// GetMatrix should return an error after close
		m, err := ctx2.GetMatrix()
		assert.Error(t, err, "GetMatrix should return error after close")
		assert.Nil(t, m, "Matrix should be nil after close")
	})

	t.Run("set_matrix_after_close", func(t *testing.T) {
		surf2, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surf2.Close()

		ctx2, err := NewContext(surf2)
		require.NoError(t, err)

		customMatrix := matrix.NewIdentityMatrix()
		defer customMatrix.Close()

		// Close the context
		err = ctx2.Close()
		require.NoError(t, err)

		// SetMatrix should be safe no-op after close
		ctx2.SetMatrix(customMatrix)

		// Status should indicate closed
		st := ctx2.Status()
		assert.Equal(t, status.NullPointer, st)
	})

	t.Run("set_nil_matrix", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Translate(5.0, 10.0)

		// SetMatrix with nil should be safe no-op
		ctx.SetMatrix(nil)

		// Matrix should be unchanged
		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		assert.InDelta(t, 5.0, m.X0, 0.001)
		assert.InDelta(t, 10.0, m.Y0, 0.001)
	})
}

// TestContextIdentityMatrix verifies identity matrix reset.
func TestContextIdentityMatrix(t *testing.T) {
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

	t.Run("reset_after_translate", func(t *testing.T) {
		ctx.Translate(10.0, 20.0)
		ctx.IdentityMatrix()

		assert.Equal(t, status.Success, ctx.Status())

		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		assert.InDelta(t, 1.0, m.XX, 0.001)
		assert.InDelta(t, 1.0, m.YY, 0.001)
		assert.InDelta(t, 0.0, m.X0, 0.001)
		assert.InDelta(t, 0.0, m.Y0, 0.001)

		// Verify coordinates are not transformed
		x, y := ctx.UserToDevice(10.0, 20.0)
		assert.InDelta(t, 10.0, x, 0.001)
		assert.InDelta(t, 20.0, y, 0.001)
	})

	t.Run("reset_after_scale", func(t *testing.T) {
		ctx.Scale(2.0, 3.0)
		ctx.IdentityMatrix()

		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		assert.InDelta(t, 1.0, m.XX, 0.001)
		assert.InDelta(t, 1.0, m.YY, 0.001)

		x, y := ctx.UserToDevice(10.0, 20.0)
		assert.InDelta(t, 10.0, x, 0.001)
		assert.InDelta(t, 20.0, y, 0.001)
	})

	t.Run("reset_after_rotate", func(t *testing.T) {
		ctx.Rotate(math.Pi / 4) // 45 degrees
		ctx.IdentityMatrix()

		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		assert.InDelta(t, 1.0, m.XX, 0.001)
		assert.InDelta(t, 0.0, m.YX, 0.001)
		assert.InDelta(t, 0.0, m.XY, 0.001)
		assert.InDelta(t, 1.0, m.YY, 0.001)
	})

	t.Run("reset_after_complex_transformations", func(t *testing.T) {
		ctx.Translate(10.0, 20.0)
		ctx.Scale(2.0, 3.0)
		ctx.Rotate(math.Pi / 6)
		ctx.IdentityMatrix()

		m, err := ctx.GetMatrix()
		require.NoError(t, err)

		// All values should be back to identity
		assert.InDelta(t, 1.0, m.XX, 0.001)
		assert.InDelta(t, 0.0, m.YX, 0.001)
		assert.InDelta(t, 0.0, m.XY, 0.001)
		assert.InDelta(t, 1.0, m.YY, 0.001)
		assert.InDelta(t, 0.0, m.X0, 0.001)
		assert.InDelta(t, 0.0, m.Y0, 0.001)
	})

	t.Run("multiple_identity_calls", func(t *testing.T) {
		ctx.Translate(5.0, 10.0)
		ctx.IdentityMatrix()
		ctx.IdentityMatrix() // Second call should be harmless

		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		assert.InDelta(t, 1.0, m.XX, 0.001)
		assert.InDelta(t, 0.0, m.X0, 0.001)
	})
}
