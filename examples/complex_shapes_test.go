package examples

import (
	"os"
	"path/filepath"
	"testing"
)

// TestComplexShapesGeneratesValidPNG tests that GenerateComplexShapes creates a valid PNG file.
// This test verifies:
//   - The function executes without error
//   - A PNG file is created at the specified location
//   - The file size is reasonable (between 5KB and 200KB for a 600x600 image)
func TestComplexShapesGeneratesValidPNG(t *testing.T) {
	// Use t.TempDir() for automatic cleanup
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "complex_test.png")

	// Generate the image
	err := GenerateComplexShapes(outputPath)
	if err != nil {
		t.Fatalf("GenerateComplexShapes failed: %v", err)
	}

	// Verify the file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Output PNG file was not created")
	}

	// Check file size is reasonable (between 5KB and 200KB)
	// A 600x600 PNG with complex shapes should be in this range
	if !CheckFileSize(t, outputPath, 5000, 200000) {
		t.Fatal("Output PNG file size is not in expected range")
	}

	t.Logf("Successfully generated complex shapes PNG at %s", outputPath)
}

// TestComplexShapesMatchesGolden tests that GenerateComplexShapes produces output
// that matches the golden reference image.
//
// This test uses the CompareImageToGolden harness to verify pixel-perfect output.
// The golden image demonstrates all 20 Context methods:
//   - Paint (background)
//   - Save/Restore (nested rectangles)
//   - SetSourceRGB/SetSourceRGBA (opaque and translucent colors)
//   - MoveTo/LineTo/ClosePath (triangle and star)
//   - GetCurrentPoint/HasCurrentPoint (path queries)
//   - Fill/FillPreserve/Stroke/StrokePreserve (rendering variations)
//   - SetLineWidth/GetLineWidth (line width control)
//   - NewPath/NewSubPath (path management)
//
// If the test fails, run with -update-golden to regenerate the reference image.
func TestComplexShapesMatchesGolden(t *testing.T) {
	goldenPath := "testdata/golden/complex_shapes.png"

	// Use the harness to compare generated image to golden reference
	match := CompareImageToGolden(t, GenerateComplexShapes, goldenPath)

	if !match {
		t.Errorf("Generated image does not match golden reference")
		t.Errorf("  Golden path: %s", goldenPath)
		t.Errorf("  To update: go test ./examples -update-golden")
		t.Errorf("")
		t.Errorf("  This test verifies that all 20 Context methods produce")
		t.Errorf("  consistent output across runs")
	}
}
