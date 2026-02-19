// ABOUTME: SurfacePattern implementation for using images/surfaces as pattern sources.
// ABOUTME: Provides texture mapping and pattern fills with extend and filter modes.
package pattern

import (
	"unsafe"

	"github.com/mikowitz/cairo/status"
)

// Surface defines the minimal interface needed for creating surface patterns.
// This interface is satisfied by all Cairo surface types from the surface package.
// unsafe.Pointer is used instead of surface.SurfacePtr to avoid circular imports.
type Surface interface {
	Ptr() unsafe.Pointer
	Status() status.Status
}

// SurfacePattern represents a pattern based on a Cairo surface (image).
//
// Surface patterns allow using existing surfaces (like images) as the source
// for drawing operations. This enables texture mapping, pattern fills, and
// using rendered content as a brush.
//
// The source surface must remain valid for the lifetime of the pattern.
// Closing the source surface before closing the pattern will result in
// undefined behavior.
//
// SurfacePattern embeds BasePattern and implements the Pattern interface,
// providing all standard pattern methods (Close, Status, SetMatrix,
// SetExtend, SetFilter, etc.).
//
// Example:
//
//	// Create a small image to use as a texture
//	surf, err := surface.NewImageSurface(surface.FormatARGB32, 10, 10)
//	if err != nil {
//	    return err
//	}
//	defer surf.Close()
//
//	// Create a pattern from the surface
//	pattern, err := pattern.NewSurfacePattern(surf)
//	if err != nil {
//	    return err
//	}
//	defer pattern.Close()
//
//	// Set pattern to repeat (tile)
//	pattern.SetExtend(pattern.ExtendRepeat)
//
//	// Use pattern for drawing
//	ctx.SetSource(pattern)
//	ctx.Rectangle(0, 0, 100, 100)
//	ctx.Fill()
type SurfacePattern struct {
	*BasePattern
}

// NewSurfacePattern creates a new pattern from a Cairo surface.
//
// The surface can be any Cairo surface type (ImageSurface, PDF, SVG, etc.).
// The pattern will paint using the contents of the surface.
//
// Important: The source surface must remain valid (not closed) for the
// entire lifetime of the pattern. Closing the surface before closing the
// pattern will result in undefined behavior.
//
// By default, the pattern uses ExtendNone (transparent outside bounds)
// and FilterGood (balanced quality/performance) settings. Use SetExtend
// and SetFilter to customize these behaviors.
//
// Parameters:
//   - surface: The Cairo surface to use as the pattern source
//
// Returns a new SurfacePattern and nil error on success.
// Returns nil and an error if pattern creation fails or if the surface
// is invalid.
//
// The returned pattern must be closed with Close() when finished to release
// Cairo resources. A finalizer is registered for safety, but explicit cleanup
// is strongly recommended.
//
// Example:
//
//	surf, err := surface.NewImageSurface(surface.FormatARGB32, 20, 20)
//	if err != nil {
//	    return err
//	}
//	defer surf.Close()
//
//	pattern, err := pattern.NewSurfacePattern(surf)
//	if err != nil {
//	    return err
//	}
//	defer pattern.Close()
func NewSurfacePattern(surface Surface) (*SurfacePattern, error) {
	if surface == nil {
		return nil, status.NullPointer
	}

	if st := surface.Status(); st != status.Success {
		return nil, st
	}

	ptr := patternCreateForSurface(surface.Ptr())
	st := patternStatus(ptr)

	if st != status.Success {
		return nil, st
	}

	basePattern := newBasePattern(ptr, PatternTypeSurface)
	return &SurfacePattern{
		BasePattern: basePattern,
	}, nil
}

