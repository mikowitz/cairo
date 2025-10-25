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
