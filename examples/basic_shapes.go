package examples

import (
	"fmt"
	"math"

	"github.com/mikowitz/cairo"
)

// GenerateBasicShapes creates a 600x400 PNG image demonstrating basic Cairo drawing operations.
//
// The image contains:
//   - A filled red rectangle
//   - A blue stroked rectangle outline
//   - A filled green circle (using Arc)
//   - An orange stroked circle
//   - A purple curved path (using CurveTo)
//
// This function demonstrates the complete Cairo workflow:
//  1. Create an image surface
//  2. Create a drawing context
//  3. Set source colors
//  4. Draw and fill shapes (rectangles, circles, curves)
//  5. Draw and stroke shapes
//  6. Save to PNG
//
// All resources are properly cleaned up using defer statements.
func GenerateBasicShapes(outputPath string) error {
	// Create a 600x400 ARGB32 image surface
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 600, 400)
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

	// Fill the background with white
	ctx.SetSourceRGB(1.0, 1.0, 1.0)
	ctx.Paint()

	// Draw a filled red rectangle at (100, 100) with dimensions 200x200
	ctx.SetSourceRGB(1.0, 0.0, 0.0) // Red
	ctx.Rectangle(100.0, 100.0, 200.0, 200.0)
	ctx.Fill()

	// Draw a blue stroked rectangle outline at (120, 120) with dimensions 160x160
	ctx.SetSourceRGB(0.0, 0.0, 1.0) // Blue
	ctx.SetLineWidth(5.0)           // 5-pixel wide line
	ctx.Rectangle(120.0, 120.0, 160.0, 160.0)
	ctx.Stroke()

	// Draw a filled green circle using Arc (at right side of canvas)
	ctx.SetSourceRGB(0.0, 0.8, 0.0) // Green
	ctx.Arc(450, 100, 60, 0, 2*math.Pi)
	ctx.Fill()

	// Draw an orange stroked circle
	ctx.SetSourceRGB(1.0, 0.5, 0.0) // Orange
	ctx.SetLineWidth(3.0)
	ctx.Arc(450, 100, 40, 0, 2*math.Pi)
	ctx.Stroke()

	// Draw a purple curved path using CurveTo
	ctx.SetSourceRGB(0.6, 0.0, 0.8) // Purple
	ctx.SetLineWidth(4.0)
	ctx.MoveTo(350, 250)
	ctx.CurveTo(400, 200, 500, 300, 550, 250)
	ctx.Stroke()

	// Draw another curve - a smooth S-shape
	ctx.SetSourceRGB(0.0, 0.6, 0.8) // Cyan
	ctx.SetLineWidth(3.0)
	ctx.MoveTo(350, 350)
	ctx.CurveTo(400, 300, 450, 380, 500, 330)
	ctx.Stroke()

	// Flush any pending operations
	surface.Flush()

	// Write the surface to a PNG file
	if err := surface.WriteToPNG(outputPath); err != nil {
		return fmt.Errorf("failed to write PNG: %w", err)
	}

	return nil
}
