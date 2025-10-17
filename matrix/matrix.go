package matrix

import "fmt"

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

// String implements the [fmt.Stringer] interface for [Matrix],
// providing a compact human-readable representation.
func (m *Matrix) String() string {
	m.RLock()
	defer m.RUnlock()

	return fmt.Sprintf(
		"Matrix {\n  %.2f %.2f %.2f %.2f %.2f %.2f\n}",
		m.XX, m.YX, m.XY, m.YY, m.X0, m.Y0,
	)
}
