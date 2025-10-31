package pattern

import "github.com/mikowitz/cairo/status"

type Gradient interface {
	AddColorStopRGB(offset, r, g, b float64)
	AddColorStopRGBA(offset, r, g, b, a float64)
	GetColorStopCount() (int, error)
	GetColorStopRGBA(index int) (float64, float64, float64, float64, float64, error)
}

type BaseGradient struct {
	*BasePattern
}

func (bg *BaseGradient) AddColorStopRGB(offset, r, g, b float64) {
	bg.Lock()
	defer bg.Unlock()

	patternAddColorStopRGB(bg.ptr, offset, r, g, b)
}

func (bg *BaseGradient) AddColorStopRGBA(offset, r, g, b, a float64) {
	bg.Lock()
	defer bg.Unlock()

	patternAddColorStopRGBA(bg.ptr, offset, r, g, b, a)
}

func (bg *BaseGradient) GetColorStopCount() (int, error) {
	count, st := patternGetColorStopCount(bg.ptr)

	if st != status.Success {
		return count, st
	}
	return count, nil
}

func (bg *BaseGradient) GetColorStopRGBA(index int) (float64, float64, float64, float64, float64, error) {
	o, r, g, b, a, st := patternGetColorStopRGBA(bg.ptr, index)

	if st != status.Success {
		return o, r, g, b, a, st
	}

	return o, r, g, b, a, nil
}

type LinearGradient struct {
	*BaseGradient
}

func NewLinearGradient(x0, y0, x1, y1 float64) (*LinearGradient, error) {
	ptr := patternCreateLinear(x0, y0, x1, y1)
	st := patternStatus(ptr)

	if st != status.Success {
		return nil, st
	}

	basePattern := newBasePattern(ptr, PatternTypeLinear)
	baseGradient := &BaseGradient{BasePattern: basePattern}
	return &LinearGradient{
		BaseGradient: baseGradient,
	}, nil
}

type RadialGradient struct {
	*BaseGradient
}

func NewRadialGradient(cx0, cy0, radius0, cx1, cy1, radius1 float64) (*RadialGradient, error) {
	ptr := patternCreateRadial(cx0, cy0, radius0, cx1, cy1, radius1)
	st := patternStatus(ptr)

	if st != status.Success {
		return nil, st
	}

	basePattern := newBasePattern(ptr, PatternTypeRadial)
	baseGradient := &BaseGradient{BasePattern: basePattern}
	return &RadialGradient{
		BaseGradient: baseGradient,
	}, nil
}
