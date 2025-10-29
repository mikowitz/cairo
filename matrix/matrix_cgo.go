package matrix

// #cgo pkg-config: cairo
// #include <cairo.h>
// #include <stdlib.h>
import "C"

import (
	"runtime"
	"sync"
	"unsafe"

	"github.com/mikowitz/cairo/status"
)

type MatrixPtr *C.cairo_matrix_t

type Matrix struct {
	XX, YX, XY, YY, X0, Y0 float64
	sync.RWMutex
	ptr *C.cairo_matrix_t
}

func (m *Matrix) toC() *C.cairo_matrix_t {
	return m.ptr
}

func matrixFromC(ptr *C.cairo_matrix_t) *Matrix {
	m := &Matrix{ptr: ptr}
	m.updateFromC()

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

func matrixInitTranslate(tx, ty float64) *Matrix {
	m := (*C.cairo_matrix_t)(C.malloc(C.sizeof_cairo_matrix_t))
	C.cairo_matrix_init_translate(m, C.double(tx), C.double(ty))
	return matrixFromC(m)
}

func matrixInitScale(sx, sy float64) *Matrix {
	m := (*C.cairo_matrix_t)(C.malloc(C.sizeof_cairo_matrix_t))
	C.cairo_matrix_init_scale(m, C.double(sx), C.double(sy))
	return matrixFromC(m)
}

func matrixInitRotate(radians float64) *Matrix {
	m := (*C.cairo_matrix_t)(C.malloc(C.sizeof_cairo_matrix_t))
	C.cairo_matrix_init_rotate(m, C.double(radians))
	return matrixFromC(m)
}

func matrixTranslate(m *Matrix, tx, ty float64) {
	C.cairo_matrix_translate(m.ptr, C.double(tx), C.double(ty))
	m.updateFromC()
}

func matrixRotate(m *Matrix, radians float64) {
	C.cairo_matrix_rotate(m.ptr, C.double(radians))
	m.updateFromC()
}

func matrixScale(m *Matrix, sx, sy float64) {
	C.cairo_matrix_scale(m.ptr, C.double(sx), C.double(sy))
	m.updateFromC()
}

func matrixTransformPoint(m *Matrix, x, y float64) (float64, float64) {
	tx := C.double(x)
	ty := C.double(y)

	C.cairo_matrix_transform_point(m.ptr, &tx, &ty)
	return float64(tx), float64(ty)
}

func matrixTransformDistance(m *Matrix, dx, dy float64) (float64, float64) {
	tx := C.double(dx)
	ty := C.double(dy)

	C.cairo_matrix_transform_distance(m.ptr, &tx, &ty)
	return float64(tx), float64(ty)
}

func matrixInvert(m *Matrix) error {
	st := status.Status(C.cairo_matrix_invert(m.ptr))
	m.updateFromC()
	return st.ToError()
}

func matrixMultiply(m, n *Matrix) *Matrix {
	rPtr := (*C.cairo_matrix_t)(C.malloc(C.sizeof_cairo_matrix_t))

	C.cairo_matrix_multiply(rPtr, m.ptr, n.ptr)
	return matrixFromC(rPtr)
}

func (m *Matrix) updateFromC() {
	m.XX = float64(m.ptr.xx)
	m.YX = float64(m.ptr.yx)
	m.XY = float64(m.ptr.xy)
	m.YY = float64(m.ptr.yy)
	m.X0 = float64(m.ptr.x0)
	m.Y0 = float64(m.ptr.y0)
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
