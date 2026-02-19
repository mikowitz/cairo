package pattern

import (
	"runtime"
	"sync"
	"unsafe"

	"github.com/mikowitz/cairo/matrix"
	"github.com/mikowitz/cairo/status"
)

// Extend defines how patterns behave outside their natural bounds.
//
// When a pattern (gradient or surface pattern) is used to paint an area
// larger than the pattern naturally covers, the extend mode determines
// what happens in the areas outside the pattern's bounds.
//
//go:generate stringer -type=Extend -trimprefix=Extend
type Extend int

const (
	// ExtendNone means the pattern is not painted outside its natural bounds.
	// Areas outside the pattern will be transparent.
	ExtendNone Extend = iota

	// ExtendRepeat means the pattern is tiled by repeating.
	// The pattern repeats infinitely in all directions.
	ExtendRepeat

	// ExtendReflect means the pattern is tiled by reflecting at the edges.
	// Creates a mirrored repetition effect.
	ExtendReflect

	// ExtendPad means the pattern extends by using the closest color from its edge.
	// The edge pixels are repeated infinitely outward.
	ExtendPad
)

// Filter defines the filtering algorithm used when sampling patterns.
//
// When a pattern is transformed (scaled, rotated), Cairo needs to resample
// the pattern pixels. The filter mode determines the quality and speed of
// this resampling operation.
//
//go:generate stringer -type=Filter -trimprefix=Filter
type Filter int

const (
	// FilterFast uses a high-performance filter with lower quality.
	// Equivalent to nearest-neighbor filtering.
	FilterFast Filter = iota

	// FilterGood balances quality and performance.
	// Uses bilinear interpolation.
	FilterGood

	// FilterBest uses the highest-quality filter available.
	// May be slower but produces the best visual results.
	FilterBest

	// FilterNearest uses nearest-neighbor sampling.
	// Fast but can produce pixelated results when scaling.
	FilterNearest

	// FilterBilinear uses bilinear interpolation.
	// Smoother than nearest-neighbor with reasonable performance.
	FilterBilinear
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

	// SetExtend sets how the pattern behaves outside its natural bounds.
	//
	// Cairo applies this to all pattern types. For gradients, it controls what
	// happens beyond the gradient's endpoints. For surface patterns, it controls
	// what happens outside the source surface.
	SetExtend(extend Extend)

	// GetExtend returns the current extend mode for the pattern.
	GetExtend() Extend

	// SetFilter sets the filtering algorithm used when sampling the pattern.
	//
	// Cairo applies this to all pattern types. It is most visually significant
	// for surface patterns and transformed patterns.
	SetFilter(filter Filter)

	// GetFilter returns the current filter mode for the pattern.
	GetFilter() Filter
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

func (b *BasePattern) SetExtend(extend Extend) {
	b.Lock()
	defer b.Unlock()

	if b.ptr == nil {
		return
	}

	patternSetExtend(b.ptr, extend)
}

func (b *BasePattern) GetExtend() Extend {
	b.RLock()
	defer b.RUnlock()

	if b.ptr == nil {
		return ExtendNone
	}

	return patternGetExtend(b.ptr)
}

func (b *BasePattern) SetFilter(filter Filter) {
	b.Lock()
	defer b.Unlock()

	if b.ptr == nil {
		return
	}

	patternSetFilter(b.ptr, filter)
}

func (b *BasePattern) GetFilter() Filter {
	b.RLock()
	defer b.RUnlock()

	if b.ptr == nil {
		return FilterGood
	}

	return patternGetFilter(b.ptr)
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
//   - PatternTypeSurface: Returns a *SurfacePattern
//
// For unsupported pattern types (Linear, Radial, Mesh, RasterSource), this
// function currently defaults to returning the base pattern implementation.
// This will be updated as additional pattern types are implemented.
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
	case PatternTypeSurface:
		return &SurfacePattern{
			BasePattern: basePattern,
		}
	// TODO: Add cases for other pattern types as implemented
	default:
		return basePattern
	}
}
