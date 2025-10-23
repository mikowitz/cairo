package surface

import (
	"runtime"
	"sync"

	"github.com/mikowitz/cairo/status"
)

type Surface interface {
	Close() error
	Status() status.Status
	Flush()
	MarkDirty()
	MarkDirtyRectangle(x, y, width, height int)
}

type BaseSurface struct {
	sync.RWMutex
	closed bool
	ptr    SurfacePtr
}

func newBaseSurface(ptr SurfacePtr) *BaseSurface {
	b := &BaseSurface{
		ptr: ptr,
	}

	runtime.SetFinalizer(b, (*BaseSurface).close)

	return b
}

// Close closes the surface, ensuring that any reserved memory related to
// the surface is cleaned up. Any subsequent methods called on the surface
// after closing will have no effect on the surface's state.
func (b *BaseSurface) Close() error {
	return b.close()
}

// Status checks whether an error has previously occurred for this surface.
func (b *BaseSurface) Status() status.Status {
	b.RLock()
	defer b.RUnlock()

	if b.ptr == nil {
		return status.NullPointer
	}
	return surfaceStatus(b.ptr)
}

// Flush does any pending drawing for the surface and also restores any
// temporary modifications [cairo] has made to the surface's state. This
// function must be called before switching from drawing on the surface
// with cairo to drawing on it directly with native APIs, or accessing
// its memory outside of Cairo. If the surface doesn't support direct
// access, then this function does nothing.
func (b *BaseSurface) Flush() {
	b.Lock()
	defer b.Unlock()

	if b.ptr == nil {
		return
	}
	surfaceFlush(b.ptr)
}

// MarkDirty tells cairo that drawing has been done to surface using means
// other than cairo, and that cairo should reread any cached areas. Note that
// you must call [Flush] before doing such drawing.
func (b *BaseSurface) MarkDirty() {
	b.Lock()
	defer b.Unlock()

	if b.ptr == nil {
		return
	}
	surfaceMarkDirty(b.ptr)
}

// MarkDirtyRectangle is like [MarkDirty], but drawing has been done only to
// the specified rectangle, so that cairo can retain cached contents for other
// parts of the surface.
//
// Any cached clip set on the surface will be reset by this function, to make
// sure that future cairo calls have the clip set that they expect.
func (b *BaseSurface) MarkDirtyRectangle(x, y, width, height int) {
	b.Lock()
	defer b.Unlock()

	if b.ptr == nil {
		return
	}
	surfaceMarkDirtyRectangle(b.ptr, x, y, width, height)
}

func (b *BaseSurface) close() error {
	b.Lock()
	defer b.Unlock()

	if b.ptr != nil {
		surfaceClose(b.ptr)
		runtime.SetFinalizer(b, nil)
		b.closed = true
		b.ptr = nil
	}

	return nil
}
