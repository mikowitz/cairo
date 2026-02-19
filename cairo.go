package cairo

import (
	"unsafe"

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
//   - Linear gradients: Color gradients along a line (NewLinearGradient)
//   - Radial gradients: Color gradients in a circular pattern (NewRadialGradient)
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

// LinearGradient represents a gradient pattern that transitions colors along a line.
//
// Linear gradients create smooth color transitions between two or more colors along
// a line defined by two points. Colors are specified using color stops, which define
// the color at specific positions along the gradient line.
//
// # Creating Linear Gradients
//
// Create a linear gradient using NewLinearGradient, specifying the start and end points:
//
//	gradient, err := cairo.NewLinearGradient(0, 0, 200, 0)  // Horizontal gradient
//	if err != nil {
//	    return err
//	}
//	defer gradient.Close()
//
// Add color stops to define the gradient colors:
//
//	gradient.AddColorStopRGB(0.0, 1.0, 0.0, 0.0)  // Red at start (0%)
//	gradient.AddColorStopRGB(0.5, 0.0, 1.0, 0.0)  // Green at middle (50%)
//	gradient.AddColorStopRGB(1.0, 0.0, 0.0, 1.0)  // Blue at end (100%)
//
// Use the gradient as a drawing source:
//
//	ctx.SetSource(gradient)
//	ctx.Rectangle(0, 0, 200, 100)
//	ctx.Fill()
//
// # Gradient Direction
//
// The gradient line defines the direction of color transition. Colors are interpolated
// perpendicular to this line, extending infinitely in both directions. Areas before
// the start point take the first color stop, and areas after the end point take the
// last color stop.
//
// Examples:
//   - Horizontal: NewLinearGradient(0, 0, 100, 0) - left to right
//   - Vertical: NewLinearGradient(0, 0, 0, 100) - top to bottom
//   - Diagonal: NewLinearGradient(0, 0, 100, 100) - top-left to bottom-right
//
// # Resource Management
//
// Linear gradients must be explicitly closed when finished to release Cairo resources:
//
//	gradient, err := cairo.NewLinearGradient(0, 0, 200, 0)
//	if err != nil {
//	    return err
//	}
//	defer gradient.Close()  // Essential
//
// For more details, see the pattern package documentation.
type LinearGradient = pattern.LinearGradient

// RadialGradient represents a gradient pattern that transitions colors between two circles.
//
// Radial gradients create smooth color transitions from one circle to another, allowing
// for effects like spotlights, glows, and radial color fades. Colors are specified using
// color stops that define the color at specific positions between the two circles.
//
// # Creating Radial Gradients
//
// Create a radial gradient using NewRadialGradient, specifying two circles:
//
//	// Gradient from small inner circle to large outer circle (spotlight effect)
//	gradient, err := cairo.NewRadialGradient(100, 100, 10, 100, 100, 100)
//	if err != nil {
//	    return err
//	}
//	defer gradient.Close()
//
// Add color stops to define the gradient colors:
//
//	gradient.AddColorStopRGB(0.0, 1.0, 1.0, 1.0)  // White at center
//	gradient.AddColorStopRGB(1.0, 0.0, 0.0, 1.0)  // Blue at edge
//
// Use the gradient as a drawing source:
//
//	ctx.SetSource(gradient)
//	ctx.Arc(100, 100, 100, 0, 2*math.Pi)  // Draw circle
//	ctx.Fill()
//
// # Gradient Effects
//
// Different circle configurations create different visual effects:
//
//   - Concentric circles (same center): Creates a uniform radial gradient
//     NewRadialGradient(50, 50, 10, 50, 50, 100)
//
//   - Offset centers: Creates directional lighting or highlight effects
//     NewRadialGradient(40, 40, 10, 60, 60, 100)
//
//   - Zero inner radius: Gradient starts from a point
//     NewRadialGradient(50, 50, 0, 50, 50, 100)
//
// # Transparency Effects
//
// Use AddColorStopRGBA to create gradients that fade in or out:
//
//	gradient.AddColorStopRGBA(0.0, 1.0, 0.5, 0.0, 1.0)  // Opaque orange center
//	gradient.AddColorStopRGBA(1.0, 1.0, 0.5, 0.0, 0.0)  // Transparent orange edge
//
// # Resource Management
//
// Radial gradients must be explicitly closed when finished to release Cairo resources:
//
//	gradient, err := cairo.NewRadialGradient(100, 100, 10, 100, 100, 100)
//	if err != nil {
//	    return err
//	}
//	defer gradient.Close()  // Essential
//
// For more details, see the pattern package documentation.
type RadialGradient = pattern.RadialGradient

// NewLinearGradient creates a new linear gradient pattern along the line from (x0, y0) to (x1, y1).
//
// The gradient coordinates are in pattern space, which initially matches user space.
// After creation, you must add color stops using AddColorStopRGB or AddColorStopRGBA
// to define the gradient colors.
//
// Parameters:
//   - x0, y0: Starting point coordinates of the gradient line
//   - x1, y1: Ending point coordinates of the gradient line
//
// The gradient line defines the direction of color transition. Color stop offset 0.0
// corresponds to the start point (x0, y0), and offset 1.0 corresponds to the end
// point (x1, y1). Cairo interpolates colors smoothly between stops.
//
// The returned gradient must be closed with Close() when finished to release
// Cairo resources. A finalizer is registered for safety, but explicit cleanup
// is strongly recommended.
//
// Example - Simple horizontal gradient:
//
//	// Create gradient from left to right
//	gradient, err := cairo.NewLinearGradient(0, 0, 200, 0)
//	if err != nil {
//	    return err
//	}
//	defer gradient.Close()
//
//	// Add color stops: red to blue
//	gradient.AddColorStopRGB(0.0, 1.0, 0.0, 0.0)  // Red at start
//	gradient.AddColorStopRGB(1.0, 0.0, 0.0, 1.0)  // Blue at end
//
//	// Use the gradient
//	ctx.SetSource(gradient)
//	ctx.Rectangle(0, 0, 200, 100)
//	ctx.Fill()
//
// Example - Multi-color rainbow gradient:
//
//	gradient, err := cairo.NewLinearGradient(0, 0, 300, 0)
//	if err != nil {
//	    return err
//	}
//	defer gradient.Close()
//
//	gradient.AddColorStopRGB(0.0, 1.0, 0.0, 0.0)   // Red
//	gradient.AddColorStopRGB(0.33, 0.0, 1.0, 0.0)  // Green
//	gradient.AddColorStopRGB(0.67, 0.0, 0.0, 1.0)  // Blue
//	gradient.AddColorStopRGB(1.0, 1.0, 1.0, 0.0)   // Yellow
//
//	ctx.SetSource(gradient)
//	ctx.Paint()
//
// Example - Gradient with transparency:
//
//	gradient, err := cairo.NewLinearGradient(0, 0, 0, 200)  // Vertical
//	if err != nil {
//	    return err
//	}
//	defer gradient.Close()
//
//	gradient.AddColorStopRGBA(0.0, 0.0, 0.0, 0.0, 1.0)  // Opaque black at top
//	gradient.AddColorStopRGBA(1.0, 0.0, 0.0, 0.0, 0.0)  // Transparent at bottom
//
//	ctx.SetSource(gradient)
//	ctx.Paint()  // Creates fade-out effect
func NewLinearGradient(x0, y0, x1, y1 float64) (*LinearGradient, error) {
	return pattern.NewLinearGradient(x0, y0, x1, y1)
}

// NewRadialGradient creates a new radial gradient pattern between two circles.
//
// The gradient interpolates colors from the start circle (cx0, cy0, radius0) to
// the end circle (cx1, cy1, radius1). After creation, you must add color stops
// using AddColorStopRGB or AddColorStopRGBA to define the gradient colors.
//
// Parameters:
//   - cx0, cy0: Center coordinates of the start circle
//   - radius0: Radius of the start circle
//   - cx1, cy1: Center coordinates of the end circle
//   - radius1: Radius of the end circle
//
// The gradient coordinates are in pattern space, which initially matches user space.
// Color stop offset 0.0 corresponds to the start circle, and offset 1.0 corresponds
// to the end circle. Cairo interpolates colors smoothly between stops.
//
// The returned gradient must be closed with Close() when finished to release
// Cairo resources. A finalizer is registered for safety, but explicit cleanup
// is strongly recommended.
//
// Example - Simple radial gradient (spotlight effect):
//
//	// Gradient from small center to large edge
//	gradient, err := cairo.NewRadialGradient(100, 100, 10, 100, 100, 100)
//	if err != nil {
//	    return err
//	}
//	defer gradient.Close()
//
//	// Add color stops: white center to blue edge
//	gradient.AddColorStopRGB(0.0, 1.0, 1.0, 1.0)  // White at center
//	gradient.AddColorStopRGB(1.0, 0.0, 0.0, 1.0)  // Blue at edge
//
//	// Draw with the gradient
//	ctx.SetSource(gradient)
//	ctx.Arc(100, 100, 100, 0, 2*math.Pi)
//	ctx.Fill()
//
// Example - Offset gradient (lighting effect):
//
//	// Offset centers create directional highlight
//	gradient, err := cairo.NewRadialGradient(80, 80, 20, 120, 120, 100)
//	if err != nil {
//	    return err
//	}
//	defer gradient.Close()
//
//	gradient.AddColorStopRGB(0.0, 1.0, 1.0, 0.8)  // Pale yellow highlight
//	gradient.AddColorStopRGB(0.5, 1.0, 0.5, 0.0)  // Orange
//	gradient.AddColorStopRGB(1.0, 0.5, 0.0, 0.0)  // Dark red
//
//	ctx.SetSource(gradient)
//	ctx.Paint()
//
// Example - Fade-out effect with transparency:
//
//	// Gradient from point (radius0=0) with transparency
//	gradient, err := cairo.NewRadialGradient(150, 150, 0, 150, 150, 100)
//	if err != nil {
//	    return err
//	}
//	defer gradient.Close()
//
//	gradient.AddColorStopRGBA(0.0, 1.0, 0.5, 0.0, 1.0)  // Opaque orange center
//	gradient.AddColorStopRGBA(0.7, 1.0, 0.0, 0.5, 0.5)  // Semi-transparent pink
//	gradient.AddColorStopRGBA(1.0, 1.0, 0.0, 0.0, 0.0)  // Transparent red edge
//
//	ctx.SetSource(gradient)
//	ctx.Arc(150, 150, 100, 0, 2*math.Pi)
//	ctx.Fill()  // Creates glow effect
func NewRadialGradient(cx0, cy0, radius0, cx1, cy1, radius1 float64) (*RadialGradient, error) {
	return pattern.NewRadialGradient(cx0, cy0, radius0, cx1, cy1, radius1)
}

// SurfacePattern represents a pattern based on a Cairo surface (image).
//
// Surface patterns allow using existing surfaces (like images) as the source
// for drawing operations. This enables texture mapping, pattern fills, and
// using rendered content as a brush.
//
// # Creating Surface Patterns
//
// Create a surface pattern from an existing surface:
//
//	// Create a small image to use as a texture
//	surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 20, 20)
//	if err != nil {
//	    return err
//	}
//	defer surf.Close()
//
//	// ... draw something on the surface ...
//
//	// Create pattern from the surface
//	pattern, err := cairo.NewSurfacePattern(surface)
//	if err != nil {
//	    return err
//	}
//	defer pattern.Close()
//
// # Extend Modes
//
// Control what happens outside the pattern bounds using SetExtend:
//
//	pattern.SetExtend(cairo.ExtendRepeat)   // Tile the pattern
//	pattern.SetExtend(cairo.ExtendReflect)  // Mirror at edges
//	pattern.SetExtend(cairo.ExtendPad)      // Extend edge colors
//	pattern.SetExtend(cairo.ExtendNone)     // Transparent outside (default)
//
// # Filter Modes
//
// Control resampling quality using SetFilter:
//
//	pattern.SetFilter(cairo.FilterNearest)   // Fast, pixelated
//	pattern.SetFilter(cairo.FilterBilinear)  // Smooth, balanced (default)
//	pattern.SetFilter(cairo.FilterBest)      // Highest quality
//
// # Important Notes
//
// The source surface must remain valid (not closed) for the entire lifetime
// of the pattern. Closing the surface before closing the pattern will result
// in undefined behavior.
//
// # Resource Management
//
// Surface patterns must be explicitly closed when finished:
//
//	pattern, err := cairo.NewSurfacePattern(surface)
//	if err != nil {
//	    return err
//	}
//	defer pattern.Close()  // Essential
//
// For more details, see the pattern package documentation.
type SurfacePattern = pattern.SurfacePattern

// NewSurfacePattern creates a new surface pattern from an existing surface.
//
// Surface patterns allow using existing surfaces (like images) as the source
// for drawing operations. This enables texture mapping, pattern fills, and
// using rendered content as a brush.
//
// The source surface must remain valid (not closed) for the entire lifetime
// of the pattern. Closing the surface before closing the pattern will result
// in undefined behavior.
//
// The returned pattern must be closed with Close() when finished to release
// Cairo resources. A finalizer is registered for safety, but explicit cleanup
// is strongly recommended.
//
// Example - Simple texture pattern:
//
//	// Create a small image to use as a texture
//	surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 20, 20)
//	if err != nil {
//	    return err
//	}
//	defer surf.Close()
//
//	// Draw something on the texture surface
//	ctx, err := cairo.NewContext(surf)
//	if err != nil {
//	    return err
//	}
//	defer ctx.Close()
//
//	ctx.SetSourceRGB(1.0, 0.0, 0.0)
//	ctx.Rectangle(0, 0, 10, 10)
//	ctx.Fill()
//	ctx.SetSourceRGB(0.0, 0.0, 1.0)
//	ctx.Rectangle(10, 10, 10, 10)
//	ctx.Fill()
//
//	// Create pattern from the surface
//	pattern, err := cairo.NewSurfacePattern(surf)
//	if err != nil {
//	    return err
//	}
//	defer pattern.Close()
//
//	// Configure pattern to repeat (tile)
//	pattern.SetExtend(cairo.ExtendRepeat)
//
//	// Use the pattern on a larger surface
//	mainSurf, err := cairo.NewImageSurface(cairo.FormatARGB32, 200, 200)
//	if err != nil {
//	    return err
//	}
//	defer mainSurf.Close()
//
//	mainCtx, err := cairo.NewContext(mainSurf)
//	if err != nil {
//	    return err
//	}
//	defer mainCtx.Close()
//
//	mainCtx.SetSource(pattern)
//	mainCtx.Paint()  // Fills entire surface with tiled pattern
//
// Example - Pattern with filter control:
//
//	pattern, err := cairo.NewSurfacePattern(imageSurface)
//	if err != nil {
//	    return err
//	}
//	defer pattern.Close()
//
//	// Use nearest-neighbor for pixelated effect
//	pattern.SetFilter(cairo.FilterNearest)
//	pattern.SetExtend(cairo.ExtendRepeat)
//
//	ctx.SetSource(pattern)
//	ctx.Rectangle(0, 0, 400, 300)
//	ctx.Fill()
func NewSurfacePattern(surface Surface) (*SurfacePattern, error) {
	return pattern.NewSurfacePattern(surfaceAdapter{surface})
}

// surfaceAdapter adapts surface.Surface to work with pattern.NewSurfacePattern.
// surface.Surface has Ptr() SurfacePtr, but pattern.Surface requires Ptr() unsafe.Pointer
// to avoid a circular import between the surface and pattern packages.
type surfaceAdapter struct {
	Surface
}

func (s surfaceAdapter) Ptr() unsafe.Pointer {
	return unsafe.Pointer(s.Surface.Ptr())
}

// Extend defines how patterns behave outside their natural bounds.
//
// When a pattern (gradient or surface pattern) is used to paint an area
// larger than the pattern naturally covers, the extend mode determines
// what happens in the areas outside the pattern's bounds.
type Extend = pattern.Extend

const (
	// ExtendNone means the pattern is not painted outside its natural bounds.
	// Areas outside the pattern will be transparent.
	ExtendNone Extend = pattern.ExtendNone

	// ExtendRepeat means the pattern is tiled by repeating.
	// The pattern repeats infinitely in all directions.
	ExtendRepeat Extend = pattern.ExtendRepeat

	// ExtendReflect means the pattern is tiled by reflecting at the edges.
	// Creates a mirrored repetition effect.
	ExtendReflect Extend = pattern.ExtendReflect

	// ExtendPad means the pattern extends by using the closest color from its edge.
	// The edge pixels are repeated infinitely outward.
	ExtendPad Extend = pattern.ExtendPad
)

// Filter defines the filtering algorithm used when sampling patterns.
//
// When a pattern is transformed (scaled, rotated), Cairo needs to resample
// the pattern pixels. The filter mode determines the quality and speed of
// this resampling operation.
type Filter = pattern.Filter

const (
	// FilterFast uses a high-performance filter with lower quality.
	// Equivalent to nearest-neighbor filtering.
	FilterFast Filter = pattern.FilterFast

	// FilterGood balances quality and performance.
	// Uses bilinear interpolation.
	FilterGood Filter = pattern.FilterGood

	// FilterBest uses the highest-quality filter available.
	// May be slower but produces the best visual results.
	FilterBest Filter = pattern.FilterBest

	// FilterNearest uses nearest-neighbor sampling.
	// Fast but can produce pixelated results when scaling.
	FilterNearest Filter = pattern.FilterNearest

	// FilterBilinear uses bilinear interpolation.
	// Smoother than nearest-neighbor with reasonable performance.
	FilterBilinear Filter = pattern.FilterBilinear

	// FilterGaussian uses gaussian interpolation.
	// Currently unimplemented.
	FilterGaussian Filter = pattern.FilterGaussian
)

// LineCap specifies how the endpoints of lines are rendered when stroking.
//
// The line cap style only affects the endpoints of lines. The appearance of
// line joins is controlled by [LineJoin].
type LineCap = context.LineCap

const (
	// LineCapButt starts and stops the line exactly at the start and end points.
	LineCapButt LineCap = iota

	// LineCapRound uses a round ending, with the center of the circle at the end point.
	LineCapRound

	// LineCapSquare uses a squared ending, with the center of the square at the end point.
	LineCapSquare
)

// LineJoin specifies how the junctions between line segments are rendered when stroking.
//
// The line join style only affects the junctions between line segments. The appearance
// of line endpoints is controlled by [LineCap].
type LineJoin = context.LineJoin

const (
	// LineJoinMiter uses a sharp (angled) corner. If the miter would extend beyond
	// the miter limit (as set by [Context.SetMiterLimit]), a bevel join is used instead.
	LineJoinMiter LineJoin = iota

	// LineJoinRound uses a rounded join, with the center of the circle at the join point.
	LineJoinRound

	// LineJoinBevel uses a cut-off join, with the join cut off at half the line width
	// from the join point.
	LineJoinBevel
)
