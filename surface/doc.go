/*
Package surface manages the abstract type representing all
different drawing targets that [cairo] can render to.
The actual drawings are performed using a cairo context.

A cairo surface is created by using backend-specific
constructors, typically of the form
cairo_backend_surface_create().

Most surface types allow accessing the surface without
using Cairo functions. If you do this, keep in mind
that it is mandatory that you call cairo_surface_flush()
before reading from or writing to the surface and
that you must use cairo_surface_mark_dirty()
after modifying it.

# Surface Types Overview

Cairo supports multiple surface types, each optimized for different use cases:

ImageSurface - Raster image buffer in memory

ImageSurface stores pixel data in memory and is the most commonly used surface
type. It supports various pixel formats (ARGB32, RGB24, A8, etc.) and can be
written to PNG files. Use ImageSurface when you need:
  - Pixel-level manipulation of image data
  - PNG export
  - A simple in-memory drawing target
  - Fast raster rendering

Example:
  surf, err := surface.NewImageSurface(surface.FormatARGB32, 640, 480)

PDFSurface - Vector output to PDF files (future)

PDFSurface writes vector graphics directly to PDF files, preserving scalability.
Text and graphics remain resolution-independent. Use PDFSurface when you need:
  - Scalable vector output
  - Multi-page documents
  - Print-quality output
  - Small file sizes for simple graphics

SVGSurface - Vector output to SVG files (future)

SVGSurface generates Scalable Vector Graphics files suitable for web use and
vector editing tools. Use SVGSurface when you need:
  - Web-compatible vector graphics
  - Editable vector output
  - CSS styling integration
  - Embedding in HTML

RecordingSurface - Capture and replay drawing operations (future)

RecordingSurface records all drawing operations for later playback to other
surfaces. Use RecordingSurface when you need:
  - Drawing operation recording
  - Render to multiple output formats
  - Drawing operation analysis
  - Deferred rendering

# Resource Management

All surfaces implement the Surface interface and must be properly closed to
release resources:

  surf, err := surface.NewImageSurface(surface.FormatARGB32, 400, 400)
  if err != nil {
      return err
  }
  defer surf.Close()

While finalizers provide a safety net, explicit Close() calls ensure deterministic
resource cleanup, which is particularly important when creating many surfaces or
in long-running applications.

# Thread Safety

All surface types are safe for concurrent use via embedded sync.RWMutex. However,
modifying a surface concurrently from multiple contexts may produce undefined
results. Best practice is to complete all drawing operations on a surface from
a single context before accessing it from another.
*/
package surface
