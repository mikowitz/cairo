package cairo_test

import (
	"testing"

	"github.com/mikowitz/cairo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFormatReexport verifies that Format type is re-exported correctly
func TestFormatReexport(t *testing.T) {
	// Verify we can use Format type from cairo package
	format := cairo.FormatARGB32
	assert.Equal(t, cairo.FormatARGB32, format)
}

// TestFormatConstants verifies all format constants are accessible
func TestFormatConstants(t *testing.T) {
	tests := []struct {
		name   string
		format cairo.Format
	}{
		{"FormatInvalid", cairo.FormatInvalid},
		{"FormatARGB32", cairo.FormatARGB32},
		{"FormatRGB24", cairo.FormatRGB24},
		{"FormatA8", cairo.FormatA8},
		{"FormatA1", cairo.FormatA1},
		{"FormatRGB16_565", cairo.FormatRGB16_565},
		{"FormatRGB30", cairo.FormatRGB30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify the constant is accessible and has a valid value
			_ = tt.format
		})
	}
}

// TestNewImageSurface verifies the re-exported function works
func TestNewImageSurface(t *testing.T) {
	// Use re-exported types and function from cairo package
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 100, 100)
	require.NoError(t, err)
	require.NotNil(t, surface)
	defer func() {
		err := surface.Close()
		require.NoError(t, err, "Surface should close without error")
	}()

	// Verify surface properties
	assert.Equal(t, cairo.FormatARGB32, surface.GetFormat())
	assert.Equal(t, 100, surface.GetWidth())
	assert.Equal(t, 100, surface.GetHeight())
}

// TestNewImageSurfaceWithDifferentFormats tests various format constants
func TestNewImageSurfaceWithDifferentFormats(t *testing.T) {
	formats := []cairo.Format{
		cairo.FormatARGB32,
		cairo.FormatRGB24,
		cairo.FormatA8,
		cairo.FormatA1,
		cairo.FormatRGB16_565,
		cairo.FormatRGB30,
	}

	for _, format := range formats {
		t.Run(format.String(), func(t *testing.T) {
			surface, err := cairo.NewImageSurface(format, 50, 50)
			require.NoError(t, err)
			require.NotNil(t, surface)
			defer func() {
				err := surface.Close()
				require.NoError(t, err, "Surface should close without error")
			}()

			assert.Equal(t, format, surface.GetFormat())
		})
	}
}

// TestNewImageSurfaceInvalidFormat tests error handling with invalid format
func TestNewImageSurfaceInvalidFormat(t *testing.T) {
	surface, err := cairo.NewImageSurface(cairo.FormatInvalid, 100, 100)
	assert.Error(t, err)
	assert.Nil(t, surface)
}

// TestSurfaceInterfaceReexport verifies that Surface interface is re-exported
func TestSurfaceInterfaceReexport(t *testing.T) {
	// Create a surface using the re-exported function
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 100, 100)
	require.NoError(t, err)
	require.NotNil(t, surface)
	defer func() {
		err := surface.Close()
		require.NoError(t, err, "Surface should close without error")
	}()

	// Verify the surface can be used as a cairo.Surface interface
	var _ cairo.Surface = surface

	// Test Surface interface methods are accessible
	surface.Flush()
	surface.MarkDirty()
	surface.MarkDirtyRectangle(0, 0, 50, 50)

	// Test Status method
	status := surface.Status()
	assert.Equal(t, 0, int(status)) // StatusSuccess should be 0
}

// TestSurfaceLifecycle demonstrates proper resource management
func TestSurfaceLifecycle(t *testing.T) {
	// Create a surface
	surface, err := cairo.NewImageSurface(cairo.FormatARGB32, 200, 200)
	require.NoError(t, err, "Should create surface successfully")
	require.NotNil(t, surface, "Surface should not be nil")

	// Verify surface is usable
	assert.Equal(t, 200, surface.GetWidth())
	assert.Equal(t, 200, surface.GetHeight())
	assert.Equal(t, cairo.FormatARGB32, surface.GetFormat())

	// Flush any pending operations
	surface.Flush()

	// Mark surface as dirty (simulating external modification)
	surface.MarkDirty()

	// Close the surface explicitly
	err = surface.Close()
	require.NoError(t, err, "Should close without error")

	// Attempting to close again should be safe (idempotent)
	err = surface.Close()
	require.NoError(t, err, "Closing again should be safe")
}
