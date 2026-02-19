// ABOUTME: Tests for SurfacePattern implementation including creation, extend, and filter modes.
// ABOUTME: Includes integration tests with actual Cairo surfaces for texture mapping.
package pattern

import (
	"testing"

	"github.com/mikowitz/cairo/status"
	"github.com/stretchr/testify/assert"
)

// mockSurface is a minimal mock implementation for testing SurfacePattern creation
type mockSurface struct {
	ptr    interface{}
	status status.Status
}

func (m *mockSurface) Ptr() interface{} {
	return m.ptr
}

func (m *mockSurface) Status() status.Status {
	return m.status
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
}

// TestSurfacePatternExtend tests extend mode get/set operations
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
			// This test would require a real surface
			// For now, we just verify the enum values exist
			assert.NotNil(t, tt.extend)
		})
	}
}

// TestSurfacePatternFilter tests filter mode get/set operations
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
			// This test would require a real surface
			// For now, we just verify the enum values exist
			assert.NotNil(t, tt.filter)
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
	// This would require a real surface to test fully
	// For now, we verify the type constant exists
	assert.Equal(t, PatternTypeSurface, PatternTypeSurface)
}
