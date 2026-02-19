// ABOUTME: Clipping operations for Cairo drawing contexts.
// ABOUTME: Provides clip, reset, extents, and point-in-clip testing in user coordinates.
package context

func (c *Context) Clip() {
	c.withLock(func() {
		contextClip(c.ptr)
	})
}

func (c *Context) ClipPreserve() {
	c.withLock(func() {
		contextClipPreserve(c.ptr)
	})
}

func (c *Context) ResetClip() {
	c.withLock(func() {
		contextResetClip(c.ptr)
	})
}

func (c *Context) ClipExtents() (float64, float64, float64, float64) {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return 0, 0, 0, 0
	}

	return contextClipExtents(c.ptr)
}

func (c *Context) InClip(x, y float64) bool {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return false
	}

	return contextInClip(c.ptr, x, y)
}
