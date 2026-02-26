package surface

// #cgo pkg-config: cairo
// #include <cairo.h>
// #include <stdlib.h>
import "C"

import (
	"unsafe"

	"github.com/mikowitz/cairo/status"
)

type SurfacePtr *C.cairo_surface_t

func surfaceClose(ptr SurfacePtr) {
	C.cairo_surface_destroy(ptr)
}

func surfaceStatus(ptr SurfacePtr) status.Status {
	return status.Status(C.cairo_surface_status(ptr))
}

func surfaceFlush(ptr SurfacePtr) {
	C.cairo_surface_flush(ptr)
}

func surfaceMarkDirty(ptr SurfacePtr) {
	C.cairo_surface_mark_dirty(ptr)
}

func surfaceMarkDirtyRectangle(ptr SurfacePtr, x, y, width, height int) {
	C.cairo_surface_mark_dirty_rectangle(
		ptr,
		C.int(x), C.int(y),
		C.int(width), C.int(height),
	)
}

func imageSurfaceCreate(format Format, width, height int) SurfacePtr {
	return SurfacePtr(
		C.cairo_image_surface_create(
			C.cairo_format_t(format),
			C.int(width),
			C.int(height),
		),
	)
}

func surfaceShowPage(ptr SurfacePtr) {
	C.cairo_surface_show_page(ptr)
}

func surfaceWriteToPNG(ptr SurfacePtr, filepath string) error {
	cFilepath := C.CString(filepath)
	defer C.free(unsafe.Pointer(cFilepath))
	st := C.cairo_surface_write_to_png(ptr, cFilepath)

	s := status.Status(st)
	if s == status.Success {
		return nil
	}
	return s
}
