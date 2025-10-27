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
// The Context lifecycle follows these stages:
//
//  1. Creation: NewContext() creates a Context bound to a Surface. The Surface
//     must be valid and not closed. The Context increments the Surface's reference
//     count, keeping it alive even if you close the Surface object (Cairo manages
//     the underlying C resources).
//
//  2. Active Use: During this phase, you can:
//     - Set drawing parameters (colors, line width, transformations, etc.)
//     - Build paths (MoveTo, LineTo, Rectangle, Arc, etc.)
//     - Render paths (Fill, Stroke, Paint)
//     - Manage graphics state via Save/Restore
//     All operations are thread-safe due to internal locking.
//
//  3. State Preservation: The graphics state includes:
//     - Current transformation matrix (CTM)
//     - Current path and current point
//     - Source pattern (color or gradient)
//     - Line width, cap, join, dash pattern
//     - Fill rule, operator, tolerance
//     - Clip region, anti-aliasing mode
//     Save/Restore operations preserve all these settings.
//
//  4. Cleanup: Close() should be called when drawing is complete to immediately
//     release Cairo resources. After Close():
//     - All drawing operations become no-ops
//     - Status() returns NullPointer
//     - The underlying Surface reference is released
//     - Double-closing is safe (subsequent calls are ignored)
//
//  5. Finalizer Safety Net: If Close() is never called, a finalizer will eventually
//     clean up the resources during garbage collection. However, relying on this
//     can delay resource release significantly. Always prefer explicit Close(),
//     typically via defer.
//
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
//	// Set up drawing parameters
//	ctx.SetSourceRGB(1.0, 0.0, 0.0)  // Red color
//	ctx.SetLineWidth(2.0)
//
//	// Build and render a path
//	ctx.Rectangle(50, 50, 100, 100)
//	ctx.Fill()
//
//	// The defer statements ensure proper cleanup in reverse order
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
