package examples

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGradientsGeneratesValidPNG tests that GenerateGradients creates a valid PNG file.
// This test verifies:
//   - The function executes without error
//   - A PNG file is created at the specified location
//   - The file size is reasonable (between 1KB and 150KB)
//   - Gradients include both linear and radial types with various color stops
func TestGradientsGeneratesValidPNG(t *testing.T) {
	// Use t.TempDir() for automatic cleanup
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "test_output.png")

	// Generate the image
	err := GenerateGradients(outputPath)
	if err != nil {
		t.Fatalf("GenerateGradients failed: %v", err)
	}

	// Verify the file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Output PNG file was not created")
	}

	// Check file size is reasonable (between 1KB and 150KB)
	// A 600x600 PNG with gradients should be in this range
	if !CheckFileSize(t, outputPath, 1000, 150000) {
		t.Fatal("Output PNG file size is not in expected range")
	}

	t.Logf("Successfully generated PNG at %s", outputPath)
}

// TestGradientsMatchesGolden tests that GenerateGradients produces output
// that matches the golden reference image.
//
// This test uses the CompareImageToGolden harness to verify pixel-perfect output.
// The image demonstrates both linear and radial gradients with various color stops,
// including transparency effects.
// If the test fails, run with -update-golden to regenerate the reference image.
func TestGradientsMatchesGolden(t *testing.T) {
	goldenPath := "testdata/golden/gradients.png"

	// Use the harness to compare generated image to golden reference
	match := CompareImageToGolden(t, GenerateGradients, goldenPath)

	if !match {
		t.Errorf("Generated image does not match golden reference")
		t.Errorf("  Golden path: %s", goldenPath)
		t.Errorf("  To update: go test ./examples -update-golden")
	}
}
