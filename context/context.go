// Package context provides the main drawing context for Cairo operations.
//
// The Context is Cairo's central object for drawing operations. It maintains
// all graphics state parameters including the current transformation matrix,
// clip region, line width, line style, colors, font properties, and more.
//
// # Drawing Pipeline
//
// The typical Cairo drawing workflow follows this pattern:
//
//  1. Create a Surface (the drawing target)
//  2. Create a Context for that Surface
//  3. Set drawing properties (colors, line width, etc.)
//  4. Construct paths (MoveTo, LineTo, Rectangle, etc.)
//  5. Render paths (Fill, Stroke, Paint, etc.)
//  6. Close the Context and Surface when done
//
// # Lifecycle and Resource Management
//
// A Context must be created for a specific Surface and holds a reference to it.
// The Context should be explicitly closed with Close() when drawing is complete
// to release Cairo resources. A finalizer is registered as a safety net, but
// explicit cleanup is strongly recommended for long-running programs.
//
// Example usage:
//
//	surface, err := surface.NewImageSurface(surface.FormatARGB32, 400, 300)
//	if err != nil {
//	    return err
//	}
//	defer surface.Close()
//
//	ctx, err := context.NewContext(surface)
//	if err != nil {
//	    return err
//	}
//	defer ctx.Close()
//
//	// Drawing operations will be added in subsequent prompts
//	// ctx.SetSourceRGB(1.0, 0.0, 0.0)  // Red color
//	// ctx.Rectangle(50, 50, 100, 100)
//	// ctx.Fill()
//
// # State Management
//
// The Context maintains a stack of graphics states. Use Save() to push the
// current state onto the stack and Restore() to pop it back. This is useful
// for temporarily changing drawing parameters:
//
//	ctx.Save()
//	// Make temporary changes
//	ctx.Restore()  // Returns to previous state
//
// # Thread Safety
//
// Context is safe for concurrent use. All methods use appropriate locking
// to protect the internal state. However, for best performance, avoid
// concurrent drawing operations on the same Context.
package context

import (
	"runtime"
	"sync"
	"unsafe"

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
