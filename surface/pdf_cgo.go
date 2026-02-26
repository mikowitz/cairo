// ABOUTME: CGO bindings for Cairo PDF surface creation and page size management.
// ABOUTME: Wraps cairo_pdf_surface_create and cairo_pdf_surface_set_size.

//go:build !nopdf

package surface

// #cgo pkg-config: cairo-pdf
// #include <cairo-pdf.h>
// #include <stdlib.h>
import "C"
import "unsafe"

func pdfSurfaceCreate(filename string, widthPt, heightPt float64) SurfacePtr {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	return SurfacePtr(C.cairo_pdf_surface_create(cFilename, C.double(widthPt), C.double(heightPt)))
}

func pdfSurfaceSetSize(ptr SurfacePtr, widthPt, heightPt float64) {
	C.cairo_pdf_surface_set_size(ptr, C.double(widthPt), C.double(heightPt))
}
