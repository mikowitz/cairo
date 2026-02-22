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
//	Row 1: Clear, Source, Over, In, Out
//	Row 2: Atop, Dest, DestOver, DestIn, DestOut
//	Row 3: DestAtop, Xor, Add, Saturate, Multiply (first blend mode)
//	Row 4: Screen, Overlay, Darken, Lighten, ColorDodge
//	Row 5: ColorBurn, HardLight, SoftLight, Difference, Exclusion
//	Row 6: HslHue, HslSaturation, HslColor, HslLuminosity
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

	for i, op := range compositingOperators() {
		col := float64(i % 5)
		row := float64(i / 5)
		ctx.Save()
		drawCompositingPanel(ctx, col*panelSize, row*panelSize, panelSize, op)
		ctx.Restore()
	}

	surface.Flush()

	if err := surface.WriteToPNG(outputPath); err != nil {
		return fmt.Errorf("failed to write PNG: %w", err)
	}

	return nil
}

// compositingOperators returns all 29 Cairo compositing operators in numeric order.
// The caller maps them onto a grid by index: column = i%5, row = i/5.
func compositingOperators() []cairo.Operator {
	return []cairo.Operator{
		// Porter-Duff operators
		cairo.OperatorClear,
		cairo.OperatorSource,
		cairo.OperatorOver,
		cairo.OperatorIn,
		cairo.OperatorOut,
		cairo.OperatorAtop,
		cairo.OperatorDest,
		cairo.OperatorDestOver,
		cairo.OperatorDestIn,
		cairo.OperatorDestOut,
		cairo.OperatorDestAtop,
		cairo.OperatorXor,
		cairo.OperatorAdd,
		cairo.OperatorSaturate,
		// Blend mode operators
		cairo.OperatorMultiply,
		cairo.OperatorScreen,
		cairo.OperatorOverlay,
		cairo.OperatorDarken,
		cairo.OperatorLighten,
		cairo.OperatorColorDodge,
		cairo.OperatorColorBurn,
		cairo.OperatorHardLight,
		cairo.OperatorSoftLight,
		cairo.OperatorDifference,
		cairo.OperatorExclusion,
		cairo.OperatorHslHue,
		cairo.OperatorHslSaturation,
		cairo.OperatorHslColor,
		cairo.OperatorHslLuminosity,
	}
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
