package pattern

import "github.com/mikowitz/cairo/status"

// Gradient defines the interface for gradient patterns in Cairo.
// Gradients support smooth color transitions using color stops, which are points
// along the gradient where a specific color is defined. Cairo interpolates colors
// between stops to create smooth transitions.
//
// There are two types of gradients:
//   - Linear gradients: transitions along a line between two points
//   - Radial gradients: transitions between two circles
//
// Color stops are defined using offset values in the range [0.0, 1.0], where 0.0
// represents the start of the gradient and 1.0 represents the end.
type Gradient interface {
	AddColorStopRGB(offset, r, g, b float64)
	AddColorStopRGBA(offset, r, g, b, a float64)
	GetColorStopCount() (int, error)
	GetColorStopRGBA(index int) (float64, float64, float64, float64, float64, error)
}

// BaseGradient provides the common implementation for all gradient pattern types.
// It embeds BasePattern and adds gradient-specific functionality for managing
// color stops. This type should not be used directly; instead use LinearGradient
// or RadialGradient.
type BaseGradient struct {
	*BasePattern
}

// AddColorStopRGB adds an opaque color stop to a gradient pattern.
// The offset specifies the location along the gradient's control vector as a value
// in the range [0.0, 1.0], where 0.0 is the start and 1.0 is the end.
//
// Parameters:
//   - offset: location along gradient (0.0 to 1.0)
//   - r: red component (0.0 to 1.0)
//   - g: green component (0.0 to 1.0)
//   - b: blue component (0.0 to 1.0)
//
// The color stop is added with full opacity (alpha = 1.0). If two or more stops
// are specified with identical offset values, they will be sorted according to
// the order in which the stops are added.
//
// Example:
//
//	gradient.AddColorStopRGB(0.0, 1.0, 0.0, 0.0)  // Red at start
//	gradient.AddColorStopRGB(1.0, 0.0, 0.0, 1.0)  // Blue at end
func (bg *BaseGradient) AddColorStopRGB(offset, r, g, b float64) {
	bg.Lock()
	defer bg.Unlock()

	patternAddColorStopRGB(bg.ptr, offset, r, g, b)
}

// AddColorStopRGBA adds a translucent color stop to a gradient pattern.
// This is identical to AddColorStopRGB but includes an alpha (transparency) component.
//
// Parameters:
//   - offset: location along gradient (0.0 to 1.0)
//   - r: red component (0.0 to 1.0)
//   - g: green component (0.0 to 1.0)
//   - b: blue component (0.0 to 1.0)
//   - a: alpha/opacity component (0.0 = transparent, 1.0 = opaque)
//
// The alpha value allows for gradients that fade in or out. Color stops with
// identical offsets are sorted by insertion order.
//
// Example:
//
//	gradient.AddColorStopRGBA(0.0, 1.0, 0.0, 0.0, 1.0)  // Opaque red at start
//	gradient.AddColorStopRGBA(1.0, 0.0, 0.0, 1.0, 0.0)  // Transparent blue at end
func (bg *BaseGradient) AddColorStopRGBA(offset, r, g, b, a float64) {
	bg.Lock()
	defer bg.Unlock()

	patternAddColorStopRGBA(bg.ptr, offset, r, g, b, a)
}

// GetColorStopCount returns the number of color stops defined in the gradient pattern.
// This is useful for iterating through color stops using GetColorStopRGBA.
//
// Returns the count and nil on success, or an error if the pattern is not a gradient
// (status.PatternTypeMismatch).
//
// Example:
//
//	count, err := gradient.GetColorStopCount()
//	if err != nil {
//	    return err
//	}
//	for i := 0; i < count; i++ {
//	    offset, r, g, b, a, _ := gradient.GetColorStopRGBA(i)
//	    // Process color stop...
//	}
func (bg *BaseGradient) GetColorStopCount() (int, error) {
	count, st := patternGetColorStopCount(bg.ptr)

	if st != status.Success {
		return count, st
	}
	return count, nil
}

// GetColorStopRGBA retrieves the color and offset information for a color stop at
// the specified index. Indices range from 0 to n-1, where n is the value returned
// by GetColorStopCount.
//
// Parameters:
//   - index: the color stop index (0-based)
//
// Returns:
//   - offset: location along gradient (0.0 to 1.0)
//   - r: red component (0.0 to 1.0)
//   - g: green component (0.0 to 1.0)
//   - b: blue component (0.0 to 1.0)
//   - a: alpha component (0.0 to 1.0)
//   - err: nil on success, or an error for invalid index or wrong pattern type
//
// The returned color values are not premultiplied by alpha. If the index is out
// of bounds, an error is returned.
//
// Example:
//
//	offset, r, g, b, a, err := gradient.GetColorStopRGBA(0)
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Stop at %.2f: rgba(%.2f, %.2f, %.2f, %.2f)\n", offset, r, g, b, a)
func (bg *BaseGradient) GetColorStopRGBA(index int) (float64, float64, float64, float64, float64, error) {
	o, r, g, b, a, st := patternGetColorStopRGBA(bg.ptr, index)

	if st != status.Success {
		return o, r, g, b, a, st
	}

	return o, r, g, b, a, nil
}

// LinearGradient represents a gradient pattern that transitions colors along a line.
// Linear gradients are defined by two points (x0, y0) and (x1, y1), with colors
// interpolated along the line connecting them.
//
// The gradient coordinates are in pattern space. For a new pattern, pattern space
// is identical to user space, but the relationship can be changed using SetMatrix.
//
// Linear gradients extend infinitely perpendicular to the gradient line. Areas
// before the start point take the color of the first stop, and areas after the
// end point take the color of the last stop (by default).
type LinearGradient struct {
	*BaseGradient
}

// NewLinearGradient creates a new linear gradient pattern along the line defined
// by the points (x0, y0) and (x1, y1).
//
// Parameters:
//   - x0, y0: starting point coordinates
//   - x1, y1: ending point coordinates
//
// After creation, color stops should be added using AddColorStopRGB or
// AddColorStopRGBA to define the gradient's colors.
//
// The coordinates are in pattern space. For a new pattern, pattern space is
// identical to user space, but this relationship can be modified using SetMatrix.
//
// Returns the new LinearGradient or an error if creation fails.
//
// Example:
//
//	// Create horizontal gradient from left to right
//	gradient, err := NewLinearGradient(0, 0, 200, 0)
//	if err != nil {
//	    return err
//	}
//	defer gradient.Close()
//
//	gradient.AddColorStopRGB(0.0, 1.0, 0.0, 0.0)  // Red at left
//	gradient.AddColorStopRGB(1.0, 0.0, 0.0, 1.0)  // Blue at right
//	ctx.SetSource(gradient)
func NewLinearGradient(x0, y0, x1, y1 float64) (*LinearGradient, error) {
	ptr := patternCreateLinear(x0, y0, x1, y1)
	st := patternStatus(ptr)

	if st != status.Success {
		return nil, st
	}

	basePattern := newBasePattern(ptr, PatternTypeLinear)
	baseGradient := &BaseGradient{BasePattern: basePattern}
	return &LinearGradient{
		BaseGradient: baseGradient,
	}, nil
}

// RadialGradient represents a gradient pattern that transitions colors between two circles.
// Radial gradients are defined by two circles: a start circle (cx0, cy0, radius0) and
// an end circle (cx1, cy1, radius1). Colors are interpolated along the cone formed by
// the two circles.
//
// The gradient coordinates are in pattern space. For a new pattern, pattern space
// is identical to user space, but the relationship can be changed using SetMatrix.
//
// Radial gradients can create various effects:
//   - Concentric circles (same center, different radii)
//   - Offset highlights (different centers, simulating light sources)
//   - Expanding/contracting effects
type RadialGradient struct {
	*BaseGradient
}

// NewRadialGradient creates a new radial gradient pattern between two circles.
// The gradient interpolates colors from the start circle to the end circle.
//
// Parameters:
//   - cx0, cy0: center coordinates of the start circle
//   - radius0: radius of the start circle
//   - cx1, cy1: center coordinates of the end circle
//   - radius1: radius of the end circle
//
// After creation, color stops should be added using AddColorStopRGB or
// AddColorStopRGBA to define the gradient's colors.
//
// The coordinates are in pattern space. For a new pattern, pattern space is
// identical to user space, but this relationship can be modified using SetMatrix.
//
// Returns the new RadialGradient or an error if creation fails.
//
// Example:
//
//	// Create radial gradient from center outward
//	gradient, err := NewRadialGradient(100, 100, 10, 100, 100, 100)
//	if err != nil {
//	    return err
//	}
//	defer gradient.Close()
//
//	gradient.AddColorStopRGB(0.0, 1.0, 1.0, 1.0)  // White at center
//	gradient.AddColorStopRGB(1.0, 0.0, 0.0, 1.0)  // Blue at edge
//	ctx.SetSource(gradient)
func NewRadialGradient(cx0, cy0, radius0, cx1, cy1, radius1 float64) (*RadialGradient, error) {
	ptr := patternCreateRadial(cx0, cy0, radius0, cx1, cy1, radius1)
	st := patternStatus(ptr)

	if st != status.Success {
		return nil, st
	}

	basePattern := newBasePattern(ptr, PatternTypeRadial)
	baseGradient := &BaseGradient{BasePattern: basePattern}
	return &RadialGradient{
		BaseGradient: baseGradient,
	}, nil
}
