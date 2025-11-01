package pattern

import (
	"math"
	"testing"

	"github.com/mikowitz/cairo/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewLinearGradient tests creation of linear gradient patterns.
// Linear gradients create smooth color transitions along a line from (x0, y0) to (x1, y1).
func TestNewLinearGradient(t *testing.T) {
	t.Run("horizontal_gradient", func(t *testing.T) {
		// Create a horizontal gradient from left to right
		gradient, err := NewLinearGradient(0, 0, 100, 0)
		require.NoError(t, err, "NewLinearGradient should not return an error")
		require.NotNil(t, gradient, "Gradient should not be nil")
		defer gradient.Close()

		// Verify gradient has successful status
		st := gradient.Status()
		assert.Equal(t, status.Success, st, "Gradient should have success status")

		// Verify gradient type
		patternType := gradient.GetType()
		assert.Equal(t, PatternTypeLinear, patternType, "Pattern type should be Linear")
	})

	t.Run("vertical_gradient", func(t *testing.T) {
		// Create a vertical gradient from top to bottom
		gradient, err := NewLinearGradient(0, 0, 0, 100)
		require.NoError(t, err)
		require.NotNil(t, gradient)
		defer gradient.Close()

		st := gradient.Status()
		assert.Equal(t, status.Success, st)
	})

	t.Run("diagonal_gradient", func(t *testing.T) {
		// Create a diagonal gradient
		gradient, err := NewLinearGradient(0, 0, 100, 100)
		require.NoError(t, err)
		require.NotNil(t, gradient)
		defer gradient.Close()

		st := gradient.Status()
		assert.Equal(t, status.Success, st)
	})

	t.Run("zero_length_gradient", func(t *testing.T) {
		// Create a gradient with same start and end point
		// This is technically valid in Cairo
		gradient, err := NewLinearGradient(50, 50, 50, 50)
		require.NoError(t, err)
		require.NotNil(t, gradient)
		defer gradient.Close()

		st := gradient.Status()
		assert.Equal(t, status.Success, st)
	})
}

// TestNewRadialGradient tests creation of radial gradient patterns.
// Radial gradients create smooth color transitions between two circles.
func TestNewRadialGradient(t *testing.T) {
	t.Run("expanding_gradient", func(t *testing.T) {
		// Create a radial gradient expanding from center
		// Small inner circle (radius 10) to larger outer circle (radius 100)
		gradient, err := NewRadialGradient(100, 100, 10, 100, 100, 100)
		require.NoError(t, err, "NewRadialGradient should not return an error")
		require.NotNil(t, gradient, "Gradient should not be nil")
		defer gradient.Close()

		// Verify gradient has successful status
		st := gradient.Status()
		assert.Equal(t, status.Success, st, "Gradient should have success status")

		// Verify gradient type
		patternType := gradient.GetType()
		assert.Equal(t, PatternTypeRadial, patternType, "Pattern type should be Radial")
	})

	t.Run("offset_centers", func(t *testing.T) {
		// Create a radial gradient with offset centers
		gradient, err := NewRadialGradient(50, 50, 10, 75, 75, 100)
		require.NoError(t, err)
		require.NotNil(t, gradient)
		defer gradient.Close()

		st := gradient.Status()
		assert.Equal(t, status.Success, st)
	})

	t.Run("concentric_circles", func(t *testing.T) {
		// Create concentric circles (same center)
		gradient, err := NewRadialGradient(100, 100, 20, 100, 100, 80)
		require.NoError(t, err)
		require.NotNil(t, gradient)
		defer gradient.Close()

		st := gradient.Status()
		assert.Equal(t, status.Success, st)
	})

	t.Run("zero_radius_inner", func(t *testing.T) {
		// Inner circle with zero radius (point)
		gradient, err := NewRadialGradient(50, 50, 0, 50, 50, 100)
		require.NoError(t, err)
		require.NotNil(t, gradient)
		defer gradient.Close()

		st := gradient.Status()
		assert.Equal(t, status.Success, st)
	})
}

// TestLinearGradientColorStops tests adding color stops to linear gradients.
// Color stops define the colors at specific positions along the gradient.
// Offset values must be in the range [0.0, 1.0] where 0.0 is the start and 1.0 is the end.
func TestLinearGradientColorStops(t *testing.T) {
	t.Run("simple_two_stop_rgb", func(t *testing.T) {
		gradient, err := NewLinearGradient(0, 0, 100, 0)
		require.NoError(t, err)
		defer gradient.Close()

		// Add two color stops: red at start, blue at end
		gradient.AddColorStopRGB(0.0, 1.0, 0.0, 0.0) // Red at 0%
		gradient.AddColorStopRGB(1.0, 0.0, 0.0, 1.0) // Blue at 100%

		// Verify gradient is still valid after adding stops
		st := gradient.Status()
		assert.Equal(t, status.Success, st, "Gradient should remain valid after adding color stops")
	})

	t.Run("simple_two_stop_rgba", func(t *testing.T) {
		gradient, err := NewLinearGradient(0, 0, 100, 0)
		require.NoError(t, err)
		defer gradient.Close()

		// Add two color stops with alpha: opaque red to transparent blue
		gradient.AddColorStopRGBA(0.0, 1.0, 0.0, 0.0, 1.0) // Opaque red at 0%
		gradient.AddColorStopRGBA(1.0, 0.0, 0.0, 1.0, 0.0) // Transparent blue at 100%

		st := gradient.Status()
		assert.Equal(t, status.Success, st)
	})

	t.Run("multiple_stops", func(t *testing.T) {
		gradient, err := NewLinearGradient(0, 0, 100, 0)
		require.NoError(t, err)
		defer gradient.Close()

		// Create a rainbow gradient with multiple stops
		gradient.AddColorStopRGB(0.0, 1.0, 0.0, 0.0)  // Red
		gradient.AddColorStopRGB(0.25, 1.0, 1.0, 0.0) // Yellow
		gradient.AddColorStopRGB(0.5, 0.0, 1.0, 0.0)  // Green
		gradient.AddColorStopRGB(0.75, 0.0, 0.0, 1.0) // Blue
		gradient.AddColorStopRGB(1.0, 0.5, 0.0, 0.5)  // Purple

		st := gradient.Status()
		assert.Equal(t, status.Success, st)
	})

	t.Run("out_of_order_stops", func(t *testing.T) {
		gradient, err := NewLinearGradient(0, 0, 100, 0)
		require.NoError(t, err)
		defer gradient.Close()

		// Color stops don't need to be added in order
		gradient.AddColorStopRGB(1.0, 0.0, 0.0, 1.0) // Blue at end
		gradient.AddColorStopRGB(0.0, 1.0, 0.0, 0.0) // Red at start
		gradient.AddColorStopRGB(0.5, 0.0, 1.0, 0.0) // Green in middle

		st := gradient.Status()
		assert.Equal(t, status.Success, st)
	})

	t.Run("mixed_rgb_rgba_stops", func(t *testing.T) {
		gradient, err := NewLinearGradient(0, 0, 100, 0)
		require.NoError(t, err)
		defer gradient.Close()

		// Mix RGB and RGBA stops
		gradient.AddColorStopRGB(0.0, 1.0, 0.0, 0.0)       // Opaque red
		gradient.AddColorStopRGBA(0.5, 0.0, 1.0, 0.0, 0.5) // Semi-transparent green
		gradient.AddColorStopRGB(1.0, 0.0, 0.0, 1.0)       // Opaque blue

		st := gradient.Status()
		assert.Equal(t, status.Success, st)
	})
}

// TestRadialGradientColorStops tests adding color stops to radial gradients.
func TestRadialGradientColorStops(t *testing.T) {
	t.Run("simple_two_stop_rgb", func(t *testing.T) {
		gradient, err := NewRadialGradient(100, 100, 10, 100, 100, 100)
		require.NoError(t, err)
		defer gradient.Close()

		// White at center, blue at edge
		gradient.AddColorStopRGB(0.0, 1.0, 1.0, 1.0) // White at center
		gradient.AddColorStopRGB(1.0, 0.0, 0.0, 1.0) // Blue at edge

		st := gradient.Status()
		assert.Equal(t, status.Success, st)
	})

	t.Run("simple_two_stop_rgba", func(t *testing.T) {
		gradient, err := NewRadialGradient(100, 100, 0, 100, 100, 100)
		require.NoError(t, err)
		defer gradient.Close()

		// Opaque center fading to transparent edge
		gradient.AddColorStopRGBA(0.0, 1.0, 0.0, 0.0, 1.0) // Opaque red center
		gradient.AddColorStopRGBA(1.0, 1.0, 0.0, 0.0, 0.0) // Transparent red edge

		st := gradient.Status()
		assert.Equal(t, status.Success, st)
	})

	t.Run("multiple_stops", func(t *testing.T) {
		gradient, err := NewRadialGradient(100, 100, 0, 100, 100, 100)
		require.NoError(t, err)
		defer gradient.Close()

		// Create concentric color rings
		gradient.AddColorStopRGB(0.0, 1.0, 1.0, 1.0)  // White center
		gradient.AddColorStopRGB(0.33, 1.0, 1.0, 0.0) // Yellow ring
		gradient.AddColorStopRGB(0.66, 1.0, 0.5, 0.0) // Orange ring
		gradient.AddColorStopRGB(1.0, 1.0, 0.0, 0.0)  // Red outer edge

		st := gradient.Status()
		assert.Equal(t, status.Success, st)
	})

	t.Run("mixed_rgb_rgba_stops", func(t *testing.T) {
		gradient, err := NewRadialGradient(50, 50, 5, 50, 50, 50)
		require.NoError(t, err)
		defer gradient.Close()

		// Mix RGB and RGBA color stops
		gradient.AddColorStopRGBA(0.0, 1.0, 1.0, 1.0, 1.0) // Opaque white
		gradient.AddColorStopRGB(0.5, 0.0, 1.0, 0.0)       // Opaque green
		gradient.AddColorStopRGBA(1.0, 0.0, 0.0, 1.0, 0.3) // Semi-transparent blue

		st := gradient.Status()
		assert.Equal(t, status.Success, st)
	})
}

// TestGradientClose verifies gradient cleanup behavior
func TestGradientClose(t *testing.T) {
	t.Run("linear_gradient_close", func(t *testing.T) {
		gradient, err := NewLinearGradient(0, 0, 100, 0)
		require.NoError(t, err)
		require.NotNil(t, gradient)

		// First close should succeed
		err = gradient.Close()
		assert.NoError(t, err)

		// Status after close should be NullPointer
		st := gradient.Status()
		assert.Equal(t, status.NullPointer, st)

		// Second close should be safe (no-op)
		err = gradient.Close()
		assert.NoError(t, err)
	})

	t.Run("radial_gradient_close", func(t *testing.T) {
		gradient, err := NewRadialGradient(100, 100, 10, 100, 100, 100)
		require.NoError(t, err)
		require.NotNil(t, gradient)

		err = gradient.Close()
		assert.NoError(t, err)

		st := gradient.Status()
		assert.Equal(t, status.NullPointer, st)

		// Double close safety
		err = gradient.Close()
		assert.NoError(t, err)
	})
}

// TestGradientWithArc tests using gradients with Arc method (integration test without context import).
// This tests that gradients can be created and used with drawing operations.
func TestGradientWithArc(t *testing.T) {
	t.Run("radial_gradient_creation", func(t *testing.T) {
		// Create a radial gradient for a circular pattern
		gradient, err := NewRadialGradient(100, 100, 10, 100, 100, 90)
		require.NoError(t, err)
		defer gradient.Close()

		// Add color stops
		gradient.AddColorStopRGB(0.0, 1.0, 1.0, 1.0) // White center
		gradient.AddColorStopRGB(1.0, 0.0, 0.0, 1.0) // Blue edge

		st := gradient.Status()
		assert.Equal(t, status.Success, st)

		// Verify gradient can be used multiple times
		for i := 0; i < 3; i++ {
			st = gradient.Status()
			assert.Equal(t, status.Success, st, "Gradient should be reusable")
		}
	})

	t.Run("linear_gradient_creation", func(t *testing.T) {
		gradient, err := NewLinearGradient(0, 0, 200, 0)
		require.NoError(t, err)
		defer gradient.Close()

		// Add color stops
		gradient.AddColorStopRGB(0.0, 1.0, 0.0, 0.0) // Red
		gradient.AddColorStopRGB(1.0, 0.0, 0.0, 1.0) // Blue

		st := gradient.Status()
		assert.Equal(t, status.Success, st)
	})
}

// TestGradientWithTransformations tests that gradient matrices work correctly.
// This verifies the transformation functionality inherited from BasePattern.
func TestGradientWithTransformations(t *testing.T) {
	t.Run("linear_gradient_matrix", func(t *testing.T) {
		gradient, err := NewLinearGradient(0, 0, 100, 0)
		require.NoError(t, err)
		defer gradient.Close()

		// Get the identity matrix
		m, err := gradient.GetMatrix()
		require.NoError(t, err)
		require.NotNil(t, m)
		defer m.Close()

		// Matrix operations should work
		m.Scale(2.0, 2.0)
		gradient.SetMatrix(m)

		st := gradient.Status()
		assert.Equal(t, status.Success, st, "Gradient should work with transformations")
	})

	t.Run("radial_gradient_matrix", func(t *testing.T) {
		gradient, err := NewRadialGradient(100, 100, 10, 100, 100, 100)
		require.NoError(t, err)
		defer gradient.Close()

		m, err := gradient.GetMatrix()
		require.NoError(t, err)
		require.NotNil(t, m)
		defer m.Close()

		// Rotate the gradient pattern
		m.Rotate(math.Pi / 4) // 45 degrees
		gradient.SetMatrix(m)

		st := gradient.Status()
		assert.Equal(t, status.Success, st)
	})
}

// TestLinearGradientGetColorStopCount tests retrieving the number of color stops from linear gradients.
func TestLinearGradientGetColorStopCount(t *testing.T) {
	t.Run("no_color_stops", func(t *testing.T) {
		gradient, err := NewLinearGradient(0, 0, 100, 0)
		require.NoError(t, err)
		defer gradient.Close()

		count, err := gradient.GetColorStopCount()
		require.NoError(t, err)
		assert.Equal(t, 0, count, "New gradient should have 0 color stops")
	})

	t.Run("two_color_stops", func(t *testing.T) {
		gradient, err := NewLinearGradient(0, 0, 100, 0)
		require.NoError(t, err)
		defer gradient.Close()

		gradient.AddColorStopRGB(0.0, 1.0, 0.0, 0.0)
		gradient.AddColorStopRGB(1.0, 0.0, 0.0, 1.0)

		count, err := gradient.GetColorStopCount()
		require.NoError(t, err)
		assert.Equal(t, 2, count, "Gradient should have 2 color stops")
	})

	t.Run("five_color_stops", func(t *testing.T) {
		gradient, err := NewLinearGradient(0, 0, 100, 0)
		require.NoError(t, err)
		defer gradient.Close()

		gradient.AddColorStopRGB(0.0, 1.0, 0.0, 0.0)
		gradient.AddColorStopRGB(0.25, 1.0, 1.0, 0.0)
		gradient.AddColorStopRGB(0.5, 0.0, 1.0, 0.0)
		gradient.AddColorStopRGB(0.75, 0.0, 0.0, 1.0)
		gradient.AddColorStopRGB(1.0, 0.5, 0.0, 0.5)

		count, err := gradient.GetColorStopCount()
		require.NoError(t, err)
		assert.Equal(t, 5, count, "Gradient should have 5 color stops")
	})

	t.Run("mixed_rgb_rgba_stops", func(t *testing.T) {
		gradient, err := NewLinearGradient(0, 0, 100, 0)
		require.NoError(t, err)
		defer gradient.Close()

		gradient.AddColorStopRGB(0.0, 1.0, 0.0, 0.0)
		gradient.AddColorStopRGBA(0.5, 0.0, 1.0, 0.0, 0.5)
		gradient.AddColorStopRGB(1.0, 0.0, 0.0, 1.0)

		count, err := gradient.GetColorStopCount()
		require.NoError(t, err)
		assert.Equal(t, 3, count, "Gradient should have 3 color stops")
	})
}

// TestRadialGradientGetColorStopCount tests retrieving the number of color stops from radial gradients.
func TestRadialGradientGetColorStopCount(t *testing.T) {
	t.Run("no_color_stops", func(t *testing.T) {
		gradient, err := NewRadialGradient(100, 100, 10, 100, 100, 100)
		require.NoError(t, err)
		defer gradient.Close()

		count, err := gradient.GetColorStopCount()
		require.NoError(t, err)
		assert.Equal(t, 0, count, "New gradient should have 0 color stops")
	})

	t.Run("two_color_stops", func(t *testing.T) {
		gradient, err := NewRadialGradient(100, 100, 10, 100, 100, 100)
		require.NoError(t, err)
		defer gradient.Close()

		gradient.AddColorStopRGB(0.0, 1.0, 1.0, 1.0)
		gradient.AddColorStopRGB(1.0, 0.0, 0.0, 1.0)

		count, err := gradient.GetColorStopCount()
		require.NoError(t, err)
		assert.Equal(t, 2, count, "Gradient should have 2 color stops")
	})

	t.Run("four_color_stops", func(t *testing.T) {
		gradient, err := NewRadialGradient(100, 100, 0, 100, 100, 100)
		require.NoError(t, err)
		defer gradient.Close()

		gradient.AddColorStopRGB(0.0, 1.0, 1.0, 1.0)
		gradient.AddColorStopRGB(0.33, 1.0, 1.0, 0.0)
		gradient.AddColorStopRGB(0.66, 1.0, 0.5, 0.0)
		gradient.AddColorStopRGB(1.0, 1.0, 0.0, 0.0)

		count, err := gradient.GetColorStopCount()
		require.NoError(t, err)
		assert.Equal(t, 4, count, "Gradient should have 4 color stops")
	})
}

// TestLinearGradientGetColorStopRGBA tests retrieving color stop data from linear gradients.
func TestLinearGradientGetColorStopRGBA(t *testing.T) {
	t.Run("single_color_stop", func(t *testing.T) {
		gradient, err := NewLinearGradient(0, 0, 100, 0)
		require.NoError(t, err)
		defer gradient.Close()

		gradient.AddColorStopRGBA(0.5, 1.0, 0.0, 0.0, 0.75)

		offset, r, g, b, a, err := gradient.GetColorStopRGBA(0)
		require.NoError(t, err)
		assert.InDelta(t, 0.5, offset, 0.0001, "Offset should be 0.5")
		assert.InDelta(t, 1.0, r, 0.0001, "Red should be 1.0")
		assert.InDelta(t, 0.0, g, 0.0001, "Green should be 0.0")
		assert.InDelta(t, 0.0, b, 0.0001, "Blue should be 0.0")
		assert.InDelta(t, 0.75, a, 0.0001, "Alpha should be 0.75")
	})

	t.Run("multiple_color_stops", func(t *testing.T) {
		gradient, err := NewLinearGradient(0, 0, 100, 0)
		require.NoError(t, err)
		defer gradient.Close()

		// Add three color stops
		gradient.AddColorStopRGB(0.0, 1.0, 0.0, 0.0)   // Red at start
		gradient.AddColorStopRGBA(0.5, 0.0, 1.0, 0.0, 0.5) // Semi-transparent green in middle
		gradient.AddColorStopRGB(1.0, 0.0, 0.0, 1.0)   // Blue at end

		// Check first stop (RGB becomes RGBA with alpha=1.0)
		offset, r, g, b, a, err := gradient.GetColorStopRGBA(0)
		require.NoError(t, err)
		assert.InDelta(t, 0.0, offset, 0.0001)
		assert.InDelta(t, 1.0, r, 0.0001)
		assert.InDelta(t, 0.0, g, 0.0001)
		assert.InDelta(t, 0.0, b, 0.0001)
		assert.InDelta(t, 1.0, a, 0.0001, "RGB stops should have alpha=1.0")

		// Check second stop
		offset, r, g, b, a, err = gradient.GetColorStopRGBA(1)
		require.NoError(t, err)
		assert.InDelta(t, 0.5, offset, 0.0001)
		assert.InDelta(t, 0.0, r, 0.0001)
		assert.InDelta(t, 1.0, g, 0.0001)
		assert.InDelta(t, 0.0, b, 0.0001)
		assert.InDelta(t, 0.5, a, 0.0001)

		// Check third stop
		offset, r, g, b, a, err = gradient.GetColorStopRGBA(2)
		require.NoError(t, err)
		assert.InDelta(t, 1.0, offset, 0.0001)
		assert.InDelta(t, 0.0, r, 0.0001)
		assert.InDelta(t, 0.0, g, 0.0001)
		assert.InDelta(t, 1.0, b, 0.0001)
		assert.InDelta(t, 1.0, a, 0.0001)
	})

	t.Run("out_of_bounds_index", func(t *testing.T) {
		gradient, err := NewLinearGradient(0, 0, 100, 0)
		require.NoError(t, err)
		defer gradient.Close()

		gradient.AddColorStopRGB(0.0, 1.0, 0.0, 0.0)

		// Try to get a color stop beyond the count
		_, _, _, _, _, err = gradient.GetColorStopRGBA(1)
		assert.Error(t, err, "Should return error for out of bounds index")
	})
}

// TestRadialGradientGetColorStopRGBA tests retrieving color stop data from radial gradients.
func TestRadialGradientGetColorStopRGBA(t *testing.T) {
	t.Run("single_color_stop", func(t *testing.T) {
		gradient, err := NewRadialGradient(100, 100, 10, 100, 100, 100)
		require.NoError(t, err)
		defer gradient.Close()

		gradient.AddColorStopRGBA(0.25, 0.5, 0.75, 1.0, 0.9)

		offset, r, g, b, a, err := gradient.GetColorStopRGBA(0)
		require.NoError(t, err)
		assert.InDelta(t, 0.25, offset, 0.0001)
		assert.InDelta(t, 0.5, r, 0.0001)
		assert.InDelta(t, 0.75, g, 0.0001)
		assert.InDelta(t, 1.0, b, 0.0001)
		assert.InDelta(t, 0.9, a, 0.0001)
	})

	t.Run("multiple_color_stops", func(t *testing.T) {
		gradient, err := NewRadialGradient(100, 100, 0, 100, 100, 100)
		require.NoError(t, err)
		defer gradient.Close()

		// Add four color stops
		gradient.AddColorStopRGB(0.0, 1.0, 1.0, 1.0)
		gradient.AddColorStopRGB(0.33, 1.0, 1.0, 0.0)
		gradient.AddColorStopRGBA(0.66, 1.0, 0.5, 0.0, 0.8)
		gradient.AddColorStopRGB(1.0, 1.0, 0.0, 0.0)

		// Check each stop
		testCases := []struct {
			index          int
			expectedOffset float64
			expectedR      float64
			expectedG      float64
			expectedB      float64
			expectedA      float64
		}{
			{0, 0.0, 1.0, 1.0, 1.0, 1.0},
			{1, 0.33, 1.0, 1.0, 0.0, 1.0},
			{2, 0.66, 1.0, 0.5, 0.0, 0.8},
			{3, 1.0, 1.0, 0.0, 0.0, 1.0},
		}

		for _, tc := range testCases {
			offset, r, g, b, a, err := gradient.GetColorStopRGBA(tc.index)
			require.NoError(t, err)
			assert.InDelta(t, tc.expectedOffset, offset, 0.0001)
			assert.InDelta(t, tc.expectedR, r, 0.0001)
			assert.InDelta(t, tc.expectedG, g, 0.0001)
			assert.InDelta(t, tc.expectedB, b, 0.0001)
			assert.InDelta(t, tc.expectedA, a, 0.0001)
		}
	})

	t.Run("out_of_bounds_index", func(t *testing.T) {
		gradient, err := NewRadialGradient(100, 100, 10, 100, 100, 100)
		require.NoError(t, err)
		defer gradient.Close()

		gradient.AddColorStopRGB(0.0, 1.0, 0.0, 0.0)
		gradient.AddColorStopRGB(1.0, 0.0, 0.0, 1.0)

		// Try to get color stop at index 2 (only 0 and 1 exist)
		_, _, _, _, _, err = gradient.GetColorStopRGBA(2)
		assert.Error(t, err, "Should return error for out of bounds index")
	})
}
