// ABOUTME: Example demonstrating Cairo clipping operations.
// ABOUTME: Shows rectangular, circular, nested, ClipPreserve, and Save/Restore clipping techniques.
package examples

import (
	"fmt"
	"math"

	"github.com/mikowitz/cairo"
)

// GenerateClipping creates a 600x600 PNG image demonstrating Cairo clipping operations.
//
// The image shows six panels in a 3x2 grid, each demonstrating a different clipping technique:
//   - Top left: Rectangular clip restricting a colorful striped background
//   - Top middle: Circular clip using Arc on the same striped background
//   - Top right: ClipPreserve stroking the clip boundary before filling the interior
//   - Bottom left: Nested clips showing intersection of two overlapping regions
//   - Bottom middle: Save/Restore showing that clip state is part of the graphics state
//   - Bottom right: Clip applied after a rotation transformation
//
// This demonstrates:
//   - ctx.Clip() - establish a clip region from the current path, consuming the path
//   - ctx.ClipPreserve() - clip while retaining the path for further drawing operations
//   - ctx.ResetClip() - clear the current clip to the full surface
//   - Clip intersection via sequential Clip() calls
//   - Save/Restore preserving and restoring clip state
//   - Clipping in a transformed (rotated) coordinate space
//
// All resources are properly cleaned up using defer statements.
func GenerateClipping(outputPath string) error {
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 600, 600)
	if err != nil {
		return fmt.Errorf("failed to create surface: %w", err)
	}
	defer func() {
		_ = surface.Close()
	}()

	ctx, err := cairo.NewContext(surface)
	if err != nil {
		return fmt.Errorf("failed to create context: %w", err)
	}
	defer func() {
		_ = ctx.Close()
	}()

	// Light gray background separates the panels visually
	ctx.SetSourceRGB(0.75, 0.75, 0.75)
	ctx.Paint()

	const panelSize = 200.0

	// 1. Top-left: Basic rectangular clip
	ctx.Save()
	drawRectangularClipPanel(ctx, 0, 0, panelSize)
	ctx.Restore()

	// 2. Top-middle: Circular clip
	ctx.Save()
	drawCircularClipPanel(ctx, panelSize, 0, panelSize)
	ctx.Restore()

	// 3. Top-right: ClipPreserve
	ctx.Save()
	drawClipPreservePanel(ctx, 2*panelSize, 0, panelSize)
	ctx.Restore()

	// 4. Bottom-left: Nested clip intersection
	ctx.Save()
	drawNestedClipsPanel(ctx, 0, panelSize, panelSize)
	ctx.Restore()

	// 5. Bottom-middle: Save/Restore with clip
	ctx.Save()
	drawSaveRestoreClipPanel(ctx, panelSize, panelSize, panelSize)
	ctx.Restore()

	// 6. Bottom-right: Clip in transformed coordinate space
	ctx.Save()
	drawTransformedClipPanel(ctx, 2*panelSize, panelSize, panelSize)
	ctx.Restore()

	surface.Flush()

	if err := surface.WriteToPNG(outputPath); err != nil {
		return fmt.Errorf("failed to write PNG: %w", err)
	}

	return nil
}

// drawClippingStripes fills a rectangular area with alternating colored horizontal stripes.
// The stripes provide a visually distinctive background that makes clipped regions obvious.
func drawClippingStripes(ctx *cairo.Context, x, y, width, height float64) {
	colors := [][3]float64{
		{1.0, 0.6, 0.6}, // coral
		{0.6, 0.8, 1.0}, // sky blue
		{0.8, 1.0, 0.6}, // lime
		{1.0, 0.9, 0.5}, // amber
	}
	stripeH := 20.0
	n := int(height/stripeH) + 1
	for i := 0; i < n; i++ {
		c := colors[i%len(colors)]
		ctx.SetSourceRGB(c[0], c[1], c[2])
		ctx.Rectangle(x, y+float64(i)*stripeH, width, stripeH)
		ctx.Fill()
	}
}

// drawClippingPanelBase draws a white background and gray border for a panel.
func drawClippingPanelBase(ctx *cairo.Context, x, y, size float64) {
	ctx.SetSourceRGB(1.0, 1.0, 1.0)
	ctx.Rectangle(x, y, size, size)
	ctx.Fill()

	ctx.SetSourceRGB(0.5, 0.5, 0.5)
	ctx.SetLineWidth(1.0)
	ctx.Rectangle(x, y, size, size)
	ctx.Stroke()
}

// drawRectangularClipPanel demonstrates Clip() with a rectangular path.
// A colorful striped background is drawn but only the area inside the clip rectangle is visible.
func drawRectangularClipPanel(ctx *cairo.Context, ox, oy, size float64) {
	drawClippingPanelBase(ctx, ox, oy, size)

	pad := 25.0

	// Establish a rectangular clip region inset from the panel edges
	ctx.Rectangle(ox+pad, oy+pad, size-2*pad, size-2*pad)
	ctx.Clip()

	// Stripes cover the full panel but only the clipped rectangle shows
	drawClippingStripes(ctx, ox, oy, size, size)

	ctx.ResetClip()

	// Redraw border on top of the clipped content
	ctx.SetSourceRGB(0.5, 0.5, 0.5)
	ctx.SetLineWidth(1.0)
	ctx.Rectangle(ox, oy, size, size)
	ctx.Stroke()
}

// drawCircularClipPanel demonstrates Clip() with a circular path created via Arc.
// The striped background is visible only through a circular window.
func drawCircularClipPanel(ctx *cairo.Context, ox, oy, size float64) {
	drawClippingPanelBase(ctx, ox, oy, size)

	cx := ox + size/2
	cy := oy + size/2
	radius := size/2 - 20

	// Establish a circular clip region
	ctx.Arc(cx, cy, radius, 0, 2*math.Pi)
	ctx.Clip()

	// Stripes cover the full panel but only the circle shows
	drawClippingStripes(ctx, ox, oy, size, size)

	ctx.ResetClip()

	ctx.SetSourceRGB(0.5, 0.5, 0.5)
	ctx.SetLineWidth(1.0)
	ctx.Rectangle(ox, oy, size, size)
	ctx.Stroke()
}

// drawClipPreservePanel demonstrates ClipPreserve(), which sets the clip while
// retaining the current path. The preserved path is then stroked to show the clip boundary.
func drawClipPreservePanel(ctx *cairo.Context, ox, oy, size float64) {
	drawClippingPanelBase(ctx, ox, oy, size)

	cx := ox + size/2
	cy := oy + size/2
	radius := size/2 - 20

	// ClipPreserve sets the clip AND keeps the path for subsequent drawing
	ctx.Arc(cx, cy, radius, 0, 2*math.Pi)
	ctx.ClipPreserve()

	// Fill the interior (clip is active, only the circle region is painted)
	drawClippingStripes(ctx, ox, oy, size, size)

	// Stroke the preserved path to draw the clip boundary
	ctx.SetSourceRGB(0.1, 0.1, 0.7)
	ctx.SetLineWidth(3.0)
	ctx.Stroke()

	ctx.ResetClip()

	ctx.SetSourceRGB(0.5, 0.5, 0.5)
	ctx.SetLineWidth(1.0)
	ctx.Rectangle(ox, oy, size, size)
	ctx.Stroke()
}

// drawNestedClipsPanel demonstrates sequential Clip() calls, which intersect.
// Each Clip() call intersects the new path with the existing clip region.
// Only the intersection of both rectangles receives paint.
func drawNestedClipsPanel(ctx *cairo.Context, ox, oy, size float64) {
	drawClippingPanelBase(ctx, ox, oy, size)

	// First clip: left 65% of the panel (with padding)
	ctx.Rectangle(ox+15, oy+15, size*0.65, size-30)
	ctx.Clip()

	// Second clip: top 65% of the panel (with padding)
	// Intersects with first; only the top-left corner region receives paint.
	ctx.Rectangle(ox+15, oy+15, size-30, size*0.65)
	ctx.Clip()

	drawClippingStripes(ctx, ox, oy, size, size)

	ctx.ResetClip()

	// Draw outlines showing each clip region independently
	ctx.SetLineWidth(1.5)
	ctx.SetSourceRGBA(0.7, 0.0, 0.0, 0.7)
	ctx.Rectangle(ox+15, oy+15, size*0.65, size-30)
	ctx.Stroke()

	ctx.SetSourceRGBA(0.0, 0.0, 0.7, 0.7)
	ctx.Rectangle(ox+15, oy+15, size-30, size*0.65)
	ctx.Stroke()

	ctx.SetSourceRGB(0.5, 0.5, 0.5)
	ctx.SetLineWidth(1.0)
	ctx.Rectangle(ox, oy, size, size)
	ctx.Stroke()
}

// drawSaveRestoreClipPanel demonstrates that the clip region is part of the graphics
// state managed by Save/Restore. A clip applied inside Save/Restore does not persist
// after Restore.
func drawSaveRestoreClipPanel(ctx *cairo.Context, ox, oy, size float64) {
	drawClippingPanelBase(ctx, ox, oy, size)

	// Draw stripes across the full panel
	drawClippingStripes(ctx, ox, oy, size, size)

	cx := ox + size/2
	cy := oy + size/2
	innerRadius := size/2 - 35

	// Inside Save/Restore: clip to a circle and paint white to erase the stripes inside it.
	// After Restore, the clip is gone, demonstrating that clip state is saved/restored.
	ctx.Save()
	ctx.Arc(cx, cy, innerRadius, 0, 2*math.Pi)
	ctx.Clip()
	ctx.SetSourceRGB(1.0, 1.0, 1.0)
	ctx.Paint()
	ctx.Restore()

	// After Restore the clip is cleared: this stroke is drawn unclipped
	ctx.SetSourceRGB(0.2, 0.2, 0.2)
	ctx.SetLineWidth(2.0)
	ctx.Arc(cx, cy, innerRadius, 0, 2*math.Pi)
	ctx.Stroke()

	ctx.SetSourceRGB(0.5, 0.5, 0.5)
	ctx.SetLineWidth(1.0)
	ctx.Rectangle(ox, oy, size, size)
	ctx.Stroke()
}

// drawTransformedClipPanel demonstrates that clip regions are defined in the current
// user coordinate space. A rotation is applied before clipping, so the clip rectangle
// is rotated relative to the surface.
func drawTransformedClipPanel(ctx *cairo.Context, ox, oy, size float64) {
	drawClippingPanelBase(ctx, ox, oy, size)

	cx := ox + size/2
	cy := oy + size/2
	halfW := size/2 - 25
	halfH := size/2 - 45

	// Translate to panel center and rotate 30 degrees before clipping
	ctx.Translate(cx, cy)
	ctx.Rotate(math.Pi / 6)

	// Clip to a rectangle in the rotated coordinate space
	ctx.Rectangle(-halfW, -halfH, 2*halfW, 2*halfH)
	ctx.Clip()

	// Draw vertical stripes in rotated space; the clip is a rotated rectangle
	for i := -int(size); i < int(size); i += 15 {
		if i%2 == 0 {
			ctx.SetSourceRGB(1.0, 0.7, 0.5)
		} else {
			ctx.SetSourceRGB(0.5, 0.7, 1.0)
		}
		ctx.Rectangle(float64(i), -size, 10, 2*size)
		ctx.Fill()
	}

	ctx.ResetClip()

	// Stroke the clip boundary in the same rotated coordinate space
	ctx.SetSourceRGB(0.1, 0.1, 0.7)
	ctx.SetLineWidth(2.0)
	ctx.Rectangle(-halfW, -halfH, 2*halfW, 2*halfH)
	ctx.Stroke()

	// Reset transform before drawing the panel border in surface coordinates
	ctx.IdentityMatrix()

	ctx.SetSourceRGB(0.5, 0.5, 0.5)
	ctx.SetLineWidth(1.0)
	ctx.Rectangle(ox, oy, size, size)
	ctx.Stroke()
}
