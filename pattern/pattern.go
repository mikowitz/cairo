package pattern

import (
	"runtime"
	"sync"
	"unsafe"

	"github.com/mikowitz/cairo/matrix"
	"github.com/mikowitz/cairo/status"
)

type Pattern interface {
	Close() error
	Status() status.Status
	SetMatrix(m *matrix.Matrix)
	GetMatrix() (*matrix.Matrix, error)
	Ptr() unsafe.Pointer
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
		return &SolidPattern{
			BasePattern: basePattern,
		}
	}
}
