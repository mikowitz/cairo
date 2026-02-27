// ABOUTME: PDF variant of the dashboard example using Cairo's PDF surface.
// ABOUTME: Reuses drawDashboard to render the same layout to a vector PDF file.

//go:build !nopdf

package examples

import (
	"fmt"

	"github.com/mikowitz/cairo"
)

// GenerateDashboardPDF renders the data dashboard to a PDF file.
//
// The page uses A4 landscape dimensions (841.89 × 595.28 points, i.e., 297 × 210 mm).
// The drawing content is identical to GenerateDashboard but output as vector graphics,
// making the PDF suitable for printing at any resolution.
func GenerateDashboardPDF(outputPath string) error {
	const (
		pageW = 841.89 // A4 landscape width in points
		pageH = 595.28 // A4 landscape height in points
	)

	surface, err := cairo.NewPDFSurface(outputPath, pageW, pageH)
	if err != nil {
		return fmt.Errorf("failed to create PDF surface: %w", err)
	}
	defer func() { _ = surface.Close() }()

	ctx, err := cairo.NewContext(surface)
	if err != nil {
		return fmt.Errorf("failed to create context: %w", err)
	}
	defer func() { _ = ctx.Close() }()

	drawDashboard(ctx, pageW, pageH)
	surface.Flush()
	return nil
}
