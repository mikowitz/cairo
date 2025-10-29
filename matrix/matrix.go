package matrix

import (
	"fmt"
	"unsafe"
)

// NewMatrix returns a matrix with the affine transformation given by
// xx , yx , xy , yy , x0 , y0. The transformation is given by:
//
//	x_new = xx * x + xy * y + x0;
//	y_new = yx * x + yy * y + y0;
func NewMatrix(xx, yx, xy, yy, x0, y0 float64) *Matrix {
	return matrixInit(xx, yx, xy, yy, x0, y0)
}

// NewIdentityMatrix returns a matrix with the identity transformation.
func NewIdentityMatrix() *Matrix {
	return matrixInitIdentity()
}

// NewTranslationMatrix returns a matrix with a transformation
// that translates by tx and ty in the X and Y dimensions,
// respectively.
func NewTranslationMatrix(tx, ty float64) *Matrix {
	return matrixInitTranslate(tx, ty)
}

// NewScalingMatrix returns a matrix with a transformation
// that scales by sx and sy in the X and Y dimensions,
// respectively.
func NewScalingMatrix(sx, sy float64) *Matrix {
	return matrixInitScale(sx, sy)
}

// NewRotationMatrix returns a matrix with a transformation
// that rotates by radians.
func NewRotationMatrix(radians float64) *Matrix {
	return matrixInitRotate(radians)
}

func (m *Matrix) Ptr() unsafe.Pointer {
	m.RLock()
	defer m.RUnlock()

	return unsafe.Pointer(m.ptr) //nolint:gosec
}

// FromPointer wraps a C cairo_matrix_t pointer in a Go [Matrix]. This
// function is primarily used for internal use by other Cairo packages.
// The caller is responsible for ensuring the pointer is valid.
func FromPointer(ptr unsafe.Pointer) *Matrix {
	return matrixFromC((MatrixPtr)(ptr))
}

// Translate applies a translation by tx, ty to the transformation in m.
// The effect of the new transformation is to first translate the
// coordinates by tx and ty, then apply the original transformation
// to the coordinates.
func (m *Matrix) Translate(tx, ty float64) {
	m.withLock(func() {
		matrixTranslate(m, tx, ty)
	})
}

// Rotate applies rotation by radians to the transformation in m.
// The effect of the new transformation is to first rotate the
// coordinates by radians, then apply the original transformation
// to the coordinates.
func (m *Matrix) Rotate(radians float64) {
	m.withLock(func() {
		matrixRotate(m, radians)
	})
}

// Scale applies scaling by sx, sy to the transformation in m.
// The effect of the new transformation is to first scale
// the coordinates by sx and sy, then apply the original
// transformation to the coordinates.
func (m *Matrix) Scale(sx, sy float64) {
	m.withLock(func() {
		matrixScale(m, sx, sy)
	})
}

// TransformPoint transforms the point (x, y) by m.
func (m *Matrix) TransformPoint(x, y float64) (float64, float64) {
	m.RLock()
	defer m.RUnlock()

	return matrixTransformPoint(m, x, y)
}

// TransformDistance transforms the distance vector (dx, dy) by m.
// This is similar to TransformPoint except that the translation
// components of the transformation are ignored. The calculation of
// the returned vector is as follows:
//
//	dx2 = xx * dx + xy * dy;
//	dy2 = yx * dx + yy * dy;
func (m *Matrix) TransformDistance(dx, dy float64) (float64, float64) {
	m.RLock()
	defer m.RUnlock()

	return matrixTransformDistance(m, dx, dy)
}

// Invert changes m to be the inverse of its original value.
// Not all transformation matrices have inverses; if the matrix
// collapses points together (it is degenerate), then it has
// no inverse and this function will fail.
func (m *Matrix) Invert() error {
	m.Lock()
	defer m.Unlock()

	return matrixInvert(m)
}

// Multiply multiplies the affine transformations in m and n together
// and returns a new matrix containing the result. The effect of the
// resulting transformation is to first apply the transformation in m
// to the coordinates and then apply the transformation in n to the
// coordinates.
//
// Example, if m in a translation and n is a rotation, the result will
// translate first, and then rotate around the new origin.
func (m *Matrix) Multiply(n *Matrix) *Matrix {
	m.RLock()
	n.RLock()
	defer m.RUnlock()
	defer n.RUnlock()

	return matrixMultiply(m, n)
}

// String implements the [fmt.Stringer] interface, providing
// a compact human-readable representation for [Matrix].
func (m *Matrix) String() string {
	m.RLock()
	defer m.RUnlock()

	return fmt.Sprintf(
		"Matrix {\n  %.2f %.2f %.2f %.2f %.2f %.2f\n}",
		m.XX, m.YX, m.XY, m.YY, m.X0, m.Y0,
	)
}

// Close releases the C memory associated with this matrix. After calling
// Close, the matrix should not be used for any operations.
//
// While matrices have finalizers that will eventually free resources during
// garbage collection, calling Close explicitly ensures deterministic resource
// cleanup. This is particularly important in long-running applications or
// when creating many matrices.
//
// Close is safe to call multiple times. Subsequent calls after the first
// will have no effect.
//
// Example:
//
//	m := matrix.NewIdentityMatrix()
//	defer m.Close()
//
//	// Use the matrix...
//	m.Translate(10, 20)
func (m *Matrix) Close() error {
	return m.destroy()
}

func (m *Matrix) withLock(f func()) {
	m.Lock()
	defer m.Unlock()

	f()
}
