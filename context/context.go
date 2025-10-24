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

func (c *Context) Save() {
	c.Lock()
	defer c.Unlock()

	if c.ptr == nil {
		return
	}
	contextSave(c.ptr)
}

func (c *Context) Restore() {
	c.Lock()
	defer c.Unlock()

	if c.ptr == nil {
		return
	}
	contextRestore(c.ptr)
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
