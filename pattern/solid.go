package pattern

import "github.com/mikowitz/cairo/status"

// SolidPattern represents a pattern with a single, uniform color.
//
// Solid patterns are the simplest pattern type in Cairo. They paint with
// a single color that can optionally include transparency via the alpha channel.
//
// Create solid patterns using NewSolidPatternRGB for opaque colors or
// NewSolidPatternRGBA for colors with transparency.
//
// SolidPattern embeds BasePattern and implements the Pattern interface,
// providing all standard pattern methods (Close, Status, SetMatrix, etc.).
//
// Example:
//
//	// Create an opaque red pattern
//	red, err := pattern.NewSolidPatternRGB(1.0, 0.0, 0.0)
//	if err != nil {
//	    return err
//	}
//	defer red.Close()
//
//	// Create a semi-transparent blue pattern
//	blue, err := pattern.NewSolidPatternRGBA(0.0, 0.0, 1.0, 0.5)
//	if err != nil {
//	    return err
//	}
//	defer blue.Close()
type SolidPattern struct {
	*BasePattern
}

// NewSolidPatternRGB creates a new solid pattern with an opaque RGB color.
//
// The color components (r, g, b) should be in the range [0.0, 1.0]:
//   - 0.0 represents no intensity (black for that channel)
//   - 1.0 represents full intensity (maximum brightness for that channel)
//
// The alpha channel is implicitly set to 1.0 (fully opaque).
//
// Parameters:
//   - r: Red component (0.0 to 1.0)
//   - g: Green component (0.0 to 1.0)
//   - b: Blue component (0.0 to 1.0)
//
// Returns a new SolidPattern and nil error on success.
// Returns nil and an error if pattern creation fails.
//
// The returned pattern must be closed with Close() when finished to release
// Cairo resources. A finalizer is registered for safety, but explicit cleanup
// is strongly recommended.
//
// Example:
//
//	// Pure red
//	red, err := pattern.NewSolidPatternRGB(1.0, 0.0, 0.0)
//	if err != nil {
//	    return err
//	}
//	defer red.Close()
//
//	// Gray (50% intensity)
//	gray, err := pattern.NewSolidPatternRGB(0.5, 0.5, 0.5)
//	if err != nil {
//	    return err
//	}
//	defer gray.Close()
func NewSolidPatternRGB(r, g, b float64) (*SolidPattern, error) {
	ptr := patternCreateRGB(r, g, b)
	st := patternStatus(ptr)

	if st != status.Success {
		return nil, st
	}

	basePattern := newBasePattern(ptr, PatternTypeSolid)
	return &SolidPattern{
		BasePattern: basePattern,
	}, nil
}

// NewSolidPatternRGBA creates a new solid pattern with an RGBA color including transparency.
//
// The color components (r, g, b, a) should be in the range [0.0, 1.0]:
//   - r, g, b: Color channels where 0.0 = no intensity, 1.0 = full intensity
//   - a: Alpha (transparency) where 0.0 = fully transparent, 1.0 = fully opaque
//
// Alpha compositing in Cairo uses premultiplied alpha. This function handles
// the premultiplication internally, so you should provide unpremultiplied values.
//
// Parameters:
//   - r: Red component (0.0 to 1.0)
//   - g: Green component (0.0 to 1.0)
//   - b: Blue component (0.0 to 1.0)
//   - a: Alpha component (0.0 = transparent, 1.0 = opaque)
//
// Returns a new SolidPattern and nil error on success.
// Returns nil and an error if pattern creation fails.
//
// The returned pattern must be closed with Close() when finished to release
// Cairo resources. A finalizer is registered for safety, but explicit cleanup
// is strongly recommended.
//
// Example:
//
//	// Fully opaque red (same as NewSolidPatternRGB)
//	red, err := pattern.NewSolidPatternRGBA(1.0, 0.0, 0.0, 1.0)
//	if err != nil {
//	    return err
//	}
//	defer red.Close()
//
//	// 50% transparent blue
//	blue, err := pattern.NewSolidPatternRGBA(0.0, 0.0, 1.0, 0.5)
//	if err != nil {
//	    return err
//	}
//	defer blue.Close()
//
//	// Fully transparent (invisible)
//	transparent, err := pattern.NewSolidPatternRGBA(1.0, 1.0, 1.0, 0.0)
//	if err != nil {
//	    return err
//	}
//	defer transparent.Close()
func NewSolidPatternRGBA(r, g, b, a float64) (*SolidPattern, error) {
	ptr := patternCreateRGBA(r, g, b, a)
	st := patternStatus(ptr)

	if st != status.Success {
		return nil, st
	}

	basePattern := newBasePattern(ptr, PatternTypeSolid)
	return &SolidPattern{
		BasePattern: basePattern,
	}, nil
}
