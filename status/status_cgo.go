package status

// #cgo pkg-config: cairo
// #include <cairo.h>
import "C"

func statusFromC(cStatus C.cairo_status_t) Status {
	return Status(cStatus)
}

func (s Status) toC() C.cairo_status_t {
	return C.cairo_status_t(s)
}

func (s Status) toString() string {
	return C.GoString(C.cairo_status_to_string(s.toC()))
}
