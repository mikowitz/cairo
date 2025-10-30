package context

import (
	"math"
	"runtime"
	"testing"

	"github.com/mikowitz/cairo/matrix"
	"github.com/mikowitz/cairo/status"
	"github.com/mikowitz/cairo/surface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestContextGetSetMatrix verifies getting and setting the transformation matrix.
func TestContextGetSetMatrix(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 200, 200)
	require.NoError(t, err, "Failed to create surface")
	defer func() {
		err := surf.Close()
		assert.NoError(t, err, "Failed to close surface")
	}()

	ctx, err := NewContext(surf)
	require.NoError(t, err, "Failed to create context")
	defer func() {
		err := ctx.Close()
		assert.NoError(t, err, "Failed to close context")
	}()

	t.Run("get_identity_matrix", func(t *testing.T) {
		ctx.IdentityMatrix()

		m, err := ctx.GetMatrix()
		require.NoError(t, err, "GetMatrix should succeed")
		require.NotNil(t, m, "Matrix should not be nil")
		defer m.Close()

		// Verify identity matrix values
		assert.InDelta(t, 1.0, m.XX, 0.001)
		assert.InDelta(t, 0.0, m.YX, 0.001)
		assert.InDelta(t, 0.0, m.XY, 0.001)
		assert.InDelta(t, 1.0, m.YY, 0.001)
		assert.InDelta(t, 0.0, m.X0, 0.001)
		assert.InDelta(t, 0.0, m.Y0, 0.001)
	})

	t.Run("get_matrix_after_transformations", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Translate(10.0, 20.0)
		ctx.Scale(2.0, 3.0)

		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		defer m.Close()

		assert.InDelta(t, 2.0, m.XX, 0.001)
		assert.InDelta(t, 3.0, m.YY, 0.001)
		assert.InDelta(t, 10.0, m.X0, 0.001)
		assert.InDelta(t, 20.0, m.Y0, 0.001)
	})

	t.Run("set_matrix", func(t *testing.T) {
		ctx.IdentityMatrix()

		// Create a custom matrix
		customMatrix := matrix.NewMatrix(2.0, 0.0, 0.0, 3.0, 10.0, 20.0)
		defer customMatrix.Close()

		ctx.SetMatrix(customMatrix)

		assert.Equal(t, status.Success, ctx.Status())

		// Verify matrix was set
		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		defer m.Close()

		assert.InDelta(t, 2.0, m.XX, 0.001)
		assert.InDelta(t, 3.0, m.YY, 0.001)
		assert.InDelta(t, 10.0, m.X0, 0.001)
		assert.InDelta(t, 20.0, m.Y0, 0.001)
	})

	t.Run("round_trip_get_set_matrix", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Translate(15.0, 25.0)
		ctx.Scale(1.5, 2.5)

		// Get the matrix
		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		defer m.Close()

		// Reset to identity
		ctx.IdentityMatrix()

		// Set the matrix back
		ctx.SetMatrix(m)

		// Verify it matches original
		m2, err := ctx.GetMatrix()
		require.NoError(t, err)
		defer m2.Close()

		assert.InDelta(t, m.XX, m2.XX, 0.001)
		assert.InDelta(t, m.YY, m2.YY, 0.001)
		assert.InDelta(t, m.X0, m2.X0, 0.001)
		assert.InDelta(t, m.Y0, m2.Y0, 0.001)
	})

	t.Run("get_matrix_after_close", func(t *testing.T) {
		surf2, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surf2.Close()

		ctx2, err := NewContext(surf2)
		require.NoError(t, err)

		// Close the context
		err = ctx2.Close()
		require.NoError(t, err)

		// GetMatrix should return an error after close
		m, err := ctx2.GetMatrix()
		assert.Error(t, err, "GetMatrix should return error after close")
		assert.Nil(t, m, "Matrix should be nil after close")
	})

	t.Run("set_matrix_after_close", func(t *testing.T) {
		surf2, err := surface.NewImageSurface(surface.FormatARGB32, 100, 100)
		require.NoError(t, err)
		defer surf2.Close()

		ctx2, err := NewContext(surf2)
		require.NoError(t, err)

		customMatrix := matrix.NewIdentityMatrix()
		defer customMatrix.Close()

		// Close the context
		err = ctx2.Close()
		require.NoError(t, err)

		// SetMatrix should be safe no-op after close
		ctx2.SetMatrix(customMatrix)

		// Status should indicate closed
		st := ctx2.Status()
		assert.Equal(t, status.NullPointer, st)
	})

	t.Run("set_nil_matrix", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Translate(5.0, 10.0)

		// SetMatrix with nil should be safe no-op
		ctx.SetMatrix(nil)

		// Matrix should be unchanged
		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		defer m.Close()

		assert.InDelta(t, 5.0, m.X0, 0.001)
		assert.InDelta(t, 10.0, m.Y0, 0.001)
	})
}

// TestContextIdentityMatrix verifies identity matrix reset.
func TestContextIdentityMatrix(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 200, 200)
	require.NoError(t, err, "Failed to create surface")
	defer func() {
		err := surf.Close()
		assert.NoError(t, err, "Failed to close surface")
	}()

	ctx, err := NewContext(surf)
	require.NoError(t, err, "Failed to create context")
	defer func() {
		err := ctx.Close()
		assert.NoError(t, err, "Failed to close context")
	}()

	t.Run("reset_after_translate", func(t *testing.T) {
		ctx.Translate(10.0, 20.0)
		ctx.IdentityMatrix()

		assert.Equal(t, status.Success, ctx.Status())

		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		defer m.Close()

		assert.InDelta(t, 1.0, m.XX, 0.001)
		assert.InDelta(t, 1.0, m.YY, 0.001)
		assert.InDelta(t, 0.0, m.X0, 0.001)
		assert.InDelta(t, 0.0, m.Y0, 0.001)

		// Verify coordinates are not transformed
		x, y := ctx.UserToDevice(10.0, 20.0)
		assert.InDelta(t, 10.0, x, 0.001)
		assert.InDelta(t, 20.0, y, 0.001)
	})

	t.Run("reset_after_scale", func(t *testing.T) {
		ctx.Scale(2.0, 3.0)
		ctx.IdentityMatrix()

		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		defer m.Close()

		assert.InDelta(t, 1.0, m.XX, 0.001)
		assert.InDelta(t, 1.0, m.YY, 0.001)

		x, y := ctx.UserToDevice(10.0, 20.0)
		assert.InDelta(t, 10.0, x, 0.001)
		assert.InDelta(t, 20.0, y, 0.001)
	})

	t.Run("reset_after_rotate", func(t *testing.T) {
		ctx.Rotate(math.Pi / 4) // 45 degrees
		ctx.IdentityMatrix()

		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		defer m.Close()

		assert.InDelta(t, 1.0, m.XX, 0.001)
		assert.InDelta(t, 0.0, m.YX, 0.001)
		assert.InDelta(t, 0.0, m.XY, 0.001)
		assert.InDelta(t, 1.0, m.YY, 0.001)
	})

	t.Run("reset_after_complex_transformations", func(t *testing.T) {
		ctx.Translate(10.0, 20.0)
		ctx.Scale(2.0, 3.0)
		ctx.Rotate(math.Pi / 6)
		ctx.IdentityMatrix()

		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		defer m.Close()

		// All values should be back to identity
		assert.InDelta(t, 1.0, m.XX, 0.001)
		assert.InDelta(t, 0.0, m.YX, 0.001)
		assert.InDelta(t, 0.0, m.XY, 0.001)
		assert.InDelta(t, 1.0, m.YY, 0.001)
		assert.InDelta(t, 0.0, m.X0, 0.001)
		assert.InDelta(t, 0.0, m.Y0, 0.001)
	})

	t.Run("multiple_identity_calls", func(t *testing.T) {
		ctx.Translate(5.0, 10.0)
		ctx.IdentityMatrix()
		ctx.IdentityMatrix() // Second call should be harmless

		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		defer m.Close()

		assert.InDelta(t, 1.0, m.XX, 0.001)
		assert.InDelta(t, 0.0, m.X0, 0.001)
	})
}

// TestContextGetMatrixMemorySafety verifies that GetMatrix returns a matrix
// with properly allocated memory that remains valid after the stack is reused.
//
// This test surfaces a critical bug where contextGetMatrix() passes a stack-allocated
// cairo_matrix_t pointer to matrix.FromPointer(), which stores the pointer directly
// and sets a finalizer to free it. This causes:
// 1. The pointer to become invalid when the stack is reused
// 2. The finalizer to attempt C.free() on stack memory (undefined behavior)
func TestContextGetMatrixMemorySafety(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 200, 200)
	require.NoError(t, err, "Failed to create surface")
	defer func() {
		err := surf.Close()
		assert.NoError(t, err, "Failed to close surface")
	}()

	ctx, err := NewContext(surf)
	require.NoError(t, err, "Failed to create context")
	defer func() {
		err := ctx.Close()
		assert.NoError(t, err, "Failed to close context")
	}()

	t.Run("matrix_pointer_remains_valid_after_stack_reuse", func(t *testing.T) {
		// Set up a non-identity transformation
		ctx.IdentityMatrix()
		ctx.Translate(10.0, 20.0)
		ctx.Scale(2.0, 3.0)

		// Get the matrix - if bug exists, this stores a stack pointer
		m, err := ctx.GetMatrix()
		require.NoError(t, err)
		require.NotNil(t, m)
		defer m.Close()

		// Verify initial values
		initialXX := m.XX
		initialYY := m.YY
		initialX0 := m.X0
		initialY0 := m.Y0

		// Perform operations that consume stack space and potentially overwrite
		// the stack memory where the matrix was allocated
		var stackConsumingOperations func(depth int)
		stackConsumingOperations = func(depth int) {
			if depth > 0 {
				// Create temporary matrices to consume stack
				temp := matrix.NewIdentityMatrix()
				temp.Translate(float64(depth), float64(depth))
				temp.Scale(2.0, 2.0)

				// Create large local variables to consume stack
				largeLocal := make([]float64, 256)
				for i := range largeLocal {
					largeLocal[i] = float64(i)
				}

				// Recursive call to consume more stack
				stackConsumingOperations(depth - 1)

				_ = temp.Close()
			}
		}
		stackConsumingOperations(10)

		// Perform other context operations that use stack
		for i := 0; i < 5; i++ {
			ctx.Save()
			ctx.Translate(float64(i), float64(i))
			ctx.Restore()
		}

		// Try to use the matrix's internal pointer via operations that access m.ptr
		// These operations will crash or produce corrupted results if the pointer is invalid
		t.Run("multiply_uses_internal_pointer", func(t *testing.T) {
			identity := matrix.NewIdentityMatrix()
			defer identity.Close()

			// Multiply uses m.ptr - will crash if pointer is invalid
			result := m.Multiply(identity)
			defer result.Close()

			// Results should match original matrix
			assert.InDelta(t, initialXX, result.XX, 0.001, "Matrix data corrupted after stack reuse")
			assert.InDelta(t, initialYY, result.YY, 0.001)
			assert.InDelta(t, initialX0, result.X0, 0.001)
			assert.InDelta(t, initialY0, result.Y0, 0.001)
		})

		t.Run("transform_point_uses_internal_pointer", func(t *testing.T) {
			// TransformPoint uses m.ptr internally
			x, y := m.TransformPoint(1.0, 1.0)

			// With scale(2,3) and translate(10,20):
			// x_new = 2*1 + 10 = 12
			// y_new = 3*1 + 20 = 23
			assert.InDelta(t, 12.0, x, 0.001, "TransformPoint failed - pointer may be invalid")
			assert.InDelta(t, 23.0, y, 0.001)
		})

		t.Run("cached_values_still_accessible", func(t *testing.T) {
			// The Go fields should still match (they were copied)
			assert.InDelta(t, initialXX, m.XX, 0.001)
			assert.InDelta(t, initialYY, m.YY, 0.001)
			assert.InDelta(t, initialX0, m.X0, 0.001)
			assert.InDelta(t, initialY0, m.Y0, 0.001)
		})

		// Note: defer m.Close() at top of function ensures cleanup
		// We also test that GC doesn't crash with proper heap allocation
		runtime.GC()
		runtime.Gosched()
		runtime.GC()
	})

	t.Run("matrix_operations_after_multiple_stack_frames", func(t *testing.T) {
		matrices := make([]*matrix.Matrix, 5)

		// Get multiple matrices through different stack frames
		for i := 0; i < 5; i++ {
			ctx.IdentityMatrix()
			ctx.Translate(float64(i*10), float64(i*20))

			m, err := ctx.GetMatrix()
			require.NoError(t, err)
			matrices[i] = m
		}

		// Now use all the matrices - if they have stack pointers, they'll be corrupted
		for i, m := range matrices {
			x, y := m.TransformPoint(1.0, 1.0)

			expectedX := float64(i*10) + 1.0
			expectedY := float64(i*20) + 1.0

			assert.InDelta(t, expectedX, x, 0.001, "Matrix %d corrupted", i)
			assert.InDelta(t, expectedY, y, 0.001, "Matrix %d corrupted", i)
		}

		// Clean up
		for _, m := range matrices {
			_ = m.Close()
		}
	})
}

// TestContextGetMatrixStackCorruption demonstrates the stack corruption bug
// through assertion failures rather than crashes.
//
// This test shows that GetMatrix() returns a matrix with an invalid stack pointer.
// After the stack is reused, operations that access m.ptr will produce corrupted
// results, causing the assertions to fail.
//
// Expected behavior:
//   - Current buggy implementation: FAILS with incorrect values
//   - After fix (heap allocation): PASSES with correct values
func TestContextGetMatrixStackCorruption(t *testing.T) {
	surf, err := surface.NewImageSurface(surface.FormatARGB32, 200, 200)
	require.NoError(t, err, "Failed to create surface")
	defer surf.Close()

	ctx, err := NewContext(surf)
	require.NoError(t, err, "Failed to create context")
	defer ctx.Close()

	// Set up a known transformation: translate(100, 200) then scale(2, 3)
	ctx.IdentityMatrix()
	ctx.Translate(100.0, 200.0)
	ctx.Scale(2.0, 3.0)

	// Get the matrix - if bug exists, this stores a STACK pointer
	m, err := ctx.GetMatrix()
	require.NoError(t, err)
	require.NotNil(t, m)
	defer m.Close()

	// Expected values for the transformation matrix
	// Scale comes first in matrix, then translation offset
	expectedXX := 2.0
	expectedYY := 3.0
	expectedX0 := 100.0
	expectedY0 := 200.0

	// Verify initial cached Go values are correct
	assert.InDelta(t, expectedXX, m.XX, 0.001, "Initial XX value should be correct")
	assert.InDelta(t, expectedYY, m.YY, 0.001, "Initial YY value should be correct")
	assert.InDelta(t, expectedX0, m.X0, 0.001, "Initial X0 value should be correct")
	assert.InDelta(t, expectedY0, m.Y0, 0.001, "Initial Y0 value should be correct")

	// Now corrupt the stack by calling a function that will reuse the same stack space
	// and overwrite it with known values that will cause assertions to fail
	overwriteStack := func() {
		// Call GetMatrix again - this will allocate a new cairo_matrix_t on the stack
		// in potentially the same location as the previous one
		ctx.IdentityMatrix()
		ctx.Translate(999.0, 888.0) // Very different values
		ctx.Scale(77.0, 66.0)        // Very different values

		m2, _ := ctx.GetMatrix()
		defer m2.Close()

		// Create more matrices to thoroughly overwrite stack
		for i := 0; i < 50; i++ {
			temp := matrix.NewIdentityMatrix()
			temp.Translate(float64(i)*111, float64(i)*222)
			temp.Scale(float64(i)+7, float64(i)+13)
			_ = temp.Close()
		}

		// Use large stack allocations
		largeArray := make([]byte, 4096)
		for i := range largeArray {
			largeArray[i] = byte(i % 256)
		}
		_ = largeArray
	}

	// Execute stack corruption
	overwriteStack()

	// Perform more context operations that use different stack frames
	for i := 0; i < 10; i++ {
		ctx.Save()
		ctx.Translate(float64(i)*5, float64(i)*10)
		ctx.Scale(1.5, 2.5)
		ctx.Restore()
	}

	// NOW: Try to use the matrix's internal C pointer
	// With the bug, m.ptr points to corrupted/overwritten stack memory
	// This will produce incorrect results

	t.Run("transform_point_with_corrupted_pointer", func(t *testing.T) {
		// TransformPoint uses m.ptr internally via C.cairo_matrix_transform_point
		// Test point: (10, 20)
		// Expected: x_new = 2*10 + 100 = 120, y_new = 3*20 + 200 = 260
		x, y := m.TransformPoint(10.0, 20.0)

		// These assertions WILL FAIL with the buggy implementation
		// because m.ptr points to corrupted stack memory
		assert.InDelta(t, 120.0, x, 0.001,
			"TransformPoint X failed - matrix pointer is corrupted! Expected 120.0, got %v", x)
		assert.InDelta(t, 260.0, y, 0.001,
			"TransformPoint Y failed - matrix pointer is corrupted! Expected 260.0, got %v", y)
	})

	t.Run("multiply_with_corrupted_pointer", func(t *testing.T) {
		// Multiply uses m.ptr via C.cairo_matrix_multiply
		identity := matrix.NewIdentityMatrix()
		defer identity.Close()

		result := m.Multiply(identity)
		defer result.Close()

		// Result should equal original matrix since we're multiplying by identity
		// These assertions WILL FAIL with the buggy implementation
		assert.InDelta(t, expectedXX, result.XX, 0.001,
			"Multiply XX failed - matrix pointer is corrupted! Expected %v, got %v", expectedXX, result.XX)
		assert.InDelta(t, expectedYY, result.YY, 0.001,
			"Multiply YY failed - matrix pointer is corrupted! Expected %v, got %v", expectedYY, result.YY)
		assert.InDelta(t, expectedX0, result.X0, 0.001,
			"Multiply X0 failed - matrix pointer is corrupted! Expected %v, got %v", expectedX0, result.X0)
		assert.InDelta(t, expectedY0, result.Y0, 0.001,
			"Multiply Y0 failed - matrix pointer is corrupted! Expected %v, got %v", expectedY0, result.Y0)
	})

	t.Run("transform_distance_with_corrupted_pointer", func(t *testing.T) {
		// TransformDistance uses m.ptr via C.cairo_matrix_transform_distance
		// Distance (10, 20) should be scaled but NOT translated
		// Expected: dx = 2*10 = 20, dy = 3*20 = 60
		dx, dy := m.TransformDistance(10.0, 20.0)

		// These assertions WILL FAIL with the buggy implementation
		assert.InDelta(t, 20.0, dx, 0.001,
			"TransformDistance DX failed - matrix pointer is corrupted! Expected 20.0, got %v", dx)
		assert.InDelta(t, 60.0, dy, 0.001,
			"TransformDistance DY failed - matrix pointer is corrupted! Expected 60.0, got %v", dy)
	})

	// Note: We deliberately DO NOT call runtime.GC() here
	// We want to see the assertion failures, not the finalizer crash
}
