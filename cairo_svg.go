// ABOUTME: Re-exports SVGSurface type, SVGUnit type and constants, and NewSVGSurface
// ABOUTME: constructor from the surface package for use via the root cairo package.

//go:build !nosvg

package cairo

import "github.com/mikowitz/cairo/surface"

// SVGSurface is a surface that writes drawing operations to an SVG file.
// Dimensions are specified in points, where 1 point equals 1/72 of an inch.
// The coordinate origin is at the top-left corner of the image.
//
// Use NewSVGSurface to create an SVG surface. Close the surface when finished
// to flush and finalize the SVG file.
//
// Requires Cairo's SVG backend (cairo-svg pkg-config entry).
type SVGSurface = surface.SVGSurface

// SVGUnit specifies the unit for coordinates in an SVG document.
// Pass one of the SVGUnit constants to SVGSurface.SetDocumentUnit.
type SVGUnit = surface.SVGUnit

const (
	// SVGUnitUser uses user-space units (the default Cairo coordinate system).
	SVGUnitUser SVGUnit = surface.SVGUnitUser
	// SVGUnitEm uses em units, relative to the current font size.
	SVGUnitEm SVGUnit = surface.SVGUnitEm
	// SVGUnitEx uses ex units, relative to the x-height of the current font.
	SVGUnitEx SVGUnit = surface.SVGUnitEx
	// SVGUnitPx uses CSS pixel units (1px = 1/96 inch).
	SVGUnitPx SVGUnit = surface.SVGUnitPx
	// SVGUnitIn uses inch units.
	SVGUnitIn SVGUnit = surface.SVGUnitIn
	// SVGUnitCm uses centimeter units.
	SVGUnitCm SVGUnit = surface.SVGUnitCm
	// SVGUnitMm uses millimeter units.
	SVGUnitMm SVGUnit = surface.SVGUnitMm
	// SVGUnitPt uses point units (1pt = 1/72 inch).
	SVGUnitPt SVGUnit = surface.SVGUnitPt
	// SVGUnitPc uses pica units (1pc = 12pt).
	SVGUnitPc SVGUnit = surface.SVGUnitPc
	// SVGUnitPercent uses percentage units relative to the SVG viewport.
	SVGUnitPercent SVGUnit = surface.SVGUnitPercent
)

// NewSVGSurface creates a new SVG surface writing to filename.
// widthPt and heightPt set the dimensions in points (1/72 inch).
// Returns an error if Cairo cannot create the surface (e.g., invalid path).
//
// Requires the Cairo SVG backend. On Debian/Ubuntu: libcairo2-dev.
// On macOS: brew install cairo (includes SVG support by default).
func NewSVGSurface(filename string, widthPt, heightPt float64) (*SVGSurface, error) {
	return surface.NewSVGSurface(filename, widthPt, heightPt)
}
