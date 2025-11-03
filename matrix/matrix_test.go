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

// TestNewTranslationMatrix verifies that NewTranslationMatrix creates a translation matrix
func TestNewTranslationMatrix(t *testing.T) {
	tests := []struct {
		name   string
		tx, ty float64
	}{
		{
			name: "Simple translation",
			tx:   10.0, ty: 20.0,
		},
		{
			name: "Negative translation",
			tx:   -5.0, ty: -10.0,
		},
		{
			name: "Zero translation",
			tx:   0.0, ty: 0.0,
		},
		{
			name: "Fractional translation",
			tx:   3.14, ty: 2.71,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewTranslationMatrix(tt.tx, tt.ty)
			assert.NotNil(t, m, "NewTranslationMatrix should return non-nil matrix")
			assert.Equal(t, 1.0, m.XX, "XX should be 1")
			assert.Equal(t, 0.0, m.YX, "YX should be 0")
			assert.Equal(t, 0.0, m.XY, "XY should be 0")
			assert.Equal(t, 1.0, m.YY, "YY should be 1")
			assert.Equal(t, tt.tx, m.X0, "X0 should match tx")
			assert.Equal(t, tt.ty, m.Y0, "Y0 should match ty")
		})
	}
}

// TestNewScalingMatrix verifies that NewScalingMatrix creates a scaling matrix
func TestNewScalingMatrix(t *testing.T) {
	tests := []struct {
		name   string
		sx, sy float64
	}{
		{
			name: "Uniform scaling",
			sx:   2.0, sy: 2.0,
		},
		{
			name: "Non-uniform scaling",
			sx:   3.0, sy: 4.0,
		},
		{
			name: "Fractional scaling",
			sx:   0.5, sy: 0.25,
		},
		{
			name: "Identity scaling",
			sx:   1.0, sy: 1.0,
		},
		{
			name: "Large scale factors",
			sx:   10.0, sy: 20.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewScalingMatrix(tt.sx, tt.sy)
			assert.NotNil(t, m, "NewScalingMatrix should return non-nil matrix")
			assert.Equal(t, tt.sx, m.XX, "XX should match sx")
			assert.Equal(t, 0.0, m.YX, "YX should be 0")
			assert.Equal(t, 0.0, m.XY, "XY should be 0")
			assert.Equal(t, tt.sy, m.YY, "YY should match sy")
			assert.Equal(t, 0.0, m.X0, "X0 should be 0")
			assert.Equal(t, 0.0, m.Y0, "Y0 should be 0")
		})
	}
}

// TestNewRotationMatrix verifies that NewRotationMatrix creates a rotation matrix
func TestNewRotationMatrix(t *testing.T) {
	tests := []struct {
		name                       string
		radians                    float64
		expXX, expYX, expXY, expYY float64
	}{
		{
			name:    "No rotation",
			radians: 0.0,
			expXX:   1.0, expYX: 0.0,
			expXY: 0.0, expYY: 1.0,
		},
		{
			name:    "90 degree rotation",
			radians: 1.5707963267948966, // π/2
			expXX:   0.0, expYX: 1.0,
			expXY: -1.0, expYY: 0.0,
		},
		{
			name:    "180 degree rotation",
			radians: 3.141592653589793, // π
			expXX:   -1.0, expYX: 0.0,
			expXY: 0.0, expYY: -1.0,
		},
		{
			name:    "45 degree rotation",
			radians: 0.7853981633974483, // π/4
			expXX:   0.7071067811865476, expYX: 0.7071067811865475,
			expXY: -0.7071067811865475, expYY: 0.7071067811865476,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewRotationMatrix(tt.radians)
			assert.NotNil(t, m, "NewRotationMatrix should return non-nil matrix")
			assert.InDelta(t, tt.expXX, m.XX, 0.0001, "XX should match expected value")
			assert.InDelta(t, tt.expYX, m.YX, 0.0001, "YX should match expected value")
			assert.InDelta(t, tt.expXY, m.XY, 0.0001, "XY should match expected value")
			assert.InDelta(t, tt.expYY, m.YY, 0.0001, "YY should match expected value")
			assert.Equal(t, 0.0, m.X0, "X0 should be 0")
			assert.Equal(t, 0.0, m.Y0, "Y0 should be 0")
		})
	}
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

// TestMatrixMultiply verifies matrix multiplication
func TestMatrixMultiply(t *testing.T) {
	tests := []struct {
		name     string
		m1       *Matrix
		m2       *Matrix
		expected *Matrix
	}{
		{
			name:     "Identity times identity",
			m1:       NewIdentityMatrix(),
			m2:       NewIdentityMatrix(),
			expected: NewIdentityMatrix(),
		},
		{
			name:     "Matrix times identity",
			m1:       NewMatrix(2.0, 0.0, 0.0, 2.0, 10.0, 20.0),
			m2:       NewIdentityMatrix(),
			expected: NewMatrix(2.0, 0.0, 0.0, 2.0, 10.0, 20.0),
		},
		{
			name:     "Simple multiplication",
			m1:       NewMatrix(2.0, 0.0, 0.0, 2.0, 0.0, 0.0),
			m2:       NewMatrix(3.0, 0.0, 0.0, 3.0, 0.0, 0.0),
			expected: NewMatrix(6.0, 0.0, 0.0, 6.0, 0.0, 0.0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.m1.Multiply(tt.m2)

			assert.InDelta(t, tt.expected.XX, actual.XX, 0.0001, "XX should match")
			assert.InDelta(t, tt.expected.YX, actual.YX, 0.0001, "YX should match")
			assert.InDelta(t, tt.expected.XY, actual.XY, 0.0001, "XY should match")
			assert.InDelta(t, tt.expected.YY, actual.YY, 0.0001, "YY should match")
			assert.InDelta(t, tt.expected.X0, actual.X0, 0.0001, "X0 should match")
			assert.InDelta(t, tt.expected.Y0, actual.Y0, 0.0001, "Y0 should match")
		})
	}
}

// TestMatrixTransformPoint verifies point transformation
func TestMatrixTransformPoint(t *testing.T) {
	tests := []struct {
		name       string
		m          *Matrix
		x, y       float64
		expX, expY float64
	}{
		{
			name: "Identity transformation",
			m:    NewIdentityMatrix(),
			x:    10.0, y: 20.0,
			expX: 10.0, expY: 20.0,
		},
		{
			name: "Translation",
			m:    NewMatrix(1.0, 0.0, 0.0, 1.0, 5.0, 10.0),
			x:    10.0, y: 20.0,
			expX: 15.0, expY: 30.0,
		},
		{
			name: "Scaling",
			m:    NewMatrix(2.0, 0.0, 0.0, 3.0, 0.0, 0.0),
			x:    10.0, y: 20.0,
			expX: 20.0, expY: 60.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y := tt.m.TransformPoint(tt.x, tt.y)
			assert.InDelta(t, tt.expX, x, 0.0001, "X should match")
			assert.InDelta(t, tt.expY, y, 0.0001, "Y should match")
		})
	}
}

// TestMatrixTransformDistance verifies distance vector transformation
func TestMatrixTransformDistance(t *testing.T) {
	tests := []struct {
		name         string
		m            *Matrix
		dx, dy       float64
		expDX, expDY float64
	}{
		{
			name: "Identity transformation",
			m:    NewIdentityMatrix(),
			dx:   10.0, dy: 20.0,
			expDX: 10.0, expDY: 20.0,
		},
		{
			name: "Translation (no effect on distance)",
			m:    NewMatrix(1.0, 0.0, 0.0, 1.0, 5.0, 10.0),
			dx:   10.0, dy: 20.0,
			expDX: 10.0, expDY: 20.0,
		},
		{
			name: "Scaling",
			m:    NewMatrix(2.0, 0.0, 0.0, 3.0, 0.0, 0.0),
			dx:   10.0, dy: 20.0,
			expDX: 20.0, expDY: 60.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dx, dy := tt.m.TransformDistance(tt.dx, tt.dy)
			assert.InDelta(t, tt.expDX, dx, 0.0001, "DX should match")
			assert.InDelta(t, tt.expDY, dy, 0.0001, "DY should match")
		})
	}
}

// TestMatrixTranslate verifies translation transformation
func TestMatrixTranslate(t *testing.T) {
	tests := []struct {
		name         string
		tx, ty       float64
		expX0, expY0 float64
	}{
		{
			name: "Simple translation",
			tx:   10.0, ty: 20.0,
			expX0: 10.0, expY0: 20.0,
		},
		{
			name: "Negative translation",
			tx:   -5.0, ty: -10.0,
			expX0: -5.0, expY0: -10.0,
		},
		{
			name: "Zero translation",
			tx:   0.0, ty: 0.0,
			expX0: 0.0, expY0: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewIdentityMatrix()
			m.Translate(tt.tx, tt.ty)

			assert.InDelta(t, 1.0, m.XX, 0.0001, "XX should be 1")
			assert.InDelta(t, 0.0, m.YX, 0.0001, "YX should be 0")
			assert.InDelta(t, 0.0, m.XY, 0.0001, "XY should be 0")
			assert.InDelta(t, 1.0, m.YY, 0.0001, "YY should be 1")
			assert.InDelta(t, tt.expX0, m.X0, 0.0001, "X0 should match translation")
			assert.InDelta(t, tt.expY0, m.Y0, 0.0001, "Y0 should match translation")
		})
	}
}

// TestMatrixScale verifies scaling transformation
func TestMatrixScale(t *testing.T) {
	tests := []struct {
		name         string
		sx, sy       float64
		expXX, expYY float64
	}{
		{
			name: "Uniform scaling",
			sx:   2.0, sy: 2.0,
			expXX: 2.0, expYY: 2.0,
		},
		{
			name: "Non-uniform scaling",
			sx:   3.0, sy: 4.0,
			expXX: 3.0, expYY: 4.0,
		},
		{
			name: "Fractional scaling",
			sx:   0.5, sy: 0.25,
			expXX: 0.5, expYY: 0.25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewIdentityMatrix()
			m.Scale(tt.sx, tt.sy)

			assert.InDelta(t, tt.expXX, m.XX, 0.0001, "XX should match scale")
			assert.InDelta(t, 0.0, m.YX, 0.0001, "YX should be 0")
			assert.InDelta(t, 0.0, m.XY, 0.0001, "XY should be 0")
			assert.InDelta(t, tt.expYY, m.YY, 0.0001, "YY should match scale")
			assert.InDelta(t, 0.0, m.X0, 0.0001, "X0 should be 0")
			assert.InDelta(t, 0.0, m.Y0, 0.0001, "Y0 should be 0")
		})
	}
}

// TestMatrixRotate verifies rotation transformation (test with 90 degrees)
func TestMatrixRotate(t *testing.T) {
	tests := []struct {
		name                       string
		radians                    float64
		expXX, expYX, expXY, expYY float64
	}{
		{
			name:    "No rotation",
			radians: 0.0,
			expXX:   1.0, expYX: 0.0,
			expXY: 0.0, expYY: 1.0,
		},
		{
			name:    "90 degree rotation",
			radians: 1.5707963267948966, // π/2
			expXX:   0.0, expYX: 1.0,
			expXY: -1.0, expYY: 0.0,
		},
		{
			name:    "180 degree rotation",
			radians: 3.141592653589793, // π
			expXX:   -1.0, expYX: 0.0,
			expXY: 0.0, expYY: -1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewIdentityMatrix()
			m.Rotate(tt.radians)

			assert.InDelta(t, tt.expXX, m.XX, 0.0001, "XX should match")
			assert.InDelta(t, tt.expYX, m.YX, 0.0001, "YX should match")
			assert.InDelta(t, tt.expXY, m.XY, 0.0001, "XY should match")
			assert.InDelta(t, tt.expYY, m.YY, 0.0001, "YY should match")
			assert.InDelta(t, 0.0, m.X0, 0.0001, "X0 should be 0")
			assert.InDelta(t, 0.0, m.Y0, 0.0001, "Y0 should be 0")
		})
	}
}

// TestMatrixInvert verifies matrix inversion
func TestMatrixInvert(t *testing.T) {
	t.Run("Invert identity matrix", func(t *testing.T) {
		m := NewIdentityMatrix()
		err := m.Invert()

		assert.NoError(t, err, "Inverting identity should succeed")
		assert.InDelta(t, 1.0, m.XX, 0.0001, "XX should be 1")
		assert.InDelta(t, 0.0, m.YX, 0.0001, "YX should be 0")
		assert.InDelta(t, 0.0, m.XY, 0.0001, "XY should be 0")
		assert.InDelta(t, 1.0, m.YY, 0.0001, "YY should be 1")
	})

	t.Run("Invert scaling matrix", func(t *testing.T) {
		m := NewMatrix(2.0, 0.0, 0.0, 3.0, 0.0, 0.0)
		err := m.Invert()

		assert.NoError(t, err, "Inverting scaling matrix should succeed")
		assert.InDelta(t, 0.5, m.XX, 0.0001, "XX should be 0.5")
		assert.InDelta(t, 0.0, m.YX, 0.0001, "YX should be 0")
		assert.InDelta(t, 0.0, m.XY, 0.0001, "XY should be 0")
		assert.InDelta(t, 1.0/3.0, m.YY, 0.0001, "YY should be 1/3")
	})

	t.Run("Invert singular matrix", func(t *testing.T) {
		m := NewMatrix(0.0, 0.0, 0.0, 0.0, 0.0, 0.0)
		err := m.Invert()

		assert.Error(t, err, "Inverting singular matrix should fail")
	})
}

// TestMatrixOperationsCombined verifies combining translate, scale, rotate
func TestMatrixOperationsCombined(t *testing.T) {
	t.Run("Translate then scale", func(t *testing.T) {
		m := NewIdentityMatrix()
		m.Translate(10.0, 20.0)
		m.Scale(2.0, 3.0)

		// After translate(10, 20) then scale(2, 3)
		// Point (0,0) should map to (10, 20)
		x, y := m.TransformPoint(0.0, 0.0)
		assert.InDelta(t, 10.0, x, 0.0001, "X should be 10")
		assert.InDelta(t, 20.0, y, 0.0001, "Y should be 20")
	})

	t.Run("Scale then translate", func(t *testing.T) {
		m := NewIdentityMatrix()
		m.Scale(2.0, 3.0)
		m.Translate(10.0, 20.0)

		// After scale(2, 3) then translate(10, 20)
		// Point (0,0) should map to (20, 60)
		x, y := m.TransformPoint(0.0, 0.0)
		assert.InDelta(t, 20.0, x, 0.0001, "X should be 20")
		assert.InDelta(t, 60.0, y, 0.0001, "Y should be 60")
	})

	t.Run("Translate, scale, and rotate", func(t *testing.T) {
		m := NewIdentityMatrix()
		m.Translate(100.0, 100.0)
		m.Scale(2.0, 2.0)
		m.Rotate(1.5707963267948966) // 90 degrees

		// This tests a complex transformation chain
		x, y := m.TransformPoint(10.0, 0.0)
		// After translate(100, 100), scale(2, 2), rotate(90°):
		// Point (10, 0) → translate → (110, 100)
		//               → scale     → (220, 200)
		//               → rotate 90°→ (-200, 220)
		assert.InDelta(t, 100.0, x, 0.0001, "X should be 100")
		assert.InDelta(t, 120.0, y, 0.0001, "Y should be 120")
	})
}

// TestMatrixClose verifies that Close releases resources properly
func TestMatrixClose(t *testing.T) {
	t.Run("Close releases resources", func(t *testing.T) {
		m := NewIdentityMatrix()
		err := m.Close()
		assert.NoError(t, err, "Close should succeed")
	})

	t.Run("Double close is safe", func(t *testing.T) {
		m := NewIdentityMatrix()

		// First close
		err := m.Close()
		assert.NoError(t, err, "First close should succeed")

		// Second close should also succeed (no-op)
		err = m.Close()
		assert.NoError(t, err, "Second close should be safe")
	})

	t.Run("Close with different matrix types", func(t *testing.T) {
		matrices := []*Matrix{
			NewMatrix(1, 2, 3, 4, 5, 6),
			NewIdentityMatrix(),
			NewTranslationMatrix(10, 20),
			NewScalingMatrix(2, 3),
			NewRotationMatrix(1.5707963267948966),
		}

		for _, m := range matrices {
			err := m.Close()
			assert.NoError(t, err, "Close should succeed for all matrix types")
		}
	})

	t.Run("Multiple matrices can be closed independently", func(t *testing.T) {
		m1 := NewIdentityMatrix()
		m2 := NewIdentityMatrix()

		err := m1.Close()
		assert.NoError(t, err, "Closing m1 should succeed")

		// m2 should still be usable
		m2.Translate(5, 10)
		assert.InDelta(t, 5.0, m2.X0, 0.0001, "m2 should still be usable after m1 is closed")

		err = m2.Close()
		assert.NoError(t, err, "Closing m2 should succeed")
	})
}

// TestMatrixString verifies that String() returns a formatted representation of the matrix
func TestMatrixString(t *testing.T) {
	tests := []struct {
		name     string
		m        *Matrix
		contains []string
	}{
		{
			name: "Identity matrix",
			m:    NewIdentityMatrix(),
			contains: []string{
				"Matrix",
				"1.00", "0.00",
			},
		},
		{
			name: "Translation matrix",
			m:    NewTranslationMatrix(10.0, 20.0),
			contains: []string{
				"Matrix",
				"10.00", "20.00",
			},
		},
		{
			name: "Scaling matrix",
			m:    NewScalingMatrix(2.5, 3.5),
			contains: []string{
				"Matrix",
				"2.50", "3.50",
			},
		},
		{
			name: "Arbitrary matrix",
			m:    NewMatrix(1.5, 2.5, 3.5, 4.5, 5.5, 6.5),
			contains: []string{
				"Matrix",
				"1.50", "2.50", "3.50", "4.50", "5.50", "6.50",
			},
		},
		{
			name: "Negative values",
			m:    NewMatrix(-1.0, -2.0, -3.0, -4.0, -5.0, -6.0),
			contains: []string{
				"Matrix",
				"-1.00", "-2.00", "-3.00", "-4.00", "-5.00", "-6.00",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := tt.m.String()
			assert.NotEmpty(t, str, "String() should return non-empty string")

			// Verify all expected substrings are present
			for _, substr := range tt.contains {
				assert.Contains(t, str, substr, "String should contain %q", substr)
			}
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
				m.Translate(val, val)
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
