package pattern

import (
	"testing"

	"github.com/mikowitz/cairo/matrix"
	"github.com/mikowitz/cairo/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewSolidPatternRGB tests creation of solid RGB patterns
func TestNewSolidPatternRGB(t *testing.T) {
	tests := []struct {
		name    string
		r, g, b float64
	}{
		{"red", 1.0, 0.0, 0.0},
		{"green", 0.0, 1.0, 0.0},
		{"blue", 0.0, 0.0, 1.0},
		{"white", 1.0, 1.0, 1.0},
		{"black", 0.0, 0.0, 0.0},
		{"gray", 0.5, 0.5, 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern, err := NewSolidPatternRGB(tt.r, tt.g, tt.b)
			require.NoError(t, err, "NewSolidPatternRGB should not return an error")
			require.NotNil(t, pattern, "Pattern should not be nil")

			// Verify pattern has successful status
			st := pattern.Status()
			assert.Equal(t, status.Success, st, "Pattern should have success status")

			// Clean up
			err = pattern.Close()
			assert.NoError(t, err, "Close should not return an error")
		})
	}
}

// TestNewSolidPatternRGBA tests creation of solid RGBA patterns with alpha
func TestNewSolidPatternRGBA(t *testing.T) {
	tests := []struct {
		name       string
		r, g, b, a float64
	}{
		{"opaque_red", 1.0, 0.0, 0.0, 1.0},
		{"transparent_red", 1.0, 0.0, 0.0, 0.0},
		{"semi_transparent_blue", 0.0, 0.0, 1.0, 0.5},
		{"semi_transparent_white", 1.0, 1.0, 1.0, 0.75},
		{"almost_transparent_black", 0.0, 0.0, 0.0, 0.1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern, err := NewSolidPatternRGBA(tt.r, tt.g, tt.b, tt.a)
			require.NoError(t, err, "NewSolidPatternRGBA should not return an error")
			require.NotNil(t, pattern, "Pattern should not be nil")

			// Verify pattern has successful status
			st := pattern.Status()
			assert.Equal(t, status.Success, st, "Pattern should have success status")

			// Clean up
			err = pattern.Close()
			assert.NoError(t, err, "Close should not return an error")
		})
	}
}

// TestPatternClose verifies close behavior and double-close safety
func TestPatternClose(t *testing.T) {
	t.Run("close_once", func(t *testing.T) {
		pattern, err := NewSolidPatternRGB(1.0, 0.0, 0.0)
		require.NoError(t, err)
		require.NotNil(t, pattern)

		err = pattern.Close()
		assert.NoError(t, err, "First close should succeed")
	})

	t.Run("double_close_safe", func(t *testing.T) {
		pattern, err := NewSolidPatternRGB(0.0, 1.0, 0.0)
		require.NoError(t, err)
		require.NotNil(t, pattern)

		// First close
		err = pattern.Close()
		assert.NoError(t, err, "First close should succeed")

		// Second close should be safe (no-op)
		err = pattern.Close()
		assert.NoError(t, err, "Second close should be safe (no-op)")
	})

	t.Run("operations_after_close", func(t *testing.T) {
		pattern, err := NewSolidPatternRGB(0.0, 0.0, 1.0)
		require.NoError(t, err)
		require.NotNil(t, pattern)

		// Close the pattern
		err = pattern.Close()
		require.NoError(t, err)

		// Operations after close should be safe (no-ops or return appropriate errors)
		// Status should still work
		st := pattern.Status()
		assert.NotEqual(t, status.Success, st, "Status after close might not be Success")
	})
}

// TestPatternStatus verifies status reporting
func TestPatternStatus(t *testing.T) {
	t.Run("valid_pattern_has_success_status", func(t *testing.T) {
		pattern, err := NewSolidPatternRGB(1.0, 0.5, 0.0)
		require.NoError(t, err)
		require.NotNil(t, pattern)
		defer pattern.Close()

		st := pattern.Status()
		assert.Equal(t, status.Success, st, "Valid pattern should have Success status")
	})

	t.Run("status_multiple_calls", func(t *testing.T) {
		pattern, err := NewSolidPatternRGBA(0.5, 0.5, 0.5, 0.5)
		require.NoError(t, err)
		require.NotNil(t, pattern)
		defer pattern.Close()

		// Multiple status calls should return consistent results
		st1 := pattern.Status()
		st2 := pattern.Status()
		assert.Equal(t, st1, st2, "Status should be consistent across multiple calls")
	})
}

// TestPatternMatrix verifies matrix get/set operations
func TestPatternMatrix(t *testing.T) {
	t.Run("get_default_matrix", func(t *testing.T) {
		pattern, err := NewSolidPatternRGB(1.0, 0.0, 0.0)
		require.NoError(t, err)
		require.NotNil(t, pattern)
		defer pattern.Close()

		m, err := pattern.GetMatrix()
		require.NoError(t, err, "GetMatrix should not return an error")
		require.NotNil(t, m, "Matrix should not be nil")

		// Default matrix should be identity
		identity := matrix.NewIdentityMatrix()
		assert.Equal(t, identity.XX, m.XX, "XX should match identity")
		assert.Equal(t, identity.YX, m.YX, "YX should match identity")
		assert.Equal(t, identity.XY, m.XY, "XY should match identity")
		assert.Equal(t, identity.YY, m.YY, "YY should match identity")
		assert.Equal(t, identity.X0, m.X0, "X0 should match identity")
		assert.Equal(t, identity.Y0, m.Y0, "Y0 should match identity")
	})

	t.Run("set_and_get_matrix", func(t *testing.T) {
		pattern, err := NewSolidPatternRGB(0.0, 1.0, 0.0)
		require.NoError(t, err)
		require.NotNil(t, pattern)
		defer pattern.Close()

		// Create a translation matrix
		m := matrix.NewTranslationMatrix(10.0, 20.0)

		// Set the matrix
		pattern.SetMatrix(m)

		// Get it back
		retrieved, err := pattern.GetMatrix()
		require.NoError(t, err, "GetMatrix should not return an error")
		require.NotNil(t, retrieved, "Retrieved matrix should not be nil")

		// Verify the matrix values match
		assert.InDelta(t, m.XX, retrieved.XX, 0.0001, "XX should match")
		assert.InDelta(t, m.YX, retrieved.YX, 0.0001, "YX should match")
		assert.InDelta(t, m.XY, retrieved.XY, 0.0001, "XY should match")
		assert.InDelta(t, m.YY, retrieved.YY, 0.0001, "YY should match")
		assert.InDelta(t, m.X0, retrieved.X0, 0.0001, "X0 should match")
		assert.InDelta(t, m.Y0, retrieved.Y0, 0.0001, "Y0 should match")
	})

	t.Run("set_scaling_matrix", func(t *testing.T) {
		pattern, err := NewSolidPatternRGBA(0.0, 0.0, 1.0, 0.8)
		require.NoError(t, err)
		require.NotNil(t, pattern)
		defer pattern.Close()

		// Create a scaling matrix
		m := matrix.NewScalingMatrix(2.0, 3.0)

		// Set the matrix
		pattern.SetMatrix(m)

		// Get it back
		retrieved, err := pattern.GetMatrix()
		require.NoError(t, err)
		require.NotNil(t, retrieved)

		// Verify scaling values
		assert.InDelta(t, 2.0, retrieved.XX, 0.0001, "XX should be 2.0")
		assert.InDelta(t, 3.0, retrieved.YY, 0.0001, "YY should be 3.0")
	})

	t.Run("set_rotation_matrix", func(t *testing.T) {
		pattern, err := NewSolidPatternRGB(1.0, 1.0, 0.0)
		require.NoError(t, err)
		require.NotNil(t, pattern)
		defer pattern.Close()

		// Create a rotation matrix (90 degrees)
		m := matrix.NewRotationMatrix(1.5707963267948966) // π/2

		// Set the matrix
		pattern.SetMatrix(m)

		// Get it back
		retrieved, err := pattern.GetMatrix()
		require.NoError(t, err)
		require.NotNil(t, retrieved)

		// Verify the matrix represents a rotation
		// For 90 degree rotation: XX ≈ 0, YX ≈ 1, XY ≈ -1, YY ≈ 0
		assert.InDelta(t, 0.0, retrieved.XX, 0.0001, "XX should be ~0 for 90° rotation")
		assert.InDelta(t, 1.0, retrieved.YX, 0.0001, "YX should be ~1 for 90° rotation")
		assert.InDelta(t, -1.0, retrieved.XY, 0.0001, "XY should be ~-1 for 90° rotation")
		assert.InDelta(t, 0.0, retrieved.YY, 0.0001, "YY should be ~0 for 90° rotation")
	})

	t.Run("matrix_after_close", func(t *testing.T) {
		pattern, err := NewSolidPatternRGB(0.5, 0.5, 0.5)
		require.NoError(t, err)
		require.NotNil(t, pattern)

		// Close the pattern
		err = pattern.Close()
		require.NoError(t, err)

		// GetMatrix after close should return an error or handle gracefully
		m, err := pattern.GetMatrix()
		// Either error is returned or m is nil
		if err == nil {
			assert.Nil(t, m, "Matrix should be nil after close if no error")
		}
	})
}

// TestPatternInterfaceCompleteness verifies the Pattern interface can be satisfied
func TestPatternInterfaceCompleteness(t *testing.T) {
	// Verify that BasePattern implements Pattern interface
	var _ Pattern = (*BasePattern)(nil)

	// Verify that SolidPattern implements Pattern interface
	var _ Pattern = (*SolidPattern)(nil)
}

// TestPatternWithDifferentColors tests patterns work with various color values
func TestPatternWithDifferentColors(t *testing.T) {
	tests := []struct {
		name       string
		r, g, b, a float64
		useAlpha   bool
	}{
		{"pure_red", 1.0, 0.0, 0.0, 1.0, false},
		{"pure_green", 0.0, 1.0, 0.0, 1.0, false},
		{"pure_blue", 0.0, 0.0, 1.0, 1.0, false},
		{"yellow", 1.0, 1.0, 0.0, 1.0, false},
		{"cyan", 0.0, 1.0, 1.0, 1.0, false},
		{"magenta", 1.0, 0.0, 1.0, 1.0, false},
		{"semi_transparent_red", 1.0, 0.0, 0.0, 0.5, true},
		{"very_transparent_white", 1.0, 1.0, 1.0, 0.1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pattern Pattern
			var err error

			if tt.useAlpha {
				pattern, err = NewSolidPatternRGBA(tt.r, tt.g, tt.b, tt.a)
			} else {
				pattern, err = NewSolidPatternRGB(tt.r, tt.g, tt.b)
			}

			require.NoError(t, err)
			require.NotNil(t, pattern)
			defer pattern.Close()

			// Verify status is success
			assert.Equal(t, status.Success, pattern.Status())
		})
	}
}

// TestPatternThreadSafety verifies patterns are thread-safe
func TestPatternThreadSafety(t *testing.T) {
	pattern, err := NewSolidPatternRGB(0.5, 0.5, 0.5)
	require.NoError(t, err)
	require.NotNil(t, pattern)
	defer pattern.Close()

	// Number of goroutines
	const numGoroutines = 10
	const numOperations = 100

	done := make(chan bool, numGoroutines)

	// Multiple goroutines performing concurrent operations
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer func() { done <- true }()
			for j := 0; j < numOperations; j++ {
				// Concurrent reads
				_ = pattern.Status()

				// Concurrent matrix operations
				m, err := pattern.GetMatrix()
				if err == nil && m != nil {
					pattern.SetMatrix(m)
				}
			}
		}()
	}

	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Pattern should still be valid
	assert.Equal(t, status.Success, pattern.Status())
}

// TestSolidPatternGetType verifies that solid patterns return PatternTypeSolid
func TestSolidPatternGetType(t *testing.T) {
	t.Run("rgb_pattern_type", func(t *testing.T) {
		pattern, err := NewSolidPatternRGB(1.0, 0.0, 0.0)
		require.NoError(t, err)
		require.NotNil(t, pattern)
		defer pattern.Close()

		patternType := pattern.GetType()
		assert.Equal(t, PatternTypeSolid, patternType, "RGB solid pattern should have PatternTypeSolid type")
	})

	t.Run("rgba_pattern_type", func(t *testing.T) {
		pattern, err := NewSolidPatternRGBA(0.0, 1.0, 0.0, 0.5)
		require.NoError(t, err)
		require.NotNil(t, pattern)
		defer pattern.Close()

		patternType := pattern.GetType()
		assert.Equal(t, PatternTypeSolid, patternType, "RGBA solid pattern should have PatternTypeSolid type")
	})

	t.Run("type_after_close", func(t *testing.T) {
		pattern, err := NewSolidPatternRGB(0.5, 0.5, 0.5)
		require.NoError(t, err)
		require.NotNil(t, pattern)

		// Get type before close
		typeBefore := pattern.GetType()
		assert.Equal(t, PatternTypeSolid, typeBefore)

		// Close the pattern
		err = pattern.Close()
		require.NoError(t, err)

		// Type should still be retrievable after close
		typeAfter := pattern.GetType()
		assert.Equal(t, PatternTypeSolid, typeAfter, "Pattern type should still be Solid even after close")
	})
}
