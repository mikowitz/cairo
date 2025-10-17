package status

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestStatusSuccess verifies that Success equals 0
func TestStatusSuccess(t *testing.T) {
	assert.Equal(t, 0, int(Success), "Success should equal 0")
}

// TestStatusError verifies that Error() returns non-empty strings for all non-success statuses
func TestStatusError(t *testing.T) {
	tests := []struct {
		name   string
		status Status
	}{
		{"NoMemory", NoMemory},
		{"InvalidRestore", InvalidRestore},
		{"InvalidPopGroup", InvalidPopGroup},
		{"NoCurrentPoint", NoCurrentPoint},
		{"InvalidMatrix", InvalidMatrix},
		{"InvalidStatus", InvalidStatus},
		{"NullPointer", NullPointer},
		{"InvalidString", InvalidString},
		{"InvalidPathData", InvalidPathData},
		{"ReadError", ReadError},
		{"WriteError", WriteError},
		{"SurfaceFinished", SurfaceFinished},
		{"SurfaceTypeMismatch", SurfaceTypeMismatch},
		{"PatternTypeMismatch", PatternTypeMismatch},
		{"InvalidContent", InvalidContent},
		{"InvalidFormat", InvalidFormat},
		{"InvalidVisual", InvalidVisual},
		{"FileNotFound", FileNotFound},
		{"InvalidDash", InvalidDash},
		{"InvalidDscComment", InvalidDscComment},
		{"InvalidIndex", InvalidIndex},
		{"ClipNotRepresentable", ClipNotRepresentable},
		{"TempFileError", TempFileError},
		{"InvalidStride", InvalidStride},
		{"FontTypeMismatch", FontTypeMismatch},
		{"UserFontImmutable", UserFontImmutable},
		{"UserFontError", UserFontError},
		{"NegativeCount", NegativeCount},
		{"InvalidClusters", InvalidClusters},
		{"InvalidSlant", InvalidSlant},
		{"InvalidWeight", InvalidWeight},
		{"InvalidSize", InvalidSize},
		{"UserFontNotImplemented", UserFontNotImplemented},
		{"DeviceTypeMismatch", DeviceTypeMismatch},
		{"DeviceError", DeviceError},
		{"InvalidMeshConstruction", InvalidMeshConstruction},
		{"DeviceFinished", DeviceFinished},
		{"JBig2GlobalMissing", JBig2GlobalMissing},
		{"PngError", PngError},
		{"FreetypeError", FreetypeError},
		{"Win32GdiError", Win32GdiError},
		{"TagError", TagError},
		{"DWriteError", DWriteError},
		{"SvgFontError", SvgFontError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errMsg := tt.status.Error()
			assert.NotEmpty(t, errMsg, "Error() should return non-empty string")
		})
	}
}

// TestStatusToError verifies that ToError() returns nil for Success and error for others
func TestStatusToError(t *testing.T) {
	// Success should return nil
	t.Run("Success", func(t *testing.T) {
		err := Success.ToError()
		assert.NoError(t, err, "Success.ToError() should return nil")
	})

	// All other statuses should return an error
	tests := []struct {
		name   string
		status Status
	}{
		{"NoMemory", NoMemory},
		{"InvalidRestore", InvalidRestore},
		{"InvalidPopGroup", InvalidPopGroup},
		{"NoCurrentPoint", NoCurrentPoint},
		{"InvalidMatrix", InvalidMatrix},
		{"InvalidStatus", InvalidStatus},
		{"NullPointer", NullPointer},
		{"InvalidString", InvalidString},
		{"InvalidPathData", InvalidPathData},
		{"ReadError", ReadError},
		{"WriteError", WriteError},
		{"SurfaceFinished", SurfaceFinished},
		{"SurfaceTypeMismatch", SurfaceTypeMismatch},
		{"PatternTypeMismatch", PatternTypeMismatch},
		{"InvalidContent", InvalidContent},
		{"InvalidFormat", InvalidFormat},
		{"InvalidVisual", InvalidVisual},
		{"FileNotFound", FileNotFound},
		{"InvalidDash", InvalidDash},
		{"InvalidDscComment", InvalidDscComment},
		{"InvalidIndex", InvalidIndex},
		{"ClipNotRepresentable", ClipNotRepresentable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.status.ToError()
			assert.Error(t, err, "ToError() should return non-nil error")
		})
	}
}

// TestCGOStatusConversion verifies round-trip conversion between C and Go status codes
func TestCGOStatusConversion(t *testing.T) {
	tests := []struct {
		name   string
		status Status
	}{
		{"Success", Success},
		{"NoMemory", NoMemory},
		{"InvalidRestore", InvalidRestore},
		{"InvalidPopGroup", InvalidPopGroup},
		{"NoCurrentPoint", NoCurrentPoint},
		{"InvalidMatrix", InvalidMatrix},
		{"InvalidStatus", InvalidStatus},
		{"NullPointer", NullPointer},
		{"InvalidString", InvalidString},
		{"InvalidPathData", InvalidPathData},
		{"ReadError", ReadError},
		{"WriteError", WriteError},
		{"SurfaceFinished", SurfaceFinished},
		{"SurfaceTypeMismatch", SurfaceTypeMismatch},
		{"PatternTypeMismatch", PatternTypeMismatch},
		{"InvalidContent", InvalidContent},
		{"InvalidFormat", InvalidFormat},
		{"InvalidVisual", InvalidVisual},
		{"FileNotFound", FileNotFound},
		{"InvalidDash", InvalidDash},
		{"InvalidDscComment", InvalidDscComment},
		{"InvalidIndex", InvalidIndex},
		{"ClipNotRepresentable", ClipNotRepresentable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test round-trip: Go -> C -> Go
			cStatus := tt.status.toC()
			goStatus := statusFromC(cStatus)
			assert.Equal(t, tt.status, goStatus, "Round-trip conversion should preserve status value")

			// Verify that C.cairo_status_to_string returns non-empty string for all statuses
			statusString := tt.status.toString()
			assert.NotEmpty(t, statusString, "Cairo C library should provide status string")
		})
	}
}
