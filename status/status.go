//go:generate stringer -type=Status
package status

import "fmt"

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

func (s Status) Error() string {
	return s.toString()
}

func (s Status) ToError() error {
	if s == Success {
		return nil
	}
	return fmt.Errorf("%v", s)
}
