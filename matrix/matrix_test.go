package matrix

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewMatrix verifies that NewMatrix creates a matrix with the given values
func TestNewMatrix(t *testing.T) {
	tests := []struct {
		name   string
		xx, yx float64
		xy, yy float64
		x0, y0 float64
	}{
		{
			name: "Identity values",
			xx:   1.0, yx: 0.0,
			xy: 0.0, yy: 1.0,
			x0: 0.0, y0: 0.0,
		},
		{
			name: "Arbitrary values",
			xx:   2.0, yx: 3.0,
			xy: 4.0, yy: 5.0,
			x0: 6.0, y0: 7.0,
		},
		{
			name: "Negative values",
			xx:   -1.0, yx: -2.0,
			xy: -3.0, yy: -4.0,
			x0: -5.0, y0: -6.0,
		},
		{
			name: "Fractional values",
			xx:   0.5, yx: 0.25,
			xy: 0.125, yy: 0.0625,
			x0: 0.03125, y0: 0.015625,
		},
		{
			name: "Zero matrix",
			xx:   0.0, yx: 0.0,
			xy: 0.0, yy: 0.0,
			x0: 0.0, y0: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMatrix(tt.xx, tt.yx, tt.xy, tt.yy, tt.x0, tt.y0)
			assert.NotNil(t, m, "NewMatrix should return non-nil matrix")
			assert.Equal(t, tt.xx, m.XX, "XX should match")
			assert.Equal(t, tt.yx, m.YX, "YX should match")
			assert.Equal(t, tt.xy, m.XY, "XY should match")
			assert.Equal(t, tt.yy, m.YY, "YY should match")
			assert.Equal(t, tt.x0, m.X0, "X0 should match")
			assert.Equal(t, tt.y0, m.Y0, "Y0 should match")
		})
	}
}

// TestNewIdentityMatrix verifies that NewIdentityMatrix returns an identity matrix
// Identity matrix has diagonal values of 1 and all others 0
func TestNewIdentityMatrix(t *testing.T) {
	m := NewIdentityMatrix()
	assert.NotNil(t, m, "NewIdentityMatrix should return non-nil matrix")
	assert.Equal(t, 1.0, m.XX, "XX should be 1")
	assert.Equal(t, 0.0, m.YX, "YX should be 0")
	assert.Equal(t, 0.0, m.XY, "XY should be 0")
	assert.Equal(t, 1.0, m.YY, "YY should be 1")
	assert.Equal(t, 0.0, m.X0, "X0 should be 0")
	assert.Equal(t, 0.0, m.Y0, "Y0 should be 0")
}

// TestMatrixCGOConversion verifies round-trip conversion between Go and C matrix representations
func TestMatrixCGOConversion(t *testing.T) {
	tests := []struct {
		name   string
		xx, yx float64
		xy, yy float64
		x0, y0 float64
	}{
		{
			name: "Identity matrix",
			xx:   1.0, yx: 0.0,
			xy: 0.0, yy: 1.0,
			x0: 0.0, y0: 0.0,
		},
		{
			name: "Arbitrary values",
			xx:   1.5, yx: 2.5,
			xy: 3.5, yy: 4.5,
			x0: 5.5, y0: 6.5,
		},
		{
			name: "Negative values",
			xx:   -1.0, yx: -2.0,
			xy: -3.0, yy: -4.0,
			x0: -5.0, y0: -6.0,
		},
		{
			name: "Fractional values",
			xx:   0.125, yx: 0.25,
			xy: 0.5, yy: 0.75,
			x0: 1.125, y0: 1.625,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create original matrix
			original := NewMatrix(tt.xx, tt.yx, tt.xy, tt.yy, tt.x0, tt.y0)

			// Convert to C and back
			cMatrix := original.toC()
			converted := matrixFromC(cMatrix)

			// Verify all fields match
			assert.Equal(t, original.XX, converted.XX, "XX should match after round-trip")
			assert.Equal(t, original.YX, converted.YX, "YX should match after round-trip")
			assert.Equal(t, original.XY, converted.XY, "XY should match after round-trip")
			assert.Equal(t, original.YY, converted.YY, "YY should match after round-trip")
			assert.Equal(t, original.X0, converted.X0, "X0 should match after round-trip")
			assert.Equal(t, original.Y0, converted.Y0, "Y0 should match after round-trip")

			// Verify C pointer is the same (we're working with the same C struct)
			assert.Equal(t, original.toC(), converted.toC(), "C pointers should be equal")
		})
	}
}

// TestMatrixThreadSafety verifies that concurrent reads/writes don't race
// TODO: implement this once we have methods to update matrices
func TestMatrixThreadSafety(t *testing.T) {
	t.Skip()

	m := NewIdentityMatrix()

	var wg sync.WaitGroup
	iterations := 100

	// Start multiple goroutines writing to the matrix
	for i := range 10 {
		wg.Add(1)
		go func(val float64) {
			defer wg.Done()
			for range iterations {
				// m.NewMatrix(val, val, val, val, val, val)
			}
		}(float64(i))
	}

	// Start multiple goroutines reading from the matrix
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range iterations {
				// Read all fields
				_ = m.XX
				_ = m.YX
				_ = m.XY
				_ = m.YY
				_ = m.X0
				_ = m.Y0
			}
		}()
	}

	wg.Wait()
	// If we get here without a race condition, the test passes
	// Run with `go test -race` to verify
}
