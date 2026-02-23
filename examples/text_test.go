// ABOUTME: Tests for the text rendering example using structural pixel checks.
// ABOUTME: Uses structural tests instead of golden comparison due to platform-dependent rendering.
package examples

import (
	"image"
	"image/color"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestTextGeneratesValidPNG tests that GenerateText creates a valid PNG file.
// This test verifies:
//   - The function executes without error
//   - A PNG file is created at the specified location
//   - The file size is reasonable (between 1KB and 200KB for a 400x300 image)
func TestTextGeneratesValidPNG(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "text_test.png")

	err := GenerateText(outputPath)
	require.NoError(t, err, "GenerateText failed")

	require.FileExists(t, outputPath, "Output PNG file was not created")
	require.True(t, CheckFileSize(t, outputPath, 1000, 200000), "Output PNG file size is not in expected range")

	t.Logf("Successfully generated text PNG at %s", outputPath)
}

// TestTextRendersVisibleContent verifies that each text row contains non-background
// pixels, confirming that Cairo rendered actual text glyphs.
//
// This test uses structural pixel checks instead of golden image comparison
// because text rendering is platform-dependent: different operating systems
// use different fonts, hinting engines, and antialiasing strategies, producing
// differences far too large for pixel-level comparison.
func TestTextRendersVisibleContent(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "text_structural.png")

	err := GenerateText(outputPath)
	require.NoError(t, err, "GenerateText failed")

	img, err := decodePNG(outputPath)
	require.NoError(t, err, "failed to decode generated PNG")

	rows := []struct {
		name   string
		y0, y1 int
	}{
		{"normal sans-serif (row 1)", 30, 60},
		{"bold sans-serif (row 2)", 80, 110},
		{"italic serif (row 3)", 130, 160},
		{"oblique monospace (row 4)", 180, 210},
		{"large bold path (row 5)", 240, 285},
	}

	for _, row := range rows {
		row := row
		t.Run(row.name, func(t *testing.T) {
			require.True(t,
				regionHasNonBackgroundPixels(img, 20, row.y0, 350, row.y1),
				"text region %q should contain rendered glyphs (non-white pixels)", row.name,
			)
		})
	}
}

// regionHasNonBackgroundPixels returns true if any pixel in the rectangle
// [x0,x1] × [y0,y1] differs significantly from white (the background color).
func regionHasNonBackgroundPixels(img image.Image, x0, y0, x1, y1 int) bool {
	bounds := img.Bounds()
	for y := y0; y <= y1 && y < bounds.Max.Y; y++ {
		for x := x0; x <= x1 && x < bounds.Max.X; x++ {
			px := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			if px.R < 250 || px.G < 250 || px.B < 250 {
				return true
			}
		}
	}
	return false
}
