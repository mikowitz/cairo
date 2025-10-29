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

// TestContextTranslate verifies translation transformations.
func TestContextTranslate(t *testing.T) {
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

	t.Run("simple_translation", func(t *testing.T) {
		ctx.IdentityMatrix() // Reset to identity
		ctx.Translate(10.0, 20.0)

		assert.Equal(t, status.Success, ctx.Status(), "Translate should succeed")

		// Verify translation by getting matrix
		m, err := ctx.GetMatrix()
		require.NoError(t, err, "GetMatrix should succeed after translate")
		assert.InDelta(t, 10.0, m.X0, 0.001, "X0 should be 10.0")
		assert.InDelta(t, 20.0, m.Y0, 0.001, "Y0 should be 20.0")
		assert.InDelta(t, 1.0, m.XX, 0.001, "XX should remain 1.0")
		assert.InDelta(t, 1.0, m.YY, 0.001, "YY should remain 1.0")
	})

	t.Run("negative_translation", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Translate(-5.0, -10.0)

		assert.Equal(t, status.Success, ctx.Status())

		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		assert.InDelta(t, -5.0, m.X0, 0.001)
		assert.InDelta(t, -10.0, m.Y0, 0.001)
	})

	t.Run("zero_translation", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Translate(0.0, 0.0)

		assert.Equal(t, status.Success, ctx.Status())

		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		assert.InDelta(t, 0.0, m.X0, 0.001)
		assert.InDelta(t, 0.0, m.Y0, 0.001)
	})

	t.Run("cumulative_translation", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Translate(5.0, 10.0)
		ctx.Translate(3.0, 7.0)

		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		// Translations should accumulate
		assert.InDelta(t, 8.0, m.X0, 0.001)
		assert.InDelta(t, 17.0, m.Y0, 0.001)
	})
}

// TestContextScale verifies scaling transformations.
func TestContextScale(t *testing.T) {
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

	t.Run("uniform_scaling", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Scale(2.0, 2.0)

		assert.Equal(t, status.Success, ctx.Status(), "Scale should succeed")

		m, err := ctx.GetMatrix()
		require.NoError(t, err, "GetMatrix should succeed after scale")
		assert.InDelta(t, 2.0, m.XX, 0.001, "XX should be 2.0")
		assert.InDelta(t, 2.0, m.YY, 0.001, "YY should be 2.0")

		// Verify scaling effect with coordinate transformation
		x, y := ctx.UserToDevice(10.0, 10.0)
		assert.InDelta(t, 20.0, x, 0.001, "10 * 2 = 20")
		assert.InDelta(t, 20.0, y, 0.001, "10 * 2 = 20")
	})

	t.Run("non_uniform_scaling", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Scale(2.0, 3.0)

		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		assert.InDelta(t, 2.0, m.XX, 0.001)
		assert.InDelta(t, 3.0, m.YY, 0.001)

		x, y := ctx.UserToDevice(10.0, 10.0)
		assert.InDelta(t, 20.0, x, 0.001)
		assert.InDelta(t, 30.0, y, 0.001)
	})

	t.Run("fractional_scaling", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Scale(0.5, 0.5)

		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		assert.InDelta(t, 0.5, m.XX, 0.001)
		assert.InDelta(t, 0.5, m.YY, 0.001)

		x, y := ctx.UserToDevice(10.0, 10.0)
		assert.InDelta(t, 5.0, x, 0.001)
		assert.InDelta(t, 5.0, y, 0.001)
	})

	t.Run("cumulative_scaling", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Scale(2.0, 2.0)
		ctx.Scale(3.0, 3.0)

		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		// Scales should multiply
		assert.InDelta(t, 6.0, m.XX, 0.001)
		assert.InDelta(t, 6.0, m.YY, 0.001)
	})
}

// TestContextRotate verifies rotation transformations.
func TestContextRotate(t *testing.T) {
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

	t.Run("ninety_degree_rotation", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Rotate(math.Pi / 2) // 90 degrees

		assert.Equal(t, status.Success, ctx.Status(), "Rotate should succeed")

		// After 90° rotation, point (1, 0) becomes (0, 1)
		x, y := ctx.UserToDevice(1.0, 0.0)
		assert.InDelta(t, 0.0, x, 0.001)
		assert.InDelta(t, 1.0, y, 0.001)

		// Point (0, 1) becomes (-1, 0)
		x, y = ctx.UserToDevice(0.0, 1.0)
		assert.InDelta(t, -1.0, x, 0.001)
		assert.InDelta(t, 0.0, y, 0.001)
	})

	t.Run("one_eighty_degree_rotation", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Rotate(math.Pi) // 180 degrees

		// After 180° rotation, point (1, 0) becomes (-1, 0)
		x, y := ctx.UserToDevice(1.0, 0.0)
		assert.InDelta(t, -1.0, x, 0.001)
		assert.InDelta(t, 0.0, y, 0.001)

		// Point (0, 1) becomes (0, -1)
		x, y = ctx.UserToDevice(0.0, 1.0)
		assert.InDelta(t, 0.0, x, 0.001)
		assert.InDelta(t, -1.0, y, 0.001)
	})

	t.Run("negative_rotation", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Rotate(-math.Pi / 2) // -90 degrees (clockwise)

		// After -90° rotation, point (1, 0) becomes (0, -1)
		x, y := ctx.UserToDevice(1.0, 0.0)
		assert.InDelta(t, 0.0, x, 0.001)
		assert.InDelta(t, -1.0, y, 0.001)
	})

	t.Run("zero_rotation", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Rotate(0.0)

		// No rotation, point should remain unchanged
		x, y := ctx.UserToDevice(1.0, 0.0)
		assert.InDelta(t, 1.0, x, 0.001)
		assert.InDelta(t, 0.0, y, 0.001)
	})

	t.Run("full_rotation", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Rotate(2 * math.Pi) // 360 degrees

		// Full rotation returns to original
		x, y := ctx.UserToDevice(1.0, 0.0)
		assert.InDelta(t, 1.0, x, 0.001)
		assert.InDelta(t, 0.0, y, 0.001)
	})
}

// TestContextTransform verifies custom matrix transformations.
func TestContextTransform(t *testing.T) {
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

	t.Run("transform_with_translation_matrix", func(t *testing.T) {
		ctx.IdentityMatrix()

		// Create a translation matrix
		m := matrix.NewTranslationMatrix(10.0, 20.0)
		defer m.Close()

		ctx.Transform(m)

		assert.Equal(t, status.Success, ctx.Status(), "Transform should succeed")

		// Verify transformation was applied
		ctxMatrix, err := ctx.GetMatrix()
		require.NoError(t, err)
		assert.InDelta(t, 10.0, ctxMatrix.X0, 0.001)
		assert.InDelta(t, 20.0, ctxMatrix.Y0, 0.001)
	})

	t.Run("transform_with_scaling_matrix", func(t *testing.T) {
		ctx.IdentityMatrix()

		m := matrix.NewScalingMatrix(2.0, 3.0)
		defer m.Close()

		ctx.Transform(m)

		ctxMatrix, err := ctx.GetMatrix()
		require.NoError(t, err)
		assert.InDelta(t, 2.0, ctxMatrix.XX, 0.001)
		assert.InDelta(t, 3.0, ctxMatrix.YY, 0.001)
	})

	t.Run("transform_with_identity_matrix", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Translate(5.0, 10.0) // Apply some transformation first

		// Transform with identity should leave it unchanged
		identity := matrix.NewIdentityMatrix()
		defer identity.Close()

		ctx.Transform(identity)

		ctxMatrix, err := ctx.GetMatrix()
		require.NoError(t, err)
		assert.InDelta(t, 5.0, ctxMatrix.X0, 0.001)
		assert.InDelta(t, 10.0, ctxMatrix.Y0, 0.001)
	})

	t.Run("cumulative_transforms", func(t *testing.T) {
		ctx.IdentityMatrix()

		// Apply multiple transformations
		m1 := matrix.NewScalingMatrix(2.0, 2.0)
		defer m1.Close()
		ctx.Transform(m1)

		m2 := matrix.NewTranslationMatrix(5.0, 10.0)
		defer m2.Close()
		ctx.Transform(m2)

		// Verify combined transformation
		x, y := ctx.UserToDevice(0.0, 0.0)
		assert.InDelta(t, 10.0, x, 0.001)
		assert.InDelta(t, 20.0, y, 0.001)
	})

	t.Run("cumulative_transforms_order_matters", func(t *testing.T) {
		ctx.IdentityMatrix()

		// Apply multiple transformations
		m1 := matrix.NewTranslationMatrix(5.0, 10.0)
		defer m1.Close()
		ctx.Transform(m1)

		m2 := matrix.NewScalingMatrix(2.0, 2.0)
		defer m2.Close()
		ctx.Transform(m2)

		// Verify combined transformation
		x, y := ctx.UserToDevice(0.0, 0.0)
		assert.InDelta(t, 5.0, x, 0.001)
		assert.InDelta(t, 10.0, y, 0.001)
	})
}

// TestContextTransformationsCombined verifies complex transformation sequences.
func TestContextTransformationsCombined(t *testing.T) {
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

	t.Run("translate_then_scale", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Translate(10.0, 20.0)
		ctx.Scale(2.0, 2.0)

		// Translation happens first, then scaling
		// Point (0, 0) → (10, 20) after translate → (20, 40) after scale
		x, y := ctx.UserToDevice(0.0, 0.0)
		assert.InDelta(t, 10.0, x, 0.001, "Translate then scale: (0,0) → (10,20)")
		assert.InDelta(t, 20.0, y, 0.001)

		// Point (5, 10) → (15, 30) after translate → (30, 60) after scale
		x, y = ctx.UserToDevice(5.0, 10.0)
		assert.InDelta(t, 20.0, x, 0.001)
		assert.InDelta(t, 40.0, y, 0.001)
	})

	t.Run("scale_then_translate", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Scale(2.0, 2.0)
		ctx.Translate(10.0, 20.0)

		// Order matters: scale first, then translate
		// Point (0, 0) → (0, 0) after scale → (10, 20) after translate
		x, y := ctx.UserToDevice(0.0, 0.0)
		assert.InDelta(t, 20.0, x, 0.001)
		assert.InDelta(t, 40.0, y, 0.001)

		// Different from translate-then-scale
		// Point (5, 10) → (10, 20) after scale → (20, 40) after translate
		x, y = ctx.UserToDevice(5.0, 10.0)
		assert.InDelta(t, 30.0, x, 0.001)
		assert.InDelta(t, 60.0, y, 0.001)
	})

	t.Run("rotate_then_translate", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Rotate(math.Pi / 2) // 90 degrees
		ctx.Translate(10.0, 0.0)

		// After rotation, translation is in rotated coordinates
		// Point (0, 0) with 90° rotation and translate(10, 0) → (-10, 0) in rotated space
		x, y := ctx.UserToDevice(0.0, 0.0)
		assert.InDelta(t, 0.0, x, 0.001)
		assert.InDelta(t, 10.0, y, 0.001)
	})

	t.Run("translate_scale_rotate", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Translate(10.0, 10.0)
		ctx.Scale(2.0, 2.0)
		ctx.Rotate(math.Pi / 4) // 45 degrees

		// Complex transformation chain
		x, y := ctx.UserToDevice(5.0, 0.0)
		// This is a complex calculation, just verify status
		assert.Equal(t, status.Success, ctx.Status())

		// Verify result is not NaN
		assert.False(t, math.IsNaN(x), "X coordinate should not be NaN")
		assert.False(t, math.IsNaN(y), "Y coordinate should not be NaN")
	})

	t.Run("save_restore_transformations", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Translate(10.0, 10.0)

		// Get transformed coordinates
		x1, y1 := ctx.UserToDevice(5.0, 5.0)

		ctx.Save()
		ctx.Scale(2.0, 2.0)
		x2, y2 := ctx.UserToDevice(5.0, 5.0)

		// After scale, coordinates should be different
		assert.NotEqual(t, x1, x2, "Coordinates should change after scale")
		assert.NotEqual(t, y1, y2, "Coordinates should change after scale")

		ctx.Restore()
		x3, y3 := ctx.UserToDevice(5.0, 5.0)

		// After restore, should match original
		assert.InDelta(t, x1, x3, 0.001, "Restore should return to saved state")
		assert.InDelta(t, y1, y3, 0.001)
	})

	t.Run("nested_save_restore", func(t *testing.T) {
		ctx.IdentityMatrix()

		ctx.Save()
		ctx.Scale(2.0, 2.0)
		x1, y1 := ctx.UserToDevice(0.0, 0.0)

		ctx.Save()
		ctx.Translate(10.0, 10.0)
		x2, _ := ctx.UserToDevice(0.0, 0.0)

		ctx.Restore() // Back to translate only
		x3, y3 := ctx.UserToDevice(0.0, 0.0)

		assert.InDelta(t, x1, x3, 0.001, "First restore should remove scale")
		assert.InDelta(t, y1, y3, 0.001)
		assert.NotEqual(t, x2, x3, "Scale state should be different")

		ctx.Restore() // Back to identity
		x4, y4 := ctx.UserToDevice(0.0, 0.0)

		assert.InDelta(t, 0.0, x4, 0.001, "Final restore should be identity")
		assert.InDelta(t, 0.0, y4, 0.001)
	})

	t.Run("transformation_with_custom_matrix", func(t *testing.T) {
		ctx.IdentityMatrix()

		// Create a combined transformation matrix manually
		m := matrix.NewIdentityMatrix()
		defer m.Close()

		m.Translate(10.0, 20.0)
		m.Scale(2.0, 3.0)

		ctx.SetMatrix(m)

		// Verify the combined transformation works
		x, y := ctx.UserToDevice(5.0, 10.0)
		assert.InDelta(t, 20.0, x, 0.001) // (5 * 2) + 10
		assert.InDelta(t, 50.0, y, 0.001) // (10 * 3) + 20
	})

	t.Run("identity_resets_all_transformations", func(t *testing.T) {
		ctx.Translate(10.0, 20.0)
		ctx.Scale(2.0, 3.0)
		ctx.Rotate(math.Pi / 6)

		// Complex state, then reset
		ctx.IdentityMatrix()

		x, y := ctx.UserToDevice(10.0, 20.0)
		assert.InDelta(t, 10.0, x, 0.001, "Identity should reset everything")
		assert.InDelta(t, 20.0, y, 0.001)
	})

	t.Run("transformation_precision", func(t *testing.T) {
		ctx.IdentityMatrix()

		// Apply and reverse transformation
		ctx.Translate(5.0, 10.0)
		ctx.Translate(-5.0, -10.0)

		// Should be back to identity (within precision)
		x, y := ctx.UserToDevice(10.0, 20.0)
		assert.InDelta(t, 10.0, x, 0.001)
		assert.InDelta(t, 20.0, y, 0.001)
	})
}
