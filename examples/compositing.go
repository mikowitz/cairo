// ABOUTME: Example demonstrating Cairo compositing operators for blending drawing operations.
// ABOUTME: Shows OperatorOver, OperatorAdd, OperatorMultiply, and OperatorXor in a visual grid.
package examples

import (
	"fmt"
	"math"

	"github.com/mikowitz/cairo"
)

// GenerateCompositing creates a 400x400 PNG image demonstrating Cairo compositing operators.
//
// The image shows four panels in a 2x2 grid, each demonstrating a different operator
// applied when drawing a red circle over a blue circle:
//   - Top left: OperatorOver - default alpha compositing (red drawn on top of blue)
//   - Top right: OperatorAdd - additive blending (overlap brightens toward white)
//   - Bottom left: OperatorMultiply - multiplicative blending (overlap darkens)
//   - Bottom right: OperatorXor - exclusive-or (overlap becomes transparent)
//
// This demonstrates:
//   - ctx.SetOperator() to change the compositing operator
//   - Porter-Duff operators: Over, Add, Xor
//   - Blend mode operators: Multiply
//   - How operators interact with opaque source colors
//
// All resources are properly cleaned up using defer statements.
func GenerateCompositing(outputPath string) error {
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 400, 400)
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

	type panel struct {
		ox, oy float64
		op     cairo.Operator
	}

	panels := []panel{
		{0, 0, cairo.OperatorOver},
		{panelSize, 0, cairo.OperatorAdd},
		{0, panelSize, cairo.OperatorMultiply},
		{panelSize, panelSize, cairo.OperatorXor},
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
