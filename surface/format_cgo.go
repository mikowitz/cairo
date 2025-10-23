package surface

// #cgo pkg-config: cairo
// #include <cairo.h>
import "C"

func formatStrideForWidth(format Format, width int) int {
	return int(C.cairo_format_stride_for_width(C.cairo_format_t(format), C.int(width)))
}
