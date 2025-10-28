package pattern

import (
	"runtime"
	"sync"
	"unsafe"

	"github.com/mikowitz/cairo/matrix"
	"github.com/mikowitz/cairo/status"
)

type Pattern interface {
	Close() error
	Status() status.Status
	SetMatrix(m *matrix.Matrix)
	GetMatrix() (*matrix.Matrix, error)
	Ptr() unsafe.Pointer
	GetType() PatternType
}

type BasePattern struct {
	sync.RWMutex
	ptr         PatternPtr
	closed      bool
	patternType PatternType
}

func newBasePattern(ptr PatternPtr, patternType PatternType) *BasePattern {
	b := &BasePattern{
		ptr:         ptr,
		patternType: patternType,
	}

	runtime.SetFinalizer(b, (*BasePattern).close)

	return b
}

func (b *BasePattern) Close() error {
	return b.close()
}

func (b *BasePattern) Status() status.Status {
	b.RLock()
	defer b.RUnlock()

	if b.ptr == nil {
		return status.NullPointer
	}
	return patternStatus(b.ptr)
}

func (b *BasePattern) SetMatrix(m *matrix.Matrix) {
	b.Lock()
	m.RLock()

	defer b.Unlock()
	defer m.RUnlock()

	patternSetMatrix(b.ptr, m.Ptr())
}

func (b *BasePattern) GetMatrix() (*matrix.Matrix, error) {
	b.RLock()
	defer b.RUnlock()

	if b.ptr == nil {
		return nil, status.NullPointer
	}

	return patternGetMatrix(b.ptr)
}

func (b *BasePattern) Ptr() unsafe.Pointer {
	b.RLock()
	defer b.RUnlock()

	return unsafe.Pointer(b.ptr)
}

func (b *BasePattern) GetType() PatternType {
	b.RLock()
	defer b.RUnlock()

	return b.patternType
}

func (b *BasePattern) close() error {
	b.Lock()
	defer b.Unlock()

	if b.ptr != nil {
		patternClose(b.ptr)
		runtime.SetFinalizer(b, nil)
		b.closed = true
		b.ptr = nil
	}

	return nil
}
