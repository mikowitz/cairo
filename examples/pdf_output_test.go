// ABOUTME: Tests for the PDF output example verifying file creation and PDF validity.
// ABOUTME: Uses structural tests since PDF content cannot be decoded like PNG images.

//go:build !nopdf

package examples

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestPDFOutputGeneratesFile verifies that GeneratePDFOutput creates a file at the
// specified path and that the file size is within a reasonable range for a 3-page PDF.
func TestPDFOutputGeneratesFile(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "output.pdf")

	err := GeneratePDFOutput(outputPath)
	require.NoError(t, err, "GeneratePDFOutput failed")

	require.FileExists(t, outputPath, "Output PDF file was not created")
	require.True(t,
		CheckFileSize(t, outputPath, 1000, 5_000_000),
		"Output PDF file size is not in expected range",
	)

	t.Logf("Successfully generated PDF at %s", outputPath)
}

// TestPDFOutputValidPDFHeader verifies that the generated file is a valid PDF by
// checking for the required "%PDF-" magic bytes at the start of the file.
func TestPDFOutputValidPDFHeader(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "output.pdf")

	err := GeneratePDFOutput(outputPath)
	require.NoError(t, err, "GeneratePDFOutput failed")

	data, err := os.ReadFile(outputPath)
	require.NoError(t, err, "failed to read PDF file")

	require.GreaterOrEqual(t, len(data), 5, "PDF file is too small to contain magic bytes")
	require.Equal(t, "%PDF-", string(data[:5]), "file does not begin with PDF magic bytes")
}

// TestPDFOutputSubstantialSize verifies that the generated PDF is large enough to
// represent 3 pages of drawing content. Cairo uses compressed object streams
// (PDF 1.5+ ObjStm), so page dictionaries are not visible as raw text and page
// count cannot be verified without a PDF parsing library. File size is used as a
// proxy: a 3-page document with shapes, gradients, and text should be at least 10KB.
func TestPDFOutputSubstantialSize(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "output.pdf")

	err := GeneratePDFOutput(outputPath)
	require.NoError(t, err, "GeneratePDFOutput failed")

	data, err := os.ReadFile(outputPath)
	require.NoError(t, err, "failed to read PDF file")

	require.GreaterOrEqual(t, len(data), 10000,
		"PDF file is too small for a 3-page document with drawing content",
	)
}
