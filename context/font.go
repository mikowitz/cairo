// ABOUTME: Implements the toy font API methods on Context.
// ABOUTME: Provides SelectFontFace, SetFontSize, ShowText, and TextPath.

package context

import "github.com/mikowitz/cairo/font"

// SelectFontFace selects a font face for the context using a font family name,
// slant, and weight. This is part of Cairo's toy font API, which provides a
// simple interface for basic text rendering without full font management.
//
// The family parameter specifies the font family (e.g., "serif", "sans-serif",
// "monospace") using the platform's font system. Results are platform-dependent.
//
// Text is positioned at the current point. Call [Context.MoveTo] before
// drawing text to establish the baseline origin.
//
// Example:
//
//	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightBold)
func (c *Context) SelectFontFace(family string, slant font.Slant, weight font.Weight) {
	c.withLock(func() {
		contextSelectFontFace(c.ptr, family, slant, weight)
	})
}

// SetFontSize sets the current font size for text rendering, specified in
// user-space units. The default font size is 10.
//
// Text is positioned at the current point. Call [Context.MoveTo] before
// drawing text to establish the baseline origin.
//
// Example:
//
//	ctx.SetFontSize(14.0)
func (c *Context) SetFontSize(size float64) {
	c.withLock(func() {
		contextSetFontSize(c.ptr, size)
	})
}

// ShowText renders the given UTF-8 string at the current point using the
// current font face and size. After rendering, the current point is advanced
// to the end of the text.
//
// ShowText is part of Cairo's toy font API. For advanced text rendering,
// use Pango integration (not yet available in this library).
//
// Example:
//
//	ctx.MoveTo(10, 50)
//	ctx.ShowText("Hello, world!")
func (c *Context) ShowText(text string) {
	c.withLock(func() {
		contextShowText(c.ptr, text)
	})
}

// TextPath appends the outline of the given UTF-8 string to the current path.
// The outlines are positioned at the current point. After the call, the current
// point is at the end of the text.
//
// Unlike [Context.ShowText], TextPath does not render immediately; the caller
// must call [Context.Fill] or [Context.Stroke] to paint the text outlines.
//
// Example:
//
//	ctx.MoveTo(10, 50)
//	ctx.TextPath("Outlined text")
//	ctx.Fill()
func (c *Context) TextPath(text string) {
	c.withLock(func() {
		contextTextPath(c.ptr, text)
	})
}
