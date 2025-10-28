// Package matrix is used throughout [cairo] to convert between different coordinate spaces.
// A [Matrix] holds an affine transformation, such as a scale, rotation, shear, or a combination of these.
// The transformation of a point (x,y) is given by:
//
//	x_new = xx * x + xy * y + x0;
//	y_new = yx * x + yy * y + y0;
//
// The current transformation matrix of a [Context], represented as a [Matrix],
// defines the transformation from user-space coordinates to device-space
// coordinates. See [context.GetMatrix()] and [context.SetMatrix(m)].
//
// # Affine Transformation Mathematics
//
// An affine transformation preserves points, straight lines, and planes. It maps
// parallel lines to parallel lines and preserves ratios of distances along lines.
// The matrix representation uses six floating-point values (xx, yx, xy, yy, x0, y0)
// that can be thought of as a 3x3 matrix in homogeneous coordinates:
//
//	[ xx  xy  x0 ]
//	[ yx  yy  y0 ]
//	[  0   0   1 ]
//
// Common transformations and their matrix representations:
//
// Identity (no transformation):
//	xx=1, yx=0, xy=0, yy=1, x0=0, y0=0
//
// Translation by (tx, ty):
//	xx=1, yx=0, xy=0, yy=1, x0=tx, y0=ty
//
// Scaling by (sx, sy):
//	xx=sx, yx=0, xy=0, yy=sy, x0=0, y0=0
//
// Rotation by angle θ (radians, counter-clockwise):
//	xx=cos(θ), yx=sin(θ), xy=-sin(θ), yy=cos(θ), x0=0, y0=0
//
// # Combining Transformations
//
// Transformations are combined through matrix multiplication. The order matters:
// transforming by matrix A then matrix B is equivalent to multiplying B × A
// (note the reverse order). This is because transformations are applied right-to-left.
//
// Example - translate then scale:
//	m := matrix.NewTranslationMatrix(10, 20)  // Move right 10, down 20
//	m.Scale(2, 2)                             // Then double the size
//	// This scales around the new origin at (10, 20)
//
// Example - scale then translate:
//	m := matrix.NewScalingMatrix(2, 2)        // Double the size
//	m.Translate(10, 20)                       // Then move right 10, down 20
//	// This moves by (20, 40) in original coordinates
//
// # User Space vs Device Space
//
// In Cairo, user-space coordinates are the coordinates you use when drawing
// (e.g., Rectangle(10, 20, 30, 40)). The current transformation matrix (CTM)
// converts these to device-space coordinates (actual pixels on the output device).
//
// The default CTM is usually the identity matrix, making user-space units equal
// to device-space units (typically pixels). Modifying the CTM lets you work in
// a more convenient coordinate system:
//
//	// Make (0,0) the center of a 400x400 surface
//	ctx.Translate(200, 200)
//	// Now Rectangle(-50, -50, 100, 100) draws centered
//
// # TransformPoint vs TransformDistance
//
// TransformPoint applies the full transformation including translation:
//	x_new = xx*x + xy*y + x0
//	y_new = yx*x + yy*y + y0
//
// TransformDistance ignores translation, useful for vectors and sizes:
//	dx_new = xx*dx + xy*dy
//	dy_new = yx*dx + yy*dy
//
// Example:
//	m := matrix.NewTranslationMatrix(100, 100)
//	px, py := m.TransformPoint(10, 10)      // Returns (110, 110)
//	dx, dy := m.TransformDistance(10, 10)   // Returns (10, 10) - translation ignored
//
// # Matrix Inversion
//
// Some operations require the inverse of a matrix to transform coordinates
// in the reverse direction. Not all matrices are invertible:
//
//   - A matrix with determinant 0 is singular (non-invertible)
//   - Singular matrices "collapse" dimensions (e.g., scale by 0)
//   - Calling Invert() on a singular matrix returns an error
//
// The determinant is calculated as: xx*yy - yx*xy
//
// Example:
//	m := matrix.NewScalingMatrix(2, 3)
//	err := m.Invert()
//	// m is now scaling by (0.5, 0.333...)
//	// Applying original then inverse returns to starting point
package matrix
