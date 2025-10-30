package context

// #cgo pkg-config: cairo
// #include <cairo.h>
// #include <stdlib.h>
import "C"

import (
	"unsafe"

	"github.com/mikowitz/cairo/matrix"
	"github.com/mikowitz/cairo/pattern"
	"github.com/mikowitz/cairo/status"
)

type ContextPtr *C.cairo_t

func contextCreate(sPtr unsafe.Pointer) ContextPtr {
	return ContextPtr(C.cairo_create((*C.cairo_surface_t)(sPtr)))
}

func contextStatus(ptr ContextPtr) status.Status {
	return status.Status(C.cairo_status(ptr))
}

func contextClose(ptr ContextPtr) {
	C.cairo_destroy(ptr)
}

func contextSave(ptr ContextPtr) {
	C.cairo_save(ptr)
}

func contextRestore(ptr ContextPtr) {
	C.cairo_restore(ptr)
}

func contextSetSourceRGB(ptr ContextPtr, r, g, b float64) {
	C.cairo_set_source_rgb(
		ptr,
		C.double(r), C.double(g), C.double(b),
	)
}

func contextSetSourceRGBA(ptr ContextPtr, r, g, b, a float64) {
	C.cairo_set_source_rgba(
		ptr,
		C.double(r), C.double(g), C.double(b), C.double(a),
	)
}

func contextMoveTo(ptr ContextPtr, x, y float64) {
	C.cairo_move_to(ptr, C.double(x), C.double(y))
}

func contextLineTo(ptr ContextPtr, x, y float64) {
	C.cairo_line_to(ptr, C.double(x), C.double(y))
}

func contextRectangle(ptr ContextPtr, x, y, width, height float64) {
	C.cairo_rectangle(
		ptr, C.double(x), C.double(y),
		C.double(width), C.double(height),
	)
}

func contextGetCurrentPoint(ptr ContextPtr) (float64, float64, error) {
	if !contextHasCurrentPoint(ptr) {
		return 0, 0, status.NoCurrentPoint
	}
	var x C.double
	var y C.double

	C.cairo_get_current_point(ptr, &x, &y)

	return float64(x), float64(y), nil
}

func contextHasCurrentPoint(ptr ContextPtr) bool {
	return int(C.cairo_has_current_point(ptr)) != 0
}

func contextNewPath(ptr ContextPtr) {
	C.cairo_new_path(ptr)
}

func contextClosePath(ptr ContextPtr) {
	C.cairo_close_path(ptr)
}

func contextNewSubPath(ptr ContextPtr) {
	C.cairo_new_sub_path(ptr)
}

func contextStroke(ptr ContextPtr) {
	C.cairo_stroke(ptr)
}

func contextStrokePreserve(ptr ContextPtr) {
	C.cairo_stroke_preserve(ptr)
}

func contextFill(ptr ContextPtr) {
	C.cairo_fill(ptr)
}

func contextFillPreserve(ptr ContextPtr) {
	C.cairo_fill_preserve(ptr)
}

func contextPaint(ptr ContextPtr) {
	C.cairo_paint(ptr)
}

func contextSetLineWidth(ptr ContextPtr, width float64) {
	C.cairo_set_line_width(ptr, C.double(width))
}

func contextGetLineWidth(ptr ContextPtr) float64 {
	return float64(C.cairo_get_line_width(ptr))
}

func contextGetSource(ptr ContextPtr) (pattern.Pattern, error) {
	patternPtr := C.cairo_get_source(ptr)
	C.cairo_pattern_reference(patternPtr)
	return pattern.PatternFromC(unsafe.Pointer(patternPtr)), nil
}

func contextSetSource(ptr ContextPtr, patternPtr unsafe.Pointer) {
	C.cairo_set_source(ptr, (*C.cairo_pattern_t)(patternPtr))
}

func contextIdentityMatrix(ptr ContextPtr) {
	C.cairo_identity_matrix(ptr)
}

func contextTranslate(ptr ContextPtr, tx, ty float64) {
	C.cairo_translate(ptr, C.double(tx), C.double(ty))
}

func contextScale(ptr ContextPtr, sx, sy float64) {
	C.cairo_scale(ptr, C.double(sx), C.double(sy))
}

func contextRotate(ptr ContextPtr, radians float64) {
	C.cairo_rotate(ptr, C.double(radians))
}

func contextTransform(ptr ContextPtr, mPtr unsafe.Pointer) {
	C.cairo_transform(ptr, (*C.cairo_matrix_t)(mPtr))
}

func contextGetMatrix(ptr ContextPtr) *matrix.Matrix {
	var mStack C.cairo_matrix_t
	C.cairo_get_matrix(ptr, &mStack)

	mHeap := (*C.cairo_matrix_t)(C.malloc(C.sizeof_cairo_matrix_t))
	*mHeap = mStack

	return matrix.FromPointer(unsafe.Pointer(mHeap))
}

func contextSetMatrix(ptr ContextPtr, mPtr unsafe.Pointer) {
	C.cairo_set_matrix(ptr, (*C.cairo_matrix_t)(mPtr))
}

func contextUserToDevice(ptr ContextPtr, x, y float64) (float64, float64) {
	rx := C.double(x)
	ry := C.double(y)

	C.cairo_user_to_device(ptr, &rx, &ry)

	return float64(rx), float64(ry)
}

func contextUserToDeviceDistance(ptr ContextPtr, dx, dy float64) (float64, float64) {
	rx := C.double(dx)
	ry := C.double(dy)

	C.cairo_user_to_device_distance(ptr, &rx, &ry)

	return float64(rx), float64(ry)
}

func contextDeviceToUser(ptr ContextPtr, x, y float64) (float64, float64) {
	rx := C.double(x)
	ry := C.double(y)

	C.cairo_device_to_user(ptr, &rx, &ry)

	return float64(rx), float64(ry)
}

func contextDeviceToUserDistance(ptr ContextPtr, dx, dy float64) (float64, float64) {
	rx := C.double(dx)
	ry := C.double(dy)

	C.cairo_device_to_user_distance(ptr, &rx, &ry)

	return float64(rx), float64(ry)
}
