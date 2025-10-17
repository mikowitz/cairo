package matrix

// #cgo pkg-config: cairo
// #include <cairo.h>
// #include <stdlib.h>
import "C"

import (
	"runtime"
	"sync"
	"unsafe"
)

type Matrix struct {
	XX, YX, XY, YY, X0, Y0 float64
	sync.RWMutex
	ptr *C.cairo_matrix_t
}

func (m *Matrix) toC() *C.cairo_matrix_t {
	return m.ptr
}

func matrixFromC(ptr *C.cairo_matrix_t) *Matrix {
	m := &Matrix{
		XX:  float64(ptr.xx),
		YX:  float64(ptr.yx),
		XY:  float64(ptr.xy),
		YY:  float64(ptr.yy),
		X0:  float64(ptr.x0),
		Y0:  float64(ptr.y0),
		ptr: ptr,
	}

	runtime.SetFinalizer(m, (*Matrix).destroy)

	return m
}

func matrixInit(xx, yx, xy, yy, x0, y0 float64) *Matrix {
	m := (*C.cairo_matrix_t)(C.malloc(C.sizeof_cairo_matrix_t))

	C.cairo_matrix_init(
		m,
		C.double(xx), C.double(yx),
		C.double(xy), C.double(yy),
		C.double(x0), C.double(y0),
	)

	return matrixFromC(m)
}

func matrixInitIdentity() *Matrix {
	m := (*C.cairo_matrix_t)(C.malloc(C.sizeof_cairo_matrix_t))
	C.cairo_matrix_init_identity(m)

	return matrixFromC(m)
}

func (m *Matrix) destroy() error {
	m.Lock()
	defer m.Unlock()

	if m.ptr != nil {
		C.free(unsafe.Pointer(m.ptr))
		runtime.SetFinalizer(m, nil)
		m.ptr = nil
	}

	return nil
}
