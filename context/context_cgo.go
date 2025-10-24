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
