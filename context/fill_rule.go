// ABOUTME: Defines the FillRule type and constants for controlling path fill behavior.
// ABOUTME: Provides SetFillRule and GetFillRule methods for the winding and even-odd rules.

package context

// FillRule controls which areas of a self-intersecting path are considered "inside"
// and therefore filled when calling [Context.Fill] or [Context.FillPreserve].
//
// The winding and even-odd rules produce different results for complex paths
// that intersect themselves, such as stars or compound shapes.
//
// # Winding Rule
//
// The winding rule considers a point "inside" a path if the path winds around
// it a non-zero number of times. This is determined by counting the number of
// times the path crosses a ray from the point in a given direction, with crossings
// in one direction counted as positive and the other as negative. If the sum is
// non-zero, the point is inside.
//
// The winding rule fills the interior of most common shapes naturally and is
// the default fill rule.
//
// # Even-Odd Rule
//
// The even-odd rule considers a point "inside" a path if a ray from the point
// crosses the path an odd number of times. This can produce "holes" in
// self-intersecting shapes (e.g., a five-pointed star drawn with a single path
// will have its center left unfilled).
//
//go:generate stringer -type=FillRule
type FillRule int

// The iota values below must match Cairo's cairo_fill_rule_t C enum exactly.
// Cairo has maintained this ordering since its initial release and documents
// it as stable. The CGO layer casts FillRule directly to cairo_fill_rule_t,
// so any divergence would silently produce incorrect fill behavior.
const (
	// FillRuleWinding uses the winding number rule. A point is considered inside
	// the path if the path winds around it a non-zero number of times.
	// This is the default fill rule.
	FillRuleWinding FillRule = iota

	// FillRuleEvenOdd uses the even-odd rule. A point is considered inside
	// the path if a ray from the point crosses the path an odd number of times.
	// This can produce holes in self-intersecting paths.
	FillRuleEvenOdd
)

// SetFillRule sets the fill rule for the context, which controls how the interior
// of self-intersecting paths is determined.
//
// The fill rule is used by [Context.Fill], [Context.FillPreserve], [Context.Clip],
// and [Context.ClipPreserve].
//
// The default fill rule is [FillRuleWinding].
func (c *Context) SetFillRule(fillRule FillRule) {
	c.withLock(func() {
		contextSetFillRule(c.ptr, fillRule)
	})
}

// GetFillRule returns the current fill rule for the context.
//
// The default fill rule is [FillRuleWinding].
//
// If called on a closed context, this method returns [FillRuleWinding].
func (c *Context) GetFillRule() FillRule {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return FillRuleWinding
	}
	return contextGetFillRule(c.ptr)
}
