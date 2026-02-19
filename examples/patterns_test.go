package examples

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratePatterns(t *testing.T) {
	generator := func(path string) error {
		return GeneratePatterns(path)
	}

	match := CompareImageToGolden(t, generator, "testdata/golden/patterns.png")
	assert.True(t, match, "Generated patterns image should match golden reference")
}

func TestGeneratePatterns_FileCreation(t *testing.T) {
	// Create temporary directory for test output
	tempDir := t.TempDir()
	outputPath := tempDir + "/patterns_test.png"

	// Generate the image
	err := GeneratePatterns(outputPath)
	require.NoError(t, err, "GeneratePatterns should not return an error")

	// Verify the file was created and has reasonable size
	// Pattern images should be substantial but not huge (expect 10KB-500KB range)
	match := CheckFileSize(t, outputPath, 10000, 500000)
	assert.True(t, match, "Generated file should have reasonable size")
}
