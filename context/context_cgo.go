package context

// #cgo pkg-config: cairo
// #include <cairo.h>
import "C"

import (
	"unsafe"

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
