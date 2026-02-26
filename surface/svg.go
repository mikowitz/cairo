// ABOUTME: SVGSurface implementation for writing vector graphics to SVG files.
// ABOUTME: Dimensions are in points (1 point = 1/72 inch); origin is at top-left.

//go:build !nosvg

package surface

import "github.com/mikowitz/cairo/status"

// SVGUnit specifies the unit for coordinates in an SVG document.
// These values correspond directly to Cairo's cairo_svg_unit_t enum.
//
//go:generate stringer -type=SVGUnit -tags "!nosvg"
type SVGUnit int

const (
	// SVGUnitUser uses user-space units (the default Cairo coordinate system).
	SVGUnitUser SVGUnit = iota
	// SVGUnitEm uses em units, relative to the current font size.
	SVGUnitEm
	// SVGUnitEx uses ex units, relative to the x-height of the current font.
	SVGUnitEx
	// SVGUnitPx uses CSS pixel units (1px = 1/96 inch).
	SVGUnitPx
	// SVGUnitIn uses inch units.
	SVGUnitIn
	// SVGUnitCm uses centimeter units.
	SVGUnitCm
	// SVGUnitMm uses millimeter units.
	SVGUnitMm
	// SVGUnitPt uses point units (1pt = 1/72 inch).
	SVGUnitPt
	// SVGUnitPc uses pica units (1pc = 12pt).
	SVGUnitPc
	// SVGUnitPercent uses percentage units relative to the SVG viewport.
	SVGUnitPercent
)

// SVGSurface is a surface that writes drawing operations to an SVG file.
// Dimensions are specified in points, where 1 point equals 1/72 of an inch.
// The coordinate origin is at the top-left corner of the image.
//
// Use NewSVGSurface to create an SVG surface. Close the surface when finished
// to flush and finalize the SVG file.
type SVGSurface struct {
	*BaseSurface
}

// NewSVGSurface creates a new SVG surface writing to filename.
// widthPt and heightPt set the dimensions in points (1/72 inch).
// Returns an error if Cairo cannot create the surface (e.g., invalid path).
func NewSVGSurface(filename string, widthPt, heightPt float64) (*SVGSurface, error) {
	ptr := svgSurfaceCreate(filename, widthPt, heightPt)
	st := surfaceStatus(ptr)
	if st != status.Success {
		surfaceClose(ptr)
		return nil, st
	}
	return &SVGSurface{BaseSurface: newBaseSurface(ptr)}, nil
}

// SetDocumentUnit sets the unit used for coordinates in the SVG document.
// This controls how SVG renderers interpret coordinate values in the output.
// Must be called before drawing operations to take effect on the whole document.
func (s *SVGSurface) SetDocumentUnit(unit SVGUnit) {
	s.Lock()
	defer s.Unlock()
	if s.ptr == nil {
		return
	}
	svgSurfaceSetDocumentUnit(s.ptr, unit)
}
