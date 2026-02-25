// ABOUTME: Re-exports PDFSurface type and NewPDFSurface constructor from the surface package.
// ABOUTME: Enables PDF surface usage through the root cairo package without a sub-package import.

//go:build !nopdf

package cairo

import "github.com/mikowitz/cairo/surface"

// PDFSurface is a surface that writes drawing operations to a PDF file.
// Dimensions are specified in points, where 1 point equals 1/72 of an inch.
// The coordinate origin is at the top-left corner of each page.
//
// Use NewPDFSurface to create a PDF surface. Call ShowPage to end one page
// and begin the next. Close the surface when finished to flush and finalize
// the PDF file.
//
// Requires Cairo's PDF backend (cairo-pdf pkg-config entry).
type PDFSurface = surface.PDFSurface

// NewPDFSurface creates a new PDF surface writing to filename.
// widthPt and heightPt set the dimensions of the first page in points (1/72 inch).
// Returns an error if Cairo cannot create the surface (e.g., invalid path).
//
// Requires the Cairo PDF backend. On Debian/Ubuntu: libcairo2-dev.
// On macOS: brew install cairo (includes PDF support by default).
func NewPDFSurface(filename string, widthPt, heightPt float64) (*PDFSurface, error) {
	return surface.NewPDFSurface(filename, widthPt, heightPt)
}
