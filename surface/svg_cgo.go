// ABOUTME: CGO bindings for Cairo SVG surface creation, document unit, and version configuration.
// ABOUTME: Wraps cairo_svg_surface_create, cairo_svg_surface_set_document_unit, and version APIs.

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

func svgSurfaceRestrictToVersion(ptr SurfacePtr, version SVGVersion) {
	C.cairo_svg_surface_restrict_to_version(ptr, C.cairo_svg_version_t(version))
}

func svgGetVersions() []SVGVersion {
	var versions *C.cairo_svg_version_t
	var numVersions C.int
	C.cairo_svg_get_versions(&versions, &numVersions)
	if numVersions == 0 || versions == nil {
		return nil
	}
	cSlice := unsafe.Slice(versions, int(numVersions))
	result := make([]SVGVersion, int(numVersions))
	for i, v := range cSlice {
		result[i] = SVGVersion(v)
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
