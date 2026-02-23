// ABOUTME: TextExtents and FontExtents methods on Context for text measurement.
// ABOUTME: These return metrics for text layout, alignment, and bounding box drawing.

package context

import "github.com/mikowitz/cairo/font"

// TextExtents measures the rendered extents of the given UTF-8 text string
// using the current font face and size. The returned [font.TextExtents] describe
// the ink bounding box and advance vector for the text, in user-space coordinates.
//
// TextExtents does not draw anything; it only measures.
//
// Common uses:
//   - Center text: move to x - extents.XBearing - extents.Width/2
//   - Right-align: move to x - extents.XAdvance
//   - Draw bounding box: use XBearing, YBearing, Width, Height relative to the current point
//
// If the context has been closed, TextExtents returns a zero-value [font.TextExtents].
func (c *Context) TextExtents(text string) *font.TextExtents {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return &font.TextExtents{}
	}
	return contextTextExtents(c.ptr, text)
}

// FontExtents returns the metrics of the current font face at the current font size.
// The returned [font.FontExtents] describe font-wide dimensions used for line spacing
// and baseline alignment in multi-line text layouts.
//
// Common uses:
//   - Line height for multi-line text: extents.Height
//   - Space above baseline: extents.Ascent
//   - Space below baseline: extents.Descent
//
// If the context has been closed, FontExtents returns a zero-value [font.FontExtents].
func (c *Context) FontExtents() *font.FontExtents {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return &font.FontExtents{}
	}
	return contextFontExtents(c.ptr)
}
