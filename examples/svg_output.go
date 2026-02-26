// ABOUTME: Example demonstrating SVG surface for web-compatible vector graphics output.
// ABOUTME: Creates an SVG with shapes, gradients, and text in a single 600×400pt document.

//go:build !nosvg

package examples

import (
	"fmt"
	"math"

	"github.com/mikowitz/cairo"
	"github.com/mikowitz/cairo/font"
	"github.com/mikowitz/cairo/surface"
)

// GenerateSVGOutput creates an SVG file demonstrating Cairo's SVG surface capabilities.
//
// The SVG (600×400 points) is divided into three sections:
//   - Left (x 0–290): Basic geometric shapes (filled rectangle, stroked rectangle,
//     filled circle, arc, semi-transparent triangle)
//   - Right (x 300–600): Linear and radial gradient patterns
//   - Bottom (y 210–400): Text rendered with various fonts, weights, and slants,
//     plus a large vector text outline via TextPath
//
// The coordinate origin is at the top-left corner of the image.
// All resources are properly cleaned up using defer statements.
func GenerateSVGOutput(outputPath string) error {
	const (
		w = 600.0
		h = 400.0
	)

	svg, err := surface.NewSVGSurface(outputPath, w, h)
	if err != nil {
		return fmt.Errorf("failed to create SVG surface: %w", err)
	}
	defer func() {
		_ = svg.Close()
	}()

	ctx, err := cairo.NewContext(svg)
	if err != nil {
		return fmt.Errorf("failed to create context: %w", err)
	}
	defer func() {
		_ = ctx.Close()
	}()

	// White background
	ctx.SetSourceRGB(1, 1, 1)
	ctx.Paint()

	drawSVGShapes(ctx)
	if err := drawSVGGradients(ctx); err != nil {
		return err
	}
	drawSVGText(ctx)
	return nil
}

// drawSVGShapes draws basic geometric shapes in the left portion of the SVG.
func drawSVGShapes(ctx *cairo.Context) {
	// Red filled rectangle
	ctx.SetSourceRGB(0.8, 0.2, 0.2)
	ctx.Rectangle(20, 20, 120, 65)
	ctx.Fill()

	// Blue stroked rectangle
	ctx.SetSourceRGB(0.2, 0.2, 0.8)
	ctx.SetLineWidth(3)
	ctx.Rectangle(20, 105, 120, 65)
	ctx.Stroke()

	// Green filled circle
	ctx.SetSourceRGB(0.2, 0.7, 0.2)
	ctx.Arc(225, 60, 42, 0, 2*math.Pi)
	ctx.Fill()

	// Orange arc (half-circle outline)
	ctx.SetSourceRGB(1.0, 0.5, 0.0)
	ctx.SetLineWidth(4)
	ctx.Arc(225, 140, 38, 0, math.Pi)
	ctx.Stroke()

	// Semi-transparent purple triangle with dark outline
	ctx.SetSourceRGBA(0.5, 0.1, 0.8, 0.65)
	ctx.MoveTo(165, 20)
	ctx.LineTo(115, 185)
	ctx.LineTo(215, 185)
	ctx.ClosePath()
	ctx.FillPreserve()
	ctx.SetSourceRGB(0.2, 0, 0.4)
	ctx.SetLineWidth(2)
	ctx.Stroke()
}

// drawSVGGradients draws linear and radial gradient patterns in the right portion.
func drawSVGGradients(ctx *cairo.Context) error {
	// Horizontal rainbow linear gradient bar
	lg, err := cairo.NewLinearGradient(310, 0, 580, 0)
	if err != nil {
		return fmt.Errorf("failed to create linear gradient: %w", err)
	}
	defer func() { _ = lg.Close() }()
	lg.AddColorStopRGB(0.0, 1.0, 0.0, 0.0)
	lg.AddColorStopRGB(0.5, 0.0, 1.0, 0.0)
	lg.AddColorStopRGB(1.0, 0.0, 0.0, 1.0)
	ctx.SetSource(lg)
	ctx.Rectangle(310, 20, 270, 65)
	ctx.Fill()

	// Warm-toned radial gradient circle
	rg, err := cairo.NewRadialGradient(400, 145, 0, 400, 145, 52)
	if err != nil {
		return fmt.Errorf("failed to create radial gradient: %w", err)
	}
	defer func() { _ = rg.Close() }()
	rg.AddColorStopRGB(0.0, 1.0, 1.0, 0.8)
	rg.AddColorStopRGB(0.55, 1.0, 0.5, 0.0)
	rg.AddColorStopRGB(1.0, 0.45, 0.0, 0.0)
	ctx.SetSource(rg)
	ctx.Arc(400, 145, 52, 0, 2*math.Pi)
	ctx.Fill()

	// Blue fade-to-transparent radial gradient
	tg, err := cairo.NewRadialGradient(520, 145, 0, 520, 145, 45)
	if err != nil {
		return fmt.Errorf("failed to create fade gradient: %w", err)
	}
	defer func() { _ = tg.Close() }()
	tg.AddColorStopRGBA(0.0, 0.2, 0.4, 1.0, 1.0)
	tg.AddColorStopRGBA(1.0, 0.2, 0.4, 1.0, 0.0)
	ctx.SetSource(tg)
	ctx.Arc(520, 145, 45, 0, 2*math.Pi)
	ctx.Fill()

	return nil
}

// drawSVGText draws text samples in the bottom portion of the SVG.
func drawSVGText(ctx *cairo.Context) {
	// Section label
	ctx.SetSourceRGB(0, 0, 0)
	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightBold)
	ctx.SetFontSize(15)
	ctx.MoveTo(20, 225)
	ctx.ShowText("Text Rendering in SVG:")

	// Normal sans-serif
	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightNormal)
	ctx.SetFontSize(13)
	ctx.MoveTo(20, 255)
	ctx.ShowText("Normal sans-serif")

	// Bold sans-serif
	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightBold)
	ctx.MoveTo(210, 255)
	ctx.ShowText("Bold sans-serif")

	// Italic serif
	ctx.SetSourceRGB(0.2, 0.2, 0.7)
	ctx.SelectFontFace("serif", font.SlantItalic, font.WeightNormal)
	ctx.MoveTo(20, 285)
	ctx.ShowText("Italic serif")

	// Oblique monospace
	ctx.SetSourceRGB(0.1, 0.5, 0.1)
	ctx.SelectFontFace("monospace", font.SlantOblique, font.WeightNormal)
	ctx.MoveTo(210, 285)
	ctx.ShowText("Oblique mono")

	// Large vector text outline via TextPath
	ctx.SetSourceRGB(0.7, 0.1, 0.1)
	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightBold)
	ctx.SetFontSize(38)
	ctx.MoveTo(20, 365)
	ctx.TextPath("Vector SVG")
	ctx.Fill()
}
