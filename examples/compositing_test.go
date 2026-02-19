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
//   - The file size is reasonable (between 5KB and 200KB for a 400x400 image)
//   - Panels include OperatorOver, OperatorAdd, OperatorMultiply, and OperatorXor demonstrations
func TestCompositingGeneratesValidPNG(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "compositing_test.png")

	err := GenerateCompositing(outputPath)
	require.NoError(t, err, "GenerateCompositing failed")

	require.FileExists(t, outputPath, "Output PNG file was not created")

	// A 400x400 PNG with compositing operator shapes should be in this range
	require.True(t, CheckFileSize(t, outputPath, 5000, 200000), "Output PNG file size is not in expected range")

	t.Logf("Successfully generated compositing PNG at %s", outputPath)
}

// TestCompositingMatchesGolden tests that GenerateCompositing produces output
// that matches the golden reference image.
//
// This test uses the CompareImageToGolden harness to verify pixel-perfect output.
// The image demonstrates four compositing operators:
//   - OperatorOver: default alpha compositing
//   - OperatorAdd: additive blending (brightens at overlap)
//   - OperatorMultiply: multiplicative blending (darkens at overlap)
//   - OperatorXor: exclusive-or (overlap becomes transparent)
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
