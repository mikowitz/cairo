package surface

import (
	"runtime"
	"sync"
	"testing"

	"github.com/mikowitz/cairo/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Temporary type declarations to allow tests to compile
// These will be removed once the actual types are implemented in surface.go

// ImageSurface is a temporary placeholder
// TODO: Remove this once ImageSurface is implemented in surface.go
// type ImageSurface struct{}
//
// // NewImageSurface is a temporary placeholder
// // TODO: Remove this once NewImageSurface is implemented in surface.go
// func NewImageSurface(format Format, width, height int) (*ImageSurface, error) {
// 	// Return a stub implementation that won't panic
// 	// This allows the test file to compile but tests are skipped anyway
// 	return &ImageSurface{}, nil
// }
//
// // Temporary methods to satisfy the interface
// // TODO: Remove these once real implementation exists
//
// func (i *ImageSurface) Close() error                               { return nil }
// func (i *ImageSurface) Status() status.Status                      { return status.Success }
// func (i *ImageSurface) Flush()                                     {}
// func (i *ImageSurface) MarkDirty()                                 {}
// func (i *ImageSurface) MarkDirtyRectangle(x, y, width, height int) {}
// func (i *ImageSurface) GetFormat() Format                          { return FormatInvalid }
// func (i *ImageSurface) GetWidth() int                              { return 0 }
// func (i *ImageSurface) GetHeight() int                             { return 0 }
// func (i *ImageSurface) GetStride() int                             { return 0 }
// func (i *ImageSurface) WriteToPNG(filename string) error           { return nil }

// TestSurfaceInterfaceCompleteness verifies the Surface interface exists
// and can be satisfied by BaseSurface
func TestSurfaceInterfaceCompleteness(t *testing.T) {
	// Verify that BaseSurface implements Surface interface
	var _ Surface = (*BaseSurface)(nil)
}

// TestBaseSurfaceCreation tests that BaseSurface can be instantiated
func TestBaseSurfaceCreation(t *testing.T) {
	// Create a test surface (will use ImageSurface once available)
	s := createTestSurface(t, FormatARGB32, 100, 100)
	defer func() {
		err := s.Close()
		require.NoError(t, err, "Surface should close without error")
	}()

	// Verify surface is not nil
	require.NotNil(t, s, "Surface should not be nil")

	// Verify initial status is success
	st := s.Status()
	assert.Equal(t, status.Success, st, "Initial status should be Success")
}

// TestBaseSurfaceClose verifies Close() method behavior
func TestBaseSurfaceClose(t *testing.T) {
	s := createTestSurface(t, FormatARGB32, 100, 100)

	// First close should succeed
	err := s.Close()
	assert.NoError(t, err, "First Close() should succeed")

	// Second close should be idempotent (no error)
	err = s.Close()
	assert.NoError(t, err, "Second Close() should be idempotent and not error")

	// Third close to ensure truly idempotent
	err = s.Close()
	assert.NoError(t, err, "Third Close() should still be safe")
}

// TestBaseSurfaceStatus verifies Status() method behavior
func TestBaseSurfaceStatus(t *testing.T) {
	s := createTestSurface(t, FormatARGB32, 100, 100)

	// Status should work before close
	st := s.Status()
	assert.Equal(t, status.Success, st, "Status before close should be Success")

	// Close the surface
	err := s.Close()
	require.NoError(t, err)

	// Status should still be callable after close
	st = s.Status()
	// The status might be Success or SurfaceFinished depending on implementation
	assert.NotEqual(t, status.InvalidStatus, st, "Status after close should still be valid")
}

// TestBaseSurfaceFlush verifies Flush() method behavior
func TestBaseSurfaceFlush(t *testing.T) {
	s := createTestSurface(t, FormatARGB32, 100, 100)
	defer func() {
		err := s.Close()
		require.NoError(t, err, "Surface should close without error")
	}()

	// Flush should complete without error on valid surface
	s.Flush()
	st := s.Status()
	assert.Equal(t, status.Success, st, "Status after Flush should be Success")

	// Multiple flushes should be safe
	s.Flush()
	s.Flush()
	st = s.Status()
	assert.Equal(t, status.Success, st, "Status after multiple Flushes should be Success")
}

// TestBaseSurfaceMarkDirty verifies MarkDirty() method behavior
func TestBaseSurfaceMarkDirty(t *testing.T) {
	s := createTestSurface(t, FormatARGB32, 100, 100)
	defer func() {
		err := s.Close()
		require.NoError(t, err, "Surface should close without error")
	}()

	// MarkDirty should be callable on valid surface
	s.MarkDirty()
	st := s.Status()
	assert.Equal(t, status.Success, st, "Status after MarkDirty should be Success")

	// Multiple MarkDirty calls should be safe
	s.MarkDirty()
	s.MarkDirty()
	st = s.Status()
	assert.Equal(t, status.Success, st, "Status after multiple MarkDirty calls should be Success")
}

// TestBaseSurfaceMarkDirtyRectangle verifies MarkDirtyRectangle() method behavior
func TestBaseSurfaceMarkDirtyRectangle(t *testing.T) {
	s := createTestSurface(t, FormatARGB32, 100, 100)
	defer func() {
		err := s.Close()
		require.NoError(t, err, "Surface should close without error")
	}()

	tests := []struct {
		name   string
		x      int
		y      int
		width  int
		height int
	}{
		{"valid rectangle", 10, 10, 50, 50},
		{"zero origin", 0, 0, 50, 50},
		{"full surface", 0, 0, 100, 100},
		{"single pixel", 50, 50, 1, 1},
		{"negative coordinates", -10, -10, 50, 50}, // Cairo should handle this
		{"zero dimensions", 10, 10, 0, 0},          // Cairo should handle this
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.MarkDirtyRectangle(tt.x, tt.y, tt.width, tt.height)
			// Cairo handles edge cases internally, just verify no crash
			st := s.Status()
			// Status should still be valid (may or may not be Success depending on parameters)
			assert.NotEqual(t, status.InvalidStatus, st, "Status should be valid after MarkDirtyRectangle")
		})
	}
}

// TestBaseSurfaceClosedState verifies operations on closed surface
func TestBaseSurfaceClosedState(t *testing.T) {
	s := createTestSurface(t, FormatARGB32, 100, 100)

	// Close the surface
	err := s.Close()
	require.NoError(t, err)

	// Status should still work on closed surface (no panic)
	st := s.Status()
	assert.NotEqual(t, status.InvalidStatus, st, "Status should work on closed surface")

	// Flush on closed surface should be safe (no panic)
	// The behavior might vary - either no-op or set an error status
	s.Flush()
	st = s.Status()
	assert.NotEqual(t, status.InvalidStatus, st, "Status after Flush on closed surface should be valid")

	// MarkDirty on closed surface should be safe (no panic)
	s.MarkDirty()
	st = s.Status()
	assert.NotEqual(t, status.InvalidStatus, st, "Status after MarkDirty on closed surface should be valid")

	// MarkDirtyRectangle on closed surface should be safe (no panic)
	s.MarkDirtyRectangle(0, 0, 10, 10)
	st = s.Status()
	assert.NotEqual(t, status.InvalidStatus, st, "Status after MarkDirtyRectangle on closed surface should be valid")
}

// TestBaseSurfaceThreadSafety verifies concurrent access is safe
// Run with: go test -race
func TestBaseSurfaceThreadSafety(t *testing.T) {
	s := createTestSurface(t, FormatARGB32, 100, 100)
	defer func() {
		err := s.Close()
		require.NoError(t, err, "Surface should close without error")
	}()

	const goroutines = 10
	const iterations = 100

	var wg sync.WaitGroup
	wg.Add(goroutines * 5) // 5 different operation types

	// Concurrent Status calls
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = s.Status()
			}
		}()
	}

	// Concurrent Flush calls
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				s.Flush()
			}
		}()
	}

	// Concurrent MarkDirty calls
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				s.MarkDirty()
			}
		}()
	}

	// Concurrent MarkDirtyRectangle calls
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				s.MarkDirtyRectangle(0, 0, 50, 50)
			}
		}()
	}

	// Mixed concurrent operations
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				switch j % 4 {
				case 0:
					_ = s.Status()
				case 1:
					s.Flush()
				case 2:
					s.MarkDirty()
				case 3:
					s.MarkDirtyRectangle(idx*10, idx*10, 10, 10)
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify surface is still in valid state
	st := s.Status()
	assert.Equal(t, status.Success, st, "Surface should still be valid after concurrent operations")
}

// TestBaseSurfaceThreadSafeClose verifies concurrent Close is safe
func TestBaseSurfaceThreadSafeClose(t *testing.T) {
	s := createTestSurface(t, FormatARGB32, 100, 100)

	const goroutines = 10
	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Multiple goroutines trying to close simultaneously
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			_ = s.Close()
		}()
	}

	wg.Wait()

	// All closes should have completed without race or panic
	// Final close should still be safe
	err := s.Close()
	assert.NoError(t, err, "Final Close after concurrent closes should be safe")
}

// TestBaseSurfaceFinalizer verifies finalizer setup
func TestBaseSurfaceFinalizer(t *testing.T) {
	// This test is inherently difficult to make deterministic
	// We can only verify that forgetting to Close doesn't leak memory over many iterations

	const iterations = 100

	for i := 0; i < iterations; i++ {
		// Create surface without calling Close
		// Finalizer should clean it up
		_ = createTestSurface(t, FormatARGB32, 100, 100)
	}

	// Force garbage collection
	runtime.GC()
	runtime.GC() // Call twice to ensure finalizers run

	// If we get here without crashing or leaking too much memory, finalizers are working
	// Note: This is not a perfect test, but it's the best we can do for finalizers

	// Now test that finalizer is safe after explicit Close
	s := createTestSurface(t, FormatARGB32, 100, 100)
	err := s.Close()
	require.NoError(t, err)

	// Clear reference and force GC - finalizer should be safe
	s = nil
	runtime.GC()
	runtime.GC()

	// If we get here, finalizer after Close is safe
}

// Placeholder tests for future ImageSurface implementation (Prompt 7)

// TestImageSurfaceCreation will test ImageSurface creation
func TestImageSurfaceCreation(t *testing.T) {
	// Test successful creation
	s, err := NewImageSurface(FormatARGB32, 100, 100)
	require.NotNil(t, s, "ImageSurface should not be nil")
	require.NoError(t, err, "Successful image surface creation should not return an error")
	defer func() {
		err := s.Close()
		require.NoError(t, err, "Surface should close without error")
	}()

	st := s.Status()
	assert.Equal(t, status.Success, st, "New surface should have Success status")
}

// TestImageSurfaceGetters will test ImageSurface getter methods
func TestImageSurfaceGetters(t *testing.T) {
	s, _ := NewImageSurface(FormatARGB32, 200, 150)
	defer func() {
		err := s.Close()
		require.NoError(t, err, "Surface should close without error")
	}()

	// Test GetFormat
	format := s.GetFormat()
	assert.Equal(t, FormatARGB32, format, "Format should match creation format")

	// Test GetWidth
	width := s.GetWidth()
	assert.Equal(t, 200, width, "Width should match creation width")

	// Test GetHeight
	height := s.GetHeight()
	assert.Equal(t, 150, height, "Height should match creation height")

	// Test GetStride
	stride := s.GetStride()
	expectedStride := FormatARGB32.StrideForWidth(200)
	assert.Equal(t, expectedStride, stride, "Stride should match expected value")
	assert.Greater(t, stride, 0, "Stride should be positive")
}

// TestImageSurfaceInvalidParameters will test error handling
func TestImageSurfaceInvalidParameters(t *testing.T) {
	tests := []struct {
		name           string
		format         Format
		width          int
		height         int
		expectedStatus status.Status
	}{
		{"invalid format", FormatInvalid, 100, 100, status.InvalidFormat},
		{"negative width", FormatARGB32, -100, 100, status.InvalidSize},
		{"negative height", FormatARGB32, 100, -100, status.InvalidSize},
		{"width too large", FormatARGB32, 1 << 15, 100, status.InvalidSize},
		{"height too large", FormatARGB32, 100, 1 << 15, status.InvalidSize},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewImageSurface(tt.format, tt.width, tt.height)

			assert.Nil(t, s, "NewImageSurface should not return a surface on error")
			assert.Equal(t, tt.expectedStatus, err, "Surface status should not be success on error")
		})
	}
}

// Placeholder tests for future PNG support (Prompt 8)

// TestImageSurfaceWriteToPNG will test PNG export
func TestImageSurfaceWriteToPNG(t *testing.T) {
	t.Skip("Will be implemented in Prompt 8 when PNG support is added")

	s, err := NewImageSurface(FormatARGB32, 100, 100)
	assert.NoError(t, err, "Should not return an error from successful image create")
	defer func() {
		err := s.Close()
		require.NoError(t, err, "Surface should close without error")
	}()

	// Flush before writing
	s.Flush()

	// Write to temporary file
	tmpDir := t.TempDir()
	filename := tmpDir + "/test.png"

	err = s.WriteToPNG(filename)
	assert.NoError(t, err, "WriteToPNG should succeed")

	// Verify file exists and has reasonable size
	// (actual verification would require reading the PNG)
}

// Helper functions for tests

// createTestSurface creates a test surface for use in tests
// This will be implemented once we have ImageSurface
func createTestSurface(t *testing.T, format Format, width, height int) Surface {
	t.Helper()

	s, err := NewImageSurface(format, width, height)
	require.NoError(t, err, "Test surface should not return error")
	require.NotNil(t, s, "Test surface should not be nil")

	return s
}

// // assertSurfaceValid verifies a surface is in valid state
// func assertSurfaceValid(t *testing.T, s Surface) {
// 	t.Helper()
// 	require.NotNil(t, s, "Surface should not be nil")
//
// 	st := s.Status()
// 	assert.Equal(t, status.Success, st, "Surface should be in success state")
// }
//
// // assertSurfaceClosed verifies a surface is properly closed
// func assertSurfaceClosed(t *testing.T, s Surface) {
// 	t.Helper()
// 	require.NotNil(t, s, "Surface should not be nil")
//
// 	// Status should still work after close, but might not be Success
// 	st := s.Status()
// 	assert.NotEqual(t, status.InvalidStatus, st, "Status should be valid even on closed surface")
// }
