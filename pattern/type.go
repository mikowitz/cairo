package pattern

//go:generate stringer -type=PatternType -trimprefix=PatternType
type PatternType int

const (
	PatternTypeSolid PatternType = iota
	PatternTypeSurface
	PatternTypeLinear
	PatternTypeRadial
	PatternTypeMesh
	PatternTypeRasterSource
)
