package context

// Arc adds a circular arc of the given radius to the current path. The arc is
// centered at (xc, yc), begins at angle1 and proceeds in the direction of
// increasing angles to end at angle2.
//
// Angles are measured in radians. An angle of 0.0 is in the direction of the
// positive X axis (in user space). An angle of π/2 radians (90 degrees) is in
// the direction of the positive Y axis (in user space). Angles increase in the
// direction from the positive X axis toward the positive Y axis. With the
// default transformation matrix, angles increase in a clockwise direction.
//
// If angle2 is less than angle1, it will be progressively increased by 2π
// until it is greater than angle1.
//
// If there is a current point, an initial line segment will be added to the
// path to connect the current point to the beginning of the arc. If you want
// to create a new sub-path with just the arc, call NewPath() or NewSubPath()
// before calling Arc.
//
// The arc is circular in user space. To achieve an elliptical arc, you can
// scale the current transformation matrix by different amounts in the X and Y
// directions. For example, to draw an ellipse in the box (x, y, width, height):
//
//	ctx.Save()
//	ctx.Translate(x+width/2.0, y+height/2.0)
//	ctx.Scale(width/2.0, height/2.0)
//	ctx.Arc(0.0, 0.0, 1.0, 0.0, 2*math.Pi)
//	ctx.Restore()
//
// Example:
//
//	// Draw a full circle
//	ctx.Arc(100, 100, 50, 0, 2*math.Pi)
//	ctx.Fill()
//
//	// Draw a quarter circle arc
//	ctx.Arc(200, 100, 30, 0, math.Pi/2)
//	ctx.Stroke()
func (c *Context) Arc(xc, yc, radius, angle1, angle2 float64) {
	c.withLock(func() {
		contextArc(c.ptr, xc, yc, radius, angle1, angle2)
	})
}

// ArcNegative adds a circular arc of the given radius to the current path. The
// arc is centered at (xc, yc), begins at angle1 and proceeds in the direction
// of decreasing angles to end at angle2.
//
// Angles are measured in radians. An angle of 0.0 is in the direction of the
// positive X axis (in user space). An angle of π/2 radians (90 degrees) is in
// the direction of the positive Y axis (in user space). With the default
// transformation matrix, angles increase in a clockwise direction.
//
// This function differs from Arc in that it proceeds in the direction of
// decreasing angles. If angle2 is greater than angle1, it will be progressively
// decreased by 2π until it is less than angle1.
//
// If there is a current point, an initial line segment will be added to the
// path to connect the current point to the beginning of the arc. If you want
// to create a new sub-path with just the arc, call NewPath() or NewSubPath()
// before calling ArcNegative.
//
// Example:
//
//	// Draw an arc in the counter-clockwise direction
//	ctx.ArcNegative(100, 100, 50, math.Pi, 0)
//	ctx.Stroke()
func (c *Context) ArcNegative(xc, yc, radius, angle1, angle2 float64) {
	c.withLock(func() {
		contextArcNegative(c.ptr, xc, yc, radius, angle1, angle2)
	})
}

// CurveTo adds a cubic Bézier spline to the path from the current point to
// position (x3, y3) in user-space coordinates, using (x1, y1) and (x2, y2) as
// the control points.
//
// After this call the current point will be (x3, y3).
//
// If there is no current point before the call to CurveTo, this function will
// behave as if preceded by a call to MoveTo(x1, y1).
//
// The curve's shape is determined by the two control points:
//   - The curve begins at the current point with a tangent in the direction of
//     the first control point (x1, y1)
//   - The curve ends at (x3, y3) with a tangent in the direction from the
//     second control point (x2, y2)
//
// Example:
//
//	// Draw a smooth S-curve
//	ctx.MoveTo(20, 100)
//	ctx.CurveTo(80, 20, 120, 180, 180, 100)
//	ctx.Stroke()
//
//	// Create a closed curved shape
//	ctx.MoveTo(100, 50)
//	ctx.CurveTo(120, 50, 150, 70, 150, 100)
//	ctx.CurveTo(150, 130, 120, 150, 100, 150)
//	ctx.ClosePath()
//	ctx.Fill()
func (c *Context) CurveTo(x1, y1, x2, y2, x3, y3 float64) {
	c.withLock(func() {
		contextCurveTo(c.ptr, x1, y1, x2, y2, x3, y3)
	})
}

// RelCurveTo adds a cubic Bézier spline to the path from the current point,
// using relative coordinates. All offsets are relative to the current point.
//
// The curve's control points are at offsets (dx1, dy1) and (dx2, dy2) from the
// current point, and the curve ends at an offset of (dx3, dy3) from the current
// point.
//
// After this call the current point will be offset by (dx3, dy3) from its
// previous position.
//
// It is an error to call this function with no current point. Doing so will
// cause the context to enter an error state.
//
// This is the relative-coordinate version of CurveTo. See CurveTo for more
// details on cubic Bézier curves.
//
// Example:
//
//	// Draw a curve using relative coordinates
//	ctx.MoveTo(50, 100)
//	// Control points at (70, 60) and (110, 60), end at (130, 100)
//	ctx.RelCurveTo(20, -40, 60, -40, 80, 0)
//	ctx.Stroke()
func (c *Context) RelCurveTo(x1, y1, x2, y2, x3, y3 float64) {
	c.withLock(func() {
		contextRelCurveTo(c.ptr, x1, y1, x2, y2, x3, y3)
	})
}

// Rectangle adds a closed rectangular sub-path to the current path.
//
// The rectangle is positioned at (x, y) in user-space with the specified
// width and height. This is equivalent to:
//
//	ctx.MoveTo(x, y)
//	ctx.LineTo(x+width, y)
//	ctx.LineTo(x+width, y+height)
//	ctx.LineTo(x, y+height)
//	ctx.ClosePath()
//
// After calling Rectangle, the current point will be at (x, y).
//
// Example:
//
//	ctx.Rectangle(20.0, 30.0, 100.0, 50.0)  // Rectangle at (20,30) sized 100x50
//	ctx.SetSourceRGB(1.0, 0.0, 0.0)         // Red
//	ctx.Fill()                               // Fill the rectangle
func (c *Context) Rectangle(x, y, width, height float64) {
	c.withLock(func() {
		contextRectangle(c.ptr, x, y, width, height)
	})
}

// LineTo adds a line segment to the path from the current point to (x, y),
// and sets the current point to (x, y).
//
// If there is no current point before the call to LineTo, this function will
// behave as if preceded by a call to MoveTo(x, y).
//
// After this call the current point will be (x, y). Coordinates are specified
// in user-space.
//
// Example:
//
//	ctx.MoveTo(10.0, 10.0)
//	ctx.LineTo(50.0, 10.0)  // Horizontal line
//	ctx.LineTo(50.0, 50.0)  // Vertical line
func (c *Context) LineTo(x, y float64) {
	c.withLock(func() {
		contextLineTo(c.ptr, x, y)
	})
}

// RelLineTo adds a line segment to the path from the current point to a point
// that is offset from the current point by (dx, dy) in user space.
//
// After this call the current point will be offset by (dx, dy) from its
// previous position.
//
// It is an error to call this function with no current point. Doing so will
// cause the context to enter an error state.
//
// This is the relative-coordinate version of LineTo. See LineTo for the
// absolute-coordinate version.
//
// Example:
//
//	// Draw a square using relative coordinates
//	ctx.MoveTo(100, 100)
//	ctx.RelLineTo(50, 0)   // Right
//	ctx.RelLineTo(0, 50)   // Down
//	ctx.RelLineTo(-50, 0)  // Left
//	ctx.RelLineTo(0, -50)  // Up (back to start)
//	ctx.Stroke()
func (c *Context) RelLineTo(x, y float64) {
	c.withLock(func() {
		contextRelLineTo(c.ptr, x, y)
	})
}

// MoveTo begins a new sub-path by setting the current point to (x, y).
//
// After this call the current point will be (x, y). Coordinates are specified
// in user-space, which is affected by the current transformation matrix (CTM).
//
// If there is no current path when MoveTo is called, this function behaves
// identically to calling NewPath() followed by MoveTo(x, y).
//
// Example:
//
//	ctx.MoveTo(50.0, 75.0)  // Start a path at (50, 75)
//	ctx.LineTo(100.0, 75.0) // Draw line to (100, 75)
func (c *Context) MoveTo(x, y float64) {
	c.withLock(func() {
		contextMoveTo(c.ptr, x, y)
	})
}

// RelMoveTo begins a new sub-path. After this call the current point will be
// offset by (dx, dy) from its previous position.
//
// It is an error to call this function with no current point. Doing so will
// cause the context to enter an error state.
//
// This is the relative-coordinate version of MoveTo. Unlike MoveTo, which sets
// the current point to absolute coordinates, RelMoveTo offsets the current
// point by the given amounts.
//
// Example:
//
//	// Draw disconnected line segments using RelMoveTo
//	ctx.MoveTo(20, 20)
//	ctx.LineTo(40, 40)
//	ctx.RelMoveTo(20, 0)  // Move right by 20, starting new sub-path
//	ctx.LineTo(80, 40)
//	ctx.Stroke()
func (c *Context) RelMoveTo(x, y float64) {
	c.withLock(func() {
		contextRelMoveTo(c.ptr, x, y)
	})
}
