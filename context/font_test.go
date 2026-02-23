// ABOUTME: Tests for the toy font API methods on Context.
// ABOUTME: Covers SelectFontFace, SetFontSize, ShowText, and TextPath.

package context

import (
	"testing"

	"github.com/mikowitz/cairo/font"
	"github.com/mikowitz/cairo/status"
	"github.com/stretchr/testify/assert"
)

// TestContextSelectFontFace tests that SelectFontFace sets the font without error.
func TestContextSelectFontFace(t *testing.T) {
	ctx := newTestContext(t, 100, 100)

	ctx.SelectFontFace("serif", font.SlantNormal, font.WeightNormal)
	assert.Equal(t, status.Success, ctx.Status())

	ctx.SelectFontFace("sans-serif", font.SlantItalic, font.WeightBold)
	assert.Equal(t, status.Success, ctx.Status())

	ctx.SelectFontFace("monospace", font.SlantOblique, font.WeightNormal)
	assert.Equal(t, status.Success, ctx.Status())
}

// TestContextSelectFontFaceClosedContext tests that SelectFontFace on a closed context is a no-op.
func TestContextSelectFontFaceClosedContext(t *testing.T) {
	ctx := newTestContext(t, 100, 100)
	_ = ctx.Close()

	// Should not panic on closed context
	ctx.SelectFontFace("serif", font.SlantNormal, font.WeightNormal)
}

// TestContextSetFontSize tests that SetFontSize sets the font size without error.
func TestContextSetFontSize(t *testing.T) {
	ctx := newTestContext(t, 100, 100)

	ctx.SetFontSize(12.0)
	assert.Equal(t, status.Success, ctx.Status())

	ctx.SetFontSize(24.0)
	assert.Equal(t, status.Success, ctx.Status())
}

// TestContextSetFontSizeClosedContext tests that SetFontSize on a closed context is a no-op.
func TestContextSetFontSizeClosedContext(t *testing.T) {
	ctx := newTestContext(t, 100, 100)
	_ = ctx.Close()

	// Should not panic on closed context
	ctx.SetFontSize(12.0)
}

// TestContextShowText tests that ShowText renders text and advances the current point.
func TestContextShowText(t *testing.T) {
	ctx := newTestContext(t, 200, 100)

	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightNormal)
	ctx.SetFontSize(12.0)
	ctx.MoveTo(10, 50)
	ctx.ShowText("Hello, Cairo!")

	assert.Equal(t, status.Success, ctx.Status())
	// ShowText advances the current point to end of the text
	assert.True(t, ctx.HasCurrentPoint())
}

// TestContextShowTextClosedContext tests that ShowText on a closed context is a no-op.
func TestContextShowTextClosedContext(t *testing.T) {
	ctx := newTestContext(t, 100, 100)
	_ = ctx.Close()

	// Should not panic on closed context
	ctx.ShowText("Hello")
}

// TestContextTextPath tests that TextPath creates a path from text without error.
func TestContextTextPath(t *testing.T) {
	ctx := newTestContext(t, 200, 100)

	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightNormal)
	ctx.SetFontSize(12.0)
	ctx.MoveTo(10, 50)
	ctx.TextPath("Hello, Cairo!")

	assert.Equal(t, status.Success, ctx.Status())
	// TextPath leaves the current point at the end of the text
	assert.True(t, ctx.HasCurrentPoint())
}

// TestContextTextPathClosedContext tests that TextPath on a closed context is a no-op.
func TestContextTextPathClosedContext(t *testing.T) {
	ctx := newTestContext(t, 100, 100)
	_ = ctx.Close()

	// Should not panic on closed context
	ctx.TextPath("Hello")
}
