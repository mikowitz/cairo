// ABOUTME: PDFSurface implementation for writing vector graphics to PDF files.
// ABOUTME: Dimensions are in points (1 point = 1/72 inch); origin is at top-left.

//go:build !nopdf

package surface

import "github.com/mikowitz/cairo/status"

// PDFSurface is a surface that writes drawing operations to a PDF file.
// Dimensions are specified in points, where 1 point equals 1/72 of an inch.
// The coordinate origin is at the top-left corner of each page.
//
// Use NewPDFSurface to create a PDF surface. Call ShowPage to end one page
// and begin the next. Close the surface when finished to flush and finalize
// the PDF file.
type PDFSurface struct {
	*BaseSurface
}

// NewPDFSurface creates a new PDF surface writing to filename.
// widthPt and heightPt set the dimensions of the first page in points (1/72 inch).
// Returns an error if Cairo cannot create the surface (e.g., invalid path).
func NewPDFSurface(filename string, widthPt, heightPt float64) (*PDFSurface, error) {
	ptr := pdfSurfaceCreate(filename, widthPt, heightPt)
	st := surfaceStatus(ptr)
	if st != status.Success {
		return nil, st
	}
	return &PDFSurface{BaseSurface: newBaseSurface(ptr)}, nil
}

// SetSize changes the page size for subsequent pages in the PDF document.
// widthPt and heightPt are in points (1/72 inch).
// This has no effect on already-emitted pages.
func (s *PDFSurface) SetSize(widthPt, heightPt float64) {
	s.Lock()
	defer s.Unlock()
	if s.ptr == nil {
		return
	}
	pdfSurfaceSetSize(s.ptr, widthPt, heightPt)
}

// ShowPage emits the current page and starts a new page in the PDF document.
// After calling ShowPage, subsequent drawing operations apply to the new page.
func (s *PDFSurface) ShowPage() {
	s.Lock()
	defer s.Unlock()
	if s.ptr == nil {
		return
	}
	surfaceShowPage(s.ptr)
}
