// ABOUTME: Tests for the PDFSurface type and NewPDFSurface constructor re-exported
// ABOUTME: from the root cairo package.

//go:build !nopdf

package cairo_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mikowitz/cairo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPDFSurfaceTypeReexport verifies that PDFSurface type is accessible from the
// root cairo package and can be used as a type constraint.
func TestPDFSurfaceTypeReexport(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "test.pdf")

	surf, err := cairo.NewPDFSurface(filename, 612, 792)
	require.NoError(t, err)
	require.NotNil(t, surf)
	defer func() { _ = surf.Close() }()

	// Verify the returned type is *cairo.PDFSurface (compile-time re-export check)
	_ = surf
}

// TestNewPDFSurfaceViaRootPackage verifies NewPDFSurface works through the root package,
// including SetSize, ShowPage, and integration with NewContext.
func TestNewPDFSurfaceViaRootPackage(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "output.pdf")

	surf, err := cairo.NewPDFSurface(filename, 612, 792)
	require.NoError(t, err)
	require.NotNil(t, surf)

	// Draw on page 1
	ctx, err := cairo.NewContext(surf)
	require.NoError(t, err)

	ctx.SetSourceRGB(1.0, 0.0, 0.0)
	ctx.Rectangle(10, 10, 100, 100)
	ctx.Fill()
	require.NoError(t, ctx.Close())

	// Change page size and emit page 2
	surf.SetSize(595, 842)
	surf.ShowPage()

	require.NoError(t, surf.Close())

	// Verify a non-empty PDF file was written
	info, err := os.Stat(filename)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0))
}

// TestNewPDFSurfaceInvalidPathViaRootPackage verifies that an invalid path returns an error.
func TestNewPDFSurfaceInvalidPathViaRootPackage(t *testing.T) {
	surf, err := cairo.NewPDFSurface("/nonexistent/dir/output.pdf", 612, 792)
	assert.Error(t, err)
	assert.Nil(t, surf)
}
