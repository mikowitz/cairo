package pattern

// #cgo pkg-config: cairo
// #include <cairo.h>
import "C"

import (
	"unsafe"

	"github.com/mikowitz/cairo/matrix"
	"github.com/mikowitz/cairo/status"
)

type PatternPtr *C.cairo_pattern_t

func PatternFromC(uPtr unsafe.Pointer) Pattern {
	ptr := PatternPtr(uPtr)
	patternType := patternGetType(ptr)
	basePattern := newBasePattern(ptr, patternType)
	return &SolidPattern{
		BasePattern: basePattern,
	}
}

func patternClose(ptr PatternPtr) {
	C.cairo_pattern_destroy(ptr)
}

func patternStatus(ptr PatternPtr) status.Status {
	return status.Status(C.cairo_pattern_status(ptr))
}

func patternSetMatrix(ptr PatternPtr, mPtr unsafe.Pointer) {
	C.cairo_pattern_set_matrix(
		ptr,
		(*C.cairo_matrix_t)(mPtr),
	)
}

func patternGetMatrix(ptr PatternPtr) (*matrix.Matrix, error) {
	m := (*C.cairo_matrix_t)(C.malloc(C.sizeof_cairo_matrix_t))

	C.cairo_pattern_get_matrix(ptr, m)

	return matrix.FromPointer(unsafe.Pointer(m)), nil
}

func patternCreateRGB(r, g, b float64) PatternPtr {
	return C.cairo_pattern_create_rgb(
		C.double(r),
		C.double(g),
		C.double(b),
	)
}

func patternCreateRGBA(r, g, b, a float64) PatternPtr {
	return C.cairo_pattern_create_rgba(
		C.double(r),
		C.double(g),
		C.double(b),
		C.double(a),
	)
}

func patternGetType(ptr PatternPtr) PatternType {
	return PatternType(C.cairo_pattern_get_type(ptr))
}
