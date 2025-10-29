package pattern

import (
	"runtime"
	"sync"
	"unsafe"

	"github.com/mikowitz/cairo/matrix"
	"github.com/mikowitz/cairo/status"
)

// Pattern is the interface that all Cairo pattern types implement.
//
// Patterns represent the "paint" that Cairo uses for drawing operations.
// They define what colors, gradients, or images to use when filling or
// stroking paths.
//
// All pattern types (solid colors, gradients, surface patterns) implement
// this interface, providing consistent methods for resource management,
// status checking, and transformation.
//
// Pattern implementations are safe for concurrent use from multiple goroutines.
type Pattern interface {
	// Close releases the Cairo resources associated with this pattern.
	//
	// After calling Close, the pattern should not be used for any operations.
	// Calling Close multiple times is safe (subsequent calls are no-ops).
	//
	// While patterns have finalizers that will eventually clean up resources,
	// explicit Close() calls are strongly recommended for predictable resource
	// management, especially in long-running programs.
	//
	// Returns an error if cleanup fails, though this is rare in practice.
	Close() error

	// Status returns the current status of the pattern.
	//
	// Returns status.Success if the pattern is valid and ready to use.
	// Returns status.NullPointer if the pattern has been closed.
	// Returns other status codes if the pattern creation failed or became invalid.
	//
	// This method is safe to call even after Close().
	Status() status.Status

	// SetMatrix sets the pattern's transformation matrix.
	//
	// The pattern matrix is used to transform the pattern coordinate space
	// before it is applied to the surface. This affects how gradients are
	// positioned and oriented, or how texture patterns are scaled and rotated.
	//
	// For solid color patterns, the matrix has no visible effect, but it can
	// still be set and retrieved.
	//
	// The transformation is independent of the Context's transformation matrix.
	// If m is nil, this method is a no-op.
	//
	// Example:
	//   m := matrix.NewScaleMatrix(2.0, 2.0)
	//   pattern.SetMatrix(m)  // Scale the pattern's coordinate space
	SetMatrix(m *matrix.Matrix)

	// GetMatrix returns the pattern's current transformation matrix.
	//
	// Returns the matrix that was previously set with SetMatrix, or the
	// identity matrix if no matrix was explicitly set.
	//
	// Returns an error (status.NullPointer) if the pattern has been closed.
	GetMatrix() (*matrix.Matrix, error)

	// Ptr returns the underlying Cairo pattern pointer.
	//
	// This method is primarily for internal use when passing patterns to
	// Cairo C API functions. Most users should not need to call this directly.
	//
	// The returned pointer is only valid while the pattern has not been closed.
	Ptr() unsafe.Pointer

	// GetType returns the pattern's type (solid, linear, radial, etc.).
	//
	// Pattern types are defined by the PatternType enumeration. This method
	// allows runtime type identification of patterns.
	//
	// Returns the PatternType value for this pattern (e.g., PatternTypeSolid).
	GetType() PatternType
}

type BasePattern struct {
	sync.RWMutex
	ptr         PatternPtr
	patternType PatternType
}

func newBasePattern(ptr PatternPtr, patternType PatternType) *BasePattern {
	b := &BasePattern{
		ptr:         ptr,
		patternType: patternType,
	}

	runtime.SetFinalizer(b, (*BasePattern).close)

	return b
}

func (b *BasePattern) Close() error {
	return b.close()
}

func (b *BasePattern) Status() status.Status {
	b.RLock()
	defer b.RUnlock()

	if b.ptr == nil {
		return status.NullPointer
	}
	return patternStatus(b.ptr)
}

func (b *BasePattern) SetMatrix(m *matrix.Matrix) {
	if m == nil {
		return
	}

	mPtr := m.Ptr() // This handles matrix un/locking internally
	b.Lock()
	defer b.Unlock()

	if b.ptr == nil {
		return
	}

	patternSetMatrix(b.ptr, mPtr)
}

func (b *BasePattern) GetMatrix() (*matrix.Matrix, error) {
	b.RLock()
	defer b.RUnlock()

	if b.ptr == nil {
		return nil, status.NullPointer
	}

	return patternGetMatrix(b.ptr)
}

func (b *BasePattern) Ptr() unsafe.Pointer {
	b.RLock()
	defer b.RUnlock()

	return unsafe.Pointer(b.ptr) //nolint:gosec
}

func (b *BasePattern) GetType() PatternType {
	b.RLock()
	defer b.RUnlock()

	return b.patternType
}

func (b *BasePattern) close() error {
	b.Lock()
	defer b.Unlock()

	if b.ptr != nil {
		patternClose(b.ptr)
		runtime.SetFinalizer(b, nil)
		b.ptr = nil
	}

	return nil
}

// PatternFromC wraps a C cairo_pattern_t pointer and returns the appropriate
// Go Pattern implementation based on the pattern's type.
//
// This function is primarily used internally when retrieving patterns from
// Cairo C API functions (e.g., cairo_get_source). It inspects the pattern's
// type and constructs the corresponding Go wrapper.
//
// Currently supported pattern types:
//   - PatternTypeSolid: Returns a *SolidPattern
//
// For unsupported pattern types (Linear, Radial, Mesh, RasterSource), this
// function currently defaults to returning a *SolidPattern as a temporary
// measure. This will be updated as additional pattern types are implemented.
//
// The returned Pattern takes ownership of the C pointer and will properly
// clean it up when Close() is called or when the finalizer runs.
func PatternFromC(uPtr unsafe.Pointer) Pattern {
	ptr := PatternPtr(uPtr)
	patternType := patternGetType(ptr)
	basePattern := newBasePattern(ptr, patternType)
	switch patternType {
	case PatternTypeSolid:
		return &SolidPattern{
			BasePattern: basePattern,
		}
	// TODO: Add cases for other pattern types as implemented
	default:
		return basePattern
	}
}
