package examples

import (
	"fmt"
	"os"

	"github.com/mikowitz/cairo"
)

// Example_drawRectangle demonstrates the basic workflow for drawing a filled rectangle.
// This shows the fundamental steps: create a surface, create a context, set a color,
// draw a rectangle path, and fill it.
func Example_drawRectangle() {
	// Create a 200x200 image surface
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 200, 200)
	if err != nil {
		panic(err)
	}
	defer surface.Close()

	// Create a drawing context
	ctx, err := cairo.NewContext(surface)
	if err != nil {
		panic(err)
	}
	defer ctx.Close()

	// Fill background with white
	ctx.SetSourceRGB(1.0, 1.0, 1.0)
	ctx.Paint()

	// Set source color to red
	ctx.SetSourceRGB(1.0, 0.0, 0.0)

	// Create a rectangle path at (50, 50) with dimensions 100x100
	ctx.Rectangle(50, 50, 100, 100)

	// Fill the rectangle
	ctx.Fill()

	// Write to a temporary file for testing
	tmpFile, err := os.CreateTemp("", "example_rectangle_*.png")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpFile.Name())

	surface.Flush()
	if err := surface.WriteToPNG(tmpFile.Name()); err != nil {
		panic(err)
	}

	fmt.Println("Rectangle drawn successfully")
	// Output: Rectangle drawn successfully
}

// Example_fillAndStroke demonstrates using both fill and stroke operations on the same path.
// This shows how to use FillPreserve to fill a shape while keeping the path for subsequent
// stroke operations.
func Example_fillAndStroke() {
	// Create a 200x200 image surface
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 200, 200)
	if err != nil {
		panic(err)
	}
	defer surface.Close()

	// Create a drawing context
	ctx, err := cairo.NewContext(surface)
	if err != nil {
		panic(err)
	}
	defer ctx.Close()

	// Fill background with white
	ctx.SetSourceRGB(1.0, 1.0, 1.0)
	ctx.Paint()

	// Create a rectangle path
	ctx.Rectangle(50, 50, 100, 100)

	// Set fill color to light blue and fill (preserving the path)
	ctx.SetSourceRGB(0.7, 0.8, 1.0)
	ctx.FillPreserve()

	// Set stroke color to dark blue and stroke the outline
	ctx.SetSourceRGB(0.0, 0.0, 0.5)
	ctx.SetLineWidth(3.0)
	ctx.Stroke()

	// Write to a temporary file for testing
	tmpFile, err := os.CreateTemp("", "example_fill_stroke_*.png")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpFile.Name())

	surface.Flush()
	if err := surface.WriteToPNG(tmpFile.Name()); err != nil {
		panic(err)
	}

	fmt.Println("Fill and stroke completed successfully")
	// Output: Fill and stroke completed successfully
}

// Example_colorBlending demonstrates drawing multiple shapes with different colors.
// This shows how to create multiple shapes by setting different source colors between
// drawing operations.
func Example_colorBlending() {
	// Create a 300x300 image surface
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 300, 300)
	if err != nil {
		panic(err)
	}
	defer surface.Close()

	// Create a drawing context
	ctx, err := cairo.NewContext(surface)
	if err != nil {
		panic(err)
	}
	defer ctx.Close()

	// Fill background with light gray
	ctx.SetSourceRGB(0.9, 0.9, 0.9)
	ctx.Paint()

	// Draw first rectangle in red
	ctx.SetSourceRGB(1.0, 0.0, 0.0)
	ctx.Rectangle(50, 50, 100, 100)
	ctx.Fill()

	// Draw second rectangle in green (overlapping with first)
	ctx.SetSourceRGB(0.0, 1.0, 0.0)
	ctx.Rectangle(100, 100, 100, 100)
	ctx.Fill()

	// Draw third rectangle in blue (overlapping with both)
	ctx.SetSourceRGB(0.0, 0.0, 1.0)
	ctx.Rectangle(75, 125, 100, 100)
	ctx.Fill()

	// Draw a semi-transparent rectangle using RGBA
	ctx.SetSourceRGBA(1.0, 1.0, 0.0, 0.5)
	ctx.Rectangle(150, 75, 100, 100)
	ctx.Fill()

	// Write to a temporary file for testing
	tmpFile, err := os.CreateTemp("", "example_color_blending_*.png")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpFile.Name())

	surface.Flush()
	if err := surface.WriteToPNG(tmpFile.Name()); err != nil {
		panic(err)
	}

	fmt.Println("Multiple colored shapes drawn successfully")
	// Output: Multiple colored shapes drawn successfully
}
