// ABOUTME: Example demonstrating text measurement with TextExtents and FontExtents.
// ABOUTME: Shows center/right alignment, multi-line spacing, and ink bounding box drawing.
package examples

import (
	"fmt"
	"math"

	"github.com/mikowitz/cairo"
	"github.com/mikowitz/cairo/font"
)

// GenerateTextExtents creates a 500x400 PNG demonstrating Cairo's text measurement API.
//
// The image is divided into three sections:
//   - Alignment: left, center, and right-aligned text relative to a vertical guide at x=250
//   - Multi-line: evenly-spaced lines using FontExtents.Height for consistent line spacing
//   - Bounding box: text with its ink bounding box drawn as a red outline and baseline dot
//
// All positioning uses TextExtents and FontExtents rather than hard-coded offsets.
// Resources are cleaned up with defer statements.
func GenerateTextExtents(outputPath string) error {
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 500, 400)
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

	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightNormal)

	// --- Section 1: Text alignment ---
	//
	// A faint vertical guide line at x=250 marks the alignment axis.
	// Left-aligned text starts at x=20; centered text is positioned so its
	// ink midpoint falls on x=250; right-aligned text ends its advance at x=480.

	ctx.SetSourceRGBA(0.75, 0.75, 0.75, 1.0)
	ctx.SetLineWidth(1.0)
	ctx.MoveTo(250, 20)
	ctx.LineTo(250, 178)
	ctx.Stroke()

	ctx.SetSourceRGB(0, 0, 0)
	ctx.SetFontSize(11.0)
	ctx.MoveTo(5, 15)
	ctx.ShowText("Alignment  (guide at x=250)")

	ctx.SetFontSize(16.0)

	// Left-aligned
	ctx.MoveTo(20, 65)
	ctx.ShowText("Left aligned")

	// Center-aligned: shift so ink midpoint lands on x=250
	centerText := "Centered text"
	te := ctx.TextExtents(centerText)
	ctx.MoveTo(250-te.XBearing-te.Width/2, 110)
	ctx.ShowText(centerText)

	// Right-aligned: shift so advance end lands at x=480
	rightText := "Right aligned"
	te = ctx.TextExtents(rightText)
	ctx.MoveTo(480-te.XAdvance, 155)
	ctx.ShowText(rightText)

	// --- Section 2: Multi-line text with FontExtents.Height ---
	//
	// FontExtents.Height gives the recommended baseline-to-baseline distance.
	// Each successive line's y is incremented by that value.

	ctx.SetSourceRGB(0, 0, 0)
	ctx.SetFontSize(11.0)
	ctx.MoveTo(5, 198)
	ctx.ShowText("Multi-line spacing  (FontExtents.Height)")

	ctx.SetFontSize(16.0)
	fe := ctx.FontExtents()

	multiLines := []string{
		"First line of text",
		"Second line of text",
		"Third line of text",
		"Fourth line of text",
	}
	startY := 222.0
	for i, line := range multiLines {
		ctx.MoveTo(20, startY+float64(i)*fe.Height)
		ctx.ShowText(line)
	}

	// --- Section 3: Ink bounding box around text ---
	//
	// TextExtents.XBearing / YBearing give the offset from the origin (MoveTo
	// point) to the top-left corner of the ink box. Width and Height give its
	// dimensions. A small red dot marks the baseline origin.

	ctx.SetSourceRGB(0, 0, 0)
	ctx.SetFontSize(11.0)
	ctx.MoveTo(5, 328)
	ctx.ShowText("Bounding box  (TextExtents ink bounds)")

	ctx.SetFontSize(22.0)
	boxText := "Bounding Box"
	bx, by := 20.0, 375.0
	te = ctx.TextExtents(boxText)
	ctx.MoveTo(bx, by)
	ctx.ShowText(boxText)

	// Red outline of the ink bounding box
	ctx.SetSourceRGB(1, 0, 0)
	ctx.SetLineWidth(1.5)
	ctx.Rectangle(bx+te.XBearing, by+te.YBearing, te.Width, te.Height)
	ctx.Stroke()

	// Red dot at the baseline origin (MoveTo point)
	ctx.Arc(bx, by, 3, 0, 2*math.Pi)
	ctx.Fill()

	surface.Flush()
	return surface.WriteToPNG(outputPath)
}
