// ABOUTME: Example demonstrating Cairo fill rules (winding vs even-odd) with a self-intersecting star.
// ABOUTME: Shows how fill rules affect the interior of self-intersecting paths.
package examples

import (
	"fmt"
	"math"

	"github.com/mikowitz/cairo"
)

// GenerateFillRules creates a 400x200 PNG image demonstrating the two Cairo fill rules.
//
// The image shows two panels side by side, each drawing an identical five-pointed star
// using a single self-intersecting path:
//   - Left panel: FillRuleWinding — the interior of the star (including the center pentagon)
//     is filled because the winding count is non-zero everywhere inside the path.
//   - Right panel: FillRuleEvenOdd — the center pentagon is left unfilled because the
//     path crosses it twice (an even number), which the even-odd rule considers "outside".
//
// This demonstrates:
//   - ctx.SetFillRule(cairo.FillRuleWinding) — default, fills entire interior
//   - ctx.SetFillRule(cairo.FillRuleEvenOdd) — creates a hole where the path self-intersects
//
// All resources are properly cleaned up using defer statements.
func GenerateFillRules(outputPath string) error {
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 400, 200)
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

	// Light gray background for the full image
	ctx.SetSourceRGB(0.85, 0.85, 0.85)
	ctx.Paint()

	// Left panel: white background
	ctx.SetSourceRGB(1, 1, 1)
	ctx.Rectangle(2, 2, 196, 196)
	ctx.Fill()

	// Right panel: white background
	ctx.SetSourceRGB(1, 1, 1)
	ctx.Rectangle(202, 2, 196, 196)
	ctx.Fill()

	// Left panel: five-pointed star with FillRuleWinding
	ctx.Save()
	ctx.SetFillRule(cairo.FillRuleWinding)
	ctx.SetSourceRGB(0.2, 0.4, 0.8)
	starPath(ctx, 100, 100, 85)
	ctx.Fill()
	ctx.Restore()

	// Right panel: same star with FillRuleEvenOdd
	ctx.Save()
	ctx.SetFillRule(cairo.FillRuleEvenOdd)
	ctx.SetSourceRGB(0.2, 0.4, 0.8)
	starPath(ctx, 300, 100, 85)
	ctx.Fill()
	ctx.Restore()

	surface.Flush()
	return surface.WriteToPNG(outputPath)
}

// starPath constructs a self-intersecting five-pointed star (pentagram) path centered
// at (cx, cy) with the given radius. It connects every other vertex of a regular
// pentagon, producing a path that crosses itself and creates a center region enclosed
// twice — which makes the even-odd and winding fill rules produce visibly different
// results. The path is not filled or stroked; the caller controls rendering.
func starPath(ctx *cairo.Context, cx, cy, radius float64) {
	const points = 5
	for i := 0; i < points; i++ {
		// Skip every other vertex to create a self-intersecting pentagram
		idx := (i * 2) % points
		angle := float64(idx)*2*math.Pi/float64(points) - math.Pi/2
		x := cx + radius*math.Cos(angle)
		y := cy + radius*math.Sin(angle)
		if i == 0 {
			ctx.MoveTo(x, y)
		} else {
			ctx.LineTo(x, y)
		}
	}
	ctx.ClosePath()
}
