package context

import (
	"github.com/mikowitz/cairo/matrix"
	"github.com/mikowitz/cairo/status"
)

func (c *Context) IdentityMatrix() {
	c.withLock(func() {
		contextIdentityMatrix(c.ptr)
	})
}

func (c *Context) Translate(tx, ty float64) {
	c.withLock(func() {
		contextTranslate(c.ptr, tx, ty)
	})
}

func (c *Context) Scale(sx, sy float64) {
	c.withLock(func() {
		contextScale(c.ptr, sx, sy)
	})
}

func (c *Context) Rotate(radians float64) {
	c.withLock(func() {
		contextRotate(c.ptr, radians)
	})
}

func (c *Context) Transform(m *matrix.Matrix) {
	c.withLock(func() {
		mPtr := m.Ptr()
		contextTransform(c.ptr, mPtr)
	})
}

func (c *Context) GetMatrix() (*matrix.Matrix, error) {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return nil, status.NullPointer
	}
	return contextGetMatrix(c.ptr), nil
}

func (c *Context) SetMatrix(m *matrix.Matrix) {
	if m == nil {
		return
	}

	c.withLock(func() {
		mPtr := m.Ptr()

		contextSetMatrix(c.ptr, mPtr)
	})
}

func (c *Context) UserToDevice(x, y float64) (float64, float64) {
	c.RLock()
	defer c.RUnlock()

	return contextUserToDevice(c.ptr, x, y)
}

func (c *Context) UserToDeviceDistance(dx, dy float64) (float64, float64) {
	c.RLock()
	defer c.RUnlock()

	return contextUserToDeviceDistance(c.ptr, dx, dy)
}

func (c *Context) DeviceToUser(x, y float64) (float64, float64) {
	c.RLock()
	defer c.RUnlock()

	return contextDeviceToUser(c.ptr, x, y)
}

func (c *Context) DeviceToUserDistance(dx, dy float64) (float64, float64) {
	c.RLock()
	defer c.RUnlock()

	return contextDeviceToUserDistance(c.ptr, dx, dy)
}
