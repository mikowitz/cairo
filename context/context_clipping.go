// ABOUTME: Clipping operations for Cairo drawing contexts.
// ABOUTME: Provides clip, reset, extents, and point-in-clip testing in user coordinates.
package context

// Clip establishes a new clip region by intersecting the current path with the
// existing clip region. After Clip, the current path is cleared.
//
// Drawing operations are restricted to the intersection of all clip regions
// established since the last ResetClip. Each successive Clip call can only
// shrink the effective drawing area — it can never expand a previously set clip.
//
// The clip region is part of the graphics state and is saved and restored by
// Save and Restore. This makes it straightforward to apply a temporary clip
// within a Save/Restore block without permanently altering the clip for
// subsequent drawing.
//
// Note: Unlike ClipPreserve, this function clears the current path after
// clipping. If you need to reuse the path (e.g., to stroke the clip boundary),
// use ClipPreserve instead.
//
// Example:
//
//	// Draw stripes but only inside a 50x50 rectangle
//	ctx.Rectangle(25, 25, 50, 50)
//	ctx.Clip()
//	ctx.SetSourceRGB(1.0, 0.0, 0.0)
//	ctx.Paint()  // Only paints within the clipped rectangle
func (c *Context) Clip() {
	c.withLock(func() {
		contextClip(c.ptr)
	})
}

// ClipPreserve establishes a new clip region by intersecting the current path
// with the existing clip region, while retaining the current path for further
// drawing operations.
//
// This is the preserve variant of Clip. The clipping behavior is identical to
// Clip — drawing is restricted to the intersection of all active clip regions —
// but the current path remains available after the call. This is useful when
// you want to both clip to a shape and then stroke its boundary.
//
// Example:
//
//	// Clip to a circle and stroke its boundary
//	ctx.Arc(cx, cy, radius, 0, 2*math.Pi)
//	ctx.ClipPreserve()
//	ctx.SetSourceRGB(0.8, 0.2, 0.2)
//	ctx.Paint()  // Fills inside the circle (clipped)
//	ctx.SetSourceRGB(0.0, 0.0, 0.8)
//	ctx.SetLineWidth(2.0)
//	ctx.Stroke()  // Strokes the preserved circle path
func (c *Context) ClipPreserve() {
	c.withLock(func() {
		contextClipPreserve(c.ptr)
	})
}

// ResetClip resets the current clip region to cover the full extent of the
// surface, effectively removing any clip regions established by prior Clip or
// ClipPreserve calls.
//
// ResetClip does not affect the current path.
//
// Note: Because the clip is part of the Cairo graphics state, an alternative
// to calling ResetClip is to bracket the clip region inside a Save/Restore
// pair. Restore will automatically undo any clip changes made since the
// corresponding Save.
//
// Example:
//
//	ctx.Rectangle(25, 25, 50, 50)
//	ctx.Clip()
//	ctx.Paint()       // Paints only inside the clipped area
//	ctx.ResetClip()   // Remove the clip
//	ctx.Paint()       // Now paints the full surface
func (c *Context) ResetClip() {
	c.withLock(func() {
		contextResetClip(c.ptr)
	})
}

// ClipExtents returns the bounding box of the current clip region in user-space
// coordinates. The return values are (x1, y1, x2, y2) representing the
// top-left and bottom-right corners of the bounding box.
//
// The returned coordinates are in user space, meaning they are expressed in the
// current user coordinate system and are affected by any active transformation
// matrix. For example, after a Translate(50, 50) the origin of user space moves,
// so the extents shift accordingly relative to the original surface coordinates.
//
// If the context has been closed, ClipExtents returns (0, 0, 0, 0).
//
// Example:
//
//	ctx.Rectangle(20, 30, 80, 60)
//	ctx.Clip()
//	x1, y1, x2, y2 := ctx.ClipExtents()
//	// x1=20, y1=30, x2=100, y2=90 (user coordinates)
func (c *Context) ClipExtents() (float64, float64, float64, float64) {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return 0, 0, 0, 0
	}

	return contextClipExtents(c.ptr)
}

// InClip reports whether the given point is inside the current clip region.
// The x and y coordinates must be given in user-space coordinates, not device
// coordinates.
//
// Because InClip operates in user space, the result depends on the current
// transformation matrix. After a Translate(50, 50), a point at (100, 100) in
// user space corresponds to (150, 150) in device space; pass the user-space
// value (100, 100) to test whether that location is within the clip.
//
// If the context has been closed, InClip returns false.
//
// Example:
//
//	ctx.Rectangle(20, 20, 60, 60)  // clip region: x=[20,80], y=[20,80]
//	ctx.Clip()
//	ctx.InClip(50, 50)  // true  — inside the clip
//	ctx.InClip(90, 90)  // false — outside the clip
func (c *Context) InClip(x, y float64) bool {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return false
	}

	return contextInClip(c.ptr, x, y)
}
