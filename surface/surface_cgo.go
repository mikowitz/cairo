package surface

// #cgo pkg-config: cairo
// #include <cairo.h>
import "C"

type SurfacePtr *C.cairo_surface_t
