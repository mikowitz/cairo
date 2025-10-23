package cairo_test

import (
	"fmt"
	"log"

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
		surface.Close()
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
	defer surface.Close() // Ensures resources are freed

	// Use the surface...
	fmt.Printf("Surface created: %dx%d\n", surface.GetWidth(), surface.GetHeight())

	// surface.Close() will be called automatically when function exits

	// Output:
	// Surface created: 100x100
}
