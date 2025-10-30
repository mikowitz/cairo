package surface

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/mikowitz/cairo/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSurfaceInterfaceCompleteness verifies the Surface interface exists
// and can be satisfied by BaseSurface
func TestSurfaceInterfaceCompleteness(t *testing.T) {
	// Verify that BaseSurface implements Surface interface
	var _ Surface = (*BaseSurface)(nil)
}

// TestBaseSurfaceCreation tests that BaseSurface can be instantiated
func TestBaseSurfaceCreation(t *testing.T) {
	// Create a test surface (will use ImageSurface once available)
	s := createTestSurface(t)
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
	s := createTestSurface(t)

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
	s := createTestSurface(t)

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
	s := createTestSurface(t)
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
	s := createTestSurface(t)
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
	s := createTestSurface(t)
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
	s := createTestSurface(t)

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
	s := createTestSurface(t)
	defer func() {
		err := s.Close()
		require.NoError(t, err, "Surface should close without error")
	}()

	const goroutines = 10
	const iterations = 100

	var wg sync.WaitGroup
	wg.Add(goroutines * 5) // 5 different operation types

	// Concurrent Status calls
	for range goroutines {
		go func() {
			defer wg.Done()
			for range iterations {
				_ = s.Status()
			}
		}()
	}

	// Concurrent Flush calls
	for range goroutines {
		go func() {
			defer wg.Done()
			for range iterations {
				s.Flush()
			}
		}()
	}

	// Concurrent MarkDirty calls
	for range goroutines {
		go func() {
			defer wg.Done()
			for range iterations {
				s.MarkDirty()
			}
		}()
	}

	// Concurrent MarkDirtyRectangle calls
	for range goroutines {
		go func() {
			defer wg.Done()
			for range iterations {
				s.MarkDirtyRectangle(0, 0, 50, 50)
			}
		}()
	}

	// Mixed concurrent operations
	for i := range goroutines {
		go func(idx int) {
			defer wg.Done()
			for j := range iterations {
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
	s := createTestSurface(t)

	const goroutines = 10
	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Multiple goroutines trying to close simultaneously
	for range goroutines {
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

	for range iterations {
		// Create surface without calling Close
		// Finalizer should clean it up
		_ = createTestSurface(t)
	}

	// Force garbage collection
	runtime.GC()
	runtime.GC() // Call twice to ensure finalizers run

	// If we get here without crashing or leaking too much memory, finalizers are working
	// Note: This is not a perfect test, but it's the best we can do for finalizers

	// Now test that finalizer is safe after explicit Close
	s := createTestSurface(t)
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

// PNG Support Tests (Prompt 8)

// TestImageSurfaceWriteToPNG tests successful PNG writing
func TestImageSurfaceWriteToPNG(t *testing.T) {
	tests := []struct {
		name   string
		format Format
		width  int
		height int
	}{
		{"ARGB32 surface", FormatARGB32, 100, 100},
		{"RGB24 surface", FormatRGB24, 50, 50},
		{"A8 surface", FormatA8, 75, 75},
		{"A1 surface", FormatA1, 64, 64}, // A1 requires dimensions that work well with 1-bit packing
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for test output
			tmpDir := t.TempDir()
			filename := tmpDir + "/test_output.png"
			filepath.Clean(filename)

			// Create surface
			surf, err := NewImageSurface(tt.format, tt.width, tt.height)
			require.NoError(t, err, "Failed to create surface")
			require.NotNil(t, surf, "Surface should not be nil")
			defer func() {
				err := surf.Close()
				assert.NoError(t, err, "surface should close without error")
			}()

			// Flush the surface before writing (as documented)
			surf.Flush()

			// Write to PNG
			err = surf.WriteToPNG(filename)
			require.NoError(t, err, "WriteToPNG should succeed")

			// Verify file exists and is non-empty
			info, err := os.Stat(filename)
			require.NoError(t, err, "PNG file should exist")
			assert.Greater(t, info.Size(), int64(0), "PNG file should not be empty")

			// Verify it's a valid PNG by checking magic bytes
			file, err := os.Open(filename)
			require.NoError(t, err, "Should be able to open PNG file")
			defer func() {
				err := surf.Close()
				assert.NoError(t, err, "surface should close without error")
			}()

			magic := make([]byte, 8)
			_, err = file.Read(magic)
			require.NoError(t, err, "Should be able to read PNG header")

			// PNG magic bytes: 89 50 4E 47 0D 0A 1A 0A
			expectedMagic := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
			assert.Equal(t, expectedMagic, magic, "File should have valid PNG magic bytes")
		})
	}
}

// TestImageSurfaceWriteToPNGInvalidPath tests error handling with invalid paths
func TestImageSurfaceWriteToPNGInvalidPath(t *testing.T) {
	// Create a surface
	surf, err := NewImageSurface(FormatARGB32, 100, 100)
	require.NoError(t, err, "Failed to create surface")
	require.NotNil(t, surf, "Surface should not be nil")
	defer func() {
		err := surf.Close()
		assert.NoError(t, err, "surface should close without error")
	}()

	// Flush the surface
	surf.Flush()

	// Try to write to an invalid path (directory that doesn't exist)
	invalidPath := "/nonexistent/directory/that/should/not/exist/test.png"
	err = surf.WriteToPNG(invalidPath)
	assert.Error(t, err, "WriteToPNG should fail with invalid path")
}

// TestImageSurfaceWriteToPNGInvalidFilename tests error handling with invalid filenames
func TestImageSurfaceWriteToPNGInvalidFilename(t *testing.T) {
	// Create a surface
	surf, err := NewImageSurface(FormatARGB32, 100, 100)
	require.NoError(t, err, "Failed to create surface")
	require.NotNil(t, surf, "Surface should not be nil")
	defer func() {
		err := surf.Close()
		assert.NoError(t, err, "surface should close without error")
	}()

	// Flush the surface
	surf.Flush()

	tests := []struct {
		name     string
		filename string
		reason   string
	}{
		{
			name:     "empty string",
			filename: "",
			reason:   "Empty filename should fail",
		},
		{
			name:     "null byte at start",
			filename: "\x00test.png",
			reason:   "Filename starting with null byte should fail (truncates to empty string)",
		},
		{
			name:     "filename too long",
			filename: "/" + string(make([]byte, 4096)) + ".png", // Most filesystems have limits around 255-4096
			reason:   "Extremely long filename should fail",
		},
		{
			name:     "invalid directory path",
			filename: "/nonexistent/deeply/nested/directory/structure/that/does/not/exist/test.png",
			reason:   "Path with non-existent directories should fail",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := surf.WriteToPNG(tt.filename)
			assert.Equal(t, err, status.WriteError, tt.reason)
		})
	}
}

// TestImageSurfaceWriteToPNGNullByteHandling tests how null bytes in filenames are handled
// Note: C.CString truncates at the first null byte, so these tests verify that behavior
func TestImageSurfaceWriteToPNGNullByteHandling(t *testing.T) {
	// Create a surface
	surf, err := NewImageSurface(FormatARGB32, 100, 100)
	require.NoError(t, err, "Failed to create surface")
	require.NotNil(t, surf, "Surface should not be nil")
	defer func() {
		err := surf.Close()
		assert.NoError(t, err, "surface should close without error")
	}()

	surf.Flush()

	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		filename    string
		shouldError bool
		reason      string
	}{
		{
			name:        "null byte in middle - creates truncated file",
			filename:    tmpDir + "/test\x00ignored.png",
			shouldError: false, // C.CString truncates at null, so this tries to write tmpDir + "/test"
			reason:      "Null byte in middle truncates filename at null byte",
		},
		{
			name:        "null byte at end - same as valid filename",
			filename:    tmpDir + "/test.png\x00",
			shouldError: false, // C.CString truncates, so this is equivalent to tmpDir + "/test.png"
			reason:      "Null byte at end is truncated, making valid filename",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := surf.WriteToPNG(tt.filename)
			if tt.shouldError {
				assert.Error(t, err, tt.reason)
			} else {
				// These may succeed or fail depending on the truncated path validity
				// We're just documenting the behavior here
				t.Logf("WriteToPNG with null byte: error=%v (reason: %s)", err, tt.reason)
			}
		})
	}
}

// TestImageSurfaceWriteToPNGAfterClose tests error handling when writing after close
func TestImageSurfaceWriteToPNGAfterClose(t *testing.T) {
	// Create a temporary directory for test output
	tmpDir := t.TempDir()
	filename := tmpDir + "/test_output.png"

	// Create surface
	surf, err := NewImageSurface(FormatARGB32, 100, 100)
	require.NoError(t, err, "Failed to create surface")
	require.NotNil(t, surf, "Surface should not be nil")

	// Close the surface
	err = surf.Close()
	require.NoError(t, err, "Close should succeed")

	// Try to write to PNG after closing
	err = surf.WriteToPNG(filename)
	assert.Error(t, err, "WriteToPNG should fail on closed surface")
}

// TestImageSurfaceWriteToPNGMultipleTimes tests writing the same surface multiple times
func TestImageSurfaceWriteToPNGMultipleTimes(t *testing.T) {
	// Create a temporary directory for test output
	tmpDir := t.TempDir()

	// Create surface
	surf, err := NewImageSurface(FormatARGB32, 50, 50)
	require.NoError(t, err, "Failed to create surface")
	require.NotNil(t, surf, "Surface should not be nil")
	defer func() {
		err := surf.Close()
		assert.NoError(t, err, "surface should close without error")
	}()

	// Write to multiple PNG files
	for i := range 3 {
		filename := tmpDir + "/test_output_" + string(rune('0'+i)) + ".png"
		surf.Flush()
		err = surf.WriteToPNG(filename)
		require.NoError(t, err, "WriteToPNG should succeed on iteration %d", i)

		// Verify file exists
		info, err := os.Stat(filename)
		require.NoError(t, err, "PNG file should exist on iteration %d", i)
		assert.Greater(t, info.Size(), int64(0), "PNG file should not be empty on iteration %d", i)
	}
}

// TestImageSurfaceWriteToPNGWithDifferentFormats tests PNG writing with all supported formats
func TestImageSurfaceWriteToPNGWithDifferentFormats(t *testing.T) {
	formats := []Format{
		FormatARGB32,
		FormatRGB24,
		FormatA8,
		FormatA1,
		FormatRGB16_565,
		FormatRGB30,
	}

	for _, format := range formats {
		t.Run(format.String(), func(t *testing.T) {
			// Skip invalid format
			if format == FormatInvalid {
				t.Skip("Skipping invalid format")
			}

			tmpDir := t.TempDir()
			filename := tmpDir + "/test_" + format.String() + ".png"

			// Create surface with format-appropriate dimensions
			width, height := 100, 100
			if format == FormatA1 {
				// A1 format works better with dimensions that align well
				width, height = 64, 64
			}

			surf, err := NewImageSurface(format, width, height)
			require.NoError(t, err, "Failed to create surface with format %s", format)
			require.NotNil(t, surf, "Surface should not be nil")
			defer func() {
				err := surf.Close()
				assert.NoError(t, err, "surface should close without error")
			}()

			surf.Flush()
			err = surf.WriteToPNG(filename)
			require.NoError(t, err, "WriteToPNG should succeed for format %s", format)

			// Verify file exists and is non-empty
			info, err := os.Stat(filename)
			require.NoError(t, err, "PNG file should exist for format %s", format)
			assert.Greater(t, info.Size(), int64(0), "PNG file should not be empty for format %s", format)
		})
	}
}

// Helper functions for tests

// createTestSurface creates a test surface for use in tests
// This will be implemented once we have ImageSurface
func createTestSurface(t *testing.T) *ImageSurface {
	t.Helper()

	s, err := NewImageSurface(FormatARGB32, 100, 100)
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
