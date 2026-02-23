// ABOUTME: Example demonstrating Cairo's toy font API for basic text rendering.
// ABOUTME: Shows text with different fonts, sizes, slants, weights, and positions.
package examples

import (
	"fmt"

	"github.com/mikowitz/cairo"
	"github.com/mikowitz/cairo/font"
)

// GenerateText creates a 400x300 PNG image demonstrating the toy font API.
//
// The image shows five rows of text, each demonstrating different font
// properties using Cairo's toy font API:
//   - Row 1: Normal weight, upright sans-serif at size 20
//   - Row 2: Bold weight sans-serif at size 20
//   - Row 3: Italic slant serif at size 20
//   - Row 4: Oblique slant monospace at size 20
//   - Row 5: Large bold text via TextPath (filled outline), size 30
//
// Text is positioned at the current point via ctx.MoveTo before each call.
// The toy font API results are platform-dependent; exact glyph rendering
// varies by operating system and installed fonts.
//
// All resources are properly cleaned up using defer statements.
func GenerateText(outputPath string) error {
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 400, 300)
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

	// White background
	ctx.SetSourceRGB(1, 1, 1)
	ctx.Paint()
	ctx.SetSourceRGB(0, 0, 0)

	// Row 1: Normal upright sans-serif at size 20
	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightNormal)
	ctx.SetFontSize(20.0)
	ctx.MoveTo(20, 50)
	ctx.ShowText("Normal sans-serif")

	// Row 2: Bold sans-serif at size 20
	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightBold)
	ctx.SetFontSize(20.0)
	ctx.MoveTo(20, 100)
	ctx.ShowText("Bold sans-serif")

	// Row 3: Italic serif at size 20
	ctx.SelectFontFace("serif", font.SlantItalic, font.WeightNormal)
	ctx.SetFontSize(20.0)
	ctx.MoveTo(20, 150)
	ctx.ShowText("Italic serif")

	// Row 4: Oblique monospace at size 20
	ctx.SelectFontFace("monospace", font.SlantOblique, font.WeightNormal)
	ctx.SetFontSize(20.0)
	ctx.MoveTo(20, 200)
	ctx.ShowText("Oblique monospace")

	// Row 5: TextPath builds glyph outlines then Fill renders them
	ctx.SetSourceRGB(0.2, 0.2, 0.8)
	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightBold)
	ctx.SetFontSize(30.0)
	ctx.MoveTo(20, 270)
	ctx.TextPath("Large bold path")
	ctx.Fill()

	surface.Flush()
	return surface.WriteToPNG(outputPath)
}
