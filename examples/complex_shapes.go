package examples

import (
	"fmt"
	"math"

	"github.com/mikowitz/cairo"
)

// GenerateComplexShapes creates a comprehensive demonstration of all currently
// implemented Cairo Context functionality.
//
// The 600x600 image demonstrates all 20 Context methods:
// - Lifecycle: Save, Restore, Status
// - Colors: SetSourceRGB, SetSourceRGBA
// - Paths: MoveTo, LineTo, Rectangle, ClosePath, NewPath, NewSubPath
// - Queries: GetCurrentPoint, HasCurrentPoint
// - Rendering: Fill, FillPreserve, Stroke, StrokePreserve, Paint
// - Line properties: SetLineWidth, GetLineWidth
//
// The image is organized into sections, each demonstrating specific features.
func GenerateComplexShapes(outputPath string) error {
	// Create a 600x600 ARGB32 image surface
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 600, 600)
	if err != nil {
		return fmt.Errorf("failed to create surface: %w", err)
	}
	defer func() {
		_ = surface.Close()
	}()

	// Create drawing context
	ctx, err := cairo.NewContext(surface)
	if err != nil {
		return fmt.Errorf("failed to create context: %w", err)
	}
	defer func() {
		_ = ctx.Close()
	}()

	// Section 1: Paint() - Fill background with light gray
	ctx.SetSourceRGB(0.95, 0.95, 0.95)
	ctx.Paint()

	// Check Status() after Paint
	if st := ctx.Status(); st != 0 {
		return fmt.Errorf("context error after Paint: %v", st)
	}

	// Section 2: Save/Restore with nested rectangles (top-left quadrant)
	// Demonstrates state management
	drawNestedRectanglesWithSaveRestore(ctx)

	// Section 3: RGB vs RGBA colors (top-right quadrant)
	// Demonstrates opaque and translucent colors
	drawColorDemonstration(ctx)

	// Section 4: Complex paths with MoveTo, LineTo, ClosePath (bottom-left)
	// Demonstrates path construction and current point queries
	drawComplexPaths(ctx)

	// Section 5: Fill vs Stroke variations (bottom-right)
	// Demonstrates Fill, FillPreserve, Stroke, StrokePreserve
	drawFillStrokeVariations(ctx)

	// Section 6: Line width variations (center vertical stripe)
	// Demonstrates SetLineWidth and GetLineWidth
	drawLineWidthVariations(ctx)

	// Section 7: NewSubPath demonstration (center)
	// Demonstrates disconnected paths in single operation
	drawMultipleSubPaths(ctx)

	// Final status check
	if st := ctx.Status(); st != 0 {
		return fmt.Errorf("context error at end: %v", st)
	}

	// Flush any pending operations
	surface.Flush()

	// Write the surface to a PNG file
	if err := surface.WriteToPNG(outputPath); err != nil {
		return fmt.Errorf("failed to write PNG: %w", err)
	}

	return nil
}

// drawNestedRectanglesWithSaveRestore demonstrates Save() and Restore()
// by drawing nested rectangles where each inner rectangle is drawn with
// a saved/restored state.
func drawNestedRectanglesWithSaveRestore(ctx *cairo.Context) {
	// Starting position: top-left quadrant
	baseX, baseY := 50.0, 50.0

	colors := []struct{ r, g, b float64 }{
		{1.0, 0.0, 0.0}, // Red
		{1.0, 0.5, 0.0}, // Orange
		{1.0, 1.0, 0.0}, // Yellow
		{0.0, 1.0, 0.0}, // Green
	}

	for i, color := range colors {
		// Save the current state
		ctx.Save()

		// Set color and draw rectangle
		ctx.SetSourceRGB(color.r, color.g, color.b)
		size := 200.0 - float64(i)*40.0
		ctx.Rectangle(baseX+float64(i)*20.0, baseY+float64(i)*20.0, size, size)
		ctx.Stroke()

		// Restore the saved state
		ctx.Restore()
	}
}

// drawColorDemonstration shows the difference between SetSourceRGB (opaque)
// and SetSourceRGBA (translucent)
func drawColorDemonstration(ctx *cairo.Context) {
	// Starting position: top-right quadrant
	baseX, baseY := 350.0, 50.0

	// Opaque rectangles with SetSourceRGB
	ctx.SetSourceRGB(1.0, 0.0, 0.0) // Red
	ctx.Rectangle(baseX, baseY, 80, 80)
	ctx.Fill()

	ctx.SetSourceRGB(0.0, 0.0, 1.0) // Blue
	ctx.Rectangle(baseX+40, baseY+40, 80, 80)
	ctx.Fill()

	// Translucent overlapping circles with SetSourceRGBA
	ctx.NewPath()
	ctx.SetSourceRGBA(1.0, 0.0, 0.0, 0.5) // Semi-transparent red
	ctx.Rectangle(baseX, baseY+130, 80, 80)
	ctx.Fill()

	ctx.SetSourceRGBA(0.0, 1.0, 0.0, 0.5) // Semi-transparent green
	ctx.Rectangle(baseX+40, baseY+170, 80, 80)
	ctx.Fill()
}

// drawComplexPaths demonstrates MoveTo, LineTo, ClosePath, and current point queries
func drawComplexPaths(ctx *cairo.Context) {
	// Starting position: bottom-left quadrant
	baseX, baseY := 80.0, 350.0

	// Draw a triangle using MoveTo, LineTo, and ClosePath
	ctx.NewPath() // Clear any existing path

	// Check if we have a current point before starting
	if ctx.HasCurrentPoint() {
		// If we do, get it (though we just called NewPath so this should be false)
		x, y, _ := ctx.GetCurrentPoint()
		fmt.Printf("Unexpected current point before triangle: (%f, %f)\n", x, y)
	}

	// Build triangle path
	ctx.MoveTo(baseX, baseY)
	ctx.LineTo(baseX+60, baseY)
	ctx.LineTo(baseX+30, baseY-60)
	ctx.ClosePath()

	// After building path, we should have a current point
	if ctx.HasCurrentPoint() {
		x, y, err := ctx.GetCurrentPoint()
		if err == nil {
			// Current point should be back at start due to ClosePath
			_ = x // Use the coordinates (they're at the starting point)
			_ = y
		}
	}

	// Stroke the triangle
	ctx.SetSourceRGB(0.2, 0.2, 0.8)
	ctx.SetLineWidth(3.0)
	ctx.Stroke()

	// Draw a star shape using multiple LineTo calls
	ctx.NewPath()
	starCenterX, starCenterY := baseX+120, baseY-30.0

	// 5-pointed star
	points := 5
	outerRadius := 40.0
	innerRadius := 15.0

	for i := 0; i < points*2; i++ {
		angle := float64(i)*math.Pi/float64(points) - math.Pi/2.0 // Start from top
		radius := outerRadius
		if i%2 == 1 {
			radius = innerRadius
		}

		x := starCenterX + radius*math.Cos(angle)
		y := starCenterY + radius*math.Sin(angle)

		if i == 0 {
			ctx.MoveTo(x, y)
		} else {
			ctx.LineTo(x, y)
		}
	}
	ctx.ClosePath()

	// Fill the star
	ctx.SetSourceRGB(1.0, 0.8, 0.0)
	ctx.Fill()
}

// drawFillStrokeVariations demonstrates Fill, FillPreserve, Stroke, and StrokePreserve
func drawFillStrokeVariations(ctx *cairo.Context) {
	// Starting position: bottom-right quadrant
	baseX, baseY := 350.0, 350.0

	// Example 1: Fill() - path is consumed
	ctx.NewPath()
	ctx.Rectangle(baseX, baseY, 60, 60)
	ctx.SetSourceRGB(0.8, 0.2, 0.2)
	ctx.Fill() // Path consumed after this

	// Example 2: FillPreserve() then Stroke() - path preserved for second operation
	ctx.NewPath()
	ctx.Rectangle(baseX+80, baseY, 60, 60)
	ctx.SetSourceRGB(0.2, 0.8, 0.2)
	ctx.FillPreserve() // Path NOT consumed
	ctx.SetSourceRGB(0.0, 0.0, 0.0)
	ctx.SetLineWidth(2.0)
	ctx.Stroke() // Now path is consumed

	// Example 3: Just Stroke() - path is consumed
	ctx.NewPath()
	ctx.Rectangle(baseX, baseY+80, 60, 60)
	ctx.SetSourceRGB(0.2, 0.2, 0.8)
	ctx.SetLineWidth(4.0)
	ctx.Stroke() // Path consumed

	// Example 4: StrokePreserve() then Stroke() again with different width
	ctx.NewPath()
	ctx.Rectangle(baseX+80, baseY+80, 60, 60)
	ctx.SetSourceRGB(0.8, 0.2, 0.8)
	ctx.SetLineWidth(2.0)
	ctx.StrokePreserve() // Path preserved
	ctx.SetSourceRGB(1.0, 1.0, 1.0)
	ctx.SetLineWidth(1.0)
	ctx.Stroke() // Inner stroke with different color
}

// drawLineWidthVariations demonstrates SetLineWidth and GetLineWidth
func drawLineWidthVariations(ctx *cairo.Context) {
	// Center vertical stripe
	startX := 280.0
	startY := 150.0

	widths := []float64{1.0, 3.0, 5.0, 10.0, 15.0}

	for i, width := range widths {
		// Set the line width
		ctx.SetLineWidth(width)

		// Verify it was set correctly using GetLineWidth
		actualWidth := ctx.GetLineWidth()
		if actualWidth != width {
			fmt.Printf("Warning: expected width %f, got %f\n", width, actualWidth)
		}

		// Draw a horizontal line
		y := startY + float64(i)*20.0
		ctx.NewPath()
		ctx.MoveTo(startX, y)
		ctx.LineTo(startX+40.0, y)
		ctx.SetSourceRGB(0.3, 0.3, 0.3)
		ctx.Stroke()
	}
}

// drawMultipleSubPaths demonstrates NewSubPath() for creating disconnected paths
// that are filled/stroked in a single operation
func drawMultipleSubPaths(ctx *cairo.Context) {
	// Center area
	centerX, centerY := 300.0, 280.0

	// Create multiple disconnected rectangles using NewSubPath
	ctx.NewPath()

	// First sub-path
	ctx.Rectangle(centerX-60, centerY, 20, 20)

	// Start a new sub-path (disconnected from first)
	ctx.NewSubPath()
	ctx.Rectangle(centerX-20, centerY, 20, 20)

	// Another sub-path
	ctx.NewSubPath()
	ctx.Rectangle(centerX+20, centerY, 20, 20)

	// Fill all three rectangles in one operation
	ctx.SetSourceRGB(0.5, 0.0, 0.5)
	ctx.Fill()

	// Demonstrate with stroke as well
	ctx.NewPath()
	ctx.Rectangle(centerX-60, centerY+30, 20, 20)
	ctx.NewSubPath()
	ctx.Rectangle(centerX-20, centerY+30, 20, 20)
	ctx.NewSubPath()
	ctx.Rectangle(centerX+20, centerY+30, 20, 20)

	ctx.SetSourceRGB(0.0, 0.5, 0.5)
	ctx.SetLineWidth(2.0)
	ctx.Stroke()
}
