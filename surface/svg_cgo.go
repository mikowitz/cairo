// ABOUTME: CGO bindings for Cairo SVG surface creation, document unit, and version configuration.
// ABOUTME: Wraps cairo_svg_surface_create, cairo_svg_surface_set_document_unit, and version APIs.

//go:build !nosvg

package surface

// #cgo pkg-config: cairo-svg
// #include <cairo-svg.h>
// #include <stdlib.h>
//
// // _svgGetVersionsList fills buf with supported SVG version identifiers and
// // returns the count. buf must have room for at least 16 entries.
// static int _svgGetVersionsList(cairo_svg_version_t *buf) {
//     cairo_svg_version_t const *v;
//     int n = 0;
//     cairo_svg_get_versions(&v, &n);
//     for (int i = 0; i < n; i++) { buf[i] = v[i]; }
//     return n;
// }
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

func svgSurfaceGetDocumentUnit(ptr SurfacePtr) SVGUnit {
	return SVGUnit(C.cairo_svg_surface_get_document_unit(ptr))
}

func svgSurfaceRestrictToVersion(ptr SurfacePtr, version SVGVersion) {
	C.cairo_svg_surface_restrict_to_version(ptr, C.cairo_svg_version_t(version))
}

func svgGetVersions() []SVGVersion {
	var buf [16]C.cairo_svg_version_t
	count := int(C._svgGetVersionsList(&buf[0]))
	result := make([]SVGVersion, count)
	for i := range result {
		result[i] = SVGVersion(buf[i])
	}
	return result
}

func svgVersionToString(version SVGVersion) string {
	cStr := C.cairo_svg_version_to_string(C.cairo_svg_version_t(version))
	if cStr == nil {
		return ""
	}
	return C.GoString(cStr)
}
