package context

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestContextCoordinateConversion verifies user/device space conversions.
func TestContextCoordinateConversion(t *testing.T) {
	ctx := newTestContext(t, 200, 200)

	t.Run("user_to_device_identity", func(t *testing.T) {
		ctx.IdentityMatrix()

		x, y := ctx.UserToDevice(10.0, 20.0)
		assert.InDelta(t, 10.0, x, 0.001, "With identity matrix, coordinates should not change")
		assert.InDelta(t, 20.0, y, 0.001)
	})

	t.Run("user_to_device_after_translation", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Translate(5.0, 10.0)

		x, y := ctx.UserToDevice(10.0, 20.0)
		assert.InDelta(t, 15.0, x, 0.001, "Translation should offset coordinates")
		assert.InDelta(t, 30.0, y, 0.001)
	})

	t.Run("user_to_device_after_scaling", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Scale(2.0, 3.0)

		x, y := ctx.UserToDevice(10.0, 20.0)
		assert.InDelta(t, 20.0, x, 0.001, "Scaling should multiply coordinates")
		assert.InDelta(t, 60.0, y, 0.001)
	})

	t.Run("user_to_device_distance_identity", func(t *testing.T) {
		ctx.IdentityMatrix()

		dx, dy := ctx.UserToDeviceDistance(10.0, 20.0)
		assert.InDelta(t, 10.0, dx, 0.001)
		assert.InDelta(t, 20.0, dy, 0.001)
	})

	t.Run("user_to_device_distance_ignores_translation", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Translate(100.0, 200.0)

		// Distance should not be affected by translation
		dx, dy := ctx.UserToDeviceDistance(10.0, 20.0)
		assert.InDelta(t, 10.0, dx, 0.001, "Translation should not affect distance")
		assert.InDelta(t, 20.0, dy, 0.001)
	})

	t.Run("user_to_device_distance_with_scaling", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Scale(2.0, 3.0)

		// Distance should be scaled
		dx, dy := ctx.UserToDeviceDistance(10.0, 20.0)
		assert.InDelta(t, 20.0, dx, 0.001, "Distance should be scaled")
		assert.InDelta(t, 60.0, dy, 0.001)
	})

	t.Run("device_to_user_identity", func(t *testing.T) {
		ctx.IdentityMatrix()

		x, y := ctx.DeviceToUser(10.0, 20.0)
		assert.InDelta(t, 10.0, x, 0.001)
		assert.InDelta(t, 20.0, y, 0.001)
	})

	t.Run("device_to_user_inverse_of_user_to_device", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Translate(5.0, 10.0)
		ctx.Scale(2.0, 3.0)

		// User -> Device -> User should return original
		origX, origY := 15.0, 25.0
		dx, dy := ctx.UserToDevice(origX, origY)
		ux, uy := ctx.DeviceToUser(dx, dy)

		assert.InDelta(t, origX, ux, 0.001, "Round-trip should return original")
		assert.InDelta(t, origY, uy, 0.001)
	})

	t.Run("device_to_user_distance_inverse", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Scale(2.0, 3.0)

		// User distance -> Device distance -> User distance should return original
		origDx, origDy := 10.0, 20.0
		ddx, ddy := ctx.UserToDeviceDistance(origDx, origDy)
		udx, udy := ctx.DeviceToUserDistance(ddx, ddy)

		assert.InDelta(t, origDx, udx, 0.001)
		assert.InDelta(t, origDy, udy, 0.001)
	})

	t.Run("coordinate_conversion_with_rotation", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Rotate(math.Pi / 2) // 90 degrees

		// Point (10, 0) rotated 90° becomes (0, 10)
		x, y := ctx.UserToDevice(10.0, 0.0)
		assert.InDelta(t, 0.0, x, 0.001)
		assert.InDelta(t, 10.0, y, 0.001)

		// Verify inverse
		ux, uy := ctx.DeviceToUser(0.0, 10.0)
		assert.InDelta(t, 10.0, ux, 0.001)
		assert.InDelta(t, 0.0, uy, 0.001)
	})

	t.Run("zero_coordinates", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Translate(5.0, 10.0)

		x, y := ctx.UserToDevice(0.0, 0.0)
		assert.InDelta(t, 5.0, x, 0.001)
		assert.InDelta(t, 10.0, y, 0.001)
	})

	t.Run("negative_coordinates", func(t *testing.T) {
		ctx.IdentityMatrix()
		ctx.Scale(2.0, 2.0)

		x, y := ctx.UserToDevice(-5.0, -10.0)
		assert.InDelta(t, -10.0, x, 0.001)
		assert.InDelta(t, -20.0, y, 0.001)
	})
}
