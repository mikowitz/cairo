// ABOUTME: Tests for the PDF dashboard example.
// ABOUTME: Verifies file creation and PDF header validity.

//go:build !nopdf

package examples

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestDashboardPDFGeneratesFile verifies that GenerateDashboardPDF creates a valid PDF file.
func TestDashboardPDFGeneratesFile(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "dashboard.pdf")

	err := GenerateDashboardPDF(outputPath)
	require.NoError(t, err, "GenerateDashboardPDF failed")
	require.FileExists(t, outputPath, "output PDF was not created")
	require.True(t, CheckFileSize(t, outputPath, 1000, 5_000_000), "PDF file size out of expected range")

	// Confirm the output begins with the PDF magic bytes.
	data, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	require.True(t, len(data) >= 5 && string(data[:5]) == "%PDF-", "file does not start with PDF magic bytes")

	t.Logf("Generated dashboard PDF at %s", outputPath)
}
