package context

import "github.com/mikowitz/cairo/status"

type LineCap int

const (
	LineCapButt LineCap = iota
	LineCapRound
	LineCapSquare
)

func (c *Context) GetLineCap() LineCap {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return LineCapButt
	}
	return contextGetLineCap(c.ptr)
}

func (c *Context) SetLineCap(lineCap LineCap) {
	c.withLock(func() {
		contextSetLineCap(c.ptr, lineCap)
	})
}

type LineJoin int

const (
	LineJoinMiter LineJoin = iota
	LineJoinRound
	LineJoinBevel
)

func (c *Context) GetLineJoin() LineJoin {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return LineJoinMiter
	}
	return contextGetLineJoin(c.ptr)
}

func (c *Context) SetLineJoin(lineJoin LineJoin) {
	c.withLock(func() {
		contextSetLineJoin(c.ptr, lineJoin)
	})
}

func (c *Context) GetMiterLimit() float64 {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return 10
	}

	return contextGetMiterLimit(c.ptr)
}

func (c *Context) SetMiterLimit(limit float64) {
	c.withLock(func() {
		contextSetMiterLimit(c.ptr, limit)
	})
}

func (c *Context) SetDash(dashes []float64, offset float64) error {
	c.Lock()
	defer c.Unlock()

	if c.ptr == nil {
		return status.NullPointer
	}

	if len(dashes) > 0 && (allZeroes(dashes) || anyNegative(dashes)) {
		return status.InvalidDash
	}

	return contextSetDash(c.ptr, dashes, offset).ToError()
}

func anyNegative(s []float64) bool {
	for _, f := range s {
		if f < 0 {
			return true
		}
	}
	return false
}

func allZeroes(s []float64) bool {
	for _, f := range s {
		if f != 0 {
			return false
		}
	}
	return true
}

func (c *Context) GetDashCount() int {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return 0
	}
	return contextGetDashCount(c.ptr)
}

func (c *Context) GetDash() ([]float64, float64, error) {
	c.RLock()
	defer c.RUnlock()

	if c.ptr == nil {
		return []float64{}, 0, status.NullPointer
	}
	dashes, offset := contextGetDash(c.ptr)
	return dashes, offset, nil
}
