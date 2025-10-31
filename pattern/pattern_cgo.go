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
	var mStack C.cairo_matrix_t
	C.cairo_pattern_get_matrix(ptr, &mStack)

	mHeap := (*C.cairo_matrix_t)(C.malloc(C.sizeof_cairo_matrix_t))
	*mHeap = mStack

	return matrix.FromPointer(unsafe.Pointer(mHeap)), nil
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

func patternCreateLinear(x0, y0, x1, y1 float64) PatternPtr {
	return C.cairo_pattern_create_linear(
		C.double(x0), C.double(y0),
		C.double(x1), C.double(y1),
	)
}

func patternCreateRadial(cx0, cy0, radius0, cx1, cy1, radius1 float64) PatternPtr {
	return C.cairo_pattern_create_radial(
		C.double(cx0), C.double(cy0), C.double(radius0),
		C.double(cx1), C.double(cy1), C.double(radius1),
	)
}

func patternAddColorStopRGB(ptr PatternPtr, offset, r, g, b float64) {
	C.cairo_pattern_add_color_stop_rgb(ptr,
		C.double(offset),
		C.double(r), C.double(g), C.double(b),
	)
}

func patternAddColorStopRGBA(ptr PatternPtr, offset, r, g, b, a float64) {
	C.cairo_pattern_add_color_stop_rgba(ptr,
		C.double(offset),
		C.double(r), C.double(g), C.double(b), C.double(a),
	)
}

func patternGetColorStopCount(ptr PatternPtr) (int, status.Status) {
	var i C.int

	st := C.cairo_pattern_get_color_stop_count(ptr, &i)

	return int(i), status.Status(st)
}

func patternGetColorStopRGBA(ptr PatternPtr, index int) (float64, float64, float64, float64, float64, error) {
	var o, r, g, b, a C.double

	st := C.cairo_pattern_get_color_stop_rgba(ptr, C.int(index), &o, &r, &g, &b, &a)

	return float64(o), float64(r), float64(g), float64(b), float64(a), status.Status(st)
}

func patternGetType(ptr PatternPtr) PatternType {
	return PatternType(C.cairo_pattern_get_type(ptr))
}
