// ABOUTME: Example demonstrating a data dashboard with bar chart, line chart, and pie chart.
// ABOUTME: Uses gradients, text, and transformations to render to PNG and (via dashboard_pdf.go) PDF.

package examples

import (
	"fmt"
	"math"

	"github.com/mikowitz/cairo"
	"github.com/mikowitz/cairo/font"
)

// GenerateDashboard creates an 800×600 PNG data dashboard with three chart types.
//
// The dashboard contains:
//   - Gradient header with title
//   - Bar chart (upper-left): monthly sales with per-bar gradient fills
//   - Line chart (upper-right): weekly trend with shaded area under the line
//   - Pie chart (bottom): category breakdown with radial gradient slices and legend
//
// All Cairo resources are cleaned up with defer statements.
func GenerateDashboard(outputPath string) error {
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 800, 600)
	if err != nil {
		return fmt.Errorf("failed to create surface: %w", err)
	}
	defer func() { _ = surface.Close() }()

	ctx, err := cairo.NewContext(surface)
	if err != nil {
		return fmt.Errorf("failed to create context: %w", err)
	}
	defer func() { _ = ctx.Close() }()

	drawDashboard(ctx, 800, 600)

	surface.Flush()
	if err := surface.WriteToPNG(outputPath); err != nil {
		return fmt.Errorf("failed to write PNG: %w", err)
	}
	return nil
}

// drawDashboard renders the full dashboard onto ctx using the supplied dimensions.
// Keeping the drawing logic separate allows both PNG and PDF surfaces to reuse it.
func drawDashboard(ctx *cairo.Context, w, h float64) { //nolint:funlen
	// Light gray background
	ctx.SetSourceRGB(0.94, 0.94, 0.96)
	ctx.Paint()

	// Gradient title header: blue sweep left-to-right
	hdr, _ := cairo.NewLinearGradient(0, 0, w, 0)
	hdr.AddColorStopRGB(0, 0.10, 0.38, 0.78)
	hdr.AddColorStopRGB(1, 0.28, 0.58, 0.90)
	ctx.SetSource(hdr)
	ctx.Rectangle(0, 0, w, h*0.10)
	ctx.Fill()
	_ = hdr.Close()

	// White title text
	ctx.SetSourceRGB(1, 1, 1)
	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightBold)
	ctx.SetFontSize(h * 0.038)
	ctx.MoveTo(w*0.025, h*0.068)
	ctx.ShowText("Data Dashboard")

	// Two side-by-side panels in the upper 50% of the canvas
	drawBarChart(ctx, w*0.025, h*0.12, w*0.46, h*0.43)
	drawLineChart(ctx, w*0.525, h*0.12, w*0.45, h*0.43)

	// Pie chart centered in the lower section
	drawPieChart(ctx, w*0.32, h*0.77, h*0.18)
}

// pieSegment holds the display attributes of a single pie slice.
type pieSegment struct {
	label   string
	frac    float64
	r, g, b float64
}

// drawBarChart renders a vertical bar chart inside a white panel at (x, y, cw, ch).
// Each bar uses a vertical gradient fill from a lighter shade to the base color.
func drawBarChart(ctx *cairo.Context, x, y, cw, ch float64) { //nolint:funlen
	// White panel
	ctx.SetSourceRGB(1, 1, 1)
	ctx.Rectangle(x, y, cw, ch)
	ctx.Fill()

	ctx.SetSourceRGB(0.20, 0.20, 0.20)
	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightBold)
	ctx.SetFontSize(ch * 0.075)
	ctx.MoveTo(x+cw*0.04, y+ch*0.10)
	ctx.ShowText("Monthly Sales")

	labels := []string{"Jan", "Feb", "Mar", "Apr", "May"}
	values := []float64{0.60, 0.85, 0.45, 0.90, 0.70}
	colors := [][3]float64{
		{0.20, 0.50, 0.88},
		{0.25, 0.68, 0.40},
		{0.90, 0.60, 0.18},
		{0.80, 0.28, 0.38},
		{0.50, 0.28, 0.80},
	}

	n := float64(len(values))
	barW := cw * 0.12
	gap := (cw - barW) / n
	chartTop := y + ch*0.17
	chartH := ch * 0.70

	for i, v := range values {
		bx := x + gap*float64(i) + gap*0.25
		bh := chartH * v
		by := chartTop + chartH - bh
		c := colors[i]

		// Gradient fill: lighter highlight at top, base color at bottom
		grad, _ := cairo.NewLinearGradient(bx, by, bx, by+bh)
		grad.AddColorStopRGB(0, clamp1(c[0]+0.25), clamp1(c[1]+0.25), clamp1(c[2]+0.25))
		grad.AddColorStopRGB(1, c[0], c[1], c[2])
		ctx.SetSource(grad)
		ctx.Rectangle(bx, by, barW, bh)
		ctx.Fill()
		_ = grad.Close()

		// Dark outline for contrast
		ctx.SetSourceRGB(c[0]*0.65, c[1]*0.65, c[2]*0.65)
		ctx.SetLineWidth(1)
		ctx.Rectangle(bx, by, barW, bh)
		ctx.Stroke()

		// Month label below the bar
		ctx.SetSourceRGB(0.30, 0.30, 0.30)
		ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightNormal)
		ctx.SetFontSize(ch * 0.062)
		ctx.MoveTo(bx+barW*0.15, chartTop+chartH+ch*0.07)
		ctx.ShowText(labels[i])
	}
}

// drawLineChart renders a line chart with a semi-transparent shaded area beneath the line.
// Data points are marked with filled circles containing white centers.
func drawLineChart(ctx *cairo.Context, x, y, cw, ch float64) {
	// White panel
	ctx.SetSourceRGB(1, 1, 1)
	ctx.Rectangle(x, y, cw, ch)
	ctx.Fill()

	ctx.SetSourceRGB(0.20, 0.20, 0.20)
	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightBold)
	ctx.SetFontSize(ch * 0.075)
	ctx.MoveTo(x+cw*0.04, y+ch*0.10)
	ctx.ShowText("Weekly Trend")

	data := []float64{0.30, 0.55, 0.40, 0.75, 0.60, 0.85, 0.70}
	chartTop := y + ch*0.17
	chartH := ch * 0.70
	ox := x + cw*0.05
	ptW := (cw * 0.90) / float64(len(data)-1)

	// Shaded area: vertical gradient from semi-opaque to transparent
	shade, _ := cairo.NewLinearGradient(ox, chartTop, ox, chartTop+chartH)
	shade.AddColorStopRGBA(0, 0.18, 0.52, 0.88, 0.40)
	shade.AddColorStopRGBA(1, 0.18, 0.52, 0.88, 0.04)
	ctx.SetSource(shade)
	ctx.MoveTo(ox, chartTop+chartH)
	for i, v := range data {
		ctx.LineTo(ox+float64(i)*ptW, chartTop+chartH*(1-v))
	}
	ctx.LineTo(ox+float64(len(data)-1)*ptW, chartTop+chartH)
	ctx.ClosePath()
	ctx.Fill()
	_ = shade.Close()

	// Trend line
	ctx.SetSourceRGB(0.10, 0.40, 0.80)
	ctx.SetLineWidth(2.5)
	ctx.MoveTo(ox, chartTop+chartH*(1-data[0]))
	for i := 1; i < len(data); i++ {
		ctx.LineTo(ox+float64(i)*ptW, chartTop+chartH*(1-data[i]))
	}
	ctx.Stroke()

	// Markers: blue dot with white center
	for i, v := range data {
		px := ox + float64(i)*ptW
		py := chartTop + chartH*(1-v)
		ctx.SetSourceRGB(0.10, 0.40, 0.80)
		ctx.Arc(px, py, ch*0.022, 0, 2*math.Pi)
		ctx.Fill()
		ctx.SetSourceRGB(1, 1, 1)
		ctx.Arc(px, py, ch*0.010, 0, 2*math.Pi)
		ctx.Fill()
	}
}

// drawPieChart renders a four-segment pie chart centered at (cx, cy) with radius r.
// Each slice has a radial gradient and is bordered in white. A legend is drawn to
// the right listing each category label.
func drawPieChart(ctx *cairo.Context, cx, cy, r float64) { //nolint:funlen
	ctx.SetSourceRGB(0.20, 0.20, 0.20)
	ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightBold)
	ctx.SetFontSize(r * 0.40)
	ctx.MoveTo(cx-r*0.60, cy-r*1.20)
	ctx.ShowText("Category Split")

	segs := []pieSegment{
		{"A – 35%", 0.35, 0.20, 0.50, 0.88},
		{"B – 25%", 0.25, 0.25, 0.68, 0.40},
		{"C – 20%", 0.20, 0.90, 0.60, 0.18},
		{"D – 20%", 0.20, 0.80, 0.28, 0.38},
	}

	angle := -math.Pi / 2 // start at 12 o'clock
	for i, s := range segs {
		end := angle + 2*math.Pi*s.frac

		// Radial gradient: bright at offset highlight, solid at rim
		grad, _ := cairo.NewRadialGradient(cx-r*0.18, cy-r*0.18, 0, cx, cy, r)
		grad.AddColorStopRGB(0, clamp1(s.r+0.30), clamp1(s.g+0.30), clamp1(s.b+0.30))
		grad.AddColorStopRGB(1, s.r, s.g, s.b)
		ctx.SetSource(grad)
		ctx.MoveTo(cx, cy)
		ctx.Arc(cx, cy, r, angle, end)
		ctx.ClosePath()
		ctx.Fill()
		_ = grad.Close()

		// White slice border
		ctx.SetSourceRGB(1, 1, 1)
		ctx.SetLineWidth(2)
		ctx.MoveTo(cx, cy)
		ctx.Arc(cx, cy, r, angle, end)
		ctx.ClosePath()
		ctx.Stroke()

		// Legend: color swatch + label stacked vertically to the right
		lx := cx + r*1.30
		ly := cy - r*0.55 + float64(i)*r*0.45
		ctx.SetSourceRGB(s.r, s.g, s.b)
		ctx.Rectangle(lx, ly-r*0.14, r*0.24, r*0.24)
		ctx.Fill()

		ctx.SetSourceRGB(0.20, 0.20, 0.20)
		ctx.SelectFontFace("sans-serif", font.SlantNormal, font.WeightNormal)
		ctx.SetFontSize(r * 0.28)
		ctx.MoveTo(lx+r*0.32, ly+r*0.06)
		ctx.ShowText(s.label)

		angle = end
	}
}

// clamp1 returns v clamped to [0, 1] for use in gradient color stops.
func clamp1(v float64) float64 {
	if v > 1 {
		return 1
	}
	return v
}
