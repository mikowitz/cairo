package examples

import (
	"testing"
)

// TestGenerateTransformations verifies that the transformations example generates
// the expected output by comparing it to a golden reference image.
//
// This test demonstrates:
//   - Translation: moving shapes to different positions
//   - Scaling: making shapes larger or smaller
//   - Rotation: rotating shapes by angles in radians
//   - Combined transformations: applying multiple transformations together
//   - Save/Restore: preserving and restoring transformation state
//
// To update the golden reference image, run:
//
//	go test -update-golden ./examples
func TestGenerateTransformations(t *testing.T) {
	generator := func(path string) error {
		return GenerateTransformations(path)
	}

	goldenPath := "testdata/golden/transformations.png"
	match := CompareImageToGolden(t, generator, goldenPath)

	if !match {
		t.Error("Generated transformations image does not match golden reference")
		t.Log("Visual differences may indicate:")
		t.Log("  - Translation not working correctly")
		t.Log("  - Scaling producing wrong dimensions")
		t.Log("  - Rotation angle calculations incorrect")
		t.Log("  - Transformation order/composition issues")
		t.Log("  - Save/Restore not preserving state properly")
	}
}

// TestTransformationsFileSize performs a sanity check that the generated
// transformations image has a reasonable file size.
func TestTransformationsFileSize(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := tempDir + "/transformations.png"

	if err := GenerateTransformations(outputPath); err != nil {
		t.Fatalf("Failed to generate transformations image: %v", err)
	}

	// Expect file size between 5KB and 100KB for a 600x600 image with shapes
	// Actual size will depend on PNG compression and content complexity
	if !CheckFileSize(t, outputPath, 5000, 100000) {
		t.Error("Transformations image file size is outside expected range")
	}
}
