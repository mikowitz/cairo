package examples

import (
	"fmt"

	"github.com/mikowitz/cairo"
)

// GenerateLineStyles creates a 700x600 PNG image demonstrating Cairo line styling options.
//
// The image contains:
//   - Top section: Three line caps (Butt, Round, Square) shown with thick lines
//   - Middle section: Three line joins (Miter, Round, Bevel) shown at sharp angles
//   - Bottom left: Various dash patterns (solid, dashed, dotted, complex)
//   - Bottom right: Miter limit demonstration showing the effect on sharp angles
//
// This demonstrates:
//   - Setting line cap styles with SetLineCap
//   - Setting line join styles with SetLineJoin
//   - Creating dash patterns with SetDash
//   - Adjusting miter limits with SetMiterLimit
//   - Setting line widths with SetLineWidth
//   - Using Save/Restore to manage line style state
//
// All resources are properly cleaned up using defer statements.
func GenerateLineStyles(outputPath string) error { //nolint:funlen
	// Create a 700x600 ARGB32 image surface
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 700, 600)
	if err != nil {
		return fmt.Errorf("failed to create surface: %w", err)
	}
	defer func() {
		_ = surface.Close()
	}()

	// Create drawing context
	ctx, err := cairo.NewContext(surface)
	if err != nil {
		return fmt.Errorf("failed to create context: %w", err)
	}
	defer func() {
		_ = ctx.Close()
	}()

	// Fill the background with white
	ctx.SetSourceRGB(1.0, 1.0, 1.0)
	ctx.Paint()

	// Set up common line properties
	ctx.SetSourceRGB(0.0, 0.0, 0.0) // Black lines
	ctx.SetLineWidth(20.0)

	// === TOP SECTION: Line Caps ===
	// Draw guide lines to show line endpoints
	ctx.Save()
	ctx.SetLineWidth(1.0)
	ctx.SetSourceRGB(0.8, 0.8, 0.8) // Light gray guides

	// Guide lines for caps section
	for x := 50.0; x <= 650.0; x += 200.0 {
		ctx.MoveTo(x, 30)
		ctx.LineTo(x, 160)
		ctx.Stroke()
	}
	ctx.Restore()

	// 1. Line Cap: Butt (default)
	ctx.Save()
	ctx.SetLineCap(cairo.LineCapButt)
	ctx.SetLineWidth(20.0)
	ctx.SetSourceRGB(0.8, 0.0, 0.0) // Dark red
	ctx.MoveTo(50, 60)
	ctx.LineTo(50, 130)
	ctx.Stroke()
	ctx.Restore()

	// 2. Line Cap: Round
	ctx.Save()
	ctx.SetLineCap(cairo.LineCapRound)
	ctx.SetLineWidth(20.0)
	ctx.SetSourceRGB(0.0, 0.6, 0.0) // Dark green
	ctx.MoveTo(250, 60)
	ctx.LineTo(250, 130)
	ctx.Stroke()
	ctx.Restore()

	// 3. Line Cap: Square
	ctx.Save()
	ctx.SetLineCap(cairo.LineCapSquare)
	ctx.SetLineWidth(20.0)
	ctx.SetSourceRGB(0.0, 0.0, 0.8) // Dark blue
	ctx.MoveTo(450, 60)
	ctx.LineTo(450, 130)
	ctx.Stroke()
	ctx.Restore()

	// === MIDDLE SECTION: Line Joins ===
	yOffset := 180.0

	// 4. Line Join: Miter
	ctx.Save()
	ctx.SetLineJoin(cairo.LineJoinMiter)
	ctx.SetLineCap(cairo.LineCapButt)
	ctx.SetLineWidth(20.0)
	ctx.SetSourceRGB(0.8, 0.0, 0.0) // Dark red
	ctx.MoveTo(20, yOffset+70)
	ctx.LineTo(80, yOffset+20)
	ctx.LineTo(140, yOffset+70)
	ctx.Stroke()
	ctx.Restore()

	// 5. Line Join: Round
	ctx.Save()
	ctx.SetLineJoin(cairo.LineJoinRound)
	ctx.SetLineCap(cairo.LineCapButt)
	ctx.SetLineWidth(20.0)
	ctx.SetSourceRGB(0.0, 0.6, 0.0) // Dark green
	ctx.MoveTo(220, yOffset+70)
	ctx.LineTo(280, yOffset+20)
	ctx.LineTo(340, yOffset+70)
	ctx.Stroke()
	ctx.Restore()

	// 6. Line Join: Bevel
	ctx.Save()
	ctx.SetLineJoin(cairo.LineJoinBevel)
	ctx.SetLineCap(cairo.LineCapButt)
	ctx.SetLineWidth(20.0)
	ctx.SetSourceRGB(0.0, 0.0, 0.8) // Dark blue
	ctx.MoveTo(420, yOffset+70)
	ctx.LineTo(480, yOffset+20)
	ctx.LineTo(540, yOffset+70)
	ctx.Stroke()
	ctx.Restore()

	// === BOTTOM LEFT: Dash Patterns ===
	yOffset = 340.0

	// 7. Solid line (no dash)
	ctx.Save()
	ctx.SetLineWidth(6.0)
	ctx.SetSourceRGB(0.0, 0.0, 0.0)
	err = ctx.SetDash(nil, 0.0) // Solid
	if err != nil {
		return fmt.Errorf("failed to set dash: %w", err)
	}
	ctx.MoveTo(20, yOffset)
	ctx.LineTo(300, yOffset)
	ctx.Stroke()
	ctx.Restore()

	// 8. Simple dash pattern
	yOffset += 30
	ctx.Save()
	ctx.SetLineWidth(6.0)
	ctx.SetSourceRGB(0.0, 0.0, 0.0)
	err = ctx.SetDash([]float64{20.0, 10.0}, 0.0)
	if err != nil {
		return fmt.Errorf("failed to set dash: %w", err)
	}
	ctx.MoveTo(20, yOffset)
	ctx.LineTo(300, yOffset)
	ctx.Stroke()
	ctx.Restore()

	// 9. Dotted pattern with round caps
	yOffset += 30
	ctx.Save()
	ctx.SetLineCap(cairo.LineCapRound)
	ctx.SetLineWidth(6.0)
	ctx.SetSourceRGB(0.0, 0.0, 0.0)
	err = ctx.SetDash([]float64{0.0, 10.0}, 0.0)
	if err != nil {
		return fmt.Errorf("failed to set dash: %w", err)
	}
	ctx.MoveTo(20, yOffset)
	ctx.LineTo(300, yOffset)
	ctx.Stroke()
	ctx.Restore()

	// 10. Complex dash pattern
	yOffset += 30
	ctx.Save()
	ctx.SetLineWidth(6.0)
	ctx.SetSourceRGB(0.0, 0.0, 0.0)
	err = ctx.SetDash([]float64{20.0, 5.0, 5.0, 5.0}, 0.0)
	if err != nil {
		return fmt.Errorf("failed to set dash: %w", err)
	}
	ctx.MoveTo(20, yOffset)
	ctx.LineTo(300, yOffset)
	ctx.Stroke()
	ctx.Restore()

	// 11. Dash with offset
	yOffset += 30
	ctx.Save()
	ctx.SetLineWidth(6.0)
	ctx.SetSourceRGB(0.6, 0.0, 0.6) // Purple
	err = ctx.SetDash([]float64{20.0, 10.0}, 15.0) // Offset by 15
	if err != nil {
		return fmt.Errorf("failed to set dash: %w", err)
	}
	ctx.MoveTo(20, yOffset)
	ctx.LineTo(300, yOffset)
	ctx.Stroke()
	ctx.Restore()

	// === BOTTOM RIGHT: Miter Limit ===
	xOffset := 380.0
	yOffset = 340.0

	// 12. Miter limit = 10 (default) - miter preserved
	ctx.Save()
	ctx.SetLineJoin(cairo.LineJoinMiter)
	ctx.SetMiterLimit(10.0)
	ctx.SetLineCap(cairo.LineCapButt)
	ctx.SetLineWidth(15.0)
	ctx.SetSourceRGB(0.0, 0.6, 0.0) // Green
	ctx.MoveTo(xOffset, yOffset+30)
	ctx.LineTo(xOffset+80, yOffset+10)
	ctx.LineTo(xOffset+160, yOffset+30)
	ctx.Stroke()
	ctx.Restore()

	// 13. Miter limit = 2 - converted to bevel at sharper angles
	yOffset += 60
	ctx.Save()
	ctx.SetLineJoin(cairo.LineJoinMiter)
	ctx.SetMiterLimit(2.0)
	ctx.SetLineCap(cairo.LineCapButt)
	ctx.SetLineWidth(15.0)
	ctx.SetSourceRGB(0.8, 0.5, 0.0) // Orange
	ctx.MoveTo(xOffset, yOffset+30)
	ctx.LineTo(xOffset+80, yOffset+10)
	ctx.LineTo(xOffset+160, yOffset+30)
	ctx.Stroke()
	ctx.Restore()

	// 14. Very sharp angle with high miter limit
	yOffset += 60
	ctx.Save()
	ctx.SetLineJoin(cairo.LineJoinMiter)
	ctx.SetMiterLimit(50.0) // High limit to show extreme miter
	ctx.SetLineCap(cairo.LineCapButt)
	ctx.SetLineWidth(15.0)
	ctx.SetSourceRGB(0.8, 0.0, 0.0) // Red
	ctx.MoveTo(xOffset, yOffset+30)
	ctx.LineTo(xOffset+80, yOffset+5) // Very sharp angle
	ctx.LineTo(xOffset+160, yOffset+30)
	ctx.Stroke()
	ctx.Restore()

	// Flush any pending operations
	surface.Flush()

	// Write the surface to a PNG file
	if err := surface.WriteToPNG(outputPath); err != nil {
		return fmt.Errorf("failed to write PNG: %w", err)
	}

	return nil
}
