package cairo_test

import (
	"fmt"
	"log"
	"math"

	"github.com/mikowitz/cairo"
)

// Example demonstrates basic usage of the cairo package with re-exported types.
func Example() {
	// Create a 200x200 ARGB32 image surface using re-exported types
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 200, 200)
	if err != nil {
		log.Fatal(err)
	}
	defer surface.Close()

	// Get surface properties
	fmt.Printf("Format: %v\n", surface.GetFormat())
	fmt.Printf("Width: %d\n", surface.GetWidth())
	fmt.Printf("Height: %d\n", surface.GetHeight())
	fmt.Printf("Stride: %d\n", surface.GetStride())

	// Check surface status
	fmt.Printf("Status: %v\n", surface.Status())

	// Output:
	// Format: ARGB32
	// Width: 200
	// Height: 200
	// Stride: 800
	// Status: no error has occurred
}

// Example_formatTypes demonstrates the different pixel formats available.
func Example_formatTypes() {
	formats := []struct {
		format cairo.Format
		desc   string
	}{
		{cairo.FormatARGB32, "32-bit ARGB with alpha"},
		{cairo.FormatRGB24, "24-bit RGB (no alpha)"},
		{cairo.FormatA8, "8-bit alpha only"},
		{cairo.FormatA1, "1-bit alpha only"},
	}

	for _, f := range formats {
		surface, err := cairo.NewImageSurface(f.format, 100, 100)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: stride=%d\n", f.format, surface.GetStride())
		err = surface.Close()
		if err != nil {
			log.Fatal(err)
		}
	}

	// Output:
	// ARGB32: stride=400
	// RGB24: stride=400
	// A8: stride=100
	// A1: stride=16
}

// Example_memoryManagement demonstrates proper cleanup patterns.
func Example_memoryManagement() {
	// Best practice: use defer to ensure cleanup
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 100, 100)
	if err != nil {
		log.Fatal(err)
	}
	defer surface.Close()

	// Use the surface...
	fmt.Printf("Surface created: %dx%d\n", surface.GetWidth(), surface.GetHeight())

	// surface.Close() will be called automatically when function exits

	// Output:
	// Surface created: 100x100
}

// ExampleNewContext demonstrates creating a drawing context from an image surface.
func ExampleNewContext() {
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 200, 200)
	if err != nil {
		log.Fatal(err)
	}
	defer surface.Close()

	ctx, err := cairo.NewContext(surface)
	if err != nil {
		log.Fatal(err)
	}
	defer ctx.Close()

	fmt.Printf("Context status: %v\n", ctx.Status())

	// Output:
	// Context status: no error has occurred
}

// ExampleContext_Arc demonstrates drawing a circle using Arc.
func ExampleContext_Arc() {
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 200, 200)
	if err != nil {
		log.Fatal(err)
	}
	defer surface.Close()

	ctx, err := cairo.NewContext(surface)
	if err != nil {
		log.Fatal(err)
	}
	defer ctx.Close()

	// Paint a white background
	ctx.SetSourceRGB(1, 1, 1)
	ctx.Paint()

	// Draw a filled red circle centered at (100, 100) with radius 50
	ctx.SetSourceRGB(1, 0, 0)
	ctx.Arc(100, 100, 50, 0, 2*math.Pi)
	ctx.Fill()

	fmt.Printf("Draw status: %v\n", ctx.Status())

	// Output:
	// Draw status: no error has occurred
}

// ExampleContext_Save demonstrates the graphics state stack with Save and Restore.
func ExampleContext_Save() {
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 100, 100)
	if err != nil {
		log.Fatal(err)
	}
	defer surface.Close()

	ctx, err := cairo.NewContext(surface)
	if err != nil {
		log.Fatal(err)
	}
	defer ctx.Close()

	ctx.SetLineWidth(2)
	fmt.Printf("Before Save: line width = %.0f\n", ctx.GetLineWidth())

	ctx.Save()
	ctx.SetLineWidth(10) // Temporarily change line width
	fmt.Printf("Inside Save: line width = %.0f\n", ctx.GetLineWidth())
	ctx.Restore()

	fmt.Printf("After Restore: line width = %.0f\n", ctx.GetLineWidth())

	// Output:
	// Before Save: line width = 2
	// Inside Save: line width = 10
	// After Restore: line width = 2
}

// ExampleNewLinearGradient demonstrates creating a linear gradient pattern.
func ExampleNewLinearGradient() {
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 300, 100)
	if err != nil {
		log.Fatal(err)
	}
	defer surface.Close()

	ctx, err := cairo.NewContext(surface)
	if err != nil {
		log.Fatal(err)
	}
	defer ctx.Close()

	// Create a horizontal gradient from red (left) to blue (right)
	grad, err := cairo.NewLinearGradient(0, 0, 300, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer grad.Close()

	grad.AddColorStopRGB(0.0, 1, 0, 0) // Red at start
	grad.AddColorStopRGB(1.0, 0, 0, 1) // Blue at end

	ctx.SetSource(grad)
	ctx.Rectangle(0, 0, 300, 100)
	ctx.Fill()

	fmt.Printf("Gradient type: %v\n", grad.GetType())

	// Output:
	// Gradient type: Linear
}
