package pattern

import "github.com/mikowitz/cairo/status"

type SolidPattern struct {
	*BasePattern
}

func NewSolidPatternRGB(r, g, b float64) (*SolidPattern, error) {
	ptr := patternCreateRGB(r, g, b)
	st := patternStatus(ptr)

	if st != status.Success {
		return nil, st
	}

	basePattern := newBasePattern(ptr, PatternTypeSolid)
	return &SolidPattern{
		BasePattern: basePattern,
	}, nil
}

func NewSolidPatternRGBA(r, g, b, a float64) (*SolidPattern, error) {
	ptr := patternCreateRGBA(r, g, b, a)
	st := patternStatus(ptr)

	if st != status.Success {
		return nil, st
	}

	basePattern := newBasePattern(ptr, PatternTypeSolid)
	return &SolidPattern{
		BasePattern: basePattern,
	}, nil
}
