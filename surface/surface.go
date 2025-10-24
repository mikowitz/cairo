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

// WriteToPNG writes the contents of the surface to a new PNG file at the specified filepath.
//
// The surface should be flushed with Flush() before calling WriteToPNG to ensure all
// pending drawing operations are completed. This is particularly important if the surface
// has been modified through direct memory access or external APIs.
//
// The filepath is converted to a C string, so it will be truncated at the first null byte
// if present. Empty filepaths or paths to non-existent directories will result in an error.
//
// Returns an error if the surface is closed or if Cairo encounters an error writing the file
// (such as invalid path, insufficient permissions, or disk full).
//
// Example:
//
//	surf, err := NewImageSurface(FormatARGB32, 640, 480)
//	if err != nil {
//		return err
//	}
//	defer surf.Close()
//
//	// ... draw to surface ...
//
//	surf.Flush() // Ensure all drawing is complete
//	err = surf.WriteToPNG("output.png")
//	if err != nil {
//		return err
//	}
func (b *BaseSurface) WriteToPNG(filepath string) error {
	if b.ptr == nil {
		return status.NullPointer
	}
	return surfaceWriteToPNG(b.ptr, filepath)
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
