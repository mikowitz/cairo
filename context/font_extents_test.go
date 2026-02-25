// ABOUTME: Tests for TextExtents and FontExtents methods on Context.
// ABOUTME: Covers text measurement, empty strings, closed context, and font comparison.

package context

import (
	"testing"

	"github.com/mikowitz/cairo/font"
	"github.com/mikowitz/cairo/surface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestContextTextExtents verifies that TextExtents returns non-zero measurements for a text string.
func TestContextTextExtents(t *testing.T) {
	ctx := newTestContext(t, 200, 100)
	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightNormal)
	ctx.SetFontSize(12.0)

	extents := ctx.TextExtents("Hello")

	assert.Greater(t, extents.Width, 0.0)
	assert.Greater(t, extents.Height, 0.0)
	assert.Greater(t, extents.XAdvance, 0.0)
}

// TestContextFontExtents verifies that FontExtents returns valid font metrics.
func TestContextFontExtents(t *testing.T) {
	ctx := newTestContext(t, 200, 100)
	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightNormal)
	ctx.SetFontSize(12.0)

	extents := ctx.FontExtents()

	assert.Greater(t, extents.Ascent, 0.0)
	assert.Greater(t, extents.Descent, 0.0)
	assert.Greater(t, extents.Height, 0.0)
	// Height must be at least Ascent + Descent
	assert.GreaterOrEqual(t, extents.Height, extents.Ascent+extents.Descent-0.001)
}

// TestContextTextExtentsEmpty verifies that an empty string yields all-zero extents.
func TestContextTextExtentsEmpty(t *testing.T) {
	ctx := newTestContext(t, 200, 100)
	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightNormal)
	ctx.SetFontSize(12.0)

	extents := ctx.TextExtents("")

	assert.Equal(t, 0.0, extents.XBearing)
	assert.Equal(t, 0.0, extents.YBearing)
	assert.Equal(t, 0.0, extents.Width)
	assert.Equal(t, 0.0, extents.Height)
	assert.Equal(t, 0.0, extents.XAdvance)
	assert.Equal(t, 0.0, extents.YAdvance)
}

// TestContextExtentsWithDifferentFonts verifies that different font faces give different metrics.
func TestContextExtentsWithDifferentFonts(t *testing.T) {
	ctx := newTestContext(t, 200, 100)
	ctx.SetFontSize(12.0)

	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightNormal)
	normalExtents := ctx.TextExtents("Hello")

	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightBold)
	boldExtents := ctx.TextExtents("Hello")

	// Bold text is wider or equal to normal text
	assert.GreaterOrEqual(t, boldExtents.Width, normalExtents.Width)
}

// TestContextTextExtentsClosedContext verifies TextExtents returns zero-value on a closed context.
func TestContextTextExtentsClosedContext(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
	require.NoError(t, err)
	defer surf.Close()

	ctx, err := NewContext(surf)
	require.NoError(t, err)
	require.NoError(t, ctx.Close())

	extents := ctx.TextExtents("Hello")
	assert.Equal(t, font.TextExtents{}, extents)
}

// TestContextFontExtentsClosedContext verifies FontExtents returns zero-value on a closed context.
func TestContextFontExtentsClosedContext(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
	require.NoError(t, err)
	defer surf.Close()

	ctx, err := NewContext(surf)
	require.NoError(t, err)
	require.NoError(t, ctx.Close())

	extents := ctx.FontExtents()
	assert.Equal(t, font.FontExtents{}, extents)
}
