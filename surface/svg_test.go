// ABOUTME: Tests for SVGSurface creation and document unit configuration.
// ABOUTME: Uses t.TempDir() for test SVG files to ensure automatic cleanup.

//go:build !nosvg

package surface

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewSVGSurface verifies that an SVG surface can be created and produces a valid file.
func TestNewSVGSurface(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "test.svg")

	s, err := NewSVGSurface(filename, 400, 300)
	require.NoError(t, err)
	require.NotNil(t, s)

	s.Flush()
	err = s.Close()
	require.NoError(t, err)

	info, err := os.Stat(filename)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0), "SVG file should not be empty")

	// Verify the file starts with a proper XML header.
	data, err := os.ReadFile(filename) //nolint:gosec // filename is from t.TempDir(), not user input
	require.NoError(t, err)
	assert.Contains(t, string(data), "<?xml", "SVG file should start with XML declaration")

	// Verify the SVG file is parseable as XML.
	var v interface{}
	err = xml.Unmarshal(data, &v)
	assert.NoError(t, err, "SVG file should be valid XML")
}

// TestNewSVGSurfaceInvalidPath verifies that an invalid path returns an error.
func TestNewSVGSurfaceInvalidPath(t *testing.T) {
	_, err := NewSVGSurface("/nonexistent/dir/test.svg", 400, 300)
	require.Error(t, err)
}

// TestSVGVersions verifies that SVGVersions returns a non-empty list of supported versions.
func TestSVGVersions(t *testing.T) {
	versions := SVGVersions()
	require.NotEmpty(t, versions, "SVGVersions should return at least one version")

	// Cairo currently supports SVG 1.1 and 1.2
	assert.Contains(t, versions, SVGVersion11)
	assert.Contains(t, versions, SVGVersion12)
}

// TestSVGVersionToString verifies that SVGVersionToString returns human-readable strings.
func TestSVGVersionToString(t *testing.T) {
	assert.Equal(t, "SVG 1.1", SVGVersionToString(SVGVersion11))
	assert.Equal(t, "SVG 1.2", SVGVersionToString(SVGVersion12))
}

// TestSVGSurfaceRestrictToVersion verifies that RestrictToVersion produces a valid SVG
// for each supported version and is a no-op when called on a closed surface.
// Cairo does not add a version attribute to the <svg> element; the restriction
// controls which SVG features Cairo is permitted to emit, not the header content.
func TestSVGSurfaceRestrictToVersion(t *testing.T) {
	for _, v := range SVGVersions() {
		t.Run(v.String(), func(t *testing.T) {
			dir := t.TempDir()
			filename := filepath.Join(dir, "test.svg")

			s, err := NewSVGSurface(filename, 100, 100)
			require.NoError(t, err)
			s.RestrictToVersion(v)
			require.NoError(t, s.Close())

			data, err := os.ReadFile(filename) //nolint:gosec // filename is from t.TempDir()
			require.NoError(t, err)
			assert.Contains(t, string(data), "<svg", "output should contain SVG root element")

			var parsed interface{}
			assert.NoError(t, xml.Unmarshal(data, &parsed), "output should be well-formed XML")
		})
	}

	// RestrictToVersion on a closed surface should be a no-op, not a panic.
	dir := t.TempDir()
	s, err := NewSVGSurface(filepath.Join(dir, "test.svg"), 100, 100)
	require.NoError(t, err)
	require.NoError(t, s.Close())
	s.RestrictToVersion(SVGVersion11)
}

// TestSVGSurfaceDocumentUnit verifies that SetDocumentUnit and GetDocumentUnit round-trip correctly.
func TestSVGSurfaceDocumentUnit(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "test.svg")

	s, err := NewSVGSurface(filename, 400, 300)
	require.NoError(t, err)
	require.NotNil(t, s)

	units := []SVGUnit{
		SVGUnitUser, SVGUnitEm, SVGUnitEx,
		SVGUnitPx, SVGUnitIn, SVGUnitCm,
		SVGUnitMm, SVGUnitPt, SVGUnitPc, SVGUnitPercent,
	}
	for _, unit := range units {
		s.SetDocumentUnit(unit)
		assert.Equal(t, unit, s.GetDocumentUnit(), "GetDocumentUnit should return the unit set by SetDocumentUnit")
	}

	s.Flush()
	err = s.Close()
	require.NoError(t, err)

	// SetDocumentUnit on a closed surface should be a no-op, not a panic.
	s.SetDocumentUnit(SVGUnitPx)
}
