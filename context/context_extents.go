// ABOUTME: Fill, stroke, and path extents methods plus point-in-fill/stroke queries.
// ABOUTME: These operations query geometry without consuming or modifying the current path.

package context

// FillExtents returns the bounding box that would be affected by calling [Context.Fill]
// with the current path and fill rule. The return values (x1, y1, x2, y2) are the
// top-left and bottom-right corners of the bounding box in user-space coordinates.
//
// FillExtents does not consume the current path.
//
// If the context has been closed, FillExtents returns (0, 0, 0, 0).
func (c *Context) FillExtents() (x1, y1, x2, y2 float64) {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return 0, 0, 0, 0
	}
	return contextFillExtents(c.ptr)
}

// StrokeExtents returns the bounding box that would be affected by calling [Context.Stroke]
// with the current path and stroke parameters. The return values (x1, y1, x2, y2) are the
// top-left and bottom-right corners of the bounding box in user-space coordinates.
//
// The bounding box includes the effect of line width, line cap, and line join styles.
// StrokeExtents does not consume the current path.
//
// If the context has been closed, StrokeExtents returns (0, 0, 0, 0).
func (c *Context) StrokeExtents() (x1, y1, x2, y2 float64) {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return 0, 0, 0, 0
	}
	return contextStrokeExtents(c.ptr)
}

// PathExtents returns the bounding box of the current path in user-space coordinates.
// The return values (x1, y1, x2, y2) are the top-left and bottom-right corners.
//
// Unlike [Context.FillExtents] and [Context.StrokeExtents], PathExtents is purely
// geometric and does not account for line width or fill rule.
// PathExtents does not consume the current path.
//
// If the context has been closed or there is no path, PathExtents returns (0, 0, 0, 0).
func (c *Context) PathExtents() (x1, y1, x2, y2 float64) {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return 0, 0, 0, 0
	}
	return contextPathExtents(c.ptr)
}

// InFill reports whether the given point is inside the area that would be filled
// by calling [Context.Fill] with the current path and fill rule.
//
// The x and y coordinates are in user-space. InFill does not consume the path.
//
// If the context has been closed, InFill returns false.
func (c *Context) InFill(x, y float64) bool {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return false
	}
	return contextInFill(c.ptr, x, y)
}

// InStroke reports whether the given point is inside the area that would be affected
// by calling [Context.Stroke] with the current path and stroke parameters.
//
// The x and y coordinates are in user-space. InStroke does not consume the path.
//
// If the context has been closed, InStroke returns false.
func (c *Context) InStroke(x, y float64) bool {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return false
	}
	return contextInStroke(c.ptr, x, y)
}
