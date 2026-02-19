// ABOUTME: Tests for the clipping example image generator.
// ABOUTME: Verifies PNG generation and pixel-perfect output via golden image comparison.
package examples

import (
	"os"
	"path/filepath"
	"testing"
)

// TestClippingGeneratesValidPNG tests that GenerateClipping creates a valid PNG file.
// This test verifies:
//   - The function executes without error
//   - A PNG file is created at the specified location
//   - The file size is reasonable (between 5KB and 200KB for a 600x600 image)
//   - Clipping panels include rectangular, circular, nested, ClipPreserve, Save/Restore,
//     and transformed clip demonstrations
func TestClippingGeneratesValidPNG(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "clipping_test.png")

	err := GenerateClipping(outputPath)
	if err != nil {
		t.Fatalf("GenerateClipping failed: %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Output PNG file was not created")
	}

	// A 600x600 PNG with clipped colorful shapes should be in this range
	if !CheckFileSize(t, outputPath, 5000, 200000) {
		t.Fatal("Output PNG file size is not in expected range")
	}

	t.Logf("Successfully generated clipping PNG at %s", outputPath)
}

// TestClippingMatchesGolden tests that GenerateClipping produces output
// that matches the golden reference image.
//
// This test uses the CompareImageToGolden harness to verify pixel-perfect output.
// The image demonstrates six clipping techniques:
//   - Rectangular clip
//   - Circular clip
//   - ClipPreserve (clip with path retained for stroking)
//   - Nested clips (intersection)
//   - Save/Restore clip state management
//   - Clip in transformed coordinate space
//
// If the test fails, run with -update-golden to regenerate the reference image.
func TestClippingMatchesGolden(t *testing.T) {
	goldenPath := "testdata/golden/clipping.png"

	match := CompareImageToGolden(t, GenerateClipping, goldenPath)

	if !match {
		t.Errorf("Generated image does not match golden reference")
		t.Errorf("  Golden path: %s", goldenPath)
		t.Errorf("  To update: go test ./examples -update-golden")
	}
}
