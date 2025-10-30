package context

import (
	"github.com/mikowitz/cairo/matrix"
	"github.com/mikowitz/cairo/status"
)

// IdentityMatrix resets the current transformation matrix (CTM) to the identity matrix.
//
// After calling this method, user space and device space are aligned with a 1:1
// correspondence. This removes all scaling, rotation, and translation transformations.
//
// This is equivalent to calling SetMatrix with a new identity matrix, and is useful
// when you want to reset transformations to a known state.
//
// Example:
//
//	ctx.Translate(10, 20)
//	ctx.Scale(2, 2)
//	ctx.IdentityMatrix()  // Removes all transformations
func (c *Context) IdentityMatrix() {
	c.withLock(func() {
		contextIdentityMatrix(c.ptr)
	})
}

// Translate modifies the current transformation matrix by moving the user-space origin.
//
// The translation is applied to the current transformation matrix, so it is cumulative
// with any existing transformations. After calling Translate(tx, ty), the point (0, 0)
// in user space will correspond to the point (tx, ty) in the previous user space.
//
// Parameters:
//   - tx: Amount to translate in the X direction
//   - ty: Amount to translate in the Y direction
//
// Example:
//
//	ctx.Translate(10, 20)      // Move origin to (10, 20)
//	ctx.Translate(5, 5)        // Move origin to (15, 25) - cumulative
//	ctx.Rectangle(0, 0, 50, 50) // Rectangle actually drawn at (15, 25)
func (c *Context) Translate(tx, ty float64) {
	c.withLock(func() {
		contextTranslate(c.ptr, tx, ty)
	})
}

// Scale modifies the current transformation matrix by scaling the user-space axes.
//
// The scaling is applied to the current transformation matrix, so it is cumulative
// with any existing transformations. After calling Scale(sx, sy), user-space units
// are multiplied by sx and sy in the X and Y directions respectively.
//
// Parameters:
//   - sx: Scale factor for the X axis (2.0 doubles width, 0.5 halves it)
//   - sy: Scale factor for the Y axis (2.0 doubles height, 0.5 halves it)
//
// Note: Negative scale factors can be used to create mirror reflections.
//
// Example:
//
//	ctx.Scale(2.0, 2.0)        // Double all coordinates
//	ctx.Rectangle(0, 0, 50, 50) // Actually draws a 100x100 rectangle
//
//	ctx.Scale(1.0, -1.0)       // Flip Y axis (useful for inverting coordinate system)
func (c *Context) Scale(sx, sy float64) {
	c.withLock(func() {
		contextScale(c.ptr, sx, sy)
	})
}

// Rotate modifies the current transformation matrix by rotating the user-space axes.
//
// The rotation is applied to the current transformation matrix, so it is cumulative
// with any existing transformations. The rotation angle is specified in radians.
//
// For positive angles, the rotation direction is from the positive X axis toward
// the positive Y axis. In the default Cairo coordinate system (Y axis pointing down),
// positive rotation appears clockwise.
//
// Parameters:
//   - radians: Rotation angle in radians (use math.Pi/180 to convert from degrees)
//
// Example:
//
//	ctx.Rotate(math.Pi / 4)     // Rotate 45 degrees
//	ctx.Rectangle(0, 0, 50, 50) // Rectangle is now rotated
//
//	// To rotate around a specific point:
//	ctx.Translate(centerX, centerY)
//	ctx.Rotate(angle)
//	ctx.Translate(-centerX, -centerY)
func (c *Context) Rotate(radians float64) {
	c.withLock(func() {
		contextRotate(c.ptr, radians)
	})
}

// Transform modifies the current transformation matrix by applying an additional matrix.
//
// This method multiplies the provided matrix with the current transformation matrix,
// preserving all existing transformations. The new CTM is computed as: new_CTM = m × old_CTM
//
// This is useful for applying custom transformations that combine multiple operations
// or for transformations that can't be expressed as simple translate/scale/rotate.
//
// Parameters:
//   - m: The transformation matrix to apply
//
// Example:
//
//	// Apply a custom shear transformation
//	shear := matrix.New(1, 0.5, 0, 1, 0, 0)  // Shear in X direction
//	ctx.Transform(shear)
//
//	// Combine with existing transformations
//	ctx.Translate(10, 10)
//	ctx.Transform(customMatrix)  // Applied after translation
func (c *Context) Transform(m *matrix.Matrix) {
	if m == nil {
		return
	}

	c.withLock(func() {
		mPtr := m.Ptr()
		contextTransform(c.ptr, mPtr)
	})
}

// GetMatrix retrieves the current transformation matrix (CTM).
//
// This method returns a copy of the current transformation matrix, which represents
// the cumulative effect of all transformations applied to this context. The matrix
// can be modified without affecting the context's transformation state.
//
// Returns the current transformation matrix and nil error on success.
// Returns nil and status.NullPointer if the context has been closed.
//
// The returned matrix can be used to:
//   - Save the current transformation state for later restoration with SetMatrix
//   - Inspect the current transformation parameters
//   - Create a modified version for use with Transform
//
// Example:
//
//	// Save current transformation
//	savedMatrix, err := ctx.GetMatrix()
//	if err != nil {
//	    return err
//	}
//
//	// Apply some transformations
//	ctx.Translate(10, 20)
//	ctx.Rotate(math.Pi / 4)
//
//	// Restore saved transformation
//	ctx.SetMatrix(savedMatrix)
func (c *Context) GetMatrix() (*matrix.Matrix, error) {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return nil, status.NullPointer
	}
	return contextGetMatrix(c.ptr), nil
}

// SetMatrix replaces the current transformation matrix (CTM) with the provided matrix.
//
// Unlike Transform, which multiplies matrices, SetMatrix completely replaces the
// existing transformation matrix. This is useful for restoring a previously saved
// transformation state or setting an exact transformation.
//
// If m is nil, this method does nothing and returns immediately.
//
// Parameters:
//   - m: The new transformation matrix, or nil to do nothing
//
// Example:
//
//	// Set a specific transformation directly
//	m := matrix.New(2, 0, 0, 2, 10, 20)  // Scale 2x and translate (10, 20)
//	ctx.SetMatrix(m)
//
//	// Common pattern: save/restore transformations
//	saved, _ := ctx.GetMatrix()
//	ctx.Translate(100, 100)
//	// ... draw something ...
//	ctx.SetMatrix(saved)  // Restore original transformation
func (c *Context) SetMatrix(m *matrix.Matrix) {
	if m == nil {
		return
	}

	c.withLock(func() {
		mPtr := m.Ptr()

		contextSetMatrix(c.ptr, mPtr)
	})
}

// UserToDevice transforms a coordinate from user space to device space.
//
// This method applies the current transformation matrix to convert coordinates
// from user space (the coordinate system used by drawing commands) to device space
// (the coordinate system of the target surface, typically pixels).
//
// Both position and the effects of the transformation matrix (including translation,
// scaling, and rotation) are applied to the input coordinates.
//
// Parameters:
//   - x: X coordinate in user space
//   - y: Y coordinate in user space
//
// Returns the corresponding X and Y coordinates in device space.
//
// Example:
//
//	ctx.Translate(10, 20)
//	ctx.Scale(2, 2)
//	dx, dy := ctx.UserToDevice(5, 5)
//	// dx = 20, dy = 30  (scaled then translated: 5*2 + 10, 5*2 + 20)
func (c *Context) UserToDevice(x, y float64) (float64, float64) {
	c.RLock()
	defer c.RUnlock()

	return contextUserToDevice(c.ptr, x, y)
}

// UserToDeviceDistance transforms a distance vector from user space to device space.
//
// This method is similar to UserToDevice, but it transforms distance vectors rather
// than position coordinates. The key difference is that translation components of
// the transformation matrix are ignored, as distance vectors represent relative
// offsets, not absolute positions.
//
// This is useful for converting dimensions (width, height) or direction vectors
// while accounting for scaling and rotation but not translation.
//
// Parameters:
//   - dx: X component of distance vector in user space
//   - dy: Y component of distance vector in user space
//
// Returns the corresponding X and Y components in device space.
//
// Example:
//
//	ctx.Translate(100, 100)  // Translation is ignored
//	ctx.Scale(2, 3)
//	dx, dy := ctx.UserToDeviceDistance(10, 20)
//	// dx = 20, dy = 60  (only scaling applied: 10*2, 20*3)
func (c *Context) UserToDeviceDistance(dx, dy float64) (float64, float64) {
	c.RLock()
	defer c.RUnlock()

	return contextUserToDeviceDistance(c.ptr, dx, dy)
}

// DeviceToUser transforms a coordinate from device space to user space.
//
// This method applies the inverse of the current transformation matrix to convert
// coordinates from device space (typically pixels) to user space (the coordinate
// system used by drawing commands).
//
// This is the inverse operation of UserToDevice and is useful for converting
// input coordinates (like mouse positions) from screen/pixel coordinates back
// to the coordinate system your drawing code uses.
//
// Parameters:
//   - x: X coordinate in device space
//   - y: Y coordinate in device space
//
// Returns the corresponding X and Y coordinates in user space.
//
// Example:
//
//	ctx.Translate(10, 20)
//	ctx.Scale(2, 2)
//	ux, uy := ctx.DeviceToUser(30, 50)
//	// ux = 10, uy = 15  (inverse: (30-10)/2, (50-20)/2)
//
//	// Common use: converting mouse clicks to drawing coordinates
//	userX, userY := ctx.DeviceToUser(mouseX, mouseY)
func (c *Context) DeviceToUser(x, y float64) (float64, float64) {
	c.RLock()
	defer c.RUnlock()

	return contextDeviceToUser(c.ptr, x, y)
}

// DeviceToUserDistance transforms a distance vector from device space to user space.
//
// This method is similar to DeviceToUser, but it transforms distance vectors rather
// than position coordinates. The key difference is that translation components of
// the transformation matrix are ignored, as distance vectors represent relative
// offsets, not absolute positions.
//
// This is the inverse operation of UserToDeviceDistance and is useful for converting
// dimensions or offsets from device space back to user space.
//
// Parameters:
//   - dx: X component of distance vector in device space
//   - dy: Y component of distance vector in device space
//
// Returns the corresponding X and Y components in user space.
//
// Example:
//
//	ctx.Translate(100, 100)  // Translation is ignored
//	ctx.Scale(2, 3)
//	ux, uy := ctx.DeviceToUserDistance(20, 60)
//	// ux = 10, uy = 20  (inverse scaling: 20/2, 60/3)
func (c *Context) DeviceToUserDistance(dx, dy float64) (float64, float64) {
	c.RLock()
	defer c.RUnlock()

	return contextDeviceToUserDistance(c.ptr, dx, dy)
}
