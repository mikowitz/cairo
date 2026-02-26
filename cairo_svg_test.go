// ABOUTME: Tests for the SVGSurface type, SVGUnit type, and NewSVGSurface constructor
// ABOUTME: re-exported from the root cairo package.

//go:build !nosvg

package cairo_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mikowitz/cairo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSVGSurfaceTypeReexport verifies that SVGSurface type is accessible from the
// root cairo package and can be used as a type constraint.
func TestSVGSurfaceTypeReexport(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "test.svg")

	surf, err := cairo.NewSVGSurface(filename, 600, 400)
	require.NoError(t, err)
	require.NotNil(t, surf)
	defer func() { _ = surf.Close() }()
}

// TestNewSVGSurfaceViaRootPackage verifies NewSVGSurface works through the root package,
// including SetDocumentUnit and integration with NewContext.
func TestNewSVGSurfaceViaRootPackage(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "output.svg")

	surf, err := cairo.NewSVGSurface(filename, 600, 400)
	require.NoError(t, err)
	require.NotNil(t, surf)

	// Set document unit
	surf.SetDocumentUnit(cairo.SVGUnitPx)

	// Draw on the surface
	ctx, err := cairo.NewContext(surf)
	require.NoError(t, err)

	ctx.SetSourceRGB(1.0, 0.0, 0.0)
	ctx.Rectangle(10, 10, 100, 100)
	ctx.Fill()
	require.NoError(t, ctx.Close())

	require.NoError(t, surf.Close())

	// Verify a non-empty SVG file was written
	info, err := os.Stat(filename)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0))
}

// TestNewSVGSurfaceInvalidPathViaRootPackage verifies that an invalid path returns an error.
func TestNewSVGSurfaceInvalidPathViaRootPackage(t *testing.T) {
	surf, err := cairo.NewSVGSurface("/nonexistent/dir/output.svg", 600, 400)
	assert.Error(t, err)
	assert.Nil(t, surf)
}

// TestSVGUnitConstantsReexport verifies all SVGUnit constants are accessible via
// the root cairo package.
func TestSVGUnitConstantsReexport(t *testing.T) {
	units := []cairo.SVGUnit{
		cairo.SVGUnitUser,
		cairo.SVGUnitEm,
		cairo.SVGUnitEx,
		cairo.SVGUnitPx,
		cairo.SVGUnitIn,
		cairo.SVGUnitCm,
		cairo.SVGUnitMm,
		cairo.SVGUnitPt,
		cairo.SVGUnitPc,
		cairo.SVGUnitPercent,
	}
	assert.Len(t, units, 10)

	// Verify they are distinct values
	seen := make(map[cairo.SVGUnit]bool)
	for _, u := range units {
		assert.False(t, seen[u], "duplicate SVGUnit value: %v", u)
		seen[u] = true
	}
}
