// ABOUTME: Defines text and font measurement structs for Cairo's text extents API.
// ABOUTME: TextExtents and FontExtents correspond to cairo_text_extents_t and cairo_font_extents_t.

package font

// TextExtents holds the results of measuring a string of text.
//
// The coordinate system: X grows right, Y grows down. The origin is the
// point passed to [Context.MoveTo] before rendering.
//
// XBearing and YBearing are offsets from the origin to the top-left corner of
// the ink bounding box. YBearing is typically negative because glyphs extend
// above the baseline. Width and Height give the dimensions of the ink bounding box.
// XAdvance and YAdvance give how far the current point moves after rendering.
type TextExtents struct {
	// XBearing is the horizontal distance from the origin to the left edge of
	// the ink bounding box. Positive means the left edge is to the right of
	// the origin.
	XBearing float64

	// YBearing is the vertical distance from the baseline to the top edge of
	// the ink bounding box. Negative values indicate the box extends above
	// the baseline (typical for most glyphs).
	YBearing float64

	// Width is the width of the ink bounding box.
	Width float64

	// Height is the height of the ink bounding box.
	Height float64

	// XAdvance is the distance to advance the current point horizontally
	// after rendering this text.
	XAdvance float64

	// YAdvance is the distance to advance the current point vertically after
	// rendering this text. Typically 0 for horizontal text layouts.
	YAdvance float64
}

// FontExtents holds general metrics of a font face at the current font size.
//
// These metrics are font-wide, not glyph-specific. They describe the overall
// dimensions of the font face at the current size, suitable for line spacing
// and baseline alignment in multi-line text layouts.
type FontExtents struct {
	// Ascent is the distance from the baseline to the top of the tallest glyph
	// in the font, in user-space units.
	Ascent float64

	// Descent is the distance from the baseline to the bottom of the lowest
	// descender in the font, in user-space units. Typically positive even
	// though it represents downward distance.
	Descent float64

	// Height is the recommended vertical distance between baselines for
	// consecutive lines of text. Use this value for consistent multi-line
	// text spacing. Height >= Ascent + Descent.
	Height float64

	// MaxXAdvance is the maximum horizontal advance width of any glyph in
	// the font.
	MaxXAdvance float64

	// MaxYAdvance is the maximum vertical advance of any glyph in the font.
	// Typically 0 for horizontal text layouts.
	MaxYAdvance float64
}
