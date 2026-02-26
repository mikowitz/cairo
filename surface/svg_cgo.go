// ABOUTME: CGO bindings for Cairo SVG surface creation and document unit configuration.
// ABOUTME: Wraps cairo_svg_surface_create and cairo_svg_surface_set_document_unit.

//go:build !nosvg

package surface

// #cgo pkg-config: cairo-svg
// #include <cairo-svg.h>
// #include <stdlib.h>
import "C"
import "unsafe"

func svgSurfaceCreate(filename string, widthPt, heightPt float64) SurfacePtr {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	return SurfacePtr(C.cairo_svg_surface_create(cFilename, C.double(widthPt), C.double(heightPt)))
}

func svgSurfaceSetDocumentUnit(ptr SurfacePtr, unit SVGUnit) {
	C.cairo_svg_surface_set_document_unit(ptr, C.cairo_svg_unit_t(unit))
}
