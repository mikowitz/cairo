package examples

import (
	"fmt"
	"math"

	"github.com/mikowitz/cairo"
)

// GenerateTransformations creates a 600x600 PNG image demonstrating Cairo transformation operations.
//
// The image shows the same "house" shape (rectangle + triangle roof) drawn multiple times
// using different transformations:
//   - Top left: Original shape with no transformation
//   - Top right: Translated to a different position
//   - Middle left: Scaled (made larger)
//   - Middle right: Rotated 45 degrees
//   - Bottom: Combined transformation (translate, rotate, scale)
//
// This demonstrates:
//   - ctx.Translate(tx, ty) - move the origin
//   - ctx.Scale(sx, sy) - resize coordinates
//   - ctx.Rotate(radians) - rotate the coordinate system
//   - ctx.Save() and ctx.Restore() - preserve transformation state
//   - How transformations are cumulative
//   - How transformations affect subsequent drawing operations
//
// All resources are properly cleaned up using defer statements.
func GenerateTransformations(outputPath string) error {
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

	// Fill the background with white
	ctx.SetSourceRGB(1.0, 1.0, 1.0)
	ctx.Paint()

	// Set line width for all shapes
	ctx.SetLineWidth(2.0)

	// Draw the house shape at 5 different positions with different transformations

	// 1. Top-left: No transformation (identity)
	drawHouseWithTransformation(ctx, func(ctx *cairo.Context) {
		ctx.Translate(75, 75)
	})

	// 2. Top-right: Simple translation
	drawHouseWithTransformation(ctx, func(ctx *cairo.Context) {
		ctx.Translate(375, 75)
	})

	// 3. Middle-left: Scaling
	drawHouseWithTransformation(ctx, func(ctx *cairo.Context) {
		ctx.Translate(75, 250)
		ctx.Scale(1.5, 1.5)
	})

	// 4. Middle-right: Rotation
	drawHouseWithTransformation(ctx, func(ctx *cairo.Context) {
		ctx.Translate(450, 300)
		ctx.Rotate(math.Pi / 4)
	})

	// 5. Bottom: Combined transformations
	drawHouseWithTransformation(ctx, func(ctx *cairo.Context) {
		ctx.Translate(300, 500)
		ctx.Rotate(-math.Pi / 6)
		ctx.Scale(0.8, 1.2)
	})

	// Flush any pending operations
	surface.Flush()

	// Write the surface to a PNG file
	if err := surface.WriteToPNG(outputPath); err != nil {
		return fmt.Errorf("failed to write PNG: %w", err)
	}

	return nil
}

func drawHouseWithTransformation(ctx *cairo.Context, transformation func(*cairo.Context)) {
	ctx.Save()
	transformation(ctx)
	drawHouse(ctx)
	ctx.Restore()
}

// drawHouse draws a simple house shape (rectangle with triangular roof) at the current
// transformation origin. The house is 80 units wide and 100 units tall.
func drawHouse(ctx *cairo.Context) {
	// Draw the house body (rectangle)
	ctx.SetSourceRGB(0.8, 0.4, 0.2) // Brown
	ctx.Rectangle(0, 30, 80, 70)
	ctx.FillPreserve()
	ctx.SetSourceRGB(0.0, 0.0, 0.0) // Black outline
	ctx.Stroke()

	// Draw the roof (triangle)
	ctx.SetSourceRGB(0.8, 0.2, 0.2) // Red
	ctx.MoveTo(0, 30)               // Bottom-left of roof
	ctx.LineTo(40, 0)               // Top point (center)
	ctx.LineTo(80, 30)              // Bottom-right of roof
	ctx.ClosePath()
	ctx.FillPreserve()
	ctx.SetSourceRGB(0.0, 0.0, 0.0) // Black outline
	ctx.Stroke()

	// Draw a door
	ctx.SetSourceRGB(0.4, 0.2, 0.1) // Dark brown
	ctx.Rectangle(30, 60, 20, 40)
	ctx.FillPreserve()
	ctx.SetSourceRGB(0.0, 0.0, 0.0) // Black outline
	ctx.Stroke()

	// Draw a window
	ctx.SetSourceRGB(0.7, 0.9, 1.0) // Light blue
	ctx.Rectangle(10, 45, 15, 15)
	ctx.FillPreserve()
	ctx.SetSourceRGB(0.0, 0.0, 0.0) // Black outline
	ctx.Stroke()
}
