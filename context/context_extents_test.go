// ABOUTME: Tests for fill/stroke/path extents and point-in-fill/stroke operations.
// ABOUTME: Covers FillExtents, StrokeExtents, PathExtents, InFill, InStroke, and fill rule comparison.

package context

import (
	"testing"

	"github.com/mikowitz/cairo/surface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestContext(t *testing.T) (*Context, func()) {
	t.Helper()
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 200, 200)
	require.NoError(t, err)
	ctx, err := NewContext(surf)
	require.NoError(t, err)
	return ctx, func() {
		_ = ctx.Close()
		_ = surf.Close()
	}
}

// TestContextFillExtents verifies that FillExtents returns the bounding box of the fill area.
func TestContextFillExtents(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()

	ctx.Rectangle(10, 20, 80, 60)
	x1, y1, x2, y2 := ctx.FillExtents()

	assert.InDelta(t, 10.0, x1, 0.001)
	assert.InDelta(t, 20.0, y1, 0.001)
	assert.InDelta(t, 90.0, x2, 0.001)
	assert.InDelta(t, 80.0, y2, 0.001)
}

// TestContextStrokeExtents verifies that StrokeExtents includes stroke width in its bounding box.
func TestContextStrokeExtents(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()

	ctx.SetLineWidth(4.0)
	ctx.Rectangle(10, 20, 80, 60)
	x1, y1, x2, y2 := ctx.StrokeExtents()

	// Stroke extents extend outward by half the line width (2.0)
	assert.InDelta(t, 8.0, x1, 0.5)
	assert.InDelta(t, 18.0, y1, 0.5)
	assert.InDelta(t, 92.0, x2, 0.5)
	assert.InDelta(t, 82.0, y2, 0.5)
}

// TestContextPathExtents verifies that PathExtents returns the path bounding box.
func TestContextPathExtents(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()

	ctx.Rectangle(15, 25, 70, 50)
	x1, y1, x2, y2 := ctx.PathExtents()

	assert.InDelta(t, 15.0, x1, 0.001)
	assert.InDelta(t, 25.0, y1, 0.001)
	assert.InDelta(t, 85.0, x2, 0.001)
	assert.InDelta(t, 75.0, y2, 0.001)
}

// TestContextExtentsAfterClose verifies all extents methods return zeros on a closed context.
func TestContextExtentsAfterClose(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
	require.NoError(t, err)
	defer surf.Close()

	ctx, err := NewContext(surf)
	require.NoError(t, err)
	require.NoError(t, ctx.Close())

	x1, y1, x2, y2 := ctx.FillExtents()
	assert.Equal(t, [4]float64{0, 0, 0, 0}, [4]float64{x1, y1, x2, y2})

	x1, y1, x2, y2 = ctx.StrokeExtents()
	assert.Equal(t, [4]float64{0, 0, 0, 0}, [4]float64{x1, y1, x2, y2})

	x1, y1, x2, y2 = ctx.PathExtents()
	assert.Equal(t, [4]float64{0, 0, 0, 0}, [4]float64{x1, y1, x2, y2})
}

// TestContextInFill verifies point-in-fill detection.
func TestContextInFill(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()

	ctx.Rectangle(20, 20, 100, 100)
	assert.True(t, ctx.InFill(50, 50))
	assert.False(t, ctx.InFill(5, 5))
	assert.False(t, ctx.InFill(150, 150))
}

// TestContextInStroke verifies point-in-stroke detection.
func TestContextInStroke(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()

	ctx.SetLineWidth(10.0)
	ctx.MoveTo(10, 50)
	ctx.LineTo(150, 50)

	assert.True(t, ctx.InStroke(80, 50))
	assert.True(t, ctx.InStroke(80, 54))
	assert.False(t, ctx.InStroke(80, 80))
}

// TestContextInAfterClose verifies InFill and InStroke return false on a closed context.
func TestContextInAfterClose(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
	require.NoError(t, err)
	defer surf.Close()

	ctx, err := NewContext(surf)
	require.NoError(t, err)
	require.NoError(t, ctx.Close())

	assert.False(t, ctx.InFill(50, 50))
	assert.False(t, ctx.InStroke(50, 50))
}

// TestContextFillRuleWindingVsEvenOdd shows that two overlapping sub-paths produce
// different fill results depending on the fill rule.
func TestContextFillRuleWindingVsEvenOdd(t *testing.T) {
	ctx, cleanup := newTestContext(t)
	defer cleanup()

	// Two overlapping rectangles; intersection is at (60,60)-(100,100)
	ctx.Rectangle(20, 20, 80, 80) // (20,20) to (100,100)
	ctx.Rectangle(60, 60, 80, 80) // (60,60) to (140,140)

	// Winding: intersection is inside (winding count = 2, non-zero)
	ctx.SetFillRule(FillRuleWinding)
	assert.True(t, ctx.InFill(80, 80))

	// Even-odd: intersection is outside (crossing count = 2, even)
	ctx.SetFillRule(FillRuleEvenOdd)
	assert.False(t, ctx.InFill(80, 80))
}
