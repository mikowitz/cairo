// ABOUTME: Defines the Status type mapping Cairo error codes to Go errors.
// ABOUTME: Provides descriptive Error() messages with actionable suggestions for common statuses.
//go:generate stringer -type=Status
package status

// Status is used to indicate errors that can occur when using
// Cairo. In some cases it is returned directly by functions.
// But when using `Context`, the last error, if any, is stored
// in the context and can be retrieved with `context.Status()`.
//
// New entries may be added in future versions. Use
// `status.ToString(err)` to get a human-readable
// representation of an error message.
type Status int

const (
	Success Status = iota
	NoMemory
	InvalidRestore
	InvalidPopGroup
	NoCurrentPoint
	InvalidMatrix
	InvalidStatus
	NullPointer
	InvalidString
	InvalidPathData
	ReadError
	WriteError
	SurfaceFinished
	SurfaceTypeMismatch
	PatternTypeMismatch
	InvalidContent
	InvalidFormat
	InvalidVisual
	FileNotFound
	InvalidDash
	InvalidDscComment
	InvalidIndex
	ClipNotRepresentable
	TempFileError
	InvalidStride
	FontTypeMismatch
	UserFontImmutable
	UserFontError
	NegativeCount
	InvalidClusters
	InvalidSlant
	InvalidWeight
	InvalidSize
	UserFontNotImplemented
	DeviceTypeMismatch
	DeviceError
	InvalidMeshConstruction
	DeviceFinished
	JBig2GlobalMissing
	PngError
	FreetypeError
	Win32GdiError
	TagError
	DWriteError
	SvgFontError
	LastStatus
)

// Error returns the status name followed by an actionable suggestion for
// common statuses (e.g. "NoCurrentPoint: call MoveTo before path operations
// like LineTo or Arc"). For statuses with no suggestion, it returns only the
// status name (e.g. "NullPointer"). Do not match on the exact string returned
// by Error(); use errors.Is instead.
func (s Status) Error() string {
	base := s.toString()
	if suggestion := statusSuggestion(s); suggestion != "" {
		return base + ": " + suggestion
	}
	return base
}

// statusSuggestion returns an actionable hint for common Cairo error statuses.
// It returns an empty string for statuses that have no suggestion.
func statusSuggestion(s Status) string {
	switch s {
	case NoCurrentPoint:
		return "call MoveTo before path operations like LineTo or Arc"
	case InvalidRestore:
		return "use Save/Restore in matching pairs; Restore was called without a prior Save"
	case SurfaceFinished:
		return "the surface was already destroyed via Close; avoid calling Close before drawing is complete"
	case NoMemory:
		return "system is out of memory; try reducing surface dimensions"
	case FileNotFound:
		return "verify the file path exists and is accessible"
	case WriteError:
		return "check file permissions and available disk space"
	case ReadError:
		return "check file permissions and that the file exists at the given path"
	case InvalidFormat:
		return "use a supported pixel format such as FormatARGB32 or FormatRGB24"
	case InvalidStride:
		return "stride must be at least width*bytes-per-pixel and properly aligned"
	case SurfaceTypeMismatch:
		return "the operation requires a different surface type"
	case PatternTypeMismatch:
		return "the operation requires a different pattern type"
	case InvalidMatrix:
		return "the transformation matrix is degenerate or singular"
	}
	return ""
}
