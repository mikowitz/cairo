package context

import (
	"runtime"
	"sync"
	"unsafe"

	"github.com/mikowitz/cairo/pattern"
	"github.com/mikowitz/cairo/status"
	"github.com/mikowitz/cairo/surface"
)

type Context struct {
	sync.RWMutex
	ptr    ContextPtr
	closed bool
}

func NewContext(surface surface.Surface) (*Context, error) {
	if surface == nil || surface.Ptr() == nil {
		return nil, status.NullPointer
	}
	ptr := contextCreate((unsafe.Pointer)(surface.Ptr()))
	st := contextStatus(ptr)

	if st != status.Success {
		return nil, st
	}

	c := &Context{
		ptr: ptr,
	}

	runtime.SetFinalizer(c, (*Context).close)

	return c, nil
}

// Status checks whether an error has previously occurred for this context.
func (c *Context) Status() status.Status {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return status.NullPointer
	}

	return contextStatus(c.ptr)
}

func (c *Context) Close() error {
	return c.close()
}

// Save makes a copy of the current state of the context and saves it on an
// internal stack of saved states for the context. When [Context.Restore] is called,
// the context will be restored to the saved state. Multiple calls to Save
// and Restore can be nested; each call to Restore restores the state from
// the matching paired Save.
//
// It isn't necessary to clear all saved states before a [Context] is closed.
// If the reference count of a [Context] drops to zero in response to a call to
// [Context.Close], any saved states will be freed along with the [Context].
func (c *Context) Save() {
	c.Lock()
	defer c.Unlock()

	if c.ptr == nil {
		return
	}
	contextSave(c.ptr)
}

// Restores context to the state saved by a preceding call to [Context.Save]
// and removes that state from the stack of saved states.
func (c *Context) Restore() {
	c.Lock()
	defer c.Unlock()

	if c.ptr == nil {
		return
	}
	contextRestore(c.ptr)
}

// SetSourceRGB sets the source pattern within the context to an opaque color.
// This opaque color will then be used for any subsequent drawing operation
// until a new source pattern is set.
//
// The color components are floating point numbers in the range 0 to 1. If
// the values passed in are outside that range, they will be clamped.
//
// The default source pattern is opaque black, (that is, it is equivalent
// to context.SetSourceRGB(0, 0, 0)).
func (c *Context) SetSourceRGB(r, g, b float64) {
	c.Lock()
	defer c.Unlock()

	if c.ptr == nil {
		return
	}
	contextSetSourceRGB(c.ptr, r, g, b)
}

// SetSourceRGBA sets the source pattern within the context to a translucent
// color. This color will then be used for any subsequent drawing operation
// until a new source pattern is set.
//
// The color and alpha components are floating point numbers in the range
// 0 to 1. If the values passed in are outside that range, they will be
// clamped.
//
// Note that the color and alpha values are not premultiplied.
//
// The default source pattern is opaque black, (that is, it is equivalent
// to context.SetSourceRGBA(0, 0, 0, 1)).
func (c *Context) SetSourceRGBA(r, g, b, a float64) {
	c.Lock()
	defer c.Unlock()

	if c.ptr == nil {
		return
	}
	contextSetSourceRGBA(c.ptr, r, g, b, a)
}

func (c *Context) GetSource() (pattern.Pattern, error) {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return nil, status.NullPointer
	}

	return contextGetSource(c.ptr)
}

func (c *Context) SetSource(p pattern.Pattern) {
	c.withLock(func() {
		if p == nil {
			return
		}

		contextSetSource(c.ptr, p.Ptr())
	})
}

// MoveTo begins a new sub-path by setting the current point to (x, y).
//
// After this call the current point will be (x, y). Coordinates are specified
// in user-space, which is affected by the current transformation matrix (CTM).
//
// If there is no current path when MoveTo is called, this function behaves
// identically to calling NewPath() followed by MoveTo(x, y).
//
// Example:
//
//	ctx.MoveTo(50.0, 75.0)  // Start a path at (50, 75)
//	ctx.LineTo(100.0, 75.0) // Draw line to (100, 75)
func (c *Context) MoveTo(x, y float64) {
	c.withLock(func() {
		contextMoveTo(c.ptr, x, y)
	})
}

// LineTo adds a line segment to the path from the current point to (x, y),
// and sets the current point to (x, y).
//
// If there is no current point before the call to LineTo, this function will
// behave as if preceded by a call to MoveTo(x, y).
//
// After this call the current point will be (x, y). Coordinates are specified
// in user-space.
//
// Example:
//
//	ctx.MoveTo(10.0, 10.0)
//	ctx.LineTo(50.0, 10.0)  // Horizontal line
//	ctx.LineTo(50.0, 50.0)  // Vertical line
func (c *Context) LineTo(x, y float64) {
	c.withLock(func() {
		contextLineTo(c.ptr, x, y)
	})
}

// Rectangle adds a closed rectangular sub-path to the current path.
//
// The rectangle is positioned at (x, y) in user-space with the specified
// width and height. This is equivalent to:
//
//	ctx.MoveTo(x, y)
//	ctx.LineTo(x+width, y)
//	ctx.LineTo(x+width, y+height)
//	ctx.LineTo(x, y+height)
//	ctx.ClosePath()
//
// After calling Rectangle, the current point will be at (x, y).
//
// Example:
//
//	ctx.Rectangle(20.0, 30.0, 100.0, 50.0)  // Rectangle at (20,30) sized 100x50
//	ctx.SetSourceRGB(1.0, 0.0, 0.0)         // Red
//	ctx.Fill()                               // Fill the rectangle
func (c *Context) Rectangle(x, y, width, height float64) {
	c.withLock(func() {
		contextRectangle(c.ptr, x, y, width, height)
	})
}

// GetCurrentPoint returns the current point in user-space coordinates.
//
// The current point is the point that would be used by MoveTo or LineTo if
// called. Most path construction functions alter the current point. See the
// individual function documentation for details.
//
// Returns an error if there is no current point defined. Use HasCurrentPoint()
// to check whether a current point exists before calling this method.
//
// Example:
//
//	ctx.MoveTo(25.5, 37.75)
//	x, y, err := ctx.GetCurrentPoint()
//	if err == nil {
//	    fmt.Printf("Current point: (%f, %f)\n", x, y)
//	}
func (c *Context) GetCurrentPoint() (x, y float64, err error) {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return 0, 0, status.NullPointer
	}

	x, y, st := contextGetCurrentPoint(c.ptr)
	if st != nil {
		return 0, 0, st
	}
	return x, y, nil
}

// HasCurrentPoint returns whether a current point is defined on the current path.
//
// A current point is defined after operations like MoveTo, LineTo, or Rectangle.
// The current point becomes undefined after NewPath() or if no path operations
// have been performed yet.
//
// This is useful to check before calling GetCurrentPoint() to avoid errors.
//
// Example:
//
//	if ctx.HasCurrentPoint() {
//	    x, y, _ := ctx.GetCurrentPoint()
//	    fmt.Printf("Current point: (%f, %f)\n", x, y)
//	}
func (c *Context) HasCurrentPoint() bool {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return false
	}
	return contextHasCurrentPoint(c.ptr)
}

// NewPath clears the current path and removes any current point.
//
// After this call there will be no current path and no current point.
// This is typically called before starting a new path to ensure the path
// is empty. Without calling NewPath, subsequent path calls will append to
// the existing path.
//
// Example:
//
//	ctx.Rectangle(10, 10, 50, 50)
//	ctx.Fill()          // Fills the rectangle
//	ctx.NewPath()       // Clear the path to start fresh
//	ctx.Rectangle(70, 10, 50, 50)
//	ctx.Stroke()        // Strokes the second rectangle
func (c *Context) NewPath() {
	c.withLock(func() {
		contextNewPath(c.ptr)
	})
}

// ClosePath adds a line segment from the current point to the beginning of
// the current sub-path (the most recent point passed to MoveTo), and marks
// the current sub-path as closed.
//
// After calling ClosePath, the current point will be the beginning of the
// (now closed) sub-path.
//
// The behavior of ClosePath is distinct from simply calling LineTo with the
// coordinates of the sub-path's starting point. When a sub-path is closed,
// line joins are used at the connection point, whereas when it's merely
// connected with a line, line caps are used.
//
// If there is no current point before the call to ClosePath, this function
// has no effect.
//
// Example:
//
//	ctx.MoveTo(10, 10)
//	ctx.LineTo(50, 10)
//	ctx.LineTo(30, 40)
//	ctx.ClosePath()     // Completes the triangle back to (10, 10)
func (c *Context) ClosePath() {
	c.withLock(func() {
		contextClosePath(c.ptr)
	})
}

// NewSubPath begins a new sub-path within the current path.
//
// After this call there will be no current point. In many cases, calling
// NewSubPath is not required since MoveTo automatically begins new sub-paths.
// The primary use case is for creating compound paths with multiple disconnected
// shapes that share the same fill or stroke operation.
//
// NewSubPath is particularly useful when you want to ensure a new sub-path
// starts without implicitly closing or connecting to an existing sub-path.
//
// Example:
//
//	// Create two separate circles in the same path
//	ctx.Arc(50, 50, 20, 0, 2*math.Pi)
//	ctx.NewSubPath()
//	ctx.Arc(120, 50, 20, 0, 2*math.Pi)
//	ctx.Fill()  // Both circles filled with one operation
func (c *Context) NewSubPath() {
	c.withLock(func() {
		contextNewSubPath(c.ptr)
	})
}

// Stroke draws the current path by stroking it according to the current line
// width, line join, line cap, and dash settings. After Stroke, the current
// path will be cleared from the Context.
//
// The stroke operation draws along the outline of the path, with the line
// width determining how thick the line appears. The current source pattern
// determines the color or pattern used for the stroke.
//
// Note: Unlike StrokePreserve, this function clears the path after stroking.
// If you need to reuse the path for additional operations, use StrokePreserve
// instead.
//
// Example:
//
//	ctx.SetSourceRGB(0.0, 0.0, 1.0)  // Blue
//	ctx.SetLineWidth(5.0)
//	ctx.Rectangle(20, 20, 100, 100)
//	ctx.Stroke()  // Draws blue outline, path is now cleared
func (c *Context) Stroke() {
	c.withLock(func() {
		contextStroke(c.ptr)
	})
}

// StrokePreserve draws the current path by stroking it according to the current
// line width, line join, line cap, and dash settings. Unlike Stroke, this function
// preserves the path in the Context after stroking, allowing for additional
// operations on the same path.
//
// This is useful when you want to both stroke and fill the same path, or perform
// multiple operations with different settings on the same path.
//
// Example:
//
//	ctx.Rectangle(20, 20, 100, 100)
//	ctx.SetSourceRGB(1.0, 0.0, 0.0)  // Red fill
//	ctx.FillPreserve()
//	ctx.SetSourceRGB(0.0, 0.0, 0.0)  // Black outline
//	ctx.SetLineWidth(2.0)
//	ctx.StrokePreserve()  // Path still available after this
func (c *Context) StrokePreserve() {
	c.withLock(func() {
		contextStrokePreserve(c.ptr)
	})
}

// Fill fills the current path according to the current fill rule. After Fill,
// the current path will be cleared from the Context.
//
// The fill operation paints the interior of the path using the current source
// pattern. The fill rule (even-odd or winding) determines which areas are
// considered "inside" the path.
//
// Note: Unlike FillPreserve, this function clears the path after filling.
// If you need to reuse the path for additional operations, use FillPreserve
// instead.
//
// Example:
//
//	ctx.SetSourceRGB(1.0, 0.0, 0.0)  // Red
//	ctx.Rectangle(50, 50, 100, 100)
//	ctx.Fill()  // Fills rectangle with red, path is now cleared
func (c *Context) Fill() {
	c.withLock(func() {
		contextFill(c.ptr)
	})
}

// FillPreserve fills the current path according to the current fill rule.
// Unlike Fill, this function preserves the path in the Context after filling,
// allowing for additional operations on the same path.
//
// This is particularly useful when you want to both fill and stroke the same
// path with different colors or settings, which is a common pattern for creating
// shapes with both interior color and outline.
//
// Example:
//
//	ctx.Rectangle(50, 50, 100, 100)
//	ctx.SetSourceRGBA(0.0, 1.0, 0.0, 0.7)  // Semi-transparent green
//	ctx.FillPreserve()  // Fill the rectangle
//	ctx.SetSourceRGB(0.0, 0.0, 0.0)        // Black outline
//	ctx.SetLineWidth(2.0)
//	ctx.Stroke()  // Stroke the same rectangle
func (c *Context) FillPreserve() {
	c.withLock(func() {
		contextFillPreserve(c.ptr)
	})
}

// Paint paints the current source pattern everywhere within the current clip
// region. This is useful for setting a background color or pattern across
// the entire surface (or clipped region).
//
// Unlike Fill and Stroke, Paint does not use the current path. It simply
// applies the source pattern to all pixels in the current clip region.
//
// Example:
//
//	// Set white background
//	ctx.SetSourceRGB(1.0, 1.0, 1.0)
//	ctx.Paint()
//
//	// Now draw on top of the white background
//	ctx.SetSourceRGB(0.0, 0.0, 0.0)
//	ctx.Rectangle(50, 50, 100, 100)
//	ctx.Fill()
func (c *Context) Paint() {
	c.withLock(func() {
		contextPaint(c.ptr)
	})
}

// SetLineWidth sets the current line width for the Context. The line width
// value specifies the diameter of the pen used for stroking paths, in user-space
// units.
//
// The line width affects all subsequent Stroke and StrokePreserve operations.
// The default line width is 2.0.
//
// Note: The line width is transformed by the current transformation matrix (CTM),
// so scaling transformations will affect the actual rendered line width. However,
// the width specified here is always in user-space coordinates before
// transformation.
//
// Example:
//
//	ctx.SetLineWidth(5.0)  // Set thick line
//	ctx.MoveTo(10, 10)
//	ctx.LineTo(100, 100)
//	ctx.Stroke()  // Draws 5-pixel wide line
//
//	ctx.SetLineWidth(1.0)  // Set thin line
//	ctx.MoveTo(10, 20)
//	ctx.LineTo(100, 110)
//	ctx.Stroke()  // Draws 1-pixel wide line
func (c *Context) SetLineWidth(width float64) {
	c.withLock(func() {
		contextSetLineWidth(c.ptr, width)
	})
}

// GetLineWidth returns the current line width value for the Context. The line
// width represents the diameter of the pen used for stroking paths, specified
// in user-space units.
//
// This method is thread-safe and can be called concurrently with other read
// operations. It returns the width that was most recently set with SetLineWidth,
// or the default value of 2.0 if SetLineWidth has not been called.
//
// If called on a closed context, this method returns 0.0.
//
// Example:
//
//	ctx.SetLineWidth(5.0)
//	width := ctx.GetLineWidth()  // Returns 5.0
func (c *Context) GetLineWidth() float64 {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return 0.0
	}

	return contextGetLineWidth(c.ptr)
}

func (c *Context) close() error {
	c.Lock()
	defer c.Unlock()

	if c.ptr != nil {
		contextClose(c.ptr)
		runtime.SetFinalizer(c, nil)
		c.closed = true
		c.ptr = nil
	}

	return nil
}

func (c *Context) withLock(fn func()) {
	c.Lock()
	defer c.Unlock()

	if c.ptr == nil {
		return
	}

	fn()
}
