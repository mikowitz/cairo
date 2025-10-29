// Package pattern provides Cairo pattern sources for drawing operations.
//
// # Overview
//
// Patterns are the "paint" that Cairo uses to draw. They define what colors,
// gradients, or images to use when filling or stroking paths. Every drawing
// operation in Cairo uses a pattern as its source.
//
// Cairo supports several pattern types:
//   - Solid colors (implemented in this package)
//   - Linear gradients (planned)
//   - Radial gradients (planned)
//   - Surface patterns (for texturing with images) (planned)
//   - Mesh patterns (for complex gradients) (planned)
//
// # Pattern Types
//
// Solid Patterns:
//
// The simplest pattern type represents a single solid color with optional
// transparency. Create solid patterns with NewSolidPatternRGB or
// NewSolidPatternRGBA:
//
//	// Opaque red
//	pattern, err := pattern.NewSolidPatternRGB(1.0, 0.0, 0.0)
//	if err != nil {
//	    return err
//	}
//	defer pattern.Close()
//
//	// Semi-transparent blue
//	pattern, err := pattern.NewSolidPatternRGBA(0.0, 0.0, 1.0, 0.5)
//	if err != nil {
//	    return err
//	}
//	defer pattern.Close()
//
// # Using Patterns with Context
//
// Patterns are used as the "source" for drawing operations. Set a pattern
// as the active source using Context.SetSource:
//
//	ctx, err := cairo.NewContext(surface)
//	if err != nil {
//	    return err
//	}
//	defer ctx.Close()
//
//	// Create a pattern
//	pattern, err := pattern.NewSolidPatternRGB(1.0, 0.0, 0.0)
//	if err != nil {
//	    return err
//	}
//	defer pattern.Close()
//
//	// Use the pattern for drawing
//	ctx.SetSource(pattern)
//	ctx.Rectangle(10, 10, 50, 50)
//	ctx.Fill()  // Fills the rectangle with red
//
// Convenience Methods:
//
// For simple solid colors, Context provides convenience methods that create
// and set solid patterns internally:
//
//	ctx.SetSourceRGB(1.0, 0.0, 0.0)  // Equivalent to creating and setting a solid pattern
//	ctx.Rectangle(10, 10, 50, 50)
//	ctx.Fill()
//
// Use explicit pattern creation when you need to:
//   - Reuse the same pattern for multiple operations
//   - Apply transformations to the pattern (via SetMatrix)
//   - Query pattern properties
//   - Work with non-solid pattern types (gradients, textures)
//
// # Pattern Transformations
//
// Patterns have their own transformation matrix, independent of the Context's
// transformation. This affects how the pattern is mapped onto the surface:
//
//	pattern, err := pattern.NewSolidPatternRGB(1.0, 0.0, 0.0)
//	if err != nil {
//	    return err
//	}
//	defer pattern.Close()
//
//	// Apply a transformation to the pattern
//	m := matrix.NewScaleMatrix(2.0, 2.0)
//	pattern.SetMatrix(m)
//
//	// Query the pattern's matrix
//	currentMatrix, err := pattern.GetMatrix()
//	if err != nil {
//	    return err
//	}
//
// Pattern transformations are most useful with gradient and texture patterns,
// where they control how the pattern is positioned and oriented.
//
// # Resource Management
//
// Patterns must be explicitly closed when finished to release Cairo resources:
//
//	pattern, err := pattern.NewSolidPatternRGB(1.0, 0.0, 0.0)
//	if err != nil {
//	    return err
//	}
//	defer pattern.Close()  // Essential for cleanup
//
// A finalizer is registered as a safety net, but explicit cleanup with Close()
// is strongly recommended for predictable resource management.
//
// Important: When using Context.GetSource(), the returned pattern has proper
// reference counting and can be safely garbage collected without explicit Close().
// However, patterns you create with New* functions should always be closed explicitly.
//
// # Reference Counting
//
// Cairo uses reference counting for patterns. When you create a pattern, you own
// a reference to it. When you set it as a Context's source, the Context takes
// its own reference. This means:
//
//	pattern, err := pattern.NewSolidPatternRGB(1.0, 0.0, 0.0)
//	if err != nil {
//	    return err
//	}
//	ctx.SetSource(pattern)
//	pattern.Close()  // Safe - Context has its own reference
//
//	// Pattern is still active in the Context
//	ctx.Rectangle(10, 10, 50, 50)
//	ctx.Fill()  // Works fine
//
// When retrieving the current source from a Context, the returned pattern
// has its own reference:
//
//	source, err := ctx.GetSource()
//	if err != nil {
//	    return err
//	}
//	// No need to explicitly Close() - GC will handle it
//	// But you can if you want to: defer source.Close()
//
// # Thread Safety
//
// Pattern methods are safe for concurrent use from multiple goroutines.
// All methods are protected by appropriate read/write locking.
package pattern
