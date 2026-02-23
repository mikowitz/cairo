// ABOUTME: Defines the Slant and Weight types for Cairo's toy font API.
// ABOUTME: These types control font style selection in SelectFontFace.

package font

// Slant specifies the slant style of a font face.
//
// The slant is used with [Weight] in SelectFontFace to select a font from
// the host platform's font system.
//
//go:generate stringer -type=Slant
type Slant int

// The iota values below must match Cairo's cairo_font_slant_t C enum exactly.
// Cairo has maintained this ordering since its initial release and documents
// it as stable. The CGO layer casts Slant directly to cairo_font_slant_t,
// so any divergence would silently produce incorrect font selection.
const (
	// SlantNormal selects an upright (non-slanted) font face.
	SlantNormal Slant = iota

	// SlantItalic selects an italic font face.
	SlantItalic

	// SlantOblique selects an oblique (mechanically slanted) font face.
	SlantOblique
)

// Weight specifies the weight (boldness) of a font face.
//
// The weight is used with [Slant] in SelectFontFace to select a font from
// the host platform's font system.
//
//go:generate stringer -type=Weight
type Weight int

// The iota values below must match Cairo's cairo_font_weight_t C enum exactly.
// Cairo has maintained this ordering since its initial release and documents
// it as stable. The CGO layer casts Weight directly to cairo_font_weight_t,
// so any divergence would silently produce incorrect font selection.
const (
	// WeightNormal selects a normal-weight font face.
	WeightNormal Weight = iota

	// WeightBold selects a bold font face.
	WeightBold
)
