package context

import "github.com/mikowitz/cairo/status"

// LineCap specifies how the endpoints of lines are rendered when stroking.
//
// The line cap style only affects the endpoints of lines. The appearance of
// line joins is controlled by [LineJoin].
//
//go:generate stringer -type=LineCap
type LineCap int

const (
	// LineCapButt starts and stops the line exactly at the start and end points.
	LineCapButt LineCap = iota

	// LineCapRound uses a round ending, with the center of the circle at the end point.
	LineCapRound

	// LineCapSquare uses a squared ending, with the center of the square at the end point.
	LineCapSquare
)

// GetLineCap gets the current line cap style, as set by [Context.SetLineCap].
//
// The default line cap style is [LineCapButt].
func (c *Context) GetLineCap() LineCap {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return LineCapButt
	}
	return contextGetLineCap(c.ptr)
}

// SetLineCap sets the line cap style to be used when stroking lines.
//
// The line cap style specifies how the endpoints of stroked lines are drawn.
// The default line cap style is [LineCapButt].
//
// As with other stroke parameters, the current line cap style is examined by
// [Context.Stroke] and [Context.StrokePreserve], but does not have any effect
// during path construction.
func (c *Context) SetLineCap(lineCap LineCap) {
	c.withLock(func() {
		contextSetLineCap(c.ptr, lineCap)
	})
}

// LineJoin specifies how the junctions between line segments are rendered when stroking.
//
// The line join style only affects the junctions between line segments. The appearance
// of line endpoints is controlled by [LineCap].
//
//go:generate stringer -type=LineJoin
type LineJoin int

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

// GetLineJoin gets the current line join style, as set by [Context.SetLineJoin].
//
// The default line join style is [LineJoinMiter].
func (c *Context) GetLineJoin() LineJoin {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return LineJoinMiter
	}
	return contextGetLineJoin(c.ptr)
}

// SetLineJoin sets the line join style to be used when stroking lines.
//
// The line join style specifies how the junctions between two line segments
// are drawn. The default line join style is [LineJoinMiter].
//
// As with other stroke parameters, the current line join style is examined by
// [Context.Stroke] and [Context.StrokePreserve], but does not have any effect
// during path construction.
func (c *Context) SetLineJoin(lineJoin LineJoin) {
	c.withLock(func() {
		contextSetLineJoin(c.ptr, lineJoin)
	})
}

// GetMiterLimit gets the current miter limit, as set by [Context.SetMiterLimit].
//
// The default miter limit is 10.0.
func (c *Context) GetMiterLimit() float64 {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return 10
	}

	return contextGetMiterLimit(c.ptr)
}

// SetMiterLimit sets the miter limit for [LineJoinMiter] line joins.
//
// When two line segments meet at a sharp angle and [LineJoinMiter] is used,
// it is possible for the miter to extend far beyond the thickness of the line
// stroking the path. The miter limit is used to decide when to fall back to
// drawing a [LineJoinBevel] instead.
//
// The miter limit is a ratio of the miter length to the line width. Cairo divides
// the length of the miter by the line width. If the result is greater than the
// miter limit, the style is converted to a bevel.
//
// The default miter limit is 10.0, which will convert miters with interior angles
// less than 11 degrees to bevels. A miter limit of 2.0 will restrict miters to
// angles larger than 60 degrees. A miter limit of 1.414 will restrict miters to
// angles larger than 90 degrees.
//
// The relationship between miter limit and angle can be expressed as:
// miter_limit = 1/sin(angle/2)
func (c *Context) SetMiterLimit(limit float64) {
	c.withLock(func() {
		contextSetMiterLimit(c.ptr, limit)
	})
}

// SetDash sets the dash pattern to be used when stroking lines.
//
// A dash pattern is specified by an array of positive values. Each value provides
// the length of alternate "on" and "off" portions of the stroke. The offset specifies
// an offset into the pattern at which the stroke begins.
//
// Each "on" segment will have caps applied as if the segment were a separate sub-path.
// In particular, it is valid to use an "on" length of 0.0 with [LineCapRound] or
// [LineCapSquare] in order to create dots or squares along a path.
//
// If the array is empty or nil, dashing is disabled and the line will be drawn solid.
// If the array has a single element, the "on" and "off" lengths are both set to that value.
//
// The offset value is measured in the same units as the pattern values and indicates
// how far into the dash pattern to start the stroke.
//
// The values in the dashes array are measured in user-space units as evaluated at
// stroke time. Therefore, changing the transformation matrix will change the rendered
// dash pattern.
//
// SetDash returns an error if any value in the dashes array is negative or if all
// values are zero.
//
// Example:
//
//	// Create a simple dashed line: 10 units on, 5 units off
//	err := ctx.SetDash([]float64{10.0, 5.0}, 0.0)
//
//	// Create a dot pattern with round caps
//	ctx.SetLineCap(LineCapRound)
//	err = ctx.SetDash([]float64{0.0, 5.0}, 0.0)
//
//	// Disable dashing
//	err = ctx.SetDash(nil, 0.0)
func (c *Context) SetDash(dashes []float64, offset float64) error {
	c.Lock()
	defer c.Unlock()

	if c.ptr == nil {
		return status.NullPointer
	}

	if len(dashes) > 0 && (allZeroes(dashes) || anyNegative(dashes)) {
		return status.InvalidDash
	}

	return contextSetDash(c.ptr, dashes, offset).ToError()
}

func anyNegative(s []float64) bool {
	for _, f := range s {
		if f < 0 {
			return true
		}
	}
	return false
}

func allZeroes(s []float64) bool {
	for _, f := range s {
		if f != 0 {
			return false
		}
	}
	return true
}

// GetDashCount gets the number of dashes in the current dash pattern.
//
// This can be used to determine the size of the array returned by [Context.GetDash].
// If dashing is not currently in effect, this returns 0.
func (c *Context) GetDashCount() int {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return 0
	}
	return contextGetDashCount(c.ptr)
}

// GetDash gets the current dash pattern and offset, as set by [Context.SetDash].
//
// Returns the dash array, the current dash offset, and an error if the context
// is invalid. If dashing is not currently enabled, the returned array will be empty.
//
// Example:
//
//	ctx.SetDash([]float64{10.0, 5.0}, 2.0)
//	dashes, offset, err := ctx.GetDash()
//	if err == nil {
//	    fmt.Printf("Dashes: %v, Offset: %f\n", dashes, offset)
//	    // Output: Dashes: [10 5], Offset: 2.000000
//	}
func (c *Context) GetDash() ([]float64, float64, error) {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return []float64{}, 0, status.NullPointer
	}
	dashes, offset := contextGetDash(c.ptr)
	return dashes, offset, nil
}
