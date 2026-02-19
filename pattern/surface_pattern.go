// ABOUTME: SurfacePattern implementation for using images/surfaces as pattern sources.
// ABOUTME: Provides texture mapping and pattern fills with extend and filter modes.
package pattern

import (
	"github.com/mikowitz/cairo/status"
)

// Surface defines the minimal interface needed for creating surface patterns.
// This interface is satisfied by all Cairo surface types (ImageSurface, etc.).
type Surface interface {
	Ptr() interface{}
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
// providing all standard pattern methods (Close, Status, SetMatrix, etc.)
// plus additional methods for controlling how the surface is sampled and
// extended beyond its bounds.
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

// SetExtend sets the extend mode for the pattern.
//
// The extend mode controls what happens when the pattern is sampled outside
// its natural bounds (i.e., when you try to paint an area larger than the
// source surface):
//
//   - ExtendNone: Transparent outside bounds (default)
//   - ExtendRepeat: Pattern tiles/repeats infinitely
//   - ExtendReflect: Pattern mirrors at edges
//   - ExtendPad: Edge colors extend infinitely
//
// This is particularly useful for tiling small textures across large areas.
//
// Example:
//
//	pattern.SetExtend(pattern.ExtendRepeat)  // Tile the pattern
func (p *SurfacePattern) SetExtend(extend Extend) {
	p.Lock()
	defer p.Unlock()

	if p.ptr == nil {
		return
	}

	patternSetExtend(p.ptr, extend)
}

// GetExtend returns the current extend mode for the pattern.
//
// Returns the Extend mode that was previously set with SetExtend, or
// ExtendNone if no extend mode was explicitly set.
//
// Example:
//
//	if pattern.GetExtend() == pattern.ExtendRepeat {
//	    fmt.Println("Pattern is set to tile")
//	}
func (p *SurfacePattern) GetExtend() Extend {
	p.RLock()
	defer p.RUnlock()

	if p.ptr == nil {
		return ExtendNone
	}

	return patternGetExtend(p.ptr)
}

// SetFilter sets the filter mode for the pattern.
//
// The filter mode controls how the pattern is resampled when the pattern
// matrix (via SetMatrix) causes scaling or rotation:
//
//   - FilterFast/FilterNearest: Fast but pixelated (nearest-neighbor)
//   - FilterGood/FilterBilinear: Balanced quality/speed (bilinear, default)
//   - FilterBest: Highest quality, potentially slower
//
// Choose faster filters for performance-critical operations or when you want
// a pixelated/retro aesthetic. Choose better filters for high-quality output.
//
// Example:
//
//	pattern.SetFilter(pattern.FilterNearest)  // Pixelated scaling
func (p *SurfacePattern) SetFilter(filter Filter) {
	p.Lock()
	defer p.Unlock()

	if p.ptr == nil {
		return
	}

	patternSetFilter(p.ptr, filter)
}

// GetFilter returns the current filter mode for the pattern.
//
// Returns the Filter mode that was previously set with SetFilter, or
// FilterGood if no filter was explicitly set (Cairo's default).
//
// Example:
//
//	if pattern.GetFilter() == pattern.FilterNearest {
//	    fmt.Println("Using nearest-neighbor filtering")
//	}
func (p *SurfacePattern) GetFilter() Filter {
	p.RLock()
	defer p.RUnlock()

	if p.ptr == nil {
		return FilterGood
	}

	return patternGetFilter(p.ptr)
}
