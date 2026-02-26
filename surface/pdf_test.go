// ABOUTME: Tests for PDFSurface creation, size changes, and multi-page document output.
// ABOUTME: Uses t.TempDir() for test PDF files to ensure automatic cleanup.

//go:build !nopdf

package surface

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mikowitz/cairo/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewPDFSurface verifies that a PDF surface can be created and produces a file.
func TestNewPDFSurface(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "test.pdf")

	s, err := NewPDFSurface(filename, 595, 842) // A4 in points
	require.NoError(t, err)
	require.NotNil(t, s)

	assert.Equal(t, status.Success, s.Status())

	s.Flush()
	err = s.Close()
	require.NoError(t, err)

	info, err := os.Stat(filename)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0), "PDF file should not be empty")
}

// TestNewPDFSurfaceInvalidPath verifies that an invalid path returns an error.
func TestNewPDFSurfaceInvalidPath(t *testing.T) {
	_, err := NewPDFSurface("/nonexistent/dir/test.pdf", 595, 842)
	require.Error(t, err)
}

// TestPDFSurfaceSetSize verifies that SetSize can be called without panicking,
// including on a closed surface (where it should be a no-op).
func TestPDFSurfaceSetSize(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "test.pdf")

	s, err := NewPDFSurface(filename, 595, 842)
	require.NoError(t, err)
	require.NotNil(t, s)

	// SetSize on an open surface should update subsequent page dimensions.
	s.SetSize(612, 792) // US Letter in points
	assert.Equal(t, status.Success, s.Status())

	err = s.Close()
	require.NoError(t, err)

	// SetSize on a closed surface should be a no-op, not a panic.
	s.SetSize(100, 100)
}

// TestPDFSurfaceMultiPage verifies that multi-page PDFs can be produced using ShowPage.
func TestPDFSurfaceMultiPage(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "multipage.pdf")

	s, err := NewPDFSurface(filename, 595, 842)
	require.NoError(t, err)
	require.NotNil(t, s)

	// Emit first page.
	s.ShowPage()

	// Change size for second page and emit it.
	s.SetSize(612, 792)
	s.ShowPage()

	s.Flush()
	err = s.Close()
	require.NoError(t, err)

	info, err := os.Stat(filename)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0), "Multi-page PDF file should not be empty")
}

// TestPDFSurfaceShowPageClosed verifies that ShowPage on a closed surface is a no-op.
func TestPDFSurfaceShowPageClosed(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "test.pdf")

	s, err := NewPDFSurface(filename, 595, 842)
	require.NoError(t, err)

	err = s.Close()
	require.NoError(t, err)

	// Should not panic.
	s.ShowPage()
}
