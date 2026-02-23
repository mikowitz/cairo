// ABOUTME: Tests for the fill rules example image generator.
// ABOUTME: Verifies PNG generation and pixel-perfect output via golden image comparison.
package examples

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestFillRulesGeneratesValidPNG tests that GenerateFillRules creates a valid PNG file.
// This test verifies:
//   - The function executes without error
//   - A PNG file is created at the specified location
//   - The file size is reasonable (between 1KB and 100KB for a 400x200 image)
//   - Both winding and even-odd fill rule panels are demonstrated
func TestFillRulesGeneratesValidPNG(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "fill_rules_test.png")

	err := GenerateFillRules(outputPath)
	require.NoError(t, err, "GenerateFillRules failed")

	require.FileExists(t, outputPath, "Output PNG file was not created")

	// A 400x200 PNG with two star panels should be in this range
	require.True(t, CheckFileSize(t, outputPath, 1000, 100000), "Output PNG file size is not in expected range")

	t.Logf("Successfully generated fill rules PNG at %s", outputPath)
}

// TestFillRulesMatchesGolden tests that GenerateFillRules produces output
// that matches the golden reference image.
//
// This test uses the CompareImageToGolden harness to verify pixel-perfect output.
// The image demonstrates the two Cairo fill rules:
//   - Left panel: FillRuleWinding fills the entire star interior including the center pentagon
//   - Right panel: FillRuleEvenOdd leaves the center pentagon unfilled (even crossing count)
//
// If the test fails, run with -update-golden to regenerate the reference image.
func TestFillRulesMatchesGolden(t *testing.T) {
	goldenPath := "testdata/golden/fill_rules.png"

	match := CompareImageToGolden(t, GenerateFillRules, goldenPath)

	if !match {
		t.Errorf("Generated image does not match golden reference")
		t.Errorf("  Golden path: %s", goldenPath)
		t.Errorf("  To update: go test ./examples -update-golden")
	}
}
