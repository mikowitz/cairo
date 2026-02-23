// ABOUTME: Tests for the text extents example using structural pixel checks.
// ABOUTME: Uses structural tests instead of golden comparison due to platform-dependent text rendering.
package examples

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestTextExtentsGeneratesValidPNG tests that GenerateTextExtents creates a valid PNG file.
func TestTextExtentsGeneratesValidPNG(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "text_extents_test.png")

	err := GenerateTextExtents(outputPath)
	require.NoError(t, err, "GenerateTextExtents failed")

	require.FileExists(t, outputPath, "Output PNG file was not created")
	require.True(t, CheckFileSize(t, outputPath, 1000, 200000), "Output PNG file size is not in expected range")

	t.Logf("Successfully generated text extents PNG at %s", outputPath)
}

// TestTextExtentsRendersVisibleContent verifies that each section in the output image
// contains non-background pixels, confirming actual content was rendered.
//
// Structural pixel checks are used instead of golden image comparison because
// text rendering is platform-dependent across operating systems.
func TestTextExtentsRendersVisibleContent(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "text_extents_structural.png")

	err := GenerateTextExtents(outputPath)
	require.NoError(t, err, "GenerateTextExtents failed")

	img, err := decodePNG(outputPath)
	require.NoError(t, err, "failed to decode generated PNG")

	sections := []struct {
		name       string
		x0, y0, x1, y1 int
	}{
		// Alignment section: left-aligned text starts near x=20, baseline y≈65
		{"left-aligned text", 20, 45, 250, 75},
		// Alignment section: centered text spans around x=250, baseline y≈110
		{"centered text", 80, 90, 420, 120},
		// Alignment section: right-aligned text ends near x=480, baseline y≈155
		{"right-aligned text", 250, 135, 490, 165},
		// Multi-line section: first line near y=222
		{"multi-line row 1", 20, 202, 350, 235},
		// Multi-line section: fourth line is roughly 3 line-heights below first
		{"multi-line row 4", 20, 260, 350, 310},
		// Bounding box section: text near y=375
		{"bounding box text", 20, 350, 350, 390},
	}

	for _, sec := range sections {
		sec := sec
		t.Run(sec.name, func(t *testing.T) {
			require.True(t,
				RegionHasNonBackgroundPixels(img, sec.x0, sec.y0, sec.x1, sec.y1),
				"section %q should contain rendered content (non-white pixels)", sec.name,
			)
		})
	}
}
