package context

import (
	"testing"

	"github.com/mikowitz/cairo/surface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestContextLineCap tests setting and getting line cap styles.
// Line caps control how the ends of open paths are rendered when stroked.
func TestContextLineCap(t *testing.T) {
	tests := []struct {
		name     string
		lineCap  LineCap
		setFirst bool // if false, test getting the default value
	}{
		{
			name:     "default line cap",
			lineCap:  LineCapButt,
			setFirst: false,
		},
		{
			name:     "set line cap butt",
			lineCap:  LineCapButt,
			setFirst: true,
		},
		{
			name:     "set line cap round",
			lineCap:  LineCapRound,
			setFirst: true,
		},
		{
			name:     "set line cap square",
			lineCap:  LineCapSquare,
			setFirst: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
			require.NoError(t, err)
			defer surface.Close()

			ctx, err := NewContext(surface)
			require.NoError(t, err)
			defer ctx.Close()

			if tt.setFirst {
				ctx.SetLineCap(tt.lineCap)
			}

			cap := ctx.GetLineCap()
			assert.Equal(t, tt.lineCap, cap)
		})
	}

	t.Run("line cap after close", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)

		ctx.SetLineCap(LineCapRound)
		_ = ctx.Close()

		// Operations after close should be safe no-ops
		ctx.SetLineCap(LineCapSquare)
		// GetLineCap after close should return a reasonable default
		cap := ctx.GetLineCap()
		_ = cap // Just verify it doesn't panic
	})
}

// TestContextLineJoin tests setting and getting line join styles.
// Line joins control how corners are rendered when two path segments meet.
func TestContextLineJoin(t *testing.T) {
	tests := []struct {
		name     string
		lineJoin LineJoin
		setFirst bool // if false, test getting the default value
	}{
		{
			name:     "default line join",
			lineJoin: LineJoinMiter,
			setFirst: false,
		},
		{
			name:     "set line join miter",
			lineJoin: LineJoinMiter,
			setFirst: true,
		},
		{
			name:     "set line join round",
			lineJoin: LineJoinRound,
			setFirst: true,
		},
		{
			name:     "set line join bevel",
			lineJoin: LineJoinBevel,
			setFirst: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
			require.NoError(t, err)
			defer surface.Close()

			ctx, err := NewContext(surface)
			require.NoError(t, err)
			defer ctx.Close()

			if tt.setFirst {
				ctx.SetLineJoin(tt.lineJoin)
			}

			join := ctx.GetLineJoin()
			assert.Equal(t, tt.lineJoin, join)
		})
	}

	t.Run("line join after close", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)

		ctx.SetLineJoin(LineJoinRound)
		_ = ctx.Close()

		// Operations after close should be safe no-ops
		ctx.SetLineJoin(LineJoinBevel)
		// GetLineJoin after close should return a reasonable default
		join := ctx.GetLineJoin()
		_ = join // Just verify it doesn't panic
	})
}

// TestContextDash tests setting and getting dash patterns.
// Dash patterns create dashed or dotted lines.
func TestContextDash(t *testing.T) {
	t.Run("default dash empty", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		// Default dash pattern should be empty (solid line)
		dashes, offset, err := ctx.GetDash()
		require.NoError(t, err)
		assert.Empty(t, dashes, "Default dash pattern should be empty")
		assert.Equal(t, 0.0, offset, "Default dash offset should be 0")
	})

	t.Run("set simple dash", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		// Set simple dash pattern: 10 on, 5 off
		pattern := []float64{10.0, 5.0}
		err = ctx.SetDash(pattern, 0.0)
		require.NoError(t, err)

		dashes, offset, err := ctx.GetDash()
		require.NoError(t, err)
		assert.Equal(t, pattern, dashes, "Dash pattern should match")
		assert.Equal(t, 0.0, offset, "Dash offset should be 0")
	})

	t.Run("set dash with offset", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		// Set dash pattern with offset
		pattern := []float64{10.0, 5.0}
		offset := 3.5
		err = ctx.SetDash(pattern, offset)
		require.NoError(t, err)

		dashes, actualOffset, err := ctx.GetDash()
		require.NoError(t, err)
		assert.Equal(t, pattern, dashes, "Dash pattern should match")
		assert.InDelta(t, offset, actualOffset, 0.001, "Dash offset should match")
	})

	t.Run("set complex dash pattern", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		// Set complex dash pattern: dash, gap, dot, gap
		pattern := []float64{15.0, 5.0, 3.0, 5.0}
		err = ctx.SetDash(pattern, 0.0)
		require.NoError(t, err)

		dashes, offset, err := ctx.GetDash()
		require.NoError(t, err)
		assert.Equal(t, pattern, dashes, "Dash pattern should match")
		assert.Equal(t, 0.0, offset, "Dash offset should be 0")
	})

	t.Run("set dash single value", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		// Single value creates equal on/off pattern
		pattern := []float64{10.0}
		err = ctx.SetDash(pattern, 0.0)
		require.NoError(t, err)

		dashes, offset, err := ctx.GetDash()
		require.NoError(t, err)
		assert.NotEmpty(t, dashes, "Dash pattern should not be empty")
		assert.Equal(t, 0.0, offset, "Dash offset should be 0")
	})

	t.Run("clear dash with empty slice", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		// Set a dash pattern
		pattern := []float64{10.0, 5.0}
		err = ctx.SetDash(pattern, 0.0)
		require.NoError(t, err)

		// Clear it with empty slice
		err = ctx.SetDash([]float64{}, 0.0)
		require.NoError(t, err)

		dashes, offset, err := ctx.GetDash()
		require.NoError(t, err)
		assert.Empty(t, dashes, "Dash pattern should be empty after clearing")
		assert.Equal(t, 0.0, offset, "Dash offset should be 0")
	})

	t.Run("clear dash with nil slice", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		// Set a dash pattern
		pattern := []float64{10.0, 5.0}
		err = ctx.SetDash(pattern, 0.0)
		require.NoError(t, err)

		// Clear it with nil slice
		err = ctx.SetDash(nil, 0.0)
		require.NoError(t, err)

		dashes, offset, err := ctx.GetDash()
		require.NoError(t, err)
		assert.Empty(t, dashes, "Dash pattern should be empty after clearing")
		assert.Equal(t, 0.0, offset, "Dash offset should be 0")
	})

	t.Run("dash after close", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)

		pattern := []float64{10.0, 5.0}
		err = ctx.SetDash(pattern, 0.0)
		require.NoError(t, err)
		_ = ctx.Close()

		// Operations after close should return the appropriate error, not panic
		err = ctx.SetDash([]float64{5.0, 5.0}, 0.0)
		assert.Error(t, err, "SetDash should fail after close")

		// GetDash after close should return error or empty
		_, _, err = ctx.GetDash()
		// Either error or empty result is acceptable
		_ = err
	})
}

// TestContextDashEmpty specifically tests empty dash patterns (solid lines).
func TestContextDashEmpty(t *testing.T) {
	t.Run("empty dash creates solid line", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		// Explicitly set empty dash (solid line)
		err = ctx.SetDash([]float64{}, 0.0)
		require.NoError(t, err)

		dashes, offset, err := ctx.GetDash()
		require.NoError(t, err)
		assert.Empty(t, dashes, "Empty dash pattern creates solid line")
		assert.Equal(t, 0.0, offset)
	})

	t.Run("nil dash creates solid line", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		// Set nil dash (solid line)
		err = ctx.SetDash(nil, 0.0)
		require.NoError(t, err)

		dashes, offset, err := ctx.GetDash()
		require.NoError(t, err)
		assert.Empty(t, dashes, "Nil dash pattern creates solid line")
		assert.Equal(t, 0.0, offset)
	})

	t.Run("toggle between dashed and solid", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		// Start with dashed
		pattern := []float64{10.0, 5.0}
		err = ctx.SetDash(pattern, 0.0)
		require.NoError(t, err)
		dashes, _, err := ctx.GetDash()
		require.NoError(t, err)
		assert.NotEmpty(t, dashes, "Should be dashed")

		// Switch to solid
		err = ctx.SetDash([]float64{}, 0.0)
		require.NoError(t, err)
		dashes, _, err = ctx.GetDash()
		require.NoError(t, err)
		assert.Empty(t, dashes, "Should be solid")

		// Switch back to dashed
		err = ctx.SetDash(pattern, 0.0)
		require.NoError(t, err)
		dashes, _, err = ctx.GetDash()
		require.NoError(t, err)
		assert.NotEmpty(t, dashes, "Should be dashed again")
	})
}

// TestContextGetDashCount tests retrieving the number of dashes in the current pattern.
func TestContextGetDashCount(t *testing.T) {
	t.Run("default dash count is zero", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		count := ctx.GetDashCount()
		assert.Equal(t, 0, count, "Default dash count should be 0")
	})

	t.Run("dash count after setting simple pattern", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		pattern := []float64{10.0, 5.0}
		err = ctx.SetDash(pattern, 0.0)
		require.NoError(t, err)

		count := ctx.GetDashCount()
		assert.Equal(t, 2, count, "Dash count should match pattern length")
	})

	t.Run("dash count after setting complex pattern", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		pattern := []float64{15.0, 5.0, 3.0, 5.0}
		err = ctx.SetDash(pattern, 0.0)
		require.NoError(t, err)

		count := ctx.GetDashCount()
		assert.Equal(t, 4, count, "Dash count should match pattern length")
	})

	t.Run("dash count after clearing with empty slice", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		// Set a pattern first
		pattern := []float64{10.0, 5.0}
		err = ctx.SetDash(pattern, 0.0)
		require.NoError(t, err)
		assert.Equal(t, 2, ctx.GetDashCount(), "Should have 2 dashes initially")

		// Clear with empty slice
		err = ctx.SetDash([]float64{}, 0.0)
		require.NoError(t, err)

		count := ctx.GetDashCount()
		assert.Equal(t, 0, count, "Dash count should be 0 after clearing")
	})

	t.Run("dash count after clearing with nil", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		// Set a pattern first
		pattern := []float64{10.0, 5.0}
		err = ctx.SetDash(pattern, 0.0)
		require.NoError(t, err)
		assert.Equal(t, 2, ctx.GetDashCount(), "Should have 2 dashes initially")

		// Clear with nil
		err = ctx.SetDash(nil, 0.0)
		require.NoError(t, err)

		count := ctx.GetDashCount()
		assert.Equal(t, 0, count, "Dash count should be 0 after clearing with nil")
	})

	t.Run("dash count after context close", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)

		pattern := []float64{10.0, 5.0}
		err = ctx.SetDash(pattern, 0.0)
		require.NoError(t, err)
		_ = ctx.Close()

		// GetDashCount after close should be safe (no panic)
		count := ctx.GetDashCount()
		_ = count // Just verify it doesn't panic
	})
}

// TestContextMiterLimit tests setting and getting miter limit.
// Miter limit controls when to switch from miter to bevel joins for sharp angles.
func TestContextMiterLimit(t *testing.T) {
	t.Run("default miter limit", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		// Default miter limit should be 10.0
		limit := ctx.GetMiterLimit()
		assert.InDelta(t, 10.0, limit, 0.001, "Default miter limit should be 10.0")
	})

	t.Run("set miter limit small", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		ctx.SetMiterLimit(2.0)
		limit := ctx.GetMiterLimit()
		assert.InDelta(t, 2.0, limit, 0.001, "Miter limit should be 2.0")
	})

	t.Run("set miter limit large", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		ctx.SetMiterLimit(50.0)
		limit := ctx.GetMiterLimit()
		assert.InDelta(t, 50.0, limit, 0.001, "Miter limit should be 50.0")
	})

	t.Run("set miter limit minimum", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		// Miter limit of 1.0 (minimum)
		ctx.SetMiterLimit(1.0)
		limit := ctx.GetMiterLimit()
		assert.InDelta(t, 1.0, limit, 0.001, "Miter limit should be 1.0")
	})

	t.Run("change miter limit multiple times", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		ctx.SetMiterLimit(5.0)
		assert.InDelta(t, 5.0, ctx.GetMiterLimit(), 0.001)

		ctx.SetMiterLimit(15.0)
		assert.InDelta(t, 15.0, ctx.GetMiterLimit(), 0.001)

		ctx.SetMiterLimit(3.0)
		assert.InDelta(t, 3.0, ctx.GetMiterLimit(), 0.001)
	})

	t.Run("miter limit after close", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)

		ctx.SetMiterLimit(5.0)
		_ = ctx.Close()

		// Operations after close should be safe no-ops
		ctx.SetMiterLimit(10.0)
		// GetMiterLimit after close should return a reasonable value
		limit := ctx.GetMiterLimit()
		_ = limit // Just verify it doesn't panic
	})
}

// TestContextLineStyleCombinations tests various combinations of line style settings.
func TestContextLineStyleCombinations(t *testing.T) {
	t.Run("set all line styles", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		// Set all line style properties
		ctx.SetLineWidth(5.0)
		ctx.SetLineCap(LineCapRound)
		ctx.SetLineJoin(LineJoinBevel)
		err = ctx.SetDash([]float64{10.0, 5.0}, 2.0)
		require.NoError(t, err)
		ctx.SetMiterLimit(8.0)

		// Verify all settings
		assert.InDelta(t, 5.0, ctx.GetLineWidth(), 0.001)
		assert.Equal(t, LineCapRound, ctx.GetLineCap())
		assert.Equal(t, LineJoinBevel, ctx.GetLineJoin())

		dashes, offset, err := ctx.GetDash()
		require.NoError(t, err)
		assert.Equal(t, []float64{10.0, 5.0}, dashes)
		assert.InDelta(t, 2.0, offset, 0.001)

		assert.InDelta(t, 8.0, ctx.GetMiterLimit(), 0.001)
	})

	t.Run("line styles persist across save restore", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		// Set initial line styles
		ctx.SetLineCap(LineCapRound)
		ctx.SetLineJoin(LineJoinBevel)
		err = ctx.SetDash([]float64{10.0, 5.0}, 0.0)
		require.NoError(t, err)
		ctx.SetMiterLimit(5.0)

		// Save state
		ctx.Save()

		// Change line styles
		ctx.SetLineCap(LineCapSquare)
		ctx.SetLineJoin(LineJoinRound)
		err = ctx.SetDash([]float64{5.0, 5.0}, 1.0)
		require.NoError(t, err)
		ctx.SetMiterLimit(15.0)

		// Verify changed styles
		assert.Equal(t, LineCapSquare, ctx.GetLineCap())
		assert.Equal(t, LineJoinRound, ctx.GetLineJoin())

		// Restore state
		ctx.Restore()

		// Verify original styles are restored
		assert.Equal(t, LineCapRound, ctx.GetLineCap())
		assert.Equal(t, LineJoinBevel, ctx.GetLineJoin())

		dashes, offset, err := ctx.GetDash()
		require.NoError(t, err)
		assert.Equal(t, []float64{10.0, 5.0}, dashes)
		assert.InDelta(t, 0.0, offset, 0.001)

		assert.InDelta(t, 5.0, ctx.GetMiterLimit(), 0.001)
	})

	t.Run("line styles independent of path operations", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		// Set line styles
		ctx.SetLineCap(LineCapRound)
		ctx.SetLineJoin(LineJoinRound)
		err = ctx.SetDash([]float64{10.0, 5.0}, 0.0)
		require.NoError(t, err)

		// Perform path operations
		ctx.MoveTo(10, 10)
		ctx.LineTo(90, 10)
		ctx.LineTo(90, 90)
		ctx.Stroke()

		// Line styles should remain unchanged
		assert.Equal(t, LineCapRound, ctx.GetLineCap())
		assert.Equal(t, LineJoinRound, ctx.GetLineJoin())

		dashes, _, err := ctx.GetDash()
		require.NoError(t, err)
		assert.Equal(t, []float64{10.0, 5.0}, dashes)
	})

	t.Run("different line cap and join combinations", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		// Test various combinations
		combinations := []struct {
			cap  LineCap
			join LineJoin
		}{
			{LineCapButt, LineJoinMiter},
			{LineCapButt, LineJoinRound},
			{LineCapButt, LineJoinBevel},
			{LineCapRound, LineJoinMiter},
			{LineCapRound, LineJoinRound},
			{LineCapRound, LineJoinBevel},
			{LineCapSquare, LineJoinMiter},
			{LineCapSquare, LineJoinRound},
			{LineCapSquare, LineJoinBevel},
		}

		for _, combo := range combinations {
			ctx.SetLineCap(combo.cap)
			ctx.SetLineJoin(combo.join)

			assert.Equal(t, combo.cap, ctx.GetLineCap())
			assert.Equal(t, combo.join, ctx.GetLineJoin())
		}
	})

	t.Run("dash pattern with different widths", func(t *testing.T) {
		surface, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surface.Close()

		ctx, err := NewContext(surface)
		require.NoError(t, err)
		defer ctx.Close()

		// Set dash pattern with varying widths
		patterns := [][]float64{
			{5.0, 5.0},
			{10.0, 5.0, 3.0, 5.0},
			{20.0, 10.0},
			{1.0, 1.0},
		}

		for _, pattern := range patterns {
			err = ctx.SetDash(pattern, 0.0)
			require.NoError(t, err)

			dashes, _, err := ctx.GetDash()
			require.NoError(t, err)
			assert.Equal(t, pattern, dashes)
		}
	})
}
