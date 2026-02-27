// ABOUTME: Tests for the SVG output example verifying file creation and SVG validity.
// ABOUTME: Uses structural tests since SVG text rendering is platform-dependent.

//go:build !nosvg

package examples

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestSVGOutputGeneratesFile verifies that GenerateSVGOutput creates a file at the
// specified path and that the file size is within a reasonable range for an SVG
// containing shapes, gradients, and text.
func TestSVGOutputGeneratesFile(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "output.svg")

	err := GenerateSVGOutput(outputPath)
	require.NoError(t, err, "GenerateSVGOutput failed")

	require.FileExists(t, outputPath, "Output SVG file was not created")
	require.True(t,
		CheckFileSize(t, outputPath, 1000, 5_000_000),
		"Output SVG file size is not in expected range",
	)
}

// TestSVGOutputValidXML verifies that the generated SVG is well-formed XML that
// can be parsed by a standard XML parser.
func TestSVGOutputValidXML(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "output.svg")

	err := GenerateSVGOutput(outputPath)
	require.NoError(t, err, "GenerateSVGOutput failed")

	data, err := os.ReadFile(outputPath) //nolint:gosec // path is from t.TempDir(), not user input
	require.NoError(t, err, "failed to read SVG file")

	var v interface{}
	xmlErr := xml.Unmarshal(data, &v)
	require.NoError(t, xmlErr, "SVG file should be well-formed XML")
}

// TestSVGOutputContainsSVGElement verifies that the generated file contains the
// SVG root element, confirming Cairo produced a proper SVG document.
func TestSVGOutputContainsSVGElement(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "output.svg")

	err := GenerateSVGOutput(outputPath)
	require.NoError(t, err, "GenerateSVGOutput failed")

	data, err := os.ReadFile(outputPath) //nolint:gosec // path is from t.TempDir(), not user input
	require.NoError(t, err, "failed to read SVG file")

	require.Contains(t, string(data), "<svg", "SVG file should contain svg root element")
}
