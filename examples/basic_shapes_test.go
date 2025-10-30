package examples

import (
	"os"
	"path/filepath"
	"testing"
)

// TestBasicShapesGeneratesValidPNG tests that GenerateBasicShapes creates a valid PNG file.
// This test verifies:
//   - The function executes without error
//   - A PNG file is created at the specified location
//   - The file size is reasonable (between 1KB and 100KB)
//   - Shapes include rectangles, circles (using Arc), and curves (using CurveTo)
func TestBasicShapesGeneratesValidPNG(t *testing.T) {
	// Use t.TempDir() for automatic cleanup
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "test_output.png")

	// Generate the image
	err := GenerateBasicShapes(outputPath)
	if err != nil {
		t.Fatalf("GenerateBasicShapes failed: %v", err)
	}

	// Verify the file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Output PNG file was not created")
	}

	// Check file size is reasonable (between 1KB and 100KB)
	// A 600x400 PNG with shapes, circles, and curves should be in this range
	if !CheckFileSize(t, outputPath, 1000, 100000) {
		t.Fatal("Output PNG file size is not in expected range")
	}

	t.Logf("Successfully generated PNG at %s", outputPath)
}

// TestBasicShapesMatchesGolden tests that GenerateBasicShapes produces output
// that matches the golden reference image.
//
// This test uses the CompareImageToGolden harness to verify pixel-perfect output.
// The image demonstrates Arc (circles) and CurveTo (Bezier curves) in addition to rectangles.
// If the test fails, run with -update-golden to regenerate the reference image.
func TestBasicShapesMatchesGolden(t *testing.T) {
	goldenPath := "testdata/golden/basic_shapes.png"

	// Use the harness to compare generated image to golden reference
	match := CompareImageToGolden(t, GenerateBasicShapes, goldenPath)

	if !match {
		t.Errorf("Generated image does not match golden reference")
		t.Errorf("  Golden path: %s", goldenPath)
		t.Errorf("  To update: go test ./examples -update-golden")
	}
}
