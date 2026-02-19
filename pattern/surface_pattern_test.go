// ABOUTME: Tests for SurfacePattern implementation including creation, extend, and filter modes.
// ABOUTME: Includes integration tests with actual Cairo surfaces for texture mapping.
package pattern

import (
	"testing"
	"unsafe"

	"github.com/mikowitz/cairo/status"
	"github.com/mikowitz/cairo/surface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockSurface is a minimal mock implementation for testing SurfacePattern error paths.
type mockSurface struct {
	ptr    unsafe.Pointer
	status status.Status
}

func (m *mockSurface) Ptr() unsafe.Pointer {
	return m.ptr
}

func (m *mockSurface) Status() status.Status {
	return m.status
}

// testSurfaceAdapter wraps *surface.ImageSurface to satisfy the pattern.Surface interface.
// surface.ImageSurface.Ptr() returns surface.SurfacePtr, but pattern.Surface requires unsafe.Pointer.
type testSurfaceAdapter struct {
	*surface.ImageSurface
}

func (a testSurfaceAdapter) Ptr() unsafe.Pointer {
	return unsafe.Pointer(a.ImageSurface.Ptr())
}

// newTestSurfacePattern creates a real SurfacePattern backed by a 10x10 ImageSurface.
// Returns the pattern and a cleanup function that closes both the pattern and surface.
func newTestSurfacePattern(t *testing.T) (*SurfacePattern, func()) {
	t.Helper()
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 10, 10)
	require.NoError(t, err)
	pat, err := NewSurfacePattern(testSurfaceAdapter{surf})
	require.NoError(t, err)
	return pat, func() {
		_ = pat.Close()
		_ = surf.Close()
	}
}

// TestNewSurfacePattern tests creation of surface patterns
func TestNewSurfacePattern(t *testing.T) {
	t.Run("nil_surface", func(t *testing.T) {
		pattern, err := NewSurfacePattern(nil)
		assert.Error(t, err, "Should return error for nil surface")
		assert.Nil(t, pattern, "Pattern should be nil")
		assert.Equal(t, status.NullPointer, err, "Error should be NullPointer")
	})

	t.Run("invalid_surface", func(t *testing.T) {
		invalidSurf := &mockSurface{
			ptr:    nil,
			status: status.NoMemory,
		}
		pattern, err := NewSurfacePattern(invalidSurf)
		assert.Error(t, err, "Should return error for invalid surface")
		assert.Nil(t, pattern, "Pattern should be nil for invalid surface")
	})

	t.Run("valid_surface", func(t *testing.T) {
		pat, cleanup := newTestSurfacePattern(t)
		defer cleanup()
		assert.NotNil(t, pat)
		assert.Equal(t, status.Success, pat.Status())
	})
}

// TestSurfacePatternExtend tests extend mode round-trip get/set operations
func TestSurfacePatternExtend(t *testing.T) {
	tests := []struct {
		name   string
		extend Extend
	}{
		{"none", ExtendNone},
		{"repeat", ExtendRepeat},
		{"reflect", ExtendReflect},
		{"pad", ExtendPad},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pat, cleanup := newTestSurfacePattern(t)
			defer cleanup()

			pat.SetExtend(tt.extend)
			assert.Equal(t, tt.extend, pat.GetExtend())
		})
	}
}

// TestSurfacePatternFilter tests filter mode round-trip get/set operations
func TestSurfacePatternFilter(t *testing.T) {
	tests := []struct {
		name   string
		filter Filter
	}{
		{"fast", FilterFast},
		{"good", FilterGood},
		{"best", FilterBest},
		{"nearest", FilterNearest},
		{"bilinear", FilterBilinear},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pat, cleanup := newTestSurfacePattern(t)
			defer cleanup()

			pat.SetFilter(tt.filter)
			assert.Equal(t, tt.filter, pat.GetFilter())
		})
	}
}

// TestExtendStringer verifies Extend enum has string representation
func TestExtendStringer(t *testing.T) {
	tests := []struct {
		extend   Extend
		expected string
	}{
		{ExtendNone, "None"},
		{ExtendRepeat, "Repeat"},
		{ExtendReflect, "Reflect"},
		{ExtendPad, "Pad"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			str := tt.extend.String()
			assert.Equal(t, tt.expected, str, "Extend.String() should return correct value")
		})
	}
}

// TestFilterStringer verifies Filter enum has string representation
func TestFilterStringer(t *testing.T) {
	tests := []struct {
		filter   Filter
		expected string
	}{
		{FilterFast, "Fast"},
		{FilterGood, "Good"},
		{FilterBest, "Best"},
		{FilterNearest, "Nearest"},
		{FilterBilinear, "Bilinear"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			str := tt.filter.String()
			assert.Equal(t, tt.expected, str, "Filter.String() should return correct value")
		})
	}
}

// TestSurfacePatternInterfaceCompleteness verifies the Pattern interface is satisfied
func TestSurfacePatternInterfaceCompleteness(t *testing.T) {
	// Verify that SurfacePattern implements Pattern interface
	var _ Pattern = (*SurfacePattern)(nil)
}

// TestSurfacePatternGetType verifies that surface patterns return PatternTypeSurface
func TestSurfacePatternGetType(t *testing.T) {
	pat, cleanup := newTestSurfacePattern(t)
	defer cleanup()

	assert.Equal(t, PatternTypeSurface, pat.GetType())
}
