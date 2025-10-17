// Package matrix is used throughout cairo to convert between different coordinate spaces.
// A [Matrix] holds an affine transformation, such as a scale, rotation, shear, or a combination of these.
// The transformation of a point (x,y) is given by:
//
//	x_new = xx * x + xy * y + x0;
//	y_new = yx * x + yy * y + y0;
//
// The current transformation matrix of a [Context], represented as a [Matrix],
// defines the transformation from user-space coordinates to device-space
// coordinates. See [context.GetMatrix()] and [context.SetMatrix(m)].
package matrix
