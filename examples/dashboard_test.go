// ABOUTME: Tests for the dashboard example using structural pixel checks.
// ABOUTME: Uses structural tests instead of golden comparison due to text rendering.
package examples

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestDashboardGeneratesValidPNG tests that GenerateDashboard creates a valid PNG file.
// It verifies the function succeeds and the output is a reasonably sized PNG.
func TestDashboardGeneratesValidPNG(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "dashboard.png")

	err := GenerateDashboard(outputPath)
	require.NoError(t, err, "GenerateDashboard failed")
	require.FileExists(t, outputPath, "output PNG was not created")
	require.True(t, CheckFileSize(t, outputPath, 5000, 500000), "PNG file size out of expected range")

	t.Logf("Generated dashboard PNG at %s", outputPath)
}

// TestDashboardRendersVisibleContent verifies that each dashboard section contains
// non-background pixels, confirming that the charts were actually rendered.
//
// Structural checks are used instead of golden image comparison because the
// dashboard includes text labels whose rendering is platform-dependent.
func TestDashboardRendersVisibleContent(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "dashboard_structural.png")

	err := GenerateDashboard(outputPath)
	require.NoError(t, err, "GenerateDashboard failed")

	img, err := decodePNG(outputPath)
	require.NoError(t, err, "failed to decode PNG")

	regions := []struct {
		name       string
		x0, y0     int
		x1, y1     int
	}{
		// Header gradient band across the top
		{"header gradient", 0, 0, 800, 60},
		// Bar chart panel (upper-left quadrant)
		{"bar chart panel", 20, 72, 390, 330},
		// Line chart panel (upper-right quadrant)
		{"line chart panel", 420, 72, 790, 330},
		// Pie chart area (lower center)
		{"pie chart", 150, 370, 550, 590},
	}

	for _, r := range regions {
		r := r
		t.Run(r.name, func(t *testing.T) {
			require.True(t,
				RegionHasNonBackgroundPixels(img, r.x0, r.y0, r.x1, r.y1),
				"region %q should contain rendered content", r.name,
			)
		})
	}
}
