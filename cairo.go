package cairo

import (
	"github.com/mikowitz/cairo/context"
	"github.com/mikowitz/cairo/surface"
)

// Format is used to identify the memory format of image data.
type Format = surface.Format

// Format constants specify the pixel format for image surfaces.
const (
	// FormatInvalid represents an invalid format.
	FormatInvalid = surface.FormatInvalid

	// FormatARGB32 represents 32-bit ARGB format with alpha in the most significant byte,
	// then red, green, and blue. The 32-bit quantities are stored native-endian.
	// Pre-multiplied alpha is used (i.e. 50% transparent red is 0x80800000, not 0x80ff0000).
	FormatARGB32 = surface.FormatARGB32

	// FormatRGB24 represents 24-bit RGB format stored in 32-bit quantities. The unused
	// bits should be zero. The 32-bit quantities are stored native-endian.
	FormatRGB24 = surface.FormatRGB24

	// FormatA8 represents 8-bit alpha-only format.
	FormatA8 = surface.FormatA8

	// FormatA1 represents 1-bit alpha-only format.
	FormatA1 = surface.FormatA1

	// FormatRGB16_565 represents 16-bit RGB format with 5 bits for red, 6 for green,
	// and 5 for blue.
	FormatRGB16_565 = surface.FormatRGB16_565

	// FormatRGB30 represents 30-bit RGB format with 10 bits per color component,
	// stored in 32-bit quantities.
	FormatRGB30 = surface.FormatRGB30
)

// Surface represents a destination for drawing operations.
//
// A Surface is the abstract type representing all different drawing targets
// that Cairo can render to. The actual drawing is done using a Context
// (available in later prompts).
//
// All Surface implementations must be explicitly closed when finished to free
// Cairo resources, or the finalizer will clean them up during garbage collection.
type Surface = surface.Surface

// NewImageSurface creates an image surface of the specified format and dimensions.
// The initial contents of the surface are set to transparent black (all pixels are
// fully transparent with RGBA values of 0,0,0,0).
//
// The surface should be closed with Close() when finished to release Cairo resources.
// A finalizer is registered as a safety net, but explicit cleanup is recommended.
func NewImageSurface(format Format, width, height int) (*surface.ImageSurface, error) {
	return surface.NewImageSurface(format, width, height)
}

// Context is the main object used for drawing operations in Cairo.
//
// A Context maintains the graphics state including transformations, clip region,
// line width and style, colors, font properties, and more. All drawing operations
// are performed through a Context.
//
// # Basic Usage Pattern
//
// The standard workflow for using a Context:
//
//  1. Create a Surface (drawing target)
//  2. Create a Context for that Surface
//  3. Perform drawing operations
//  4. Close both Context and Surface
//
// Example:
//
//	// Create a 400x300 ARGB32 image surface
//	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 400, 300)
//	if err != nil {
//	    return err
//	}
//	defer surface.Close()
//
//	// Create drawing context
//	ctx, err := cairo.NewContext(surface)
//	if err != nil {
//	    return err
//	}
//	defer ctx.Close()
//
//	// Set source color to opaque red
//	ctx.SetSourceRGB(1.0, 0.0, 0.0)
//
//	// Create a rectangular path and fill it
//	ctx.Rectangle(50.0, 50.0, 200.0, 100.0)
//	ctx.Fill()
//
//	// Draw a line with semi-transparent blue
//	ctx.SetSourceRGBA(0.0, 0.0, 1.0, 0.5)
//	ctx.MoveTo(10.0, 10.0)
//	ctx.LineTo(100.0, 100.0)
//	ctx.Stroke()
//
// # Path Construction
//
// Cairo uses a path-based drawing model. You construct paths using operations
// like MoveTo, LineTo, and Rectangle, then render them with Fill or Stroke.
//
// Basic path operations:
//
//	ctx.MoveTo(x, y)       // Begin a new sub-path at (x, y)
//	ctx.LineTo(x, y)       // Add a line to (x, y)
//	ctx.Rectangle(x, y, w, h) // Add a rectangular path
//	ctx.ClosePath()        // Close the current sub-path
//	ctx.NewPath()          // Clear the current path
//
// Example - Drawing a triangle:
//
//	ctx.MoveTo(50, 10)      // Start at top
//	ctx.LineTo(90, 90)      // Line to bottom-right
//	ctx.LineTo(10, 90)      // Line to bottom-left
//	ctx.ClosePath()         // Complete the triangle
//	ctx.Fill()              // Fill with current source color
//
// # Current Point
//
// Cairo maintains a "current point" which is used as the starting point for
// path operations. The current point is set by operations like MoveTo and
// updated by LineTo. You can query it with GetCurrentPoint() or check if
// one exists with HasCurrentPoint().
//
// The current point is always in user-space coordinates, meaning it's affected
// by any transformations applied to the Context. After operations like NewPath(),
// there is no current point until one is established by MoveTo or similar operations.
//
// Example:
//
//	ctx.MoveTo(25.0, 50.0)
//	if ctx.HasCurrentPoint() {
//	    x, y, _ := ctx.GetCurrentPoint()
//	    fmt.Printf("Current point: (%f, %f)\n", x, y)  // Prints: (25.0, 50.0)
//	}
//	ctx.NewPath()  // Clears path and current point
//	if !ctx.HasCurrentPoint() {
//	    fmt.Println("No current point after NewPath")
//	}
//
// # Resource Management
//
// Always close the Context when finished to release Cairo resources:
//
//	ctx, err := cairo.NewContext(surface)
//	if err != nil {
//	    return err
//	}
//	defer ctx.Close()  // Ensures cleanup
//
// While a finalizer is registered as a safety net, explicit cleanup with Close()
// is strongly recommended, especially in long-running programs.
//
// # State Stack
//
// The Context maintains a stack of graphics states. Use Save() to push the current
// state and Restore() to pop it back:
//
//	ctx.Save()
//	// Modify state (transformations, colors, etc.)
//	// ...
//	ctx.Restore()  // Returns to saved state
//
// # Thread Safety
//
// Context is safe for concurrent use. All methods are protected by appropriate
// locking. However, for optimal performance, avoid concurrent drawing operations
// on the same Context from multiple goroutines.
type Context = context.Context

// NewContext creates a new Context for drawing on the given Surface.
//
// The Context maintains all graphics state for drawing operations. It must be
// explicitly closed with Close() when finished, or rely on the finalizer for
// cleanup during garbage collection.
//
// Returns an error if the surface is nil or if Context creation fails.
//
// Example:
//
//	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 640, 480)
//	if err != nil {
//	    return err
//	}
//	defer surface.Close()
//
//	ctx, err := cairo.NewContext(surface)
//	if err != nil {
//	    return err
//	}
//	defer ctx.Close()
//
//	// Use ctx for drawing operations...
func NewContext(surface Surface) (*Context, error) {
	return context.NewContext(surface)
}
