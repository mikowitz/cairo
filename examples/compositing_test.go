// ABOUTME: Tests for the compositing example image generator.
// ABOUTME: Verifies PNG generation and pixel-perfect output via golden image comparison.
package examples

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCompositingGeneratesValidPNG tests that GenerateCompositing creates a valid PNG file.
// This test verifies:
//   - The function executes without error
//   - A PNG file is created at the specified location
//   - The file size is reasonable (between 10KB and 500KB for a 600x720 image)
//   - All 29 compositing operators are demonstrated in a 5×6 grid
func TestCompositingGeneratesValidPNG(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "compositing_test.png")

	err := GenerateCompositing(outputPath)
	require.NoError(t, err, "GenerateCompositing failed")

	require.FileExists(t, outputPath, "Output PNG file was not created")

	// A 600x720 PNG with 29 compositing operator panels should be in this range
	require.True(t, CheckFileSize(t, outputPath, 10000, 500000), "Output PNG file size is not in expected range")

	t.Logf("Successfully generated compositing PNG at %s", outputPath)
}

// TestCompositingMatchesGolden tests that GenerateCompositing produces output
// that matches the golden reference image.
//
// This test uses the CompareImageToGolden harness to verify pixel-perfect output.
// The image demonstrates all 29 compositing operators in a 5×6 grid.
//
// Operators are shown in Cairo's numeric order. The 14 Porter-Duff operators fill
// rows 1–3 with one slot remaining, so Multiply (first blend mode) ends row 3:
//
//   Row 1: Clear, Source, Over, In, Out
//   Row 2: Atop, Dest, DestOver, DestIn, DestOut
//   Row 3: DestAtop, Xor, Add, Saturate, Multiply (first blend mode)
//   Row 4: Screen, Overlay, Darken, Lighten, ColorDodge
//   Row 5: ColorBurn, HardLight, SoftLight, Difference, Exclusion
//   Row 6: HslHue, HslSaturation, HslColor, HslLuminosity
//
// If the test fails, run with -update-golden to regenerate the reference image.
func TestCompositingMatchesGolden(t *testing.T) {
	goldenPath := "testdata/golden/compositing.png"

	match := CompareImageToGolden(t, GenerateCompositing, goldenPath)

	if !match {
		t.Errorf("Generated image does not match golden reference")
		t.Errorf("  Golden path: %s", goldenPath)
		t.Errorf("  To update: go test ./examples -update-golden")
	}
}
