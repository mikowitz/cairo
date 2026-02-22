// ABOUTME: Defines the Operator type and compositing operator constants for Cairo drawing contexts.
// ABOUTME: Provides SetOperator and GetOperator methods for controlling how drawing operations blend.

package context

// Operator controls how drawing operations combine with existing surface content.
//
// Cairo supports two classes of operators:
//
// # Porter-Duff Operators
//
// The Porter-Duff operators (Clear through Saturate) implement the standard
// alpha compositing model. They control how source and destination pixels are
// combined based on their alpha values.
//
// OperatorOver is the default and most common operator. It draws the source
// on top of the destination, respecting alpha transparency.
//
// # Blend Mode Operators
//
// The blend mode operators (Multiply through HslLuminosity) implement
// photoshop-style blending. They compute a color based on both the source
// and destination colors, then composite using OperatorOver semantics.
//
//go:generate stringer -type=Operator
type Operator int

// The iota values below must match Cairo's cairo_operator_t C enum exactly.
// Cairo has maintained this ordering since version 1.10 and documents it as
// stable. The CGO layer casts Operator directly to cairo_operator_t, so any
// divergence would silently produce incorrect compositing. If Cairo ever adds
// or reorders operators, this block and the stringer output must be updated.
const (
	// OperatorClear clears the destination layer, making it fully transparent.
	OperatorClear Operator = iota

	// OperatorSource replaces the destination with the source, ignoring the destination.
	OperatorSource

	// OperatorOver draws the source on top of the destination, respecting alpha.
	// This is the default operator.
	OperatorOver

	// OperatorIn draws the source only where the destination is opaque.
	OperatorIn

	// OperatorOut draws the source only where the destination is transparent.
	OperatorOut

	// OperatorAtop draws the source on top of the destination, but only where
	// the destination exists (i.e., has non-zero alpha).
	OperatorAtop

	// OperatorDest leaves the destination unchanged, ignoring the source.
	OperatorDest

	// OperatorDestOver draws the destination on top of the source.
	OperatorDestOver

	// OperatorDestIn leaves the destination visible only where the source is opaque.
	OperatorDestIn

	// OperatorDestOut leaves the destination visible only where the source is transparent.
	OperatorDestOut

	// OperatorDestAtop leaves the destination on top of the source, only where the source exists.
	OperatorDestAtop

	// OperatorXor shows source and destination where they do not overlap.
	OperatorXor

	// OperatorAdd adds source and destination pixel values.
	OperatorAdd

	// OperatorSaturate saturates source and destination; useful for alpha accumulation.
	OperatorSaturate

	// OperatorMultiply multiplies source and destination colors, darkening the result.
	OperatorMultiply

	// OperatorScreen inverts, multiplies, then inverts again, lightening the result.
	OperatorScreen

	// OperatorOverlay multiplies or screens depending on the destination color.
	OperatorOverlay

	// OperatorDarken uses the darker of the source and destination colors.
	OperatorDarken

	// OperatorLighten uses the lighter of the source and destination colors.
	OperatorLighten

	// OperatorColorDodge brightens the destination to reflect the source.
	OperatorColorDodge

	// OperatorColorBurn darkens the destination to reflect the source.
	OperatorColorBurn

	// OperatorHardLight applies Multiply or Screen depending on the source color.
	OperatorHardLight

	// OperatorSoftLight applies a softer version of HardLight.
	OperatorSoftLight

	// OperatorDifference subtracts source from destination or vice versa.
	OperatorDifference

	// OperatorExclusion produces an effect similar to Difference but lower in contrast.
	OperatorExclusion

	// OperatorHslHue uses the hue of the source with the saturation and luminosity of the destination.
	OperatorHslHue

	// OperatorHslSaturation uses the saturation of the source with the hue and luminosity of the destination.
	OperatorHslSaturation

	// OperatorHslColor uses the hue and saturation of the source with the luminosity of the destination.
	OperatorHslColor

	// OperatorHslLuminosity uses the luminosity of the source with the hue and saturation of the destination.
	OperatorHslLuminosity
)

// SetOperator sets the compositing operator to be used for all drawing operations.
//
// The operator controls how drawing operations combine with existing content on
// the surface. The default operator is [OperatorOver], which draws the source on
// top of the destination respecting alpha transparency.
//
// Porter-Duff operators control alpha compositing:
//   - [OperatorOver]: Source over destination (default)
//   - [OperatorSource]: Replace destination with source
//   - [OperatorClear]: Clear destination to transparent
//   - [OperatorIn], [OperatorOut], [OperatorAtop]: Mask variants
//
// Blend mode operators combine colors:
//   - [OperatorMultiply]: Darkens by multiplying colors
//   - [OperatorScreen]: Lightens using inverted multiply
//   - [OperatorOverlay]: Multiplies or screens based on destination color
//   - [OperatorDifference]: Subtracts source from destination or vice versa
func (c *Context) SetOperator(op Operator) {
	c.withLock(func() {
		contextSetOperator(c.ptr, op)
	})
}

// GetOperator returns the current compositing operator for the context.
//
// The default operator is [OperatorOver].
func (c *Context) GetOperator() Operator {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return OperatorOver
	}
	return contextGetOperator(c.ptr)
}
