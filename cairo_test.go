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

// TestOperatorReexport verifies that the Operator type and SetOperator/GetOperator
// are usable via the root cairo package.
func TestOperatorReexport(t *testing.T) {
	surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 10, 10)
	require.NoError(t, err)
	defer func() { _ = surf.Close() }()

	ctx, err := cairo.NewContext(surf)
	require.NoError(t, err)
	defer func() { _ = ctx.Close() }()

	ctx.SetOperator(cairo.OperatorAdd)
	assert.Equal(t, cairo.OperatorAdd, ctx.GetOperator())
}

// TestOperatorRoundTrip verifies all 29 re-exported operator constants round-trip
// correctly through SetOperator/GetOperator.
func TestOperatorRoundTrip(t *testing.T) {
	surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 10, 10)
	require.NoError(t, err)
	defer func() { _ = surf.Close() }()

	ctx, err := cairo.NewContext(surf)
	require.NoError(t, err)
	defer func() { _ = ctx.Close() }()

	operators := []struct {
		name string
		op   cairo.Operator
	}{
		{"OperatorClear", cairo.OperatorClear},
		{"OperatorSource", cairo.OperatorSource},
		{"OperatorOver", cairo.OperatorOver},
		{"OperatorIn", cairo.OperatorIn},
		{"OperatorOut", cairo.OperatorOut},
		{"OperatorAtop", cairo.OperatorAtop},
		{"OperatorDest", cairo.OperatorDest},
		{"OperatorDestOver", cairo.OperatorDestOver},
		{"OperatorDestIn", cairo.OperatorDestIn},
		{"OperatorDestOut", cairo.OperatorDestOut},
		{"OperatorDestAtop", cairo.OperatorDestAtop},
		{"OperatorXor", cairo.OperatorXor},
		{"OperatorAdd", cairo.OperatorAdd},
		{"OperatorSaturate", cairo.OperatorSaturate},
		{"OperatorMultiply", cairo.OperatorMultiply},
		{"OperatorScreen", cairo.OperatorScreen},
		{"OperatorOverlay", cairo.OperatorOverlay},
		{"OperatorDarken", cairo.OperatorDarken},
		{"OperatorLighten", cairo.OperatorLighten},
		{"OperatorColorDodge", cairo.OperatorColorDodge},
		{"OperatorColorBurn", cairo.OperatorColorBurn},
		{"OperatorHardLight", cairo.OperatorHardLight},
		{"OperatorSoftLight", cairo.OperatorSoftLight},
		{"OperatorDifference", cairo.OperatorDifference},
		{"OperatorExclusion", cairo.OperatorExclusion},
		{"OperatorHslHue", cairo.OperatorHslHue},
		{"OperatorHslSaturation", cairo.OperatorHslSaturation},
		{"OperatorHslColor", cairo.OperatorHslColor},
		{"OperatorHslLuminosity", cairo.OperatorHslLuminosity},
	}

	for _, tt := range operators {
		t.Run(tt.name, func(t *testing.T) {
			ctx.SetOperator(tt.op)
			assert.Equal(t, tt.op, ctx.GetOperator())
		})
	}
}

// TestOperatorDefaultIsOver verifies that OperatorOver is the default (value 2 per Cairo spec).
func TestOperatorDefaultIsOver(t *testing.T) {
	surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 10, 10)
	require.NoError(t, err)
	defer func() { _ = surf.Close() }()

	ctx, err := cairo.NewContext(surf)
	require.NoError(t, err)
	defer func() { _ = ctx.Close() }()

	assert.Equal(t, cairo.OperatorOver, ctx.GetOperator())
}

// TestSlantReexport verifies that the Slant type and constants are re-exported
// correctly from the root cairo package.
func TestSlantReexport(t *testing.T) {
	tests := []struct {
		name  string
		slant cairo.Slant
	}{
		{"SlantNormal", cairo.SlantNormal},
		{"SlantItalic", cairo.SlantItalic},
		{"SlantOblique", cairo.SlantOblique},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.slant
		})
	}
}

// TestWeightReexport verifies that the Weight type and constants are re-exported
// correctly from the root cairo package.
func TestWeightReexport(t *testing.T) {
	tests := []struct {
		name   string
		weight cairo.Weight
	}{
		{"WeightNormal", cairo.WeightNormal},
		{"WeightBold", cairo.WeightBold},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.weight
		})
	}
}

// TestSelectFontFaceViaRootPackage verifies that SelectFontFace works using
// re-exported Slant and Weight constants from the root cairo package.
func TestSelectFontFaceViaRootPackage(t *testing.T) {
	surf, err := cairo.NewImageSurface(cairo.FormatARGB32, 100, 100)
	require.NoError(t, err)
	defer func() { _ = surf.Close() }()

	ctx, err := cairo.NewContext(surf)
	require.NoError(t, err)
	defer func() { _ = ctx.Close() }()

	ctx.SelectFontFace("sans-serif", cairo.SlantNormal, cairo.WeightNormal)
	ctx.SelectFontFace("serif", cairo.SlantItalic, cairo.WeightBold)
	ctx.SelectFontFace("monospace", cairo.SlantOblique, cairo.WeightNormal)

	assert.Equal(t, 0, int(ctx.Status()))
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
