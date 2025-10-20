// ABOUTME: Tests for the surface package interface and base types.
// ABOUTME: Written using TDD - these tests verify Surface interface and BaseSurface implementation.

package surface

import (
	"testing"

	"github.com/mikowitz/cairo/status"
)

// TestBaseSurfaceImplementsSurface verifies that BaseSurface implements the Surface interface.
// This is a compile-time check.
func TestBaseSurfaceImplementsSurface(t *testing.T) {
	// This is a compile-time check. If BaseSurface doesn't implement Surface,
	// this will fail to compile.
	var _ Surface = (*BaseSurface)(nil)
}

// TestBaseSurfaceStructure verifies the BaseSurface struct has the expected fields.
func TestBaseSurfaceStructure(t *testing.T) {
	// Create a zero-value BaseSurface to verify the struct compiles
	var bs BaseSurface

	// Verify initial state
	if bs.ptr != nil {
		t.Errorf("BaseSurface.ptr should be nil initially, got %v", bs.ptr)
	}

	if bs.closed != false {
		t.Errorf("BaseSurface.closed should be false initially, got %v", bs.closed)
	}
}

// TestBaseSurfaceClosedFlagManagement tests that the closed flag is properly managed.
// This test will fail until BaseSurface Close() is implemented.
func TestBaseSurfaceClosedFlagManagement(t *testing.T) {
	t.Skip("Skipping until BaseSurface.Close() CGO implementation is complete")

	// Expected behavior (to be implemented):
	// 1. New surface should have closed = false
	// 2. After Close(), closed should be true
	// 3. Second Close() should be safe (idempotent)
	// 4. Close() on already-closed surface returns nil (not error)
}

// TestBaseSurfaceCloseIdempotent tests that Close() can be called multiple times safely.
// This follows the Go pattern where Close() should be idempotent.
func TestBaseSurfaceCloseIdempotent(t *testing.T) {
	t.Skip("Skipping until BaseSurface.Close() CGO implementation is complete")

	// Expected behavior:
	// surface := createTestSurface(t)
	// err1 := surface.Close()
	// if err1 != nil {
	//     t.Errorf("First Close() failed: %v", err1)
	// }
	//
	// err2 := surface.Close()
	// if err2 != nil {
	//     t.Errorf("Second Close() should be safe, got error: %v", err2)
	// }
}

// TestBaseSurfaceStatusMethod tests the Status() method.
func TestBaseSurfaceStatusMethod(t *testing.T) {
	t.Skip("Skipping until BaseSurface.Status() CGO implementation is complete")

	// Expected behavior:
	// 1. New surface should have StatusSuccess
	// 2. Closed surface should report appropriate status
	// 3. Status() should be callable multiple times
	// 4. Status() should use RLock (not Lock) for thread safety
}

// TestBaseSurfaceFlushMethod tests the Flush() method.
func TestBaseSurfaceFlushMethod(t *testing.T) {
	t.Skip("Skipping until BaseSurface.Flush() CGO implementation is complete")

	// Expected behavior:
	// 1. Flush() should complete without error on valid surface
	// 2. Flush() should check closed flag before calling Cairo
	// 3. Flush() on closed surface should handle gracefully
	// 4. Flush() should use Lock for thread safety
}

// TestBaseSurfaceMarkDirtyMethod tests the MarkDirty() method.
func TestBaseSurfaceMarkDirtyMethod(t *testing.T) {
	t.Skip("Skipping until BaseSurface.MarkDirty() CGO implementation is complete")

	// Expected behavior:
	// 1. MarkDirty() should complete without error on valid surface
	// 2. MarkDirty() should check closed flag
	// 3. MarkDirty() should use Lock for thread safety
}

// TestBaseSurfaceMarkDirtyRectangleMethod tests the MarkDirtyRectangle() method.
func TestBaseSurfaceMarkDirtyRectangleMethod(t *testing.T) {
	t.Skip("Skipping until BaseSurface.MarkDirtyRectangle() CGO implementation is complete")

	// Expected behavior:
	// 1. MarkDirtyRectangle() should accept x, y, width, height parameters
	// 2. Should check closed flag before calling Cairo
	// 3. Should use Lock for thread safety
	// 4. Should handle edge cases (zero width/height, negative values)
}

// TestBaseSurfaceThreadSafety tests concurrent access to BaseSurface.
func TestBaseSurfaceThreadSafety(t *testing.T) {
	t.Skip("Skipping until BaseSurface methods are implemented")

	// Expected behavior:
	// Test concurrent:
	// - Status() calls (should use RLock)
	// - Flush() calls (should use Lock)
	// - MarkDirty() calls (should use Lock)
	// - Close() calls (should use Lock)
	// - Mixed read/write operations
	// No races should be detected when run with -race flag
}

// TestSurfaceInterfaceContract documents and tests the Surface interface contract.
func TestSurfaceInterfaceContract(t *testing.T) {
	// Verify the Surface interface has all required methods
	// This is mostly a documentation test - if the interface changes,
	// this test should be updated to reflect the new contract.

	tests := []struct {
		methodName string
		signature  string
	}{
		{"Close", "Close() error"},
		{"Status", "Status() status.Status"},
		{"Flush", "Flush()"},
		{"MarkDirty", "MarkDirty()"},
		{"MarkDirtyRectangle", "MarkDirtyRectangle(x, y, width, height int)"},
	}

	// This test documents the expected interface.
	// The actual interface compliance is checked at compile time.
	for _, tt := range tests {
		t.Run(tt.methodName, func(t *testing.T) {
			t.Logf("Surface interface should have method: %s", tt.signature)
		})
	}
}

// TestBaseSurfaceStatusReturnType verifies Status() returns the correct type.
func TestBaseSurfaceStatusReturnType(t *testing.T) {
	t.Skip("Skipping until BaseSurface.Status() is implemented")

	// Expected behavior:
	// surface := createTestSurface(t)
	// defer surface.Close()
	//
	// s := surface.Status()
	// // Verify s is of type status.Status
	// if s != status.StatusSuccess {
	//     t.Errorf("New surface should have StatusSuccess, got %v", s)
	// }
}

// TestBaseSurfaceOperationsOnClosedSurface tests that operations on closed surfaces are handled properly.
func TestBaseSurfaceOperationsOnClosedSurface(t *testing.T) {
	t.Skip("Skipping until BaseSurface methods are implemented")

	// Expected behavior (to be implemented):
	// surface := createTestSurface(t)
	// surface.Close()
	//
	// // All operations should either:
	// // 1. Return an error (for methods that return errors)
	// // 2. Handle gracefully without panicking (for void methods)
	//
	// // Status should still be callable
	// _ = surface.Status()
	//
	// // Flush should handle closed surface
	// surface.Flush() // Should not panic
	//
	// // MarkDirty should handle closed surface
	// surface.MarkDirty() // Should not panic
	//
	// // MarkDirtyRectangle should handle closed surface
	// surface.MarkDirtyRectangle(0, 0, 10, 10) // Should not panic
}

// TestBaseSurfaceNilPointerSafety tests that BaseSurface methods handle nil pointers safely.
func TestBaseSurfaceNilPointerSafety(t *testing.T) {
	t.Skip("Skipping until BaseSurface methods are implemented")

	// Expected behavior:
	// var bs *BaseSurface
	//
	// // All methods should handle nil receiver gracefully
	// // Either by checking for nil or by panicking with clear error
	//
	// Methods to test:
	// - bs.Close() should not crash (or return clear error)
	// - bs.Status() should return appropriate status
	// - bs.Flush() should not crash
	// - bs.MarkDirty() should not crash
	// - bs.MarkDirtyRectangle(0, 0, 10, 10) should not crash
}

// TestBaseSurfaceRWMutexEmbedding verifies that RWMutex is properly embedded.
func TestBaseSurfaceRWMutexEmbedding(t *testing.T) {
	var bs BaseSurface

	// Verify we can call Lock/Unlock directly on BaseSurface
	// because sync.RWMutex is embedded
	bs.Lock()
	bs.Unlock()

	bs.RLock()
	bs.RUnlock()

	// If this test compiles and runs, the embedding is correct
	t.Log("BaseSurface properly embeds sync.RWMutex")
}

// TestSurfaceInterfaceNilCheck verifies that nil Surface interface behaves correctly.
func TestSurfaceInterfaceNilCheck(t *testing.T) {
	var s Surface

	if s != nil {
		t.Errorf("Uninitialized Surface interface should be nil, got %v", s)
	}

	// Verify that we can assign a nil *BaseSurface to Surface
	var bs *BaseSurface
	s = bs

	if s == nil {
		t.Error("Surface interface should not be nil when holding nil *BaseSurface")
	}

	// The interface is not nil, but the underlying pointer is nil
	// This is an important Go distinction for interface types
	t.Log("Surface interface correctly handles nil *BaseSurface assignment")
}

// Helper function placeholder for when we can create actual surfaces
// This will be implemented once ImageSurface creation is available.
func createTestSurface(t *testing.T) Surface {
	t.Helper()
	t.Fatal("createTestSurface not yet implemented - needs ImageSurface from Prompt 7")
	return nil
}
