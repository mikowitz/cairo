package examples

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLineStylesGeneratesValidPNG tests that GenerateLineStyles creates a valid PNG file.
// This test verifies:
//   - The function executes without error
//   - A PNG file is created at the specified location
//   - The file size is reasonable (between 1KB and 150KB)
//   - Line styles include caps, joins, dash patterns, and miter limits
func TestLineStylesGeneratesValidPNG(t *testing.T) {
	// Use t.TempDir() for automatic cleanup
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "test_output.png")

	// Generate the image
	err := GenerateLineStyles(outputPath)
	if err != nil {
		t.Fatalf("GenerateLineStyles failed: %v", err)
	}

	// Verify the file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Output PNG file was not created")
	}

	// Check file size is reasonable (between 1KB and 150KB)
	// A 700x600 PNG with line styles should be in this range
	if !CheckFileSize(t, outputPath, 1000, 150000) {
		t.Fatal("Output PNG file size is not in expected range")
	}

	t.Logf("Successfully generated PNG at %s", outputPath)
}

// TestLineStylesMatchesGolden tests that GenerateLineStyles produces output
// that matches the golden reference image.
//
// This test uses the CompareImageToGolden harness to verify pixel-perfect output.
// The image demonstrates line caps (Butt, Round, Square), line joins (Miter, Round, Bevel),
// various dash patterns (solid, dashed, dotted, complex, offset), and miter limit effects.
// If the test fails, run with -update-golden to regenerate the reference image.
func TestLineStylesMatchesGolden(t *testing.T) {
	goldenPath := "testdata/golden/line_styles.png"

	// Use the harness to compare generated image to golden reference
	match := CompareImageToGolden(t, GenerateLineStyles, goldenPath)

	if !match {
		t.Errorf("Generated image does not match golden reference")
		t.Errorf("  Golden path: %s", goldenPath)
		t.Errorf("  To update: go test ./examples -update-golden")
	}
}
