package surface

import (
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
	ptr    SurfacePtr
	closed bool
	sync.RWMutex
}
