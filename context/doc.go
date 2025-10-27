// Package context provides the main drawing context for Cairo operations.
//
// The Context is Cairo's central object for drawing operations. It maintains
// all graphics state parameters including the current transformation matrix,
// clip region, line width, line style, colors, font properties, and more.
//
// # Drawing Pipeline
//
// The typical Cairo drawing workflow follows this pattern:
//
//  1. Create a Surface (the drawing target)
//  2. Create a Context for that Surface
//  3. Set drawing properties (colors, line width, etc.)
//  4. Construct paths (MoveTo, LineTo, Rectangle, etc.)
//  5. Render paths (Fill, Stroke, Paint, etc.)
//  6. Close the Context and Surface when done
//
// # Lifecycle and Resource Management
//
// A Context must be created for a specific Surface and holds a reference to it.
// The Context should be explicitly closed with Close() when drawing is complete
// to release Cairo resources. A finalizer is registered as a safety net, but
// explicit cleanup is strongly recommended for long-running programs.
//
// Example usage:
//
//	surface, err := surface.NewImageSurface(surface.FormatARGB32, 400, 300)
//	if err != nil {
//	    return err
//	}
//	defer surface.Close()
//
//	ctx, err := context.NewContext(surface)
//	if err != nil {
//	    return err
//	}
//	defer ctx.Close()
//
//	// Drawing operations will be added in subsequent prompts
//	// ctx.SetSourceRGB(1.0, 0.0, 0.0)  // Red color
//	// ctx.Rectangle(50, 50, 100, 100)
//	// ctx.Fill()
//
// # State Management
//
// The Context maintains a stack of graphics states. Use Save() to push the
// current state onto the stack and Restore() to pop it back. This is useful
// for temporarily changing drawing parameters:
//
//	ctx.Save()
//	// Make temporary changes
//	ctx.Restore()  // Returns to previous state
//
// # Thread Safety
//
// Context is safe for concurrent use. All methods use appropriate locking
// to protect the internal state. However, for best performance, avoid
// concurrent drawing operations on the same Context.
package context
