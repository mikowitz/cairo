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

// TestSVGSurfaceDocumentUnit verifies that SetDocumentUnit can be called with all unit types.
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
	}

	s.Flush()
	err = s.Close()
	require.NoError(t, err)

	// SetDocumentUnit on a closed surface should be a no-op, not a panic.
	s.SetDocumentUnit(SVGUnitPx)
}
