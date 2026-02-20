// ABOUTME: Example demonstrating Cairo compositing operators for blending drawing operations.
// ABOUTME: Shows all 29 Cairo compositing operators in a 5×6 visual grid.
package examples

import (
	"fmt"
	"math"

	"github.com/mikowitz/cairo"
)

// GenerateCompositing creates a 600x720 PNG image demonstrating all Cairo compositing operators.
//
// The image shows 29 panels in a 5×6 grid, each demonstrating a different compositing
// operator applied when drawing a red circle over a blue circle.
//
// Operators are shown in Cairo's numeric order across a 5×6 grid (left to right,
// top to bottom). The 14 Porter-Duff operators fill rows 1–3 with one slot remaining,
// so the first blend mode (Multiply) appears at the end of row 3:
//
//   Row 1: Clear, Source, Over, In, Out
//   Row 2: Atop, Dest, DestOver, DestIn, DestOut
//   Row 3: DestAtop, Xor, Add, Saturate, Multiply (first blend mode)
//   Row 4: Screen, Overlay, Darken, Lighten, ColorDodge
//   Row 5: ColorBurn, HardLight, SoftLight, Difference, Exclusion
//   Row 6: HslHue, HslSaturation, HslColor, HslLuminosity
//
// This demonstrates:
//   - ctx.SetOperator() to change the compositing operator
//   - All 14 Porter-Duff operators
//   - All 15 blend mode operators
//
// All resources are properly cleaned up using defer statements.
func GenerateCompositing(outputPath string) error {
	const panelSize = 120.0

	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 600, 720)
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

	type panel struct {
		ox, oy float64
		op     cairo.Operator
	}

	panels := []panel{
		// Porter-Duff operators
		{0 * panelSize, 0 * panelSize, cairo.OperatorClear},
		{1 * panelSize, 0 * panelSize, cairo.OperatorSource},
		{2 * panelSize, 0 * panelSize, cairo.OperatorOver},
		{3 * panelSize, 0 * panelSize, cairo.OperatorIn},
		{4 * panelSize, 0 * panelSize, cairo.OperatorOut},

		{0 * panelSize, 1 * panelSize, cairo.OperatorAtop},
		{1 * panelSize, 1 * panelSize, cairo.OperatorDest},
		{2 * panelSize, 1 * panelSize, cairo.OperatorDestOver},
		{3 * panelSize, 1 * panelSize, cairo.OperatorDestIn},
		{4 * panelSize, 1 * panelSize, cairo.OperatorDestOut},

		{0 * panelSize, 2 * panelSize, cairo.OperatorDestAtop},
		{1 * panelSize, 2 * panelSize, cairo.OperatorXor},
		{2 * panelSize, 2 * panelSize, cairo.OperatorAdd},
		{3 * panelSize, 2 * panelSize, cairo.OperatorSaturate},

		// Blend mode operators
		{4 * panelSize, 2 * panelSize, cairo.OperatorMultiply},

		{0 * panelSize, 3 * panelSize, cairo.OperatorScreen},
		{1 * panelSize, 3 * panelSize, cairo.OperatorOverlay},
		{2 * panelSize, 3 * panelSize, cairo.OperatorDarken},
		{3 * panelSize, 3 * panelSize, cairo.OperatorLighten},
		{4 * panelSize, 3 * panelSize, cairo.OperatorColorDodge},

		{0 * panelSize, 4 * panelSize, cairo.OperatorColorBurn},
		{1 * panelSize, 4 * panelSize, cairo.OperatorHardLight},
		{2 * panelSize, 4 * panelSize, cairo.OperatorSoftLight},
		{3 * panelSize, 4 * panelSize, cairo.OperatorDifference},
		{4 * panelSize, 4 * panelSize, cairo.OperatorExclusion},

		{0 * panelSize, 5 * panelSize, cairo.OperatorHslHue},
		{1 * panelSize, 5 * panelSize, cairo.OperatorHslSaturation},
		{2 * panelSize, 5 * panelSize, cairo.OperatorHslColor},
		{3 * panelSize, 5 * panelSize, cairo.OperatorHslLuminosity},
	}

	for _, p := range panels {
		ctx.Save()
		drawCompositingPanel(ctx, p.ox, p.oy, panelSize, p.op)
		ctx.Restore()
	}

	surface.Flush()

	if err := surface.WriteToPNG(outputPath); err != nil {
		return fmt.Errorf("failed to write PNG: %w", err)
	}

	return nil
}

// drawCompositingPanel draws a single compositing demonstration panel.
//
// Each panel shows a blue circle drawn first with OperatorOver, followed by a red
// circle drawn with the specified operator. The overlap region reveals the operator's
// blending effect. The background is white so transparent results show clearly.
func drawCompositingPanel(ctx *cairo.Context, ox, oy, size float64, op cairo.Operator) {
	// Fill panel with white using OperatorSource to replace any existing content
	ctx.SetOperator(cairo.OperatorSource)
	ctx.SetSourceRGB(1.0, 1.0, 1.0)
	ctx.Rectangle(ox, oy, size, size)
	ctx.Fill()

	cx := ox + size/2
	cy := oy + size/2
	r := size * 0.28
	offset := size * 0.13

	// Draw blue circle (left) with OperatorOver
	ctx.SetOperator(cairo.OperatorOver)
	ctx.SetSourceRGB(0.18, 0.37, 0.87)
	ctx.Arc(cx-offset, cy, r, 0, 2*math.Pi)
	ctx.Fill()

	// Draw red circle (right) with the test operator; overlap shows the blend effect
	ctx.SetOperator(op)
	ctx.SetSourceRGB(0.87, 0.18, 0.10)
	ctx.Arc(cx+offset, cy, r, 0, 2*math.Pi)
	ctx.Fill()

	// Reset to OperatorOver for decorative elements
	ctx.SetOperator(cairo.OperatorOver)

	// Panel border
	ctx.SetSourceRGB(0.4, 0.4, 0.4)
	ctx.SetLineWidth(1.0)
	ctx.Rectangle(ox, oy, size, size)
	ctx.Stroke()
}
