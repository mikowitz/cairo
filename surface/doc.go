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
*/
package surface
