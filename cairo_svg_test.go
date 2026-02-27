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

// TestSVGVersionsViaRootPackage verifies SVGVersions returns supported versions via the root package.
func TestSVGVersionsViaRootPackage(t *testing.T) {
	versions := cairo.SVGVersions()
	require.NotEmpty(t, versions)
	assert.Contains(t, versions, cairo.SVGVersion11)
	assert.Contains(t, versions, cairo.SVGVersion12)
}

// TestSVGVersionToStringViaRootPackage verifies SVGVersionToString via the root package.
func TestSVGVersionToStringViaRootPackage(t *testing.T) {
	assert.Equal(t, "SVG 1.1", cairo.SVGVersionToString(cairo.SVGVersion11))
	assert.Equal(t, "SVG 1.2", cairo.SVGVersionToString(cairo.SVGVersion12))
}

// TestSVGSurfaceRestrictToVersionViaRootPackage verifies RestrictToVersion works via the root package.
func TestSVGSurfaceRestrictToVersionViaRootPackage(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "output.svg")

	surf, err := cairo.NewSVGSurface(filename, 600, 400)
	require.NoError(t, err)
	require.NotNil(t, surf)

	surf.RestrictToVersion(cairo.SVGVersion11)

	require.NoError(t, surf.Close())
}

// TestSVGUnitConstantsReexport verifies all SVGUnit constants are accessible via the root
// cairo package and have the exact numeric values Cairo's cairo_svg_unit_t enum defines.
func TestSVGUnitConstantsReexport(t *testing.T) {
	// These numeric values must match cairo_svg_unit_t in <cairo-svg.h>.
	assert.Equal(t, cairo.SVGUnit(0), cairo.SVGUnitUser)
	assert.Equal(t, cairo.SVGUnit(1), cairo.SVGUnitEm)
	assert.Equal(t, cairo.SVGUnit(2), cairo.SVGUnitEx)
	assert.Equal(t, cairo.SVGUnit(3), cairo.SVGUnitPx)
	assert.Equal(t, cairo.SVGUnit(4), cairo.SVGUnitIn)
	assert.Equal(t, cairo.SVGUnit(5), cairo.SVGUnitCm)
	assert.Equal(t, cairo.SVGUnit(6), cairo.SVGUnitMm)
	assert.Equal(t, cairo.SVGUnit(7), cairo.SVGUnitPt)
	assert.Equal(t, cairo.SVGUnit(8), cairo.SVGUnitPc)
	assert.Equal(t, cairo.SVGUnit(9), cairo.SVGUnitPercent)
}

// TestSVGVersionConstantsReexport verifies SVGVersion constants have the exact numeric
// values Cairo's cairo_svg_version_t enum defines.
func TestSVGVersionConstantsReexport(t *testing.T) {
	// These numeric values must match cairo_svg_version_t in <cairo-svg.h>.
	assert.Equal(t, cairo.SVGVersion(0), cairo.SVGVersion11)
	assert.Equal(t, cairo.SVGVersion(1), cairo.SVGVersion12)
}
