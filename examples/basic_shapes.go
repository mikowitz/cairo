package examples

import (
	"fmt"

	"github.com/mikowitz/cairo"
)

// GenerateBasicShapes creates a 400x400 PNG image demonstrating basic Cairo drawing operations.
//
// The image contains:
//   - A filled red rectangle at (100, 100) with dimensions 200x200
//   - A blue stroked rectangle outline at (120, 120) with dimensions 160x160
//
// This function demonstrates the complete Cairo workflow:
//  1. Create an image surface
//  2. Create a drawing context
//  3. Set source colors
//  4. Draw and fill shapes
//  5. Draw and stroke shapes
//  6. Save to PNG
//
// All resources are properly cleaned up using defer statements.
func GenerateBasicShapes(outputPath string) error {
	// Create a 400x400 ARGB32 image surface
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 400, 400)
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

	// Flush any pending operations
	surface.Flush()

	// Write the surface to a PNG file
	if err := surface.WriteToPNG(outputPath); err != nil {
		return fmt.Errorf("failed to write PNG: %w", err)
	}

	return nil
}
