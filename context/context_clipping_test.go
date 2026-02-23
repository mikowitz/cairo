package context

import (
	"math"
	"testing"

	"github.com/mikowitz/cairo/status"
	"github.com/mikowitz/cairo/surface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestContextClip verifies that clipping restricts drawing to the clip region.
func TestContextClip(t *testing.T) {
	ctx, _ := newTestContext(t, 200, 200)

	t.Run("BasicClip", func(t *testing.T) {
		// Create a rectangular clip region
		ctx.Rectangle(50, 50, 100, 100)
		ctx.Clip()

		// Verify status is success
		assert.Equal(t, status.Success, ctx.Status())

		// Path should be consumed after Clip
		assert.False(t, ctx.HasCurrentPoint(), "Path should be consumed after Clip")
	})

	t.Run("ClipConsumesPath", func(t *testing.T) {
		// Create a path
		ctx.NewPath()
		ctx.Rectangle(10, 10, 50, 50)
		assert.True(t, ctx.HasCurrentPoint(), "Should have current point before clip")

		// Clip consumes the path
		ctx.Clip()
		assert.False(t, ctx.HasCurrentPoint(), "Path should be consumed after Clip")
	})

	t.Run("ClipWithNoPath", func(t *testing.T) {
		// Clipping with no path should not error
		ctx.ResetClip()
		ctx.NewPath()
		ctx.Clip()
		assert.Equal(t, status.Success, ctx.Status())
	})
}

// TestContextClipPreserve verifies that ClipPreserve preserves the path.
func TestContextClipPreserve(t *testing.T) {
	ctx, _ := newTestContext(t, 200, 200)

	t.Run("PreservesPath", func(t *testing.T) {
		// Create a path
		ctx.Rectangle(50, 50, 100, 100)
		assert.True(t, ctx.HasCurrentPoint())

		// ClipPreserve should keep the path
		ctx.ClipPreserve()

		// Path should still exist
		assert.True(t, ctx.HasCurrentPoint(), "Path should be preserved after ClipPreserve")
		assert.Equal(t, status.Success, ctx.Status())

		// We should be able to stroke the preserved path
		ctx.Stroke()
		assert.Equal(t, status.Success, ctx.Status())
	})

	t.Run("CanFillAfterClipPreserve", func(t *testing.T) {
		ctx.NewPath()
		ctx.Rectangle(20, 20, 60, 60)

		// ClipPreserve and then fill
		ctx.ClipPreserve()
		ctx.Fill()

		assert.Equal(t, status.Success, ctx.Status())
	})
}

// TestContextClipExtents verifies extents calculation for clip region.
func TestContextClipExtents(t *testing.T) {
	ctx, _ := newTestContext(t, 400, 400)

	t.Run("InitialExtents", func(t *testing.T) {
		// Initial clip extents should be entire surface
		x1, y1, x2, y2 := ctx.ClipExtents()

		assert.Equal(t, 0.0, x1, "Initial clip x1 should be 0")
		assert.Equal(t, 0.0, y1, "Initial clip y1 should be 0")
		assert.Equal(t, 400.0, x2, "Initial clip x2 should be surface width")
		assert.Equal(t, 400.0, y2, "Initial clip y2 should be surface height")
	})

	t.Run("ExtentsAfterClip", func(t *testing.T) {
		ctx.ResetClip()
		// Set a clip region
		ctx.Rectangle(50, 60, 200, 150)
		ctx.Clip()

		// Get clip extents
		x1, y1, x2, y2 := ctx.ClipExtents()

		assert.Equal(t, 50.0, x1, "Clip x1 should be 50")
		assert.Equal(t, 60.0, y1, "Clip y1 should be 60")
		assert.Equal(t, 250.0, x2, "Clip x2 should be 50+200")
		assert.Equal(t, 210.0, y2, "Clip y2 should be 60+150")
	})

	t.Run("ExtentsWithCircularClip", func(t *testing.T) {
		ctx.ResetClip()
		ctx.NewPath()
		ctx.Arc(100, 100, 50, 0, 2*math.Pi)
		ctx.Clip()

		// Extents should be bounding box of circle
		x1, y1, x2, y2 := ctx.ClipExtents()

		// Should be approximately 50-150 for both x and y
		assert.InDelta(t, 50.0, x1, 1.0, "Circle clip x1")
		assert.InDelta(t, 50.0, y1, 1.0, "Circle clip y1")
		assert.InDelta(t, 150.0, x2, 1.0, "Circle clip x2")
		assert.InDelta(t, 150.0, y2, 1.0, "Circle clip y2")
	})
}

// TestContextInClip verifies point-in-clip testing.
func TestContextInClip(t *testing.T) {
	ctx, _ := newTestContext(t, 200, 200)

	t.Run("AllPointsInClipInitially", func(t *testing.T) {
		// Before any clipping, all points should be in clip
		assert.True(t, ctx.InClip(50, 50), "Point should be in initial clip")
		assert.True(t, ctx.InClip(0, 0), "Point should be in initial clip")
		assert.True(t, ctx.InClip(199, 199), "Point should be in initial clip")
	})

	t.Run("PointsInRectangularClip", func(t *testing.T) {
		// Set rectangular clip region: 50,50 to 150,150
		ctx.Rectangle(50, 50, 100, 100)
		ctx.Clip()

		// Points inside clip
		assert.True(t, ctx.InClip(75, 75), "Point inside clip should return true")
		assert.True(t, ctx.InClip(50, 50), "Corner point should be in clip")
		assert.True(t, ctx.InClip(149, 149), "Edge point should be in clip")

		// Points outside clip
		assert.False(t, ctx.InClip(25, 25), "Point outside clip should return false")
		assert.False(t, ctx.InClip(175, 175), "Point outside clip should return false")
		assert.False(t, ctx.InClip(0, 100), "Point outside clip should return false")
	})

	t.Run("PointsInCircularClip", func(t *testing.T) {
		ctx.NewPath()
		// Circle centered at (100, 100) with radius 50
		ctx.Arc(100, 100, 50, 0, 2*math.Pi)
		ctx.Clip()

		// Point at center
		assert.True(t, ctx.InClip(100, 100), "Center point should be in circular clip")

		// Point inside circle
		assert.True(t, ctx.InClip(110, 100), "Point inside circle should be in clip")

		// Point outside circle
		assert.False(t, ctx.InClip(160, 100), "Point outside circle should not be in clip")
		assert.False(t, ctx.InClip(50, 50), "Point outside circle should not be in clip")
	})
}

// TestContextResetClip verifies clip region clearing.
func TestContextResetClip(t *testing.T) {
	ctx, _ := newTestContext(t, 200, 200)

	t.Run("ResetAfterClip", func(t *testing.T) {
		// Set a clip region
		ctx.Rectangle(50, 50, 50, 50)
		ctx.Clip()

		// Verify clipping is active
		assert.False(t, ctx.InClip(25, 25), "Point should be outside clip")
		assert.True(t, ctx.InClip(75, 75), "Point should be inside clip")

		// Reset clip
		ctx.ResetClip()

		// All points should be in clip again
		assert.True(t, ctx.InClip(25, 25), "Point should be in clip after reset")
		assert.True(t, ctx.InClip(75, 75), "Point should be in clip after reset")
		assert.True(t, ctx.InClip(0, 0), "Point should be in clip after reset")
		assert.True(t, ctx.InClip(199, 199), "Point should be in clip after reset")
	})

	t.Run("ResetWithoutClip", func(t *testing.T) {
		// Resetting without any clip should not error
		ctx.ResetClip()
		assert.Equal(t, status.Success, ctx.Status())
	})

	t.Run("ExtentsAfterReset", func(t *testing.T) {
		// Set a clip
		ctx.Rectangle(50, 50, 50, 50)
		ctx.Clip()

		// Reset
		ctx.ResetClip()

		// Extents should be full surface again
		x1, y1, x2, y2 := ctx.ClipExtents()
		assert.Equal(t, 0.0, x1)
		assert.Equal(t, 0.0, y1)
		assert.Equal(t, 200.0, x2)
		assert.Equal(t, 200.0, y2)
	})
}

// TestContextNestedClips verifies clip intersection behavior.
func TestContextNestedClips(t *testing.T) {
	ctx, _ := newTestContext(t, 200, 200)

	t.Run("IntersectingClips", func(t *testing.T) {
		// First clip: 0,0 to 150,150
		ctx.Rectangle(0, 0, 150, 150)
		ctx.Clip()

		// Second clip: 50,50 to 200,200
		ctx.Rectangle(50, 50, 150, 150)
		ctx.Clip()

		// Effective clip should be intersection: 50,50 to 150,150
		assert.True(t, ctx.InClip(100, 100), "Point in intersection should be in clip")
		assert.False(t, ctx.InClip(25, 25), "Point outside intersection should not be in clip")
		assert.False(t, ctx.InClip(175, 175), "Point outside intersection should not be in clip")

		// Check extents of intersection
		x1, y1, x2, y2 := ctx.ClipExtents()
		assert.Equal(t, 50.0, x1, "Intersection x1")
		assert.Equal(t, 50.0, y1, "Intersection y1")
		assert.Equal(t, 150.0, x2, "Intersection x2")
		assert.Equal(t, 150.0, y2, "Intersection y2")
	})

	t.Run("ThreeNestedClips", func(t *testing.T) {
		ctx.NewPath()

		// First clip: 0,0 to 180,180
		ctx.Rectangle(0, 0, 180, 180)
		ctx.Clip()

		// Second clip: 20,20 to 160,160
		ctx.Rectangle(20, 20, 140, 140)
		ctx.Clip()

		// Third clip: 40,40 to 180,180
		ctx.Rectangle(40, 40, 140, 140)
		ctx.Clip()

		// Effective clip is intersection: 40,40 to 160,160
		assert.True(t, ctx.InClip(100, 100), "Center point in intersection")
		assert.False(t, ctx.InClip(10, 10), "Point outside all clips")
		assert.False(t, ctx.InClip(170, 170), "Point outside all clips")
		assert.False(t, ctx.InClip(30, 30), "Point in first two but not third")
	})

	t.Run("NonOverlappingClips", func(t *testing.T) {
		ctx.NewPath()

		// First clip: 0,0 to 50,50
		ctx.Rectangle(0, 0, 50, 50)
		ctx.Clip()

		// Second clip: 100,100 to 200,200 (no overlap)
		ctx.Rectangle(100, 100, 100, 100)
		ctx.Clip()

		// No points should be in clip (empty intersection)
		assert.False(t, ctx.InClip(25, 25), "Point in first clip only")
		assert.False(t, ctx.InClip(150, 150), "Point in second clip only")
		assert.False(t, ctx.InClip(100, 100), "Point at edge")
	})
}

// TestContextClipWithTransform verifies clipping with transformations.
func TestContextClipWithTransform(t *testing.T) {
	ctx, _ := newTestContext(t, 400, 400)

	t.Run("ClipWithTranslation", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.ResetClip()
		ctx.NewPath()

		// Translate coordinate system: user origin moves to device (50, 50)
		ctx.Translate(50, 50)

		// Create clip in translated user coordinates: user (0,0)-(100,100)
		ctx.Rectangle(0, 0, 100, 100)
		ctx.Clip()

		// InClip takes user coordinates. The clip spans user (0,0)-(100,100).
		assert.True(t, ctx.InClip(50, 50), "Point in translated clip (user coords)")
		assert.False(t, ctx.InClip(-25, -25), "Point outside translated clip (user coords)")
		assert.False(t, ctx.InClip(125, 125), "Point outside translated clip (user coords)")
	})

	t.Run("ClipWithScale", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.ResetClip()
		ctx.NewPath()

		// Scale by 2x: user (x,y) maps to device (2x, 2y)
		ctx.Scale(2.0, 2.0)

		// Create clip at user (0,0)-(50,50), which is device (0,0)-(100,100)
		ctx.Rectangle(0, 0, 50, 50)
		ctx.Clip()

		// InClip takes user coordinates. The clip spans user (0,0)-(50,50).
		assert.True(t, ctx.InClip(25, 25), "Point in scaled clip (user coords)")
		assert.True(t, ctx.InClip(49, 49), "Point near edge of scaled clip (user coords)")
		assert.False(t, ctx.InClip(60, 60), "Point outside scaled clip (user coords)")
	})

	t.Run("ClipWithRotation", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.ResetClip()
		ctx.NewPath()

		// Rotate 45 degrees
		ctx.Rotate(math.Pi / 4) // 45 degrees

		// Create rectangular clip in rotated user coordinates
		ctx.Rectangle(-50, -50, 100, 100)
		ctx.Clip()

		// The clip should be a rotated square
		// This is complex to test precisely, just verify it works
		assert.Equal(t, status.Success, ctx.Status())

		// User origin should definitely be in clip
		assert.True(t, ctx.InClip(0, 0), "User origin should be in rotated clip")
	})

	t.Run("ClipBeforeTransform", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.ResetClip()
		ctx.NewPath()

		// Set clip with identity CTM: clip = device (50,50)-(150,150)
		ctx.Rectangle(50, 50, 100, 100)
		ctx.Clip()

		// Then transform: user origin shifts by (25, 25)
		ctx.Translate(25, 25)

		// ClipExtents returns the clip in current user coordinates.
		// The clip (device 50,50-150,150) in the new user space (translated by 25,25)
		// is user (25,25)-(125,125). A point at user (75,75) is device (100,100),
		// which is inside the clip.
		assert.True(t, ctx.InClip(75, 75), "Point inside clip in current user coordinates")
	})

	t.Run("SaveRestoreWithClipAndTransform", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.ResetClip()
		ctx.NewPath()

		// Save state
		ctx.Save()

		// Apply transform and clip
		// After Translate(50,50): user (0,0)-(50,50) = device (50,50)-(100,100)
		ctx.Translate(50, 50)
		ctx.Rectangle(0, 0, 50, 50)
		ctx.Clip()

		// InClip takes user coordinates. Clip spans user (0,0)-(50,50).
		assert.True(t, ctx.InClip(25, 25), "Point in clip (user coords)")
		assert.False(t, ctx.InClip(75, 75), "Point outside clip (user coords)")

		// Restore state: CTM and clip return to pre-Save values
		ctx.Restore()

		// After restore, full surface is clippable again
		assert.True(t, ctx.InClip(125, 125), "Point in clip after restore")
		assert.True(t, ctx.InClip(0, 0), "Point in clip after restore")
	})
}

// TestContextClipAfterClose verifies safe behavior after context close.
func TestContextClipAfterClose(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
	require.NoError(t, err)
	defer surf.Close()

	ctx, err := NewContext(surf)
	require.NoError(t, err)

	// Close the context
	err = ctx.Close()
	require.NoError(t, err)

	t.Run("ClipAfterClose", func(t *testing.T) {
		// Clip after close should not crash (no-op)
		ctx.Rectangle(10, 10, 50, 50)
		ctx.Clip()
		// Should be safe
	})

	t.Run("ClipPreserveAfterClose", func(t *testing.T) {
		ctx.ClipPreserve()
		// Should be safe
	})

	t.Run("ClipExtentsAfterClose", func(t *testing.T) {
		// Should return zeros or handle gracefully
		x1, y1, x2, y2 := ctx.ClipExtents()
		_ = x1
		_ = y1
		_ = x2
		_ = y2
		// Should not crash
	})

	t.Run("InClipAfterClose", func(t *testing.T) {
		// Should return false or handle gracefully
		result := ctx.InClip(50, 50)
		assert.False(t, result, "InClip should return false after close")
	})

	t.Run("ResetClipAfterClose", func(t *testing.T) {
		ctx.ResetClip()
		// Should be safe
	})
}
