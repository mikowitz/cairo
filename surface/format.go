package surface

//go:generate stringer -type=Format -trimprefix=Format
type Format int

const (
	FormatARGB32 Format = iota
	FormatRGB24
	FormatA8
	FormatA1
	FormatRGB16_565
	FormatRGB30

	FormatInvalid Format = -1
)

func (f Format) StrideForWidth(width int) int {
	return formatStrideForWidth(f, width)
}
