package surface

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFormatConstants verifies that Format constants are defined correctly
func TestFormatConstants(t *testing.T) {
	tests := []struct {
		name   string
		format Format
		want   int
	}{
		{"FormatInvalid", FormatInvalid, -1},
		{"FormatARGB32", FormatARGB32, 0},
		{"FormatRGB24", FormatRGB24, 1},
		{"FormatA8", FormatA8, 2},
		{"FormatA1", FormatA1, 3},
		{"FormatRGB16_565", FormatRGB16_565, 4},
		{"FormatRGB30", FormatRGB30, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, int(tt.format), "%s should equal %d", tt.name, tt.want)
		})
	}
}

// TestFormatString verifies that Format.String() returns proper names
func TestFormatString(t *testing.T) {
	tests := []struct {
		format Format
		want   string
	}{
		{FormatInvalid, "Invalid"},
		{FormatARGB32, "ARGB32"},
		{FormatRGB24, "RGB24"},
		{FormatA8, "A8"},
		{FormatA1, "A1"},
		{FormatRGB16_565, "RGB16_565"},
		{FormatRGB30, "RGB30"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.format.String())
		})
	}
}

// TestFormatStrideForWidth verifies stride calculations for different formats and widths
func TestFormatStrideForWidth(t *testing.T) {
	tests := []struct {
		name   string
		format Format
		width  int
		want   int
	}{
		// ARGB32 is 4 bytes per pixel, must be aligned to 4-byte boundaries
		{"ARGB32_width1", FormatARGB32, 1, 4},
		{"ARGB32_width10", FormatARGB32, 10, 40},
		{"ARGB32_width100", FormatARGB32, 100, 400},

		// RGB24 is also 4 bytes per pixel (with unused byte)
		{"RGB24_width1", FormatRGB24, 1, 4},
		{"RGB24_width10", FormatRGB24, 10, 40},
		{"RGB24_width100", FormatRGB24, 100, 400},

		// A8 is 1 byte per pixel, but stride must be aligned to 4-byte boundaries
		{"A8_width1", FormatA8, 1, 4},
		{"A8_width2", FormatA8, 2, 4},
		{"A8_width3", FormatA8, 3, 4},
		{"A8_width4", FormatA8, 4, 4},
		{"A8_width5", FormatA8, 5, 8},
		{"A8_width10", FormatA8, 10, 12},

		// A1 is 1 bit per pixel, 32 pixels per 4-byte word
		{"A1_width1", FormatA1, 1, 4},
		{"A1_width8", FormatA1, 8, 4},
		{"A1_width32", FormatA1, 32, 4},
		{"A1_width33", FormatA1, 33, 8},

		// RGB16_565 is 2 bytes per pixel
		{"RGB16_565_width1", FormatRGB16_565, 1, 4},
		{"RGB16_565_width2", FormatRGB16_565, 2, 4},
		{"RGB16_565_width3", FormatRGB16_565, 3, 8},
		{"RGB16_565_width10", FormatRGB16_565, 10, 20},

		// RGB30 is 4 bytes per pixel
		{"RGB30_width1", FormatRGB30, 1, 4},
		{"RGB30_width10", FormatRGB30, 10, 40},

		// Zero width should return error (-1 or 0 depending on Cairo version)
		// We'll check it returns a non-positive value
		{"ARGB32_width0", FormatARGB32, 0, 0},

		// Invalid format should return -1
		{"Invalid_width10", FormatInvalid, 10, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.format.StrideForWidth(tt.width)

			// For zero width and invalid format, we expect non-positive values
			if tt.width == 0 || tt.format == FormatInvalid {
				assert.LessOrEqual(t, got, 0, "Expected non-positive value for width=%d, format=%v", tt.width, tt.format)
				return
			}

			assert.Equal(t, tt.want, got, "StrideForWidth(%d) mismatch", tt.width)

			// Verify stride is always a multiple of 4 (alignment requirement)
			assert.Zero(t, got%4, "Stride must be multiple of 4, got %d", got)
		})
	}
}

// TestFormatStrideForWidthNegativeWidth verifies error handling for negative widths
func TestFormatStrideForWidthNegativeWidth(t *testing.T) {
	stride := FormatARGB32.StrideForWidth(-1)
	assert.LessOrEqual(t, stride, 0, "Expected non-positive value for negative width")
}

// TestFormatStrideForWidthLargeWidth verifies handling of very large widths
func TestFormatStrideForWidthLargeWidth(t *testing.T) {
	// Cairo should return -1 for widths that would overflow
	// The exact limit depends on Cairo's internal implementation,
	// but INT_MAX / 4 is definitely too large for ARGB32
	veryLargeWidth := (1 << 30) // A very large width that should fail

	stride := FormatARGB32.StrideForWidth(veryLargeWidth)
	assert.Equal(t, -1, stride, "Expected -1 for width that would overflow")
}

// Note: Tests for Surface interface and BaseSurface will be added after
// implementing the ImageSurface in Prompt 7, as we need a concrete
// surface type to test the interface methods properly.
