package cairo

import (
	"github.com/mikowitz/cairo/context"
	"github.com/mikowitz/cairo/pattern"
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
// # Rendering Operations
//
// After constructing a path, you render it using Fill, Stroke, or Paint operations.
// Understanding the difference between these operations and their "Preserve" variants
// is crucial for effective Cairo usage.
//
// Fill vs Stroke:
//
//   - Fill(): Paints the interior of the path using the current source pattern.
//     The path is consumed (cleared) after the operation.
//
//   - Stroke(): Paints along the outline of the path according to line width,
//     line join, line cap, and dash settings. The path is consumed after the operation.
//
//   - Paint(): Applies the source pattern to the entire clipped region, independent
//     of any path. Useful for setting backgrounds.
//
// Example - Simple fill and stroke:
//
//	// Fill a rectangle
//	ctx.SetSourceRGB(1.0, 0.0, 0.0)  // Red
//	ctx.Rectangle(50, 50, 100, 100)
//	ctx.Fill()  // Path is now cleared
//
//	// Stroke a line
//	ctx.SetSourceRGB(0.0, 0.0, 1.0)  // Blue
//	ctx.SetLineWidth(3.0)
//	ctx.MoveTo(10, 10)
//	ctx.LineTo(100, 100)
//	ctx.Stroke()  // Path is now cleared
//
// Path Preservation:
//
// The standard Fill() and Stroke() operations consume the current path, clearing
// it after rendering. The "Preserve" variants (FillPreserve and StrokePreserve)
// keep the path intact, allowing multiple operations on the same path.
//
// This is particularly useful for creating shapes with both fill and outline:
//
// Example - Fill and stroke the same shape:
//
//	ctx.Rectangle(50, 50, 100, 100)
//	ctx.SetSourceRGBA(0.0, 1.0, 0.0, 0.7)  // Semi-transparent green fill
//	ctx.FillPreserve()  // Fill, but keep the path
//	ctx.SetSourceRGB(0.0, 0.0, 0.0)        // Black outline
//	ctx.SetLineWidth(2.0)
//	ctx.Stroke()  // Now the path is cleared
//
// When to use Preserve variants:
//
//   - When you need both fill and stroke on the same path
//   - When you want to apply multiple operations with different settings
//   - When you need to query path information after rendering
//
// Line Width:
//
// The line width affects how Stroke operations render. It specifies the diameter
// of the pen used for stroking, in user-space units:
//
//	ctx.SetLineWidth(5.0)   // Thick lines
//	ctx.SetLineWidth(1.0)   // Thin lines
//	width := ctx.GetLineWidth()  // Query current width
//
// The default line width is 2.0. Line width is affected by transformations,
// so scaling the Context will also scale the rendered line width.
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

// Pattern is the interface that all Cairo pattern types implement.
//
// Patterns represent the "paint" that Cairo uses for drawing operations. They define
// what colors, gradients, or images to use when filling or stroking paths.
//
// # Pattern Types
//
// Cairo supports several pattern types:
//   - Solid colors: Single uniform colors (NewSolidPatternRGB, NewSolidPatternRGBA)
//   - Linear gradients: Color gradients along a line (planned)
//   - Radial gradients: Color gradients in a circular pattern (planned)
//   - Surface patterns: Texturing with images (planned)
//   - Mesh patterns: Complex multi-point gradients (planned)
//
// # Using Patterns
//
// Create a pattern and set it as the drawing source using Context.SetSource:
//
//	pattern, err := cairo.NewSolidPatternRGB(1.0, 0.0, 0.0)
//	if err != nil {
//	    return err
//	}
//	defer pattern.Close()
//
//	ctx.SetSource(pattern)
//	ctx.Rectangle(10, 10, 100, 100)
//	ctx.Fill()  // Fills with red
//
// For simple solid colors, Context provides convenience methods:
//
//	ctx.SetSourceRGB(1.0, 0.0, 0.0)  // Equivalent to creating a solid pattern
//	ctx.Fill()
//
// # Resource Management
//
// Patterns must be explicitly closed when finished:
//
//	pattern, err := cairo.NewSolidPatternRGB(1.0, 0.0, 0.0)
//	if err != nil {
//	    return err
//	}
//	defer pattern.Close()  // Essential
//
// Exception: Patterns returned from Context.GetSource() have proper reference
// counting and can be safely garbage collected without explicit Close().
//
// For more details, see the pattern package documentation.
type Pattern = pattern.Pattern

// NewSolidPatternRGB creates a new solid pattern with an opaque RGB color.
//
// The color components should be in the range [0.0, 1.0]:
//   - 0.0 represents no intensity (black for that channel)
//   - 1.0 represents full intensity (maximum brightness)
//
// The alpha channel is implicitly set to 1.0 (fully opaque).
//
// The returned pattern must be closed with Close() when finished to release
// Cairo resources. A finalizer is registered for safety, but explicit cleanup
// is strongly recommended.
//
// Example:
//
//	// Create an opaque red pattern
//	red, err := cairo.NewSolidPatternRGB(1.0, 0.0, 0.0)
//	if err != nil {
//	    return err
//	}
//	defer red.Close()
//
//	// Use it for drawing
//	ctx.SetSource(red)
//	ctx.Rectangle(10, 10, 50, 50)
//	ctx.Fill()
//
// For simple cases, consider using Context.SetSourceRGB instead:
//
//	ctx.SetSourceRGB(1.0, 0.0, 0.0)
//	ctx.Rectangle(10, 10, 50, 50)
//	ctx.Fill()
//
// Use explicit pattern creation when you need to:
//   - Reuse the same pattern for multiple operations
//   - Apply transformations to the pattern
//   - Store patterns for later use
func NewSolidPatternRGB(r, g, b float64) (*pattern.SolidPattern, error) {
	return pattern.NewSolidPatternRGB(r, g, b)
}

// NewSolidPatternRGBA creates a new solid pattern with an RGBA color including transparency.
//
// The color components should be in the range [0.0, 1.0]:
//   - r, g, b: Color channels where 0.0 = no intensity, 1.0 = full intensity
//   - a: Alpha (transparency) where 0.0 = fully transparent, 1.0 = fully opaque
//
// Alpha compositing in Cairo uses premultiplied alpha. This function handles
// the premultiplication internally, so you should provide unpremultiplied values.
//
// The returned pattern must be closed with Close() when finished to release
// Cairo resources. A finalizer is registered for safety, but explicit cleanup
// is strongly recommended.
//
// Example:
//
//	// Create a semi-transparent blue pattern
//	blue, err := cairo.NewSolidPatternRGBA(0.0, 0.0, 1.0, 0.5)
//	if err != nil {
//	    return err
//	}
//	defer blue.Close()
//
//	// Use it for drawing
//	ctx.SetSource(blue)
//	ctx.Rectangle(10, 10, 50, 50)
//	ctx.Fill()  // Fills with 50% transparent blue
//
// For simple cases, consider using Context.SetSourceRGBA instead:
//
//	ctx.SetSourceRGBA(0.0, 0.0, 1.0, 0.5)
//	ctx.Rectangle(10, 10, 50, 50)
//	ctx.Fill()
//
// Use explicit pattern creation when you need to:
//   - Reuse the same pattern for multiple operations
//   - Apply transformations to the pattern
//   - Store patterns for later use
func NewSolidPatternRGBA(r, g, b, a float64) (*pattern.SolidPattern, error) {
	return pattern.NewSolidPatternRGBA(r, g, b, a)
}
