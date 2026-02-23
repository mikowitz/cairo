package context

import (
	"math"
	"testing"
)

// TestContextArc verifies arc path creation.
// An arc is a portion of a circle from angle1 to angle2 (in radians),
// drawn counter-clockwise from the positive x-axis.
func TestContextArc(t *testing.T) {
	ctx, _ := newTestContext(t, 200, 200)

	t.Run("full_circle", func(t *testing.T) {
		// Create a full circle using Arc with 0 to 2π
		ctx.NewPath()
		ctx.Arc(100, 100, 50, 0, 2*math.Pi)

		// Verify the path was created by checking we can stroke it
		ctx.SetLineWidth(2.0)
		ctx.SetSourceRGB(0, 0, 0)
		ctx.Stroke()

		// No panic means success
	})

	t.Run("quarter_circle", func(t *testing.T) {
		// Create a quarter circle (0 to π/2)
		ctx.NewPath()
		ctx.Arc(100, 100, 30, 0, math.Pi/2)
		ctx.Stroke()

		// No panic means success
	})

	t.Run("arc_with_line_to", func(t *testing.T) {
		// Draw arc connected to a line
		ctx.NewPath()
		ctx.MoveTo(50, 50)
		ctx.Arc(100, 100, 40, 0, math.Pi)
		ctx.LineTo(150, 150)
		ctx.Stroke()

		// No panic means success
	})
}

// TestContextArcNegative verifies negative arc (clockwise arc).
// ArcNegative draws arcs in the clockwise direction.
func TestContextArcNegative(t *testing.T) {
	ctx, _ := newTestContext(t, 200, 200)

	t.Run("clockwise_arc", func(t *testing.T) {
		// Create a clockwise arc from π to 0
		ctx.NewPath()
		ctx.ArcNegative(100, 100, 50, math.Pi, 0)
		ctx.Stroke()

		// No panic means success
	})

	t.Run("negative_vs_positive", func(t *testing.T) {
		// Draw both positive and negative arcs to show they're different
		ctx.NewPath()
		ctx.SetSourceRGB(1, 0, 0) // Red for positive
		ctx.Arc(80, 100, 30, 0, math.Pi)
		ctx.Stroke()

		ctx.NewPath()
		ctx.SetSourceRGB(0, 0, 1) // Blue for negative
		ctx.ArcNegative(120, 100, 30, 0, math.Pi)
		ctx.Stroke()

		// No panic means success
	})
}

// TestContextCurveTo verifies Bezier curve creation.
// CurveTo adds a cubic Bezier curve from the current point to (x3, y3),
// using (x1, y1) and (x2, y2) as control points.
func TestContextCurveTo(t *testing.T) {
	ctx, _ := newTestContext(t, 200, 200)

	t.Run("simple_curve", func(t *testing.T) {
		// Draw a simple S-curve
		ctx.NewPath()
		ctx.MoveTo(20, 100)
		ctx.CurveTo(80, 20, 120, 180, 180, 100)
		ctx.SetLineWidth(2.0)
		ctx.Stroke()

		// No panic means success
	})

	t.Run("multiple_curves", func(t *testing.T) {
		// Chain multiple curves together
		ctx.NewPath()
		ctx.MoveTo(10, 50)
		ctx.CurveTo(30, 10, 70, 10, 90, 50)
		ctx.CurveTo(110, 90, 150, 90, 170, 50)
		ctx.Stroke()

		// No panic means success
	})

	t.Run("closed_curve_shape", func(t *testing.T) {
		// Create a closed shape with curves
		ctx.NewPath()
		ctx.MoveTo(100, 50)
		ctx.CurveTo(120, 50, 150, 70, 150, 100)
		ctx.CurveTo(150, 130, 120, 150, 100, 150)
		ctx.CurveTo(80, 150, 50, 130, 50, 100)
		ctx.CurveTo(50, 70, 80, 50, 100, 50)
		ctx.ClosePath()
		ctx.Fill()

		// No panic means success
	})
}

// TestContextRelativeOperations verifies relative path operations work correctly.
// Relative operations (RelMoveTo, RelLineTo, RelCurveTo) use offsets from the
// current point rather than absolute coordinates.
func TestContextRelativeOperations(t *testing.T) {
	ctx, _ := newTestContext(t, 200, 200)

	t.Run("rel_move_to", func(t *testing.T) {
		// Draw disconnected lines using RelMoveTo
		ctx.NewPath()
		ctx.MoveTo(20, 20)
		ctx.LineTo(40, 40)
		ctx.RelMoveTo(20, 0) // Move right by 20
		ctx.LineTo(80, 40)
		ctx.Stroke()

		// No panic means success
	})

	t.Run("rel_line_to", func(t *testing.T) {
		// Draw a path using relative movements
		ctx.NewPath()
		ctx.MoveTo(100, 100)
		ctx.RelLineTo(30, 0)   // Right
		ctx.RelLineTo(0, 30)   // Down
		ctx.RelLineTo(-30, 0)  // Left
		ctx.RelLineTo(0, -30)  // Up (back to start)
		ctx.Stroke()

		// No panic means success
	})

	t.Run("rel_curve_to", func(t *testing.T) {
		// Draw a curve using relative coordinates
		ctx.NewPath()
		ctx.MoveTo(50, 100)
		// All coordinates are relative to current point (50, 100)
		ctx.RelCurveTo(20, -40, 60, -40, 80, 0)
		ctx.Stroke()

		// No panic means success
	})

	t.Run("combined_relative_operations", func(t *testing.T) {
		// Create a complex path with mixed relative operations
		ctx.NewPath()
		ctx.MoveTo(20, 150)
		ctx.RelLineTo(20, -20)
		ctx.RelCurveTo(10, 10, 30, 10, 40, 0)
		ctx.RelLineTo(20, 20)
		ctx.RelMoveTo(10, 0)
		ctx.RelLineTo(30, -30)
		ctx.Stroke()

		// No panic means success
	})

	t.Run("relative_with_transformations", func(t *testing.T) {
		// Test that relative operations work correctly with transformations
		ctx.Save()
		ctx.Translate(50, 50)
		ctx.Scale(1.5, 1.5)

		ctx.NewPath()
		ctx.MoveTo(0, 0)
		ctx.RelLineTo(20, 0)
		ctx.RelLineTo(0, 20)
		ctx.RelLineTo(-20, 0)
		ctx.ClosePath()
		ctx.Stroke()

		ctx.Restore()

		// No panic means success
	})
}

// TestContextCircle verifies that Arc can create a complete circle.
// This is a common use case: drawing a full circle from 0 to 2π.
func TestContextCircle(t *testing.T) {
	ctx, _ := newTestContext(t, 300, 300)

	t.Run("filled_circle", func(t *testing.T) {
		// Draw a filled circle
		ctx.NewPath()
		ctx.Arc(150, 150, 100, 0, 2*math.Pi)
		ctx.SetSourceRGB(1, 0, 0)
		ctx.Fill()

		// No panic means success
	})

	t.Run("stroked_circle", func(t *testing.T) {
		// Draw a circle outline
		ctx.NewPath()
		ctx.Arc(150, 150, 80, 0, 2*math.Pi)
		ctx.SetSourceRGB(0, 0, 1)
		ctx.SetLineWidth(3.0)
		ctx.Stroke()

		// No panic means success
	})

	t.Run("multiple_circles", func(t *testing.T) {
		// Draw multiple circles of different sizes
		radii := []float64{20, 40, 60}
		for _, radius := range radii {
			ctx.NewPath()
			ctx.Arc(150, 150, radius, 0, 2*math.Pi)
			ctx.Stroke()
		}

		// No panic means success
	})

	t.Run("circle_with_fill_preserve", func(t *testing.T) {
		// Draw a circle with both fill and stroke
		ctx.NewPath()
		ctx.Arc(150, 150, 50, 0, 2*math.Pi)
		ctx.SetSourceRGB(1, 1, 0)
		ctx.FillPreserve()
		ctx.SetSourceRGB(0, 0, 0)
		ctx.SetLineWidth(2.0)
		ctx.Stroke()

		// No panic means success
	})
}
