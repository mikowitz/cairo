// ABOUTME: Example demonstrating PDF surface for multi-page vector output.
// ABOUTME: Creates a 3-page PDF with shapes, gradients, and text.

//go:build !nopdf

package examples

import (
	"fmt"
	"math"

	"github.com/mikowitz/cairo"
	"github.com/mikowitz/cairo/font"
)

// GeneratePDFOutput creates a 3-page PDF document demonstrating Cairo's PDF surface.
//
// The document uses US Letter size (612×792 points, where 1 point = 1/72 inch):
//   - Page 1: Basic shapes (filled rectangle, stroked rectangle, circle, triangle)
//   - Page 2: Linear and radial gradient patterns
//   - Page 3: Text with different fonts, weights, and slants
//
// The coordinate origin is at the top-left corner of each page.
// All resources are properly cleaned up using defer statements.
func GeneratePDFOutput(outputPath string) error {
	const (
		pageW = 612.0 // US Letter width in points (8.5 inches)
		pageH = 792.0 // US Letter height in points (11 inches)
	)

	pdf, err := cairo.NewPDFSurface(outputPath, pageW, pageH)
	if err != nil {
		return fmt.Errorf("failed to create PDF surface: %w", err)
	}
	defer func() {
		_ = pdf.Close()
	}()

	ctx, err := cairo.NewContext(pdf)
	if err != nil {
		return fmt.Errorf("failed to create context: %w", err)
	}
	defer func() {
		_ = ctx.Close()
	}()

	// Page 1: Basic shapes
	drawPDFShapesPage(ctx, pageW, pageH)
	pdf.ShowPage()

	// Page 2: Gradients
	if err := drawPDFGradientsPage(ctx, pageW, pageH); err != nil {
		return err
	}
	pdf.ShowPage()

	// Page 3: Text
	drawPDFTextPage(ctx, pageW)
	return nil
}

// drawPDFShapesPage draws basic geometric shapes demonstrating fill, stroke, and transparency.
func drawPDFShapesPage(ctx *cairo.Context, w, h float64) {
	ctx.SetSourceRGB(1, 1, 1)
	ctx.Paint()

	ctx.SetSourceRGB(0, 0, 0)
	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightBold)
	ctx.SetFontSize(28)
	ctx.MoveTo(50, 60)
	ctx.ShowText("Page 1: Shapes")

	// Red filled rectangle
	ctx.SetSourceRGB(0.8, 0.2, 0.2)
	ctx.Rectangle(50, 100, 220, 130)
	ctx.Fill()

	// Blue stroked rectangle
	ctx.SetSourceRGB(0.2, 0.2, 0.8)
	ctx.SetLineWidth(4)
	ctx.Rectangle(340, 100, 220, 130)
	ctx.Stroke()

	// Green circle
	ctx.SetSourceRGB(0.2, 0.7, 0.2)
	ctx.Arc(w/4, 360, 100, 0, 2*math.Pi)
	ctx.Fill()

	// Orange arc (half circle)
	ctx.SetSourceRGB(1.0, 0.5, 0.0)
	ctx.SetLineWidth(6)
	ctx.Arc(3*w/4, 360, 100, 0, math.Pi)
	ctx.Stroke()

	// Semi-transparent purple triangle with black outline
	ctx.SetSourceRGBA(0.4, 0.1, 0.8, 0.7)
	ctx.MoveTo(w/2, 530)
	ctx.LineTo(w/2-120, h-100)
	ctx.LineTo(w/2+120, h-100)
	ctx.ClosePath()
	ctx.FillPreserve()
	ctx.SetSourceRGB(0, 0, 0)
	ctx.SetLineWidth(2)
	ctx.Stroke()
}

// drawPDFGradientsPage draws linear and radial gradient patterns.
func drawPDFGradientsPage(ctx *cairo.Context, w, h float64) error {
	ctx.SetSourceRGB(1, 1, 1)
	ctx.Paint()

	ctx.SetSourceRGB(0, 0, 0)
	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightBold)
	ctx.SetFontSize(28)
	ctx.MoveTo(50, 60)
	ctx.ShowText("Page 2: Gradients")

	// Horizontal rainbow linear gradient
	linear, err := cairo.NewLinearGradient(50, 0, w-50, 0)
	if err != nil {
		return fmt.Errorf("failed to create linear gradient: %w", err)
	}
	defer func() { _ = linear.Close() }()
	linear.AddColorStopRGB(0.0, 1.0, 0.0, 0.0)
	linear.AddColorStopRGB(0.5, 0.0, 1.0, 0.0)
	linear.AddColorStopRGB(1.0, 0.0, 0.0, 1.0)
	ctx.SetSource(linear)
	ctx.Rectangle(50, 100, w-100, 140)
	ctx.Fill()

	// Radial gradient: warm sphere
	radial, err := cairo.NewRadialGradient(w/4, 390, 0, w/4, 390, 120)
	if err != nil {
		return fmt.Errorf("failed to create radial gradient: %w", err)
	}
	defer func() { _ = radial.Close() }()
	radial.AddColorStopRGB(0.0, 1.0, 1.0, 0.8)
	radial.AddColorStopRGB(0.5, 1.0, 0.5, 0.0)
	radial.AddColorStopRGB(1.0, 0.4, 0.0, 0.0)
	ctx.SetSource(radial)
	ctx.Arc(w/4, 390, 120, 0, 2*math.Pi)
	ctx.Fill()

	// Radial gradient with transparency: blue glow
	radialFade, err := cairo.NewRadialGradient(3*w/4, 390, 0, 3*w/4, 390, 120)
	if err != nil {
		return fmt.Errorf("failed to create fading gradient: %w", err)
	}
	defer func() { _ = radialFade.Close() }()
	radialFade.AddColorStopRGBA(0.0, 0.2, 0.4, 1.0, 1.0)
	radialFade.AddColorStopRGBA(0.6, 0.4, 0.2, 1.0, 0.7)
	radialFade.AddColorStopRGBA(1.0, 0.8, 0.0, 0.5, 0.0)
	ctx.SetSource(radialFade)
	ctx.Arc(3*w/4, 390, 120, 0, 2*math.Pi)
	ctx.Fill()

	// Vertical gradient bar
	vGrad, err := cairo.NewLinearGradient(0, 620, 0, h-60)
	if err != nil {
		return fmt.Errorf("failed to create vertical gradient: %w", err)
	}
	defer func() { _ = vGrad.Close() }()
	vGrad.AddColorStopRGB(0.0, 0.0, 0.6, 0.3)
	vGrad.AddColorStopRGB(1.0, 0.8, 1.0, 0.2)
	ctx.SetSource(vGrad)
	ctx.Rectangle(50, 620, w-100, h-680)
	ctx.Fill()

	return nil
}

// drawPDFTextPage draws text samples showing different fonts, sizes, and styles.
func drawPDFTextPage(ctx *cairo.Context, w float64) {
	ctx.SetSourceRGB(1, 1, 1)
	ctx.Paint()

	ctx.SetSourceRGB(0, 0, 0)
	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightBold)
	ctx.SetFontSize(28)
	ctx.MoveTo(50, 60)
	ctx.ShowText("Page 3: Text")

	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightNormal)
	ctx.SetFontSize(22)
	ctx.MoveTo(50, 140)
	ctx.ShowText("Normal sans-serif, size 22")

	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightBold)
	ctx.MoveTo(50, 210)
	ctx.ShowText("Bold sans-serif, size 22")

	ctx.SetSourceRGB(0.2, 0.2, 0.8)
	ctx.SelectFontFace("serif", font.SlantItalic, font.WeightNormal)
	ctx.MoveTo(50, 280)
	ctx.ShowText("Italic serif, size 22")

	ctx.SetSourceRGB(0.1, 0.5, 0.1)
	ctx.SelectFontFace("monospace", font.SlantOblique, font.WeightNormal)
	ctx.MoveTo(50, 350)
	ctx.ShowText("Oblique monospace, size 22")

	// Large text rendered via TextPath for vector outline
	ctx.SetSourceRGB(0.7, 0.1, 0.1)
	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightBold)
	ctx.SetFontSize(44)
	ctx.MoveTo(50, 460)
	ctx.TextPath("Vector PDF")
	ctx.Fill()

	// Centered text demonstrating alignment
	ctx.SetSourceRGB(0, 0, 0)
	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightNormal)
	ctx.SetFontSize(18)
	te := ctx.TextExtents("Centered text using TextExtents")
	ctx.MoveTo(w/2-te.XBearing-te.Width/2, 560)
	ctx.ShowText("Centered text using TextExtents")
}
