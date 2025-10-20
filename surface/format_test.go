package surface

import (
	"testing"
)

func TestFormatConstants(t *testing.T) {
	tests := []struct {
		name   string
		format Format
		want   int
	}{
		{"FormatInvalid is -1", FormatInvalid, -1},
		{"FormatARGB32 is 0", FormatARGB32, 0},
		{"FormatRGB24 is 1", FormatRGB24, 1},
		{"FormatA8 is 2", FormatA8, 2},
		{"FormatA1 is 3", FormatA1, 3},
		{"FormatRGB16_565 is 4", FormatRGB16_565, 4},
		{"FormatRGB30 is 5", FormatRGB30, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.format) != tt.want {
				t.Errorf("format constant %s = %d, want %d", tt.name, int(tt.format), tt.want)
			}
		})
	}
}

func TestFormatStrideForWidth(t *testing.T) {
	tests := []struct {
		name   string
		format Format
		width  int
		want   int
	}{
		// ARGB32: 4 bytes per pixel
		{"ARGB32 width 1", FormatARGB32, 1, 4},
		{"ARGB32 width 10", FormatARGB32, 10, 40},
		{"ARGB32 width 100", FormatARGB32, 100, 400},

		// RGB24: 4 bytes per pixel (same as ARGB32 in Cairo)
		{"RGB24 width 1", FormatRGB24, 1, 4},
		{"RGB24 width 10", FormatRGB24, 10, 40},
		{"RGB24 width 100", FormatRGB24, 100, 400},

		// A8: 1 byte per pixel, but Cairo aligns to 4-byte boundaries
		{"A8 width 1", FormatA8, 1, 4},
		{"A8 width 4", FormatA8, 4, 4},
		{"A8 width 5", FormatA8, 5, 8},
		{"A8 width 8", FormatA8, 8, 8},
		{"A8 width 100", FormatA8, 100, 100},

		// A1: 1 bit per pixel, packed into bytes, aligned to 4-byte boundaries
		{"A1 width 1", FormatA1, 1, 4},
		{"A1 width 8", FormatA1, 8, 4},
		{"A1 width 9", FormatA1, 9, 4},
		{"A1 width 32", FormatA1, 32, 4},
		{"A1 width 33", FormatA1, 33, 8},
		{"A1 width 100", FormatA1, 100, 16},

		// RGB16_565: 2 bytes per pixel, aligned to 4-byte boundaries
		{"RGB16_565 width 1", FormatRGB16_565, 1, 4},
		{"RGB16_565 width 2", FormatRGB16_565, 2, 4},
		{"RGB16_565 width 3", FormatRGB16_565, 3, 8},
		{"RGB16_565 width 100", FormatRGB16_565, 100, 200},

		// RGB30: 4 bytes per pixel
		{"RGB30 width 1", FormatRGB30, 1, 4},
		{"RGB30 width 10", FormatRGB30, 10, 40},
		{"RGB30 width 100", FormatRGB30, 100, 400},

		// Edge cases
		{"ARGB32 width 0", FormatARGB32, 0, 0},
		{"A8 width 0", FormatA8, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.format.StrideForWidth(tt.width)
			if got != tt.want {
				t.Errorf("Format(%v).StrideForWidth(%d) = %d, want %d", tt.format, tt.width, got, tt.want)
			}
		})
	}
}

func TestFormatStrideForWidthInvalid(t *testing.T) {
	// Invalid format should return -1
	stride := FormatInvalid.StrideForWidth(100)
	if stride != -1 {
		t.Errorf("FormatInvalid.StrideForWidth(100) = %d, want -1", stride)
	}
}

func TestFormatStringer(t *testing.T) {
	// Verify that Format implements String() method (via go:generate stringer)
	tests := []struct {
		format Format
		want   string
	}{
		{FormatInvalid, "FormatInvalid"},
		{FormatARGB32, "FormatARGB32"},
		{FormatRGB24, "FormatRGB24"},
		{FormatA8, "FormatA8"},
		{FormatA1, "FormatA1"},
		{FormatRGB16_565, "FormatRGB16_565"},
		{FormatRGB30, "FormatRGB30"},
		{FormatRGB96F, "FormatRGB96F"},
		{FormatRGBA128F, "FormatRGBA128F"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.format.String()
			if got != tt.want {
				t.Errorf("Format(%d).String() = %q, want %q", tt.format, got, tt.want)
			}
		})
	}
}
